package config

// Настройки по умолчанию.
var (
	// адрес и порт для запуска сервера.
	defaultRunAddr = "127.0.0.1:8080"
	// адрес PostgreSQL - подробнее о контейнере в файле docker-compose/dev/docker-compose.yml
	defaultDSN = "postgres://quickon:quickon@localhost:15438/qon_dev"
)
