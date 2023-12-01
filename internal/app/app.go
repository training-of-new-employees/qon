// Package app - пакет для запуска приложения.
package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/training-of-new-employees/qon/internal/app/rest"
)

// TODO: перенести, после реализации пакета config.
var defaultRunAddr = "127.0.0.1:8080"

// StartApp запускает приложение.
func StartApp() error {
	srv := &http.Server{
		Handler: rest.New(),
		Addr:    defaultRunAddr,
	}

	log.Println(fmt.Sprintf("Running server on '%s'", defaultRunAddr))

	return srv.ListenAndServe()
}
