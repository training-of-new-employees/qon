// Package handlers - пакет содержит хендлеры.
package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/training-of-new-employees/qon/internal/app/rest/middlewares"
)

// Handler - HTTP-хендлеры для API.
type Handlers struct {
}

func NewHandler() *Handlers {
	return &Handlers{}
}

// InitRoutes инициализирует handler-ы.
func (h *Handlers) InitRoutes() *gin.Engine {
	router := gin.New()

	router.Use(middlewares.DummyMiddleware())

	router.GET("/api/v1/dummy", h.Dummy)

	return router
}
