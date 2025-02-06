package pkg

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

const (
	openTimeout     = 10 * time.Second
	maxFailureCount = 3
)

type State int

const (
	Closed State = iota
	Open
	HalfOpen
)

type CircuitBreaker struct {
	state        State
	failureCount int
	openTimeout  time.Time
	mutex        sync.Mutex
	lastRequest  time.Time
}

var cb = &CircuitBreaker{state: Closed}

func GetDataWithCircuitBreaker(url string) (string, error) {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	// Проверка состояния
	switch cb.state {
	case Open:
		if time.Now().Before(cb.openTimeout) {
			return "", fmt.Errorf("circuit breaker is OPEN")
		}
		cb.state = HalfOpen
		fmt.Println("Circuit Breaker: OPEN → HALF-OPEN")

	case HalfOpen:
		if time.Since(cb.lastRequest) < 10*time.Second {
			return "", fmt.Errorf("circuit breaker is HALF-OPEN (test request pending)")
		}
	}

	// Выполнение запроса
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	// Обработка ответа
	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("failed to read response body: %v", err)
		}

		// Успешный запрос
		if cb.state == HalfOpen {
			fmt.Println("Circuit Breaker: HALF-OPEN → CLOSED (success)")
			cb.state = Closed
			cb.failureCount = 0
		}
		return string(body), nil
	}

	// Обработка ошибок
	if isRetryable(resp.StatusCode) {
		cb.failureCount++
		if cb.failureCount >= maxFailureCount {
			cb.state = Open
			cb.openTimeout = time.Now().Add(openTimeout)
			fmt.Println("Circuit Breaker: CLOSED → OPEN (3 failures)")
		}
	} else {
		cb.failureCount = 0
	}

	return "", fmt.Errorf("status code: %d", resp.StatusCode)
}

func isRetryable(statusCode int) bool {
	switch statusCode {
	case 500, 502, 503, 504:
		return true
	default:
		return false
	}
}
