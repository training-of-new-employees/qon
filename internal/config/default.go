package config

import "go.uber.org/zap"

// Настройки по умолчанию.
var (
	// Уровень логирования Debug
	defaultLogLevel = zap.DebugLevel.CapitalString()
	// Адрес и порт для запуска сервера.
	defaultRunAddr = "127.0.0.1:8080"
	// Адрес БД PostgreSQL - подробнее о контейнере в файле docker-compose/dev/docker-compose.yml
	defaultDSN       = "postgres://quickon:quickon@localhost:15438/qon_dev"
	defaultSecretKey = "simtahtonim"
)
