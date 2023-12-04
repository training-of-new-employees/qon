package rest

// InitRoutes - инициализация роутеров.
func (s *RestServer) InitRoutes() {

	//s.router.Use(s.DummyMiddleware())
	//s.router.GET("/api/v1/dummy", s.Dummy)
	s.router.Use(s.LoggerMiddleware())
	allRoutes := s.router.Group("/api")
	mvp := allRoutes.Group("/v1")
	adminGroup := mvp.Group("/admin")
	adminGroup.POST("/", s.handlerCreateAdminInCache)
	adminGroup.POST("/verify", s.EmailVerificationAndAdminCreation)
	userGroup := mvp.Group("/users")
	userGroup.POST("/", s.handlerCreateUser)
	userGroup.POST("/login", s.handlerSignIn)
	//s.router.Use(s.IsAuthenticated())
	//s.router.Use(s.IsAdmin())
}
