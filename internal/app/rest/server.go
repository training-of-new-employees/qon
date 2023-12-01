// Package rest - пакет с реализацией rest-сервера.
package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/training-of-new-employees/qon/internal/app/rest/handlers"
)

// RestServer - реализация rest-сервера.
type RestServer struct {
	router *gin.Engine
}

// New - конструктор для RestServer.
func New() *RestServer {
	gin.SetMode(gin.ReleaseMode)
	
	s := &RestServer{
		router: handlers.NewHandler().InitRoutes(),
	}

	return s
}

// ServeHTTP используется для реализации http.Handler интерфейса.
// Замечание: необходимо для совместимости с net/http пакетом.
func (s *RestServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
