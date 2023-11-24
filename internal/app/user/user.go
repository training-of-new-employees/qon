package user

import "github.com/gin-gonic/gin"

type user struct {
	route *gin.Engine
}

type Config struct {
	Route *gin.Engine
}

func New(cfg Config) *user {
	return &user{
		route: cfg.Route,
	}
}
