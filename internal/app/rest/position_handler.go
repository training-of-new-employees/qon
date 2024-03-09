package rest

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/training-of-new-employees/qon/internal/errs"
	"github.com/training-of-new-employees/qon/internal/model"
)

// CreatePosition godoc
//
//	@Summary	Админ.Должности.Создание должности
//	@Tags		position
//	@Produce	json
//	@Param		object	body		reqCreatePosition	true	"Position Create"
//	@Success	201		{object}	model.Position
//	@Failure	400		{object}	errResponse
//	@Failure	500		{object}	errResponse
//
//	@Security	Bearer
//
//	@Router		/positions [post]
func (r *RestServer) handlerCreatePosition(c *gin.Context) {
	ctx := c.Request.Context()

	request := &reqCreatePosition{}

	if err := c.ShouldBindJSON(&request); err != nil {
		r.handleError(c, errs.ErrInvalidRequest)
		return
	}

	positionReq := model.PositionSet{
		CompanyID: request.CompanyID,
		Name:      request.Name,
	}

	if err := positionReq.Validation(); err != nil {
		r.handleError(c, err)
		return
	}

	us := r.getUserSession(c)
	if us.OrgID != positionReq.CompanyID {
		r.handleError(c, errs.ErrCompanyNotFound)
		return
	}

	position, err := r.services.Position().CreatePosition(ctx, positionReq)

	if err != nil {
		r.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, position)
}

// GetPosition godoc
//
//	@Summary	Админ.Должности.Получение данных должности
//	@Tags		position
//	@Produce	json
//	@Param		id	path		int	true	"Position ID"
//	@Success	200	{object}	model.Position
//	@Failure	404	{object}	errResponse
//	@Failure	500	{object}	errResponse
//
//	@Security	Bearer
//
//	@Router		/positions/{id} [get]
func (r *RestServer) handlerGetPosition(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 0 {
		r.handleError(c, errs.ErrBadRequest)
		return
	}

	position, err := r.services.Position().GetPosition(ctx, r.getUserSession(c).OrgID, id)
	if err != nil {
		r.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, position)
}

// GetPositionCourses godoc
//
//	@Summary	Админ.Должности.Получить список курсов для должности
//	@Tags		position
//	@Produce	json
//	@Param		id	path		int	true	"Position ID"
//	@Success	200	{object}	getPositionCoursesResponse
//	@Failure	401	{object}	errResponse
//	@Failure	403	{object}	errResponse
//	@Failure	404	{object}	errResponse
//	@Failure	500	{object}	errResponse
//
//	@Security	Bearer
//
//	@Router		/positions/{id}/courses [get]
func (r *RestServer) handlerGetPositionCourses(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 0 {
		r.handleError(c, errs.ErrBadRequest)
		return
	}

	courses, err := r.services.Position().GetPositionCourses(ctx, r.getUserSession(c).OrgID, id)
	if err != nil {
		r.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, getPositionCoursesResponse{
		PositionID: id,
		CoursesID:  courses,
	})
}

type getPositionCoursesResponse struct {
	PositionID int   `json:"position_id"`
	CoursesID  []int `json:"courses_id"`
}

// GetPositions godoc
//
//	@Summary	Админ.Должности.Получение всех должностей
//	@Tags		position
//	@Produce	json
//	@Success	200	{array}		model.Position
//	@Failure	404	{object}	errResponse
//	@Failure	500	{object}	errResponse
//
//	@Security	Bearer
//
//	@Router		/positions [get]
func (r *RestServer) handlerGetPositions(c *gin.Context) {
	ctx := c.Request.Context()

	positions, err := r.services.Position().GetPositions(ctx, r.getUserSession(c).OrgID)
	if err != nil {
		r.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, positions)
}

// UpdatePosition godoc
//
//	@Summary	Админ.Должности.Редактирование/ Архивирование должности
//	@Tags		position
//	@Produce	json
//	@Param		id		path		int					true	"Position ID"
//	@Param		object	body		model.PositionSet	true	"Position info"
//	@Success	200		{object}	model.Position
//	@Failure	400		{object}	errResponse
//	@Failure	404		{object}	errResponse
//	@Failure	500		{object}	errResponse
//
//	@Security	Bearer
//
//	@Router		/positions/update/{id} [patch]
func (r *RestServer) handlerUpdatePosition(c *gin.Context) {
	ctx := c.Request.Context()
	positionReq := model.PositionSet{}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		r.handleError(c, errs.ErrBadRequest)
		return
	}

	if err := c.ShouldBindJSON(&positionReq); err != nil {
		r.handleError(c, errs.ErrInvalidRequest)
		return
	}

	if err := positionReq.ValidationEdit(); err != nil {
		r.handleError(c, err)
		return
	}

	position, err := r.services.Position().UpdatePosition(ctx, id, positionReq)
	if err != nil {
		r.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, position)
}

// @Summary	Админ.Должности. Редактировать список курсов для должности
// @Accept		json
// @Tags		position
// @Produce	json
// @Param		id		path		int							true	"Position ID"
// @Param		object	body		model.PositionAssignCourses	true	"Courses"
// @Success	200		{object}	assignCoursesResponse
// @Failure	400		{object}	errResponse	"Неверный формат запроса"
// @Failure	401		{object}	errResponse	"Пользователь не является сотрудником компании"
// @Failure	500		{object}	errResponse	"Внутренняя ошибка сервера"
//
// @Security	Bearer
//
// @Router		/positions/{id}/courses [patch]
func (r *RestServer) handlerAssignCourses(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 0 {
		r.handleError(c, errs.ErrBadRequest)
		return
	}

	var body model.PositionAssignCourses
	if err := c.ShouldBindJSON(&body); err != nil {
		r.handleError(c, errs.ErrInvalidRequest)
		return
	}

	if err := body.Validation(); err != nil {
		r.handleError(c, err)
		return
	}

	err = r.services.Position().AssignCourses(ctx, id, body.CoursesID, r.getUserSession(c).UserID)
	if err != nil {
		fmt.Println(err)
		r.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, assignCoursesResponse{
		CoursesID:  body.CoursesID,
		PositionID: id,
	})
}

type assignCoursesResponse struct {
	PositionID int   `json:"position_id"`
	CoursesID  []int `json:"course_id"`
}
