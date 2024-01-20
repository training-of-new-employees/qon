// Package handlers - пакет содержит middlewar-ы.
package rest

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/training-of-new-employees/qon/internal/errs"
	"github.com/training-of-new-employees/qon/internal/logger"
	"github.com/training-of-new-employees/qon/internal/pkg/jwttoken"
)

type UserSession struct {
	UserID        int
	IsAdmin       bool
	OrgID         int
	HashedRefresh string
}

// IsAuthenticated - middleware для проверки авторизации.
func (r *RestServer) IsAuthenticated() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := jwttoken.GetToken(c)

		claims, err := r.tokenVal.ValidateToken(token)

		if err != nil {
			logger.Log.Warn("error invalid token: %v", zap.Error(err))
			r.handleError(c, errs.ErrUnauthorized)
			return
		}

		_, err = r.cache.GetRefreshToken(c.Request.Context(), claims.HashedRefresh)
		if err != nil {
			logger.Log.Warn("error invalid session: %v", zap.Error(err))
			r.handleError(c, errs.ErrUnauthorized)
			return
		}

		us := UserSession{
			UserID:        claims.UserID,
			IsAdmin:       claims.IsAdmin,
			OrgID:         claims.OrgID,
			HashedRefresh: claims.HashedRefresh,
		}

		c.Set("session", &us)

		c.Next()
	}
}

// IsAdmin - middleware для проверки прав администратора.
func (r *RestServer) IsAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := jwttoken.GetToken(c)

		claims, err := r.tokenVal.ValidateToken(token)
		if err != nil {
			logger.Log.Warn("error  invalid token: %v", zap.Error(err))
			r.handleError(c, errs.ErrUnauthorized)
			return
		}

		if !claims.IsAdmin {
			logger.Log.Warn("error permission denied: %v", zap.Error(err))
			r.handleError(c, errs.ErrNoAccess)
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

func (r *RestServer) getUserSession(c *gin.Context) *UserSession {
	val, _ := c.Get("session")

	us, ok := val.(*UserSession)
	if !ok {
		logger.Log.Warn("ctx without user session")

		return &UserSession{}
	}

	return us
}

func (r *RestServer) CORS() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"OPTIONS", "GET", "POST", "PATCH", "PUT"},
		AllowHeaders:     []string{"Origin", "Accept", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	})
}
