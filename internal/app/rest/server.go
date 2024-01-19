// Package rest - пакет с реализацией rest-сервера.
package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/training-of-new-employees/qon/internal/pkg/jwttoken"
	"github.com/training-of-new-employees/qon/internal/service"
	"github.com/training-of-new-employees/qon/internal/store/cache"
)

// RestServer - реализация rest-сервера.
type RestServer struct {
	router    *gin.Engine
	secretKey string
	tokenVal  jwttoken.JWTValidator
	cache     cache.Cache
	services  service.Service
}

// New - конструктор для RestServer.
func New(secretKey string, services service.Service, cache cache.Cache) *RestServer {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.RedirectTrailingSlash = false

	s := &RestServer{
		router:    router,
		secretKey: secretKey,
		tokenVal:  jwttoken.NewTokenValidator(secretKey),
		cache:     cache,
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
