package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Dummy - тестовый хендлер.
func (h *Handlers) Dummy(c *gin.Context) {
	message := "Hello!"

	c.JSON(
		http.StatusOK,
		message,
	)
}
