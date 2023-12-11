package rest

// InitRoutes - инициализация роутеров.
func (s *RestServer) InitRoutes() {

	s.router.Use(s.LoggerMiddleware())
	allRoutes := s.router.Group("/api")
	mvp := allRoutes.Group("/v1")
	login := mvp.Group("/login")
	login.POST("/", s.handlerSignIn)
	adminGroup := mvp.Group("/admin")
	adminGroup.POST("/register", s.handlerCreateAdminInCache)
	adminGroup.POST("/verify", s.handlerAdminEmailVerification)
	userGroup := mvp.Group("/users")
	userGroup.POST("/", s.handlerCreateUser)

	position := mvp.Group("/positions")
	position.Use(s.IsAuthenticated())
	position.Use(s.IsAdmin())
	position.POST("/", s.handlerCreatePosition)
	position.GET("/", s.handlerGetPositions)
	position.GET("/:id", s.handlerGetPosition)
	position.PATCH("/update/:id", s.handlerUpdatePosition)
	position.DELETE("/delete/:id", s.handlerDeletePosition)
}
