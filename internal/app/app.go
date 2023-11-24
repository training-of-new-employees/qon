package app

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type application struct {
	address string
	server  *http.Server
	handler http.Handler
	logger  *slog.Logger
}

type Config struct {
	Address string
}

func New(cfg Config) *application {
	app := &application{
		address: cfg.Address,
		logger:  slog.Default(),
	}

	app.init()

	app.server = &http.Server{
		Addr:    app.address,
		Handler: app.handler,
	}

	return app
}

func (app *application) Start() {
	app.server.ListenAndServe()
}

func (app *application) stop() {
	app.server.Shutdown(context.Background())
}

func (app *application) init() {
	route := gin.Default()
	pool := pgxpool.New()
	// TODO: Здесь инициализируем все пакеты приложения
	// u := user.New(user.Config{
	// 	Logger: app.logger,
	// 	Route:  route,
	// 	Pool: pool,
	// })


	app.handler = route
}
