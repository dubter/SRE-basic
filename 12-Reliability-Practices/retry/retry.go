package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func GetData(url string) (string, error) {
	maxRetries := 3
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		resp, err := http.Get(url)
		if err != nil {
			return "", fmt.Errorf("HTTP request failed: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return "", fmt.Errorf("failed to read response body: %w", err)
			}
			return string(body), nil
		}

		if isRetryable(resp.StatusCode) {
			lastErr = fmt.Errorf("retryable status code: %d", resp.StatusCode)
			if i < maxRetries-1 {
				delay := time.Duration(1<<i) * time.Second
				time.Sleep(delay)
			}
		} else {
			return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}
	}

	return "", fmt.Errorf("after %d attempts: %w", maxRetries, lastErr)
}

func isRetryable(statusCode int) bool {
	switch statusCode {
	case 500, 502, 503, 504:
		return true
	default:
		return false
	}
}
