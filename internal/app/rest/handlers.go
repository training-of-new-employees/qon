package rest

// InitRoutes - инициализация роутеров.
func (s *RestServer) InitRoutes() {

	s.router.Use(s.LoggerMiddleware())
	allRoutes := s.router.Group("/api")
	mvp := allRoutes.Group("/v1")
	login := mvp.Group("/login")
	login.POST("/", s.handlerSignIn)
	password := mvp.Group("/password")
	password.POST("/", s.handlerResetPassword)
	adminGroup := mvp.Group("/admin")
	adminGroup.POST("/register", s.handlerCreateAdminInCache)
	adminGroup.POST("/verify", s.handlerAdminEmailVerification)
	adminGroup.POST("/employee", s.handlerCreateUser)

	adminGroup.PATCH("/info", s.handlerAdminEditInfo)
	userGroup := mvp.Group("/users")
	userGroup.POST("/set-password", s.handlerSetPassword)

	position := mvp.Group("/positions")
	position.Use(s.IsAuthenticated())
	position.Use(s.IsAdmin())
	position.POST("/", s.handlerCreatePosition)
	position.GET("/", s.handlerGetPositions)
	position.GET("/:id", s.handlerGetPosition)
	position.PATCH("/update/:id", s.handlerUpdatePosition)
	position.DELETE("/delete/:id", s.handlerDeletePosition)

	mvp.POST("/", s.handlerAssignCourse)

}
