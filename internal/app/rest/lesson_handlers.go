package rest

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/training-of-new-employees/qon/internal/errs"
	"github.com/training-of-new-employees/qon/internal/model"
)

// @Summary	Создание урока
// @Tags	lessons
// @Produce	json
// @Param	object	body		model.LessonCreate	true	"Lesson Create"
// @Success	201		{object}	model.LessonCreate
// @Failure	400		{object}	sErr
// @Failure	404		{object}	sErr
// @Failure	500		{object}	sErr
// @Router		/lessons [post]
func (r *RestServer) handlerLessonCreate(c *gin.Context) {
	lessonCreate := &model.Lesson{}

	if err := c.ShouldBindJSON(&lessonCreate); err != nil {
		c.JSON(http.StatusBadRequest, s().SetError(err))
		return
	}

	ctx := c.Request.Context()
	us := r.getUserSession(c)

	lesson, err := r.services.Lesson().CreateLesson(ctx, *lessonCreate,
		us.UserID)
	switch {
	case errors.Is(err, errs.ErrNotFound):
		c.JSON(http.StatusNotFound, s().SetError(err))
		return
	case err != nil:
		c.JSON(http.StatusInternalServerError, s().SetError(err))

		return
	}

	c.JSON(http.StatusCreated, lesson)
}

// @Summary	Получение урока
// @Tags		lessons
// @Produce	json
// @Param		id	path	int	true	"Lesson ID"
// @Success	200
// @Failure	400	{object}	sErr
// @Failure	401	{object}	sErr
// @Failure	403	{object}	sErr
// @Failure	404	{object}	sErr
// @Failure	500	{object}	sErr
// @Router		/lessons/{id} [get]
func (r *RestServer) handlerLessonGet(c *gin.Context) {
	val := c.Param("id")

	lessonID, err := strconv.Atoi(val)
	if err != nil {
		c.JSON(http.StatusBadRequest, s().SetError(err))
		return
	}

	ctx := c.Request.Context()
	lesson, err := r.services.Lesson().GetLesson(ctx, lessonID)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrNotFound):
			c.JSON(http.StatusNotFound, s().SetError(err))
			return
		case err != nil:
			c.JSON(http.StatusInternalServerError, s().SetError(err))
			return
		}
	}
	c.JSON(http.StatusOK, lesson)
}

// @Summary	Обновление урока
// @Tags		lessons
// @Produce	json
// @Param		id		path		int					true	"Lesson ID"
// @Param		object	body		model.LessonUpdate	true	"Lesson Update"
// @Success	200		{object}	model.Lesson
// @Failure	400		{object}	sErr
// @Failure	401		{object}	sErr
// @Failure	403		{object}	sErr
// @Failure	404		{object}	sErr
// @Failure	500		{object}	sErr
// @Router		/lessons/{id} [patch]
func (r *RestServer) handlerLessonUpdate(c *gin.Context) {
	var err error
	lessonUpdate := model.LessonUpdate{}

	if err = c.ShouldBindJSON(&lessonUpdate); err != nil {
		c.JSON(http.StatusBadRequest, s().SetError(err))
		return
	}

	val := c.Param("id")

	lessonUpdate.ID, err = strconv.Atoi(val)
	if err != nil {
		c.JSON(http.StatusBadRequest, s().SetError(err))
		return
	}

	ctx := c.Request.Context()
	lesson, err := r.services.Lesson().UpdateLesson(ctx, lessonUpdate)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrNotFound):
			c.JSON(http.StatusNotFound, s().SetError(err))
			return
		case err != nil:
			c.JSON(http.StatusInternalServerError, s().SetError(err))
			return
		}
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
//	@Failure	404	{object}	sErr
//	@Failure	401	{object}	sErr
//	@Failure	403	{object}	sErr
//	@Failure	500	{object}	sErr
//	@Router		/admin/courses/{id}/lessons [get]
func (r *RestServer) handlerGetLessonsList(c *gin.Context) {
	val := c.Param("id")

	courseID, err := strconv.Atoi(val)
	if err != nil {
		c.JSON(http.StatusBadRequest, s().SetError(err))
		return
	}

	ctx := c.Request.Context()
	lessonsList, err := r.services.Lesson().GetLessonsList(ctx, courseID)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrNotFound):
			c.JSON(http.StatusNotFound, s().SetError(err))
			return
		case err != nil:
			c.JSON(http.StatusInternalServerError, s().SetError(err))
			return
		}
	}
	c.JSON(http.StatusOK, lessonsList)

}
