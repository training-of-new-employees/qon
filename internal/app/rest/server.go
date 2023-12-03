// Package rest - пакет с реализацией rest-сервера.
package rest

import (
	"github.com/training-of-new-employees/qon/internal/pkg/jwttoken"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RestServer - реализация rest-сервера.
type RestServer struct {
	router    *gin.Engine
	secretKey string
	tokenVal  jwttoken.JWTValidator
}

// New - конструктор для RestServer.
func New(secretKey string) *RestServer {
	gin.SetMode(gin.ReleaseMode)

	s := &RestServer{
		router:    gin.New(),
		secretKey: secretKey,
		tokenVal:  jwttoken.NewTokenValidator(secretKey),
	}

	s.InitRoutes()

	return s
}

// ServeHTTP используется для реализации http.Handler интерфейса.
// Замечание: необходимо для совместимости с net/http пакетом.
func (s *RestServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
