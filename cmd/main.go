package main

import (
	"log"

	"github.com/training-of-new-employees/qon/internal/app"
	"github.com/training-of-new-employees/qon/internal/config"
)

func main() {
	// Инициализация настроек приложения
	cfg := config.InitConfig()

	// Запускаем приложение
	if err := app.StartApp(cfg); err != nil {
		log.Fatal(err)
	}
}
