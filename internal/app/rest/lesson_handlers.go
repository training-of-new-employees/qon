package rest

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/training-of-new-employees/qon/internal/errs"
	"github.com/training-of-new-employees/qon/internal/model"
)

//		@Summary	Создание урока
//		@Tags		lesson
//	 	@Accept		json
//		@Produce	json
//		@Success	201		{object}	model.Lesson
//		@Failure	400		{object}	sErr
//		@Failure	404		{object}	sErr
//		@Failure	500		{object}	sErr
//		@Router		/admin/lesson [post]
func (r *RestServer) handlerLessonCreate(c *gin.Context) {
	lessonCreate := &model.LessonCreate{}

	if err := c.ShouldBindJSON(&lessonCreate); err != nil {
		c.JSON(http.StatusBadRequest, s().SetError(err))
		return
	}

	ctx := c.Request.Context()
	us := r.getUserSession(c)

	lesson, err := r.services.Lesson().CreateLesson(ctx, *lessonCreate,
		us.UserID)
	switch {
	case errors.Is(err, errs.ErrCourseNotFound):
		c.JSON(http.StatusNotFound, s().SetError(err))
		return
	case err != nil:
		c.JSON(http.StatusInternalServerError, s().SetError(err))

		return
	}

	c.JSON(http.StatusCreated, lesson)
}

//		@Summary	Удаление урока
//		@Tags		lesson
//	 	@Produce	json
//	 	@Param		id		int	true	"User ID"
//		@Success	200
//		@Failure	400		{object}	sErr
//		@Failure	500		{object}	sErr
//		@Router		/admin/lesson [delete]
func (r *RestServer) handlerLessonDelete(c *gin.Context) {
	val := c.Param("id")

	lessonID, err := strconv.Atoi(val)
	if err != nil {
		c.JSON(http.StatusBadRequest, s().SetError(err))
		return
	}

	ctx := c.Request.Context()
	if err := r.services.Lesson().DeleteLesson(ctx, lessonID); err != nil {
		c.JSON(http.StatusInternalServerError, s().SetError(err))
		return
	}

	c.Status(http.StatusOK)
}

func (r *RestServer) handlerLessonGet(c *gin.Context) {
}