package main

import (
	"fmt"
	"log"
	"time"

	"circuit-breaker/pkg"
)

func main() {
	// Тестовая последовательность: 3 ошибки → успех
	urls := []string{
		"https://httpbin.org/status/500",
		"https://httpbin.org/status/502",
		"https://httpbin.org/status/503",
		"https://httpbin.org/status/200",
	}

	for _, url := range urls {
		fmt.Printf("\nRequest: %s\n", url)
		data, err := pkg.GetDataWithCircuitBreaker(url)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			log.Println("Code: 200")
			if len(data) > 0 {
				log.Printf("Respose:\n%s\n", data)
			}
		}
		time.Sleep(1 * time.Second)
	}

	// Демонстрация Half-Open → Closed
	fmt.Println("\nWaiting 11 seconds...")
	time.Sleep(11 * time.Second)
	data, err := pkg.GetDataWithCircuitBreaker("https://httpbin.org/status/200")
	fmt.Printf("\nAfter timeout:\n")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		log.Println("Code: 200")
		if len(data) > 0 {
			log.Printf("Respose:\n%s", data)
		}
	}
}
