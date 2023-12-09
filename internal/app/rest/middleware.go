// Package handlers - пакет содержит middlewar-ы.
package rest

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/training-of-new-employees/qon/internal/logger"
	"github.com/training-of-new-employees/qon/internal/pkg/jwttoken"
)

// DummyMiddleware - тестовый middleware, используется для проверки.
func (r *RestServer) DummyMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Println("into dummyMiddleware")

		ctx.Next()
	}
}

// IsAuthenticated - middleware для проверки авторизации.
func (r *RestServer) IsAuthenticated() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := jwttoken.GetToken(c)

		_, err := r.tokenVal.ValidateToken(token)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			slog.Warn("error invalid token: %v", err)

			return
		}

		c.Next()
	}
}

// IsAdmin - middleware для проверки прав администратора.
func (r *RestServer) IsAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := jwttoken.GetToken(c)

		claims, err := r.tokenVal.ValidateToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			logger.Log.Warn("error invalid token: %v", zap.Error(err))

			return
		}

		if claims.IsAdmin != true {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			logger.Log.Warn("error permission denied: %v", zap.Error(err))
			return
		}

		c.Next()
	}

}

func (r *RestServer) LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Log.Info(
			"request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
		)
		c.Next()
	}
}
