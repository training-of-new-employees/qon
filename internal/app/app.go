// Package app - пакет для запуска приложения.
package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/training-of-new-employees/qon/internal/app/rest"
	"github.com/training-of-new-employees/qon/internal/config"
	"github.com/training-of-new-employees/qon/internal/store/pg"
)

// StartApp запускает приложение.
func StartApp(cfg *config.Config) error {
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

	log.Println(fmt.Sprintf("Running server on '%s'", cfg.Address))

	return app.ListenAndServe()
}
