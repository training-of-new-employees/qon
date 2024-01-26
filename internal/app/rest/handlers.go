package rest

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/training-of-new-employees/qon/docs"
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
	s.router.NoRoute(s.invalidHandler)

	s.router.Use(s.CORS())

	s.router.Use(s.LoggerMiddleware())
	allRoutes := s.router.Group("/api")
	mvp := allRoutes.Group("/v1")
	login := mvp.Group("/login")
	login.POST("", s.handlerSignIn)
	logout := mvp.Group("/logout")
	logout.Use(s.IsAuthenticated())
	logout.POST("", s.handlerLogOut)

	password := mvp.Group("/password")
	password.POST("", s.handlerResetPassword)
	adminGroup := mvp.Group("/admin")
	adminGroup.POST("/register", s.handlerCreateAdminInCache)
	adminGroup.POST("/verify", s.handlerAdminEmailVerification)

	restrictedAdmin := adminGroup.Group("")
	restrictedAdmin.Use(s.IsAuthenticated())
	restrictedAdmin.Use(s.IsAdmin())
	restrictedAdmin.POST("/employee", s.handlerCreateUser)
	restrictedAdmin.PATCH("/info", s.handlerEditAdmin)

	adminCourses := restrictedAdmin.Group("/courses")
	adminCourses.GET("", s.handlerGetAdminCourses)
	adminCourses.POST("", s.handlerCreateCourse)
	adminCourses.PATCH("/:id", s.handlerEditCourse)
	adminCourses.GET("/:id/lessons", s.handlerGetLessonsList)
	adminLessons := restrictedAdmin.Group("/lessons")
	adminLessons.POST("", s.handlerLessonCreate)
	adminLessons.GET("/:id", s.handlerLessonGet)
	adminLessons.PATCH("/:id", s.handlerLessonUpdate)

	invitationLinkGroup := mvp.Group("/invitation-link")
	invitationLinkGroup.Use(s.IsAuthenticated(), s.IsAdmin())
	invitationLinkGroup.GET("/:email", s.handlerGetInvitationLink)
	invitationLinkGroup.PATCH("", s.handlerRegenerationInvitationLink)

	userGroup := mvp.Group("/users")
	userGroup.POST("/set-password", s.handlerSetPassword)
	userGroup.Use(s.IsAuthenticated())
	userGroup.GET("/info", s.handlerUserInfo)
	userGroup.GET("/courses", s.handlerGetUserCourses)
	userGroup.GET("/courses/:id/lessons", s.handlerGetUserCourseLessons)
	userGroup.GET("/lessons/:id", s.handlerGetLesson)
	userGroup.PATCH("/lessons/:id", s.handlerUpdateLessonStatus)
	userGroup.Use(s.IsAdmin())
	userGroup.GET("", s.handlerGetUsers)
	userGroup.GET("/:id", s.handlerGetUser)
	userGroup.PATCH("/:id", s.handlerEditUser)
	userGroup.PATCH("/archive/:id", s.handlerArchiveUser)

	position := mvp.Group("/positions")
	position.Use(s.IsAuthenticated())
	position.Use(s.IsAdmin())
	position.POST("", s.handlerCreatePosition)
	position.POST("/course", s.handlerAssignCourse)
	position.GET("", s.handlerGetPositions)
	position.GET("/:id/courses", s.handlerGetPositionCourses)
	position.PATCH("/:id/courses", s.handlerAssignCourses)
	position.GET("/:id", s.handlerGetPosition)
	position.PATCH("/update/:id", s.handlerUpdatePosition)

	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func (s *RestServer) NotFound(err error) gin.HandlerFunc {
	return func(c *gin.Context) {
		s.handleError(c, err)
	}
}
