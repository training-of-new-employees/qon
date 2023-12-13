// Package rest - пакет с реализацией rest-сервера.
package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/training-of-new-employees/qon/internal/pkg/jwttoken"
	"github.com/training-of-new-employees/qon/internal/service"
)

// RestServer - реализация rest-сервера.
type RestServer struct {
	router    *gin.Engine
	secretKey string
	tokenVal  jwttoken.JWTValidator
	services  service.Service
}

// New - конструктор для RestServer.
func New(secretKey string, services service.Service) *RestServer {
	gin.SetMode(gin.ReleaseMode)

	s := &RestServer{
		router:    gin.New(),
		secretKey: secretKey,
		tokenVal:  jwttoken.NewTokenValidator(secretKey),
		services:  services,
	}

	s.InitRoutes()

	return s
}

// ServeHTTP используется для реализации http.Handler интерфейса.
// Замечание: необходимо для совместимости с net/http пакетом.
func (s *RestServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// httpErr - используется для корректного вывода возвращаемых значений в swagger
type httpErr struct {
	Error string `json:"error"`
}

func ginError(value string) httpErr {
	return httpErr{value}
}
