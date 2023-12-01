// Package config - пакет для работы с конфигурацией приложения.
package config

import (
	"flag"
	"os"
)

var (
	runAddr string
	dsn     string
)

// Config хранит настройки приложения.
type Config struct {
	Address     string
	DatabaseDSN string
}

// InitConfig определяет настройки приложения по флагам, переменным окружения.
func InitConfig() *Config {
	// Флаги
	flag.StringVar(&runAddr, "a", defaultRunAddr, "address and port to run server")
	flag.StringVar(&dsn, "d", defaultDSN, "db address")
	// NOTE: здесь определяем последующие флаги
	// ...

	flag.Parse()

	// Переменные окружения (ENV)
	if envRunAddr := os.Getenv("RUN_ADDR"); envRunAddr != "" {
		runAddr = envRunAddr
	}

	if envDatabaseDSN := os.Getenv("DATABASE_DSN"); envDatabaseDSN != "" {
		dsn = envDatabaseDSN
	}
	// NOTE: здесь определяем последующие ENV
	// ...

	// Определение конфига
	config := &Config{
		Address:     runAddr,
		DatabaseDSN: dsn,
	}

	return config
}
