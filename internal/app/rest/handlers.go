package rest

// InitRoutes - инициализация роутеров.
func (s *RestServer) InitRoutes() {

	s.router.Use(s.DummyMiddleware())
	s.router.Use(s.LoggerMiddleware())
	//s.router.Use(s.IsAuthenticated())
	//s.router.Use(s.IsAdmin())

	s.router.GET("/api/v1/dummy", s.Dummy)

	return
}
