package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/training-of-new-employees/qon/internal/app/rest"
	"github.com/training-of-new-employees/qon/internal/config"
	"github.com/training-of-new-employees/qon/internal/logger"
	"github.com/training-of-new-employees/qon/internal/pkg/doar"
	"github.com/training-of-new-employees/qon/internal/service/impl"
	"github.com/training-of-new-employees/qon/internal/store/cache/cacheredis"
	"github.com/training-of-new-employees/qon/internal/store/pg"
)

func main() {
	// Запускаем приложение
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	// Инициализация настроек приложения
	cfg := config.InitConfig()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGHUP, syscall.SIGTERM)
	defer cancel()

	// Инициализация логгера (logger.Log)
	if err := logger.InitLogger(cfg.LogLevel); err != nil {
		return err
	}

	logger.Log.Debug("Config Initialized", zap.Any("config", cfg))

	// Создаём хранилище
	store, err := pg.NewStore(cfg.DatabaseDSN)
	if err != nil {
		return err
	}

	defer func() {
		if err = store.Close(); err != nil {
			logger.Log.Warn("closing store: %v", zap.Error(err))
		}
	}()

	clientRedis := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisDSN,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	logger.Log.Info("Redis up")
	status := clientRedis.Ping(ctx)
	logger.Log.Info("Status up")
	if status.Err() != nil {
		logger.Log.Warn("cacheredis ping: %v", zap.Error(err))
		return status.Err()
	}

	redis := cacheredis.NewRedis(clientRedis)
	logger.Log.Info("Redis up")
	sender := doar.NewSender(
		&doar.SenderConfig{
			Mode:           doar.SenderMode(cfg.SenderMode),
			SenderEmail:    cfg.SenderMode,
			SenderPassword: cfg.SenderPassword,
			SenderApiKey:   cfg.SenderApiKey,
		},
	)

	services := impl.NewServices(store, redis, cfg.SecretKey, cfg.AccessTokenExpires, cfg.RefreshTokenExpires, sender, cfg.Domain)
	// Создаём сервер
	server := rest.New(cfg.SecretKey, services, redis)

	app := &http.Server{
		Handler: server,
		Addr:    cfg.Address,
	}
	logger.Log.Info(fmt.Sprintf("Running server with log level '%s'", cfg.LogLevel), zap.String("address", cfg.Address))

	return app.ListenAndServe()
}
