// Package app - пакет для запуска приложения.
package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/training-of-new-employees/qon/internal/app/rest"
	"github.com/training-of-new-employees/qon/internal/config"
)

// StartApp запускает приложение.
func StartApp(cfg *config.Config) error {
	srv := &http.Server{
		Handler: rest.New(),
		Addr:    cfg.Address,
	}

	log.Println(fmt.Sprintf("Running server on '%s'", cfg.Address))

	return srv.ListenAndServe()
}
