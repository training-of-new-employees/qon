// Package config - пакет для работы с конфигурацией приложения.
package config

import (
	"flag"
	"os"
	"time"
)

var (
	logLevel       string
	runAddr        string
	dsn            string
	jwtSecretKey   string
	redisDSN       string
	redisPassword  string
	redisDB        int
	senderEmail    string
	senderPassword string
)

// Config хранит настройки приложения.
type Config struct {
	LogLevel            string
	Address             string
	DatabaseDSN         string
	SecretKey           string
	RedisDSN            string
	RedisPassword       string
	RedisDB             int
	AccessTokenExpires  time.Duration
	RefreshTokenExpires time.Duration
	SenderEmail         string
	SenderPassword      string
}

// InitConfig определяет настройки приложения по флагам, переменным окружения.
func InitConfig() *Config {
	// Флаги
	flag.StringVar(&logLevel, "l", defaultLogLevel, "log level")
	flag.StringVar(&runAddr, "a", defaultRunAddr, "address and port to run server")
	flag.StringVar(&dsn, "d", defaultDSN, "db address")
	flag.StringVar(&jwtSecretKey, "j", defaultSecretKey, "jwt secret key")
	flag.StringVar(&redisDSN, "r", defaultRedisDSN, "cacheredis address")
	flag.StringVar(&redisPassword, "rp", defaultRedisPassword, "cacheredis password")
	flag.IntVar(&redisDB, "rd", defaultRedisDB, "cacheredis db")
	// NOTE: здесь определяем последующие флаги
	// ...

	flag.StringVar(&senderEmail, "se", defaultSenderEmail, "sender email")
	flag.StringVar(&senderPassword, "sp", defaultSenderPassword, "sender password")

	flag.Parse()

	// Переменные окружения (ENV)
	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		logLevel = envLogLevel
	}

	if envRunAddr := os.Getenv("RUN_ADDR"); envRunAddr != "" {
		runAddr = envRunAddr
	}

	if envDatabaseDSN := os.Getenv("DATABASE_DSN"); envDatabaseDSN != "" {
		dsn = envDatabaseDSN
	}

	if envJWT := os.Getenv("JWT_SECRET_KEY"); envJWT != "" {
		jwtSecretKey = envJWT
	}

	if envRedisDSN := os.Getenv("REDIS_DSN"); envRedisDSN != "" {
		redisDSN = envRedisDSN
	}

	if envRedisPassword := os.Getenv("REDIS_PASSWORD"); envRedisPassword != "" {
		redisPassword = envRedisPassword
	}

	if envSenderEmail := os.Getenv("SENDER_EMAIL"); envSenderEmail != "" {
		senderEmail = envSenderEmail
	}

	if envSenderPassword := os.Getenv("SENDER_PASSWORD"); envSenderPassword != "" {
		senderPassword = envSenderPassword
	}

	// NOTE: здесь определяем последующие ENV
	// ...

	// Определение конфига
	config := &Config{
		LogLevel:            logLevel,
		Address:             runAddr,
		DatabaseDSN:         dsn,
		SecretKey:           jwtSecretKey,
		RedisDSN:            redisDSN,
		RedisPassword:       redisPassword,
		RedisDB:             0,
		AccessTokenExpires:  time.Minute * 15,
		RefreshTokenExpires: time.Hour * 120,
		SenderEmail:         senderEmail,
		SenderPassword:      senderPassword,
	}

	return config
}
