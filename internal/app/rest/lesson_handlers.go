package rest

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/training-of-new-employees/qon/internal/errs"
	"github.com/training-of-new-employees/qon/internal/model"
)

// @Summary	Создание урока
// @Tags		lessons
// @Produce	json
// @Param		object	body		model.Lesson	true	"Lesson Create"
// @Success	201		{object}	model.Lesson
// @Failure	400		{object}	errResponse
// @Failure	404		{object}	errResponse
// @Failure	500		{object}	errResponse
// @Router		/admin/lessons [post]
func (r *RestServer) handlerLessonCreate(c *gin.Context) {
	ctx := c.Request.Context()
	lessonCreate := model.Lesson{}
	if err := c.ShouldBindJSON(&lessonCreate); err != nil {
		r.handleError(c, errs.ErrInvalidRequest)
		return
	}

	us := r.getUserSession(c)
	if err := lessonCreate.Validation(); err != nil {
		r.handleError(c, err)
		return
	}
	lesson, err := r.services.Lesson().CreateLesson(ctx, lessonCreate, us.UserID)
	if err != nil {
		r.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, lesson)
}

// @Summary	Получение урока
// @Tags		lessons
// @Produce	json
// @Param		id	path	int	true	"Lesson ID"
// @Success	200
// @Failure	400	{object}	errResponse
// @Failure	401	{object}	errResponse
// @Failure	403	{object}	errResponse
// @Failure	404	{object}	errResponse
// @Failure	500	{object}	errResponse
// @Router		/admin/lessons/{id} [get]
func (r *RestServer) handlerLessonGet(c *gin.Context) {
	lessonID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		r.handleError(c, errs.ErrBadRequest)
		return
	}

	ctx := c.Request.Context()
	lesson, err := r.services.Lesson().GetLesson(ctx, lessonID)
	if err != nil {
		r.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, lesson)
}

// @Summary	Обновление урока
// @Tags		lessons
// @Produce	json
// @Param		id		path		int					true	"Lesson ID"
// @Param		object	body		model.LessonUpdate	true	"Lesson Update"
// @Success	200		{object}	model.Lesson
// @Failure	400		{object}	errResponse
// @Failure	401		{object}	errResponse
// @Failure	403		{object}	errResponse
// @Failure	404		{object}	errResponse
// @Failure	500		{object}	errResponse
// @Router		/admin/lessons/{id} [patch]
func (r *RestServer) handlerLessonUpdate(c *gin.Context) {
	var err error
	lessonUpdate := model.LessonUpdate{}

	lessonUpdate.ID, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		r.handleError(c, errs.ErrBadRequest)
		return
	}

	if err = c.ShouldBindJSON(&lessonUpdate); err != nil {
		r.handleError(c, errs.ErrInvalidRequest)
		return
	}

	if err := lessonUpdate.Validation(); err != nil {
		r.handleError(c, err)
		return
	}

	ctx := c.Request.Context()
	lesson, err := r.services.Lesson().UpdateLesson(ctx, lessonUpdate)
	if err != nil {
		r.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, lesson)
}

// GetLessonsList godoc
//
//	@Summary	Получение уроков курса
//	@Tags		course
//	@Produce	json
//	@Param		id	path		int	true	"Course ID"
//	@Success	200	{object}	[]model.Lesson
//	@Failure	404	{object}	errResponse
//	@Failure	401	{object}	errResponse
//	@Failure	403	{object}	errResponse
//	@Failure	500	{object}	errResponse
//	@Router		/admin/courses/{id}/lessons [get]
func (r *RestServer) handlerGetLessonsList(c *gin.Context) {
	courseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		r.handleError(c, errs.ErrBadRequest)
		return
	}

	ctx := c.Request.Context()
	lessonsList, err := r.services.Lesson().GetLessonsList(ctx, courseID)
	if err != nil {
		r.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, lessonsList)

}

// GetUserLesson godoc
//
//	@Summary	Получение данных урока пользователя
//	@Tags		lessons
//	@Param		id	path	int	true	"Lesson ID"
//	@Produce	json
//	@Success	200	{object}	model.Lesson
//	@Failure	400	{object}	sErr
//	@Failure	401	{object}	sErr
//	@Failure	404	{object}	sErr
//	@Failure	500	{object}	sErr
//	@Security	Bearer
//	@Router		/users/lessons/{id} [get]
func (r *RestServer) handlerGetLesson(c *gin.Context) {
	lessonID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		r.handleError(c, errs.ErrBadRequest)
		return
	}

	userSession := r.getUserSession(c)
	lesson, err := r.services.Lesson().GetUserLesson(c.Request.Context(), userSession.UserID, lessonID)
	if err != nil {
		r.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, lesson)
}

// UpdateLessonStatus godoc
//
//	@Summary	Обновление статуса прогресса у урока
//	@Tags		lessons
//	@Param		id		path	int							true	"Lesson ID"
//	@Param		object	body	model.LessonStatusUpdate	true	"Lesson Status Update"
//	@Produce	json
//	@Success	200	{array}		updateLessonStatusResponse
//	@Failure	400	{object}	errResponse
//	@Failure	401	{object}	errResponse
//	@Failure	404	{object}	errResponse
//	@Failure	500	{object}	errResponse
//	@Router		/users/lessons/{id} [patch]
func (r *RestServer) handlerUpdateLessonStatus(c *gin.Context) {
	lessonID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		r.handleError(c, errs.ErrBadRequest)
		return
	}

	var body model.LessonStatusUpdate
	if err := c.BindJSON(&body); err != nil {
		r.handleError(c, errs.ErrInvalidRequest)
		return
	}

	userSession := r.getUserSession(c)
	err = r.services.Lesson().UpdateLessonStatus(c.Request.Context(), userSession.UserID, lessonID, body.Status)
	if err != nil {
		r.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, updateLessonStatusResponse{
		LessonID: lessonID,
		Status:   body.Status,
	})
}

type updateLessonStatusResponse struct {
	LessonID int    `json:"lesson_id"`
	Status   string `json:"status"`
}
