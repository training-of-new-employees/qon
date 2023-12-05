package rest

// InitRoutes - инициализация роутеров.
func (s *RestServer) InitRoutes() {

	//s.router.Use(s.DummyMiddleware())
	//s.router.GET("/api/v1/dummy", s.Dummy)
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
	//s.router.Use(s.IsAuthenticated())
	//s.router.Use(s.IsAdmin())
}
