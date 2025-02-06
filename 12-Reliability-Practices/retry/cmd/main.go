package main

import (
	"fmt"
	"log"

	"retry/pkg"
)

func main() {
	// Пример 1: Успешный запрос с телом
	fmt.Println("Пример 1: Успешный запрос с телом")
	data, err := pkg.GetData("https://httpbin.org/get")
	if err != nil {
		log.Printf("Ошибка: %v\n", err)
	} else {
		log.Println("Code: 200")
		if len(data) > 0 {
			log.Printf("Данные:\n%s\n", data)
		}
	}

	// Пример 2: Сервер возвращает 500 → 500 → 200 (с телом)
	fmt.Println("\nПример 2: Retry 500 → 500 → 200")
	data, err = pkg.GetData("https://httpbin.org/status/500,500,200")
	if err != nil {
		log.Printf("Ошибка: %v\n", err)
	} else {
		log.Println("Code: 200")
		if len(data) > 0 {
			log.Printf("Данные:\n%s\n", data)
		}
	}

	// Пример 3: Все попытки неудачны (500 → 500 → 500)
	fmt.Println("\nПример 3: Все попытки неудачны")
	data, err = pkg.GetData("https://httpbin.org/status/500,500,500")
	if err != nil {
		log.Printf("Ошибка: %v\n", err)
	}

	// Пример 4: Неповторяемая ошибка (404)
	fmt.Println("\nПример 4: Неповторяемая ошибка")
	data, err = pkg.GetData("https://httpbin.org/status/404")
	if err != nil {
		log.Printf("Ошибка: %v\n", err)
	}
}
