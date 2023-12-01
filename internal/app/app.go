// Package app - пакет для запуска приложения.
package app

import (
	"fmt"
	"net/http"

	"github.com/training-of-new-employees/qon/internal/app/rest"
	"github.com/training-of-new-employees/qon/internal/config"
	"github.com/training-of-new-employees/qon/internal/logger"
	"github.com/training-of-new-employees/qon/internal/store/pg"
	"go.uber.org/zap"
)

// StartApp запускает приложение.
func StartApp(cfg *config.Config) error {
	// Инициализация логгера (logger.Log)
	if err := logger.InitLogger(cfg.LogLevel); err != nil {
		return err
	}

	// Создаём хранилище
	store, err := pg.NewStore(cfg.DatabaseDSN)
	if err != nil {
		return err
	}
	defer store.Close()

	// Создаём сервер
	server := rest.New()

	app := &http.Server{
		Handler: server,
		Addr:    cfg.Address,
	}

	logger.Log.Info(fmt.Sprintf("Running server with log level '%s'", cfg.LogLevel), zap.String("address", cfg.Address))

	return app.ListenAndServe()
}
