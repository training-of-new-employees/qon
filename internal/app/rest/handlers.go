package rest

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/training-of-new-employees/qon/docs"
	"github.com/training-of-new-employees/qon/internal/errs"
)

//	@title			QuickOn
//	@version		0.1
//	@description	Описание API QuickOn

//	@host		localhost:8080
//	@BasePath	/api/v1

//	@securityDefinitions.apikey	Bearer
//	@in							header
//	@name						Authorization
//	@description				you can get it on login page

//	@externalDocs.description	OpenAPI
//	@externalDocs.url			https://swagger.io/resources/open-api/

// InitRoutes - инициализация роутеров.
func (s *RestServer) InitRoutes() {

	s.router.Use(s.LoggerMiddleware())
	allRoutes := s.router.Group("/api")
	mvp := allRoutes.Group("/v1")
	login := mvp.Group("/login")
	login.POST("", s.handlerSignIn)
	password := mvp.Group("/password")
	password.POST("", s.handlerResetPassword)
	adminGroup := mvp.Group("/admin")
	adminGroup.POST("/register", s.handlerCreateAdminInCache)
	adminGroup.POST("/verify", s.handlerAdminEmailVerification)
	restrictedAdmin := adminGroup.Group("")
	restrictedAdmin.Use(s.IsAuthenticated())
	restrictedAdmin.Use(s.IsAdmin())
	restrictedAdmin.POST("/employee", s.handlerCreateUser)
	restrictedAdmin.PATCH("/info", s.handlerAdminEdit)

	lessons := mvp.Group("/lesson")
	lessons.POST("/", s.handlerLessonCreate)
	lessons.DELETE("/", s.handlerLessonDelete)
	lessons.GET("/", s.handlerLessonGet)
	lessons.PATCH("/", s.handlerLessonUpdate)

	userGroup := mvp.Group("/users")
	userGroup.Use(s.IsAuthenticated())
	userGroup.GET("", s.handlerGetUsers)
	userGroup.GET("/:id", s.handlerGetUser)
	userGroup.PATCH("/:id", s.handlerEditUser)
	userGroup.POST("/set-password", s.handlerSetPassword)
	userGroup.PATCH("/archive/:id", s.handlerArchiveUser)
	userGroup.Use(s.IsAuthenticated())
	userGroup.GET("/info", s.handlerUserInfo)

	position := mvp.Group("/positions")
	position.Use(s.IsAuthenticated())
	position.Use(s.IsAdmin())
	position.POST("", s.handlerCreatePosition)
	position.POST("/course", s.handlerAssignCourse)
	position.GET("", s.handlerGetPositions)
	position.Any("/", s.NotFound(errs.ErrPositionNotFound))
	position.GET("/:id", s.handlerGetPosition)
	position.PATCH("/update/:id", s.handlerUpdatePosition)

	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func (s *RestServer) NotFound(err error) gin.HandlerFunc {
	return func(c *gin.Context) {
		s.handleError(c, err)
	}
}
