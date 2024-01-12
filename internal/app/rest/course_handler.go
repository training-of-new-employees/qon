package rest

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/training-of-new-employees/qon/internal/errs"
	"github.com/training-of-new-employees/qon/internal/model"
)

// GetAdminCourses godoc
//
//	@Summary	Получение данных о курсах администратором
//	@Tags		course
//	@Produce	json
//	@Success	200	{array}		model.Course
//	@Failure	400	{object}	errResponse
//	@Failure	401	{object}	errResponse
//	@Failure	404	{object}	errResponse
//	@Failure	500	{object}	errResponse
//	@Router		/admin/courses [get]
func (r *RestServer) handlerGetAdminCourses(c *gin.Context) {
	ctx := c.Request.Context()
	us := r.getUserSession(c)
	courses, err := r.services.Course().GetCompanyCourses(ctx, us.OrgID)
	if err != nil {
		r.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, courses)
}

// GetUserCourses godoc
//
//	@Summary	Получение данных о курсах пользователем
//	@Tags		course
//	@Produce	json
//	@Success	200	{array}		model.Course
//	@Failure	400	{object}	errResponse
//	@Failure	401	{object}	errResponse
//	@Failure	404	{object}	errResponse
//	@Failure	500	{object}	errResponse
//	@Router		/users/courses [get]
func (r *RestServer) handlerGetUserCourses(c *gin.Context) {
	ctx := c.Request.Context()
	us := r.getUserSession(c)
	courses, err := r.services.Course().GetUserCourses(ctx, us.UserID)
	if err != nil {
		r.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, courses)
}

// CreateCourse godoc
//
//	@Summary	Создание нового курса
//	@Tags		course
//	@Produce	json
//	@Param		object	body		model.CourseSet	true	"Course Create"
//	@Success	201		{object}	model.Course
//	@Failure	400		{object}	errResponse
//	@Failure	401		{object}	errResponse
//	@Failure	500		{object}	errResponse
//	@Router		/admin/courses [post]
func (r *RestServer) handlerCreateCourse(c *gin.Context) {
	ctx := c.Request.Context()
	course := model.NewCourseSet(0, r.getUserSession(c).UserID)
	err := c.BindJSON(&course)
	if err != nil {
		r.handleError(c, errs.ErrBadRequest)
		return
	}
	created, err := r.services.Course().CreateCourse(ctx, course)
	if err != nil {
		r.handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, created)

}

// EditCourse godoc
//
//	@Summary	Изменение данных курса
//	@Tags		course
//	@Produce	json
//	@Param		id		path		int				true	"Course ID"
//	@Param		object	body		model.CourseSet	true	"Course Edit"
//	@Success	200		{object}	courseResp
//	@Failure	400		{object}	errResponse
//	@Failure	401		{object}	errResponse
//	@Failure	404		{object}	errResponse
//	@Failure	500		{object}	errResponse
//	@Router		/admin/courses/{id} [patch]
func (r *RestServer) handlerEditCourse(c *gin.Context) {
	ctx := c.Request.Context()
	sID := c.Param("id")
	id, err := strconv.Atoi(sID)
	if err != nil {
		r.handleError(c, errs.ErrBadRequest)
		return
	}
	us := r.getUserSession(c)
	course := model.NewCourseSet(id, us.UserID)
	err = c.BindJSON(&course)
	if err != nil {
		r.handleError(c, errs.ErrBadRequest)
		return
	}
	edited, err := r.services.Course().EditCourse(ctx, course, us.OrgID)
	if err != nil {
		r.handleError(c, err)
		return
	}
	resp := courseResp{
		ID: edited.ID,
		CourseSet: model.CourseSet{
			Name:        edited.Name,
			Description: edited.Description,
			IsArchived:  edited.IsArchived,
		},
	}
	c.JSON(http.StatusOK, resp)

}

type courseResp struct {
	ID int `json:"id"`
	model.CourseSet
}
