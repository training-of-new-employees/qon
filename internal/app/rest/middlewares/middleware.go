// Package handlers - пакет содержит middlewar-ы.
package middlewares

import (
	"log"

	"github.com/gin-gonic/gin"
)

// dummyMiddleware - тестовый middleware, используется для проверки.
func DummyMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Println("into dummyMiddleware")

		ctx.Next()
	}
}
