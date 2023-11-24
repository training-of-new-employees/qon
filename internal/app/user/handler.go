package user

import "github.com/gin-gonic/gin"

func (u *user) initHandlers() {
	u.route.GET("/", u.createUser)
}

func (u *user) createUser(ctxc *gin.Context) {
	//
}
