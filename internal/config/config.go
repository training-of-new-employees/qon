// Package config - пакет для работы с конфигурацией приложения.
package config

import (
	"flag"
	"os"
)

var (
	runAddr string
)

// Config хранит настройки приложения.
type Config struct {
	Address string
}

// InitConfig определяет настройки приложения по флагам, переменным окружения.
func InitConfig() *Config {
	// Флаги
	flag.StringVar(&runAddr, "a", defaultRunAddr, "address and port to run server")
	// TODO: здесь определяем последующие флаги
	// ...

	flag.Parse()

	// Переменные окружения (ENV)
	if envRunAddr := os.Getenv("RUN_ADDR"); envRunAddr != "" {
		runAddr = envRunAddr
	}
	// TODO: здесь определяем последующие ENV
	// ...

	// Определение конфига
	config := &Config{
		Address: runAddr,
	}

	return config
}
