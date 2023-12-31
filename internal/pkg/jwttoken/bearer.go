package jwttoken

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func GetToken(r *gin.Context) string {
	bearer := r.GetHeader("Authorization")

	arr := strings.Split(bearer, " ")
	if len(arr) == 2 {
		return arr[1]
	}

	return ""
}
