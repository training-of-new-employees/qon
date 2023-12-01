package main

import (
	"log"

	"github.com/training-of-new-employees/qon/internal/app"
)

func main() {
	// Запускаем приложение
	if err := app.StartApp(); err != nil {
		log.Fatal(err)
	}
}
