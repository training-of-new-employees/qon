package rest

import (
	"errors"
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
//
//	@Security	Bearer
//
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

// GetAdminCourse godoc
//
//	@Summary	Получение данных курса администратором
//	@Tags		course
//	@Produce	json
//	@Param		id	path		int	true	"Course ID"
//	@Success	200	{object}	model.Course
//	@Failure	400	{object}	errResponse
//	@Failure	401	{object}	errResponse
//	@Failure	404	{object}	errResponse
//	@Failure	500	{object}	errResponse
//
//	@Security	Bearer
//
//	@Router		/admin/courses/{id} [get]
func (r *RestServer) handlerGetAdminCourse(c *gin.Context) {
	ctx := c.Request.Context()
	courseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		r.handleError(c, errs.ErrBadRequest)
		return
	}
	us := r.getUserSession(c)
	course, err := r.services.Course().GetCompanyCourse(ctx, courseID, us.OrgID)
	if err != nil {
		r.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, course)
}

// GetUserCourse godoc
//
//	@Summary	Получение данных курса пользователем
//	@Tags		course
//	@Produce	json
//	@Param		id	path		int	true	"Course ID"
//	@Success	200	{object}	model.Course
//	@Failure	400	{object}	errResponse
//	@Failure	401	{object}	errResponse
//	@Failure	404	{object}	errResponse
//	@Failure	500	{object}	errResponse
//
//	@Security	Bearer
//
//	@Router		/users/courses/{id} [get]
func (r *RestServer) handlerGetUserCourse(c *gin.Context) {
	ctx := c.Request.Context()
	courseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		r.handleError(c, errs.ErrBadRequest)
		return
	}
	us := r.getUserSession(c)
	course, err := r.services.Course().GetUserCourse(ctx, courseID, us.UserID)
	if err != nil {
		r.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, course)
}

// GetUserCourses godoc
//
//	@Summary	Получение данных о курсах пользователем
//	@Tags		course
//	@Produce	json
//	@Success	200	{array}		model.CoursePreview
//	@Failure	400	{object}	sErr
//	@Failure	401	{object}	sErr
//	@Failure	404	{object}	sErr
//	@Failure	500	{object}	sErr
//
//	@Security	Bearer
//
//	@Router		/users/courses [get]
func (r *RestServer) handlerGetUserCourses(c *gin.Context) {
	ctx := c.Request.Context()
	us := r.getUserSession(c)
	courses, err := r.services.Course().GetUserCourses(ctx, us.UserID)
	if err != nil {
		r.handleError(c, err)
		return
	}

	previews := make([]model.CoursePreview, 0, len(courses))
	for _, course := range courses {
		previews = append(previews, model.CoursePreview{
			CourseID:    course.ID,
			Name:        course.Name,
			Description: course.Description,
			Status:      course.Status,
		})
	}

	c.JSON(http.StatusOK, previews)
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
//
//	@Security	Bearer
//
//	@Router		/admin/courses [post]
func (r *RestServer) handlerCreateCourse(c *gin.Context) {
	ctx := c.Request.Context()
	course := model.NewCourseSet(0, r.getUserSession(c).UserID)
	err := c.BindJSON(&course)
	if err != nil {
		r.handleError(c, errs.ErrBadRequest)
		return
	}
	err = course.Validation()
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
//
//	@Security	Bearer
//
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
	err = course.Validation()
	if !errors.Is(err, errs.ErrCourseNameIsEmpty) && err != nil {
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

// GetUserCourseLessons godoc
//
//	@Summary	Получение данных о уроках курса
//	@Tags		course
//	@Param		id	path	int	true	"Course ID"
//	@Produce	json
//	@Success	200	{array}		model.LessonPreview
//	@Failure	400	{object}	sErr
//	@Failure	401	{object}	sErr
//	@Failure	404	{object}	sErr
//	@Failure	500	{object}	sErr
//
//	@Security	Bearer
//
//	@Router		/users/courses/{id}/lessons [get]
func (r *RestServer) handlerGetUserCourseLessons(c *gin.Context) {
	ctx := c.Request.Context()
	us := r.getUserSession(c)
	courseIDStr := c.Param("id")
	courseID, err := strconv.Atoi(courseIDStr)
	if err != nil {
		r.handleError(c, errs.ErrBadRequest)
		return
	}

	lessons, err := r.services.Course().GetUserCourseLessons(ctx, us.UserID, courseID)
	if err != nil {
		r.handleError(c, err)
		return
	}

	previews := make([]model.LessonPreview, 0, len(lessons))
	for _, lesson := range lessons {
		previews = append(previews, model.LessonPreview{
			CourseID: lesson.CourseID,
			Name:     lesson.Name,
			LessonID: lesson.ID,
			Status:   lesson.Status,
		})
	}

	c.JSON(http.StatusOK, previews)
}
