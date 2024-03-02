package rest

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/training-of-new-employees/qon/internal/errs"
	"github.com/training-of-new-employees/qon/internal/model"
)

// @Summary	Админ.Урок.Создание урока в рамках курса
// @Tags		lessons
// @Produce	json
// @Param		object	body		reqCreateLesson	true	"Lesson Create"
// @Success	201		{object}	model.Lesson
// @Failure	400		{object}	errResponse
// @Failure	404		{object}	errResponse
// @Failure	500		{object}	errResponse
//
// @Security	Bearer
//
// @Router		/admin/lessons [post]
func (r *RestServer) handlerLessonCreate(c *gin.Context) {
	ctx := c.Request.Context()
	reqLesson := reqCreateLesson{}
	if err := c.ShouldBindJSON(&reqLesson); err != nil {
		r.handleError(c, errs.ErrInvalidRequest)
		return
	}

	lessonCreate := model.Lesson{
		Name:       reqLesson.Name,
		Content:    reqLesson.Content,
		CourseID:   reqLesson.CourseID,
		URLPicture: reqLesson.URLPicture,
	}

	if err := lessonCreate.Validation(); err != nil {
		r.handleError(c, err)
		return
	}
	us := r.getUserSession(c)
	lesson, err := r.services.Lesson().CreateLesson(ctx, lessonCreate, us.UserID)
	if err != nil {
		r.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, lesson)
}

// @Summary	Админ.Урок.Получение урока курса
// @Tags		lessons
// @Produce	json
// @Param		id	path	int	true	"Lesson ID"
// @Success	200
// @Failure	400	{object}	errResponse
// @Failure	401	{object}	errResponse
// @Failure	403	{object}	errResponse
// @Failure	404	{object}	errResponse
// @Failure	500	{object}	errResponse
//
// @Security	Bearer
//
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

// @Summary	Админ.Урок.Редактирование/Архивирование урока курса
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
//
// @Security	Bearer
//
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
//	@Summary	Админ.Уроки.Получение уроков в рамках курса.
//	@Tags		course
//	@Produce	json
//	@Param		id	path		int	true	"Course ID"
//	@Success	200	{array}		model.Lesson
//	@Failure	404	{object}	errResponse
//	@Failure	401	{object}	errResponse
//	@Failure	403	{object}	errResponse
//	@Failure	500	{object}	errResponse
//
//	@Security	Bearer
//
//	@Router		/admin/courses/{id}/lessons [get]
func (r *RestServer) handlerGetLessonsList(c *gin.Context) {
	courseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		r.handleError(c, errs.ErrBadRequest)
		return
	}
	us := r.getUserSession(c)

	ctx := c.Request.Context()
	lessonsList, err := r.services.Lesson().GetLessonsList(ctx, courseID, us.OrgID)
	if err != nil {
		r.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, lessonsList)

}

// GetUserLesson godoc
//
//	@Summary	Сотрудник.Урок. Получение данных урока по id
//	@Tags		lessons
//	@Param		id	path	int	true	"Lesson ID"
//	@Produce	json
//	@Success	200	{object}	model.Lesson
//	@Failure	400	{object}	errResponse
//	@Failure	401	{object}	errResponse
//	@Failure	404	{object}	errResponse
//	@Failure	500	{object}	errResponse
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
//	@Summary	Сотрудник. Урок. Прогресс по уроку
//	@Tags		lessons
//	@Param		id		path	int							true	"Lesson ID"
//	@Param		object	body	model.LessonStatusUpdate	true	"Lesson Status Update"
//	@Produce	json
//	@Success	200	{object}	updateLessonStatusResponse
//	@Failure	400	{object}	errResponse
//	@Failure	401	{object}	errResponse
//	@Failure	404	{object}	errResponse
//	@Failure	500	{object}	errResponse
//
//	@Security	Bearer
//
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
