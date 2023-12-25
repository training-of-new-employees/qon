package rest

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/training-of-new-employees/qon/internal/errs"
	"github.com/training-of-new-employees/qon/internal/model"
)

// CreatePosition godoc
//
//	@Summary	Создание новой должности
//	@Tags		position
//	@Produce	json
//	@Param		object	body		model.PositionSet	true	"Position Create"
//	@Success	201		{object}	model.Position
//	@Failure	400		{object}	sErr
//	@Failure	500		{object}	sErr
//	@Router		/positions [post]
func (r *RestServer) handlerCreatePosition(c *gin.Context) {
	ctx := c.Request.Context()
	positionReq := model.PositionSet{}

	if err := c.ShouldBindJSON(&positionReq); err != nil {
		c.JSON(http.StatusBadRequest, s().SetError(err))
		return
	}

	if err := positionReq.Validation(); err != nil {
		c.JSON(http.StatusBadRequest, s().SetError(err))
		return
	}

	us := r.getUserSession(c)
	if us.OrgID != positionReq.CompanyID {
		c.JSON(http.StatusBadRequest, s().SetError(errs.ErrCompanyNotFound))
		return
	}

	position, err := r.services.Position().CreatePosition(ctx, positionReq)
	switch {
	case errors.Is(err, model.ErrCompanyIDNotFound):
		c.JSON(http.StatusBadRequest, s().SetError(err))
		return
	case err != nil:
		c.JSON(http.StatusInternalServerError, s().SetError(err))
		return
	}

	c.JSON(http.StatusCreated, position)
}

// GetPosition godoc
//
//	@Summary	Получение всех должностей
//	@Tags		position
//	@Produce	json
//	@Param		id	path		int	true	"Position ID"
//	@Success	200	{object}	model.Position
//	@Failure	404	{object}	sErr
//	@Failure	500	{object}	sErr
//	@Router		/positions/{id} [get]
func (r *RestServer) handlerGetPosition(c *gin.Context) {
	ctx := c.Request.Context()

	val := c.Param("id")

	id, err := strconv.Atoi(val)
	if err != nil {
		c.JSON(http.StatusBadRequest, s().SetError(err))
		return
	}

	us := r.getUserSession(c)

	position, err := r.services.Position().GetPosition(ctx, us.OrgID, id)
	switch {
	case errors.Is(err, model.ErrPositionNotFound):
		c.JSON(http.StatusNotFound, s().SetError(err))
		return
	case err != nil:
		c.JSON(http.StatusInternalServerError, s().SetError(err))

		return
	}

	c.JSON(http.StatusOK, position)
}

// GetPositions godoc
//
//	@Summary	Получение всех должностей
//	@Tags		position
//	@Produce	json
//	@Success	200	{array}		model.Position
//	@Failure	404	{object}	sErr
//	@Failure	500	{object}	sErr
//	@Router		/positions [get]
func (r *RestServer) handlerGetPositions(c *gin.Context) {
	ctx := c.Request.Context()
	us := r.getUserSession(c)

	positions, err := r.services.Position().GetPositions(ctx, us.OrgID)

	switch {
	case errors.Is(err, model.ErrPositionsNotFound):
		c.JSON(http.StatusNotFound, s().SetError(err))
		return
	case err != nil:
		c.JSON(http.StatusInternalServerError, s().SetError(err))
		return
	}

	c.JSON(http.StatusOK, positions)
}

// UpdatePosition godoc
//
//	@Summary	Обновление данных о должности
//	@Tags		position
//	@Produce	json
//	@Param		id		path		int					true	"Position ID"
//	@Param		object	body		model.PositionSet	true	"Position info"
//	@Success	200		{object}	model.Position
//	@Failure	400		{object}	sErr
//	@Failure	404		{object}	sErr
//	@Failure	500		{object}	sErr
//	@Router		/positions/update/{id} [patch]
func (r *RestServer) handlerUpdatePosition(c *gin.Context) {
	ctx := c.Request.Context()
	positionReq := model.PositionSet{}

	val := c.Param("id")

	id, err := strconv.Atoi(val)
	if err != nil {
		c.JSON(http.StatusBadRequest, s().SetError(err))
		return
	}

	if err := c.ShouldBindJSON(&positionReq); err != nil {
		c.JSON(http.StatusBadRequest, s().SetError(err))
		return
	}

	if err = positionReq.Validation(); err != nil {
		c.JSON(http.StatusBadRequest, s().SetError(err))
		return
	}

	position, err := r.services.Position().UpdatePosition(ctx, id, positionReq)
	switch {
	case errors.Is(err, model.ErrPositionNotFound):
		c.JSON(http.StatusNotFound, s().SetError(err))
		return
	case err != nil:
		c.JSON(http.StatusInternalServerError, s().SetError(err))
		return
	}

	c.JSON(http.StatusOK, position)
}

// @Summary	Присвоение курса к должности
// @Accept		json
// @Tags		position
// @Produce	json
// @Success	200
// @Failure	400	{object}	error	"Неверный формат запроса"
// @Failure	401	{object}	error	"Пользователь не является сотрудником компании"
// @Failure	500	{object}	error	"Внутренняя ошибка сервера"
// @Router		/positions/course [post]
func (r *RestServer) handlerAssignCourse(c *gin.Context) {
	positionCourse := model.PositionCourse{}
	if err := c.ShouldBindJSON(&positionCourse); err != nil {
		c.JSON(http.StatusBadRequest, s().SetError(err))
		return
	}

	ctx := c.Request.Context()
	us := r.getUserSession(c)

	if err := r.services.Position().AssignCourse(ctx, positionCourse.PositionID,
		positionCourse.CourseID, us.UserID); err != nil {
		if errors.Is(err, model.ErrNoAuthorized) {
			c.JSON(http.StatusUnauthorized, s().SetError(err))
			return
		}
		c.JSON(http.StatusInternalServerError, s().SetError(err))
		return
	}
	c.Status(http.StatusOK)
}
