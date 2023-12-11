package rest

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/training-of-new-employees/qon/internal/model"
	"net/http"
	"strconv"
)

func (r *RestServer) handlerCreatePosition(c *gin.Context) {
	ctx := c.Request.Context()
	positionReq := model.PositionCreate{}

	if err := c.ShouldBindJSON(&positionReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := positionReq.Validation(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	us := r.getUserSession(c)
	if us.OrgID != positionReq.CompanyID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "StatusBadRequest"})
		return
	}

	position, err := r.services.Position().CreatePosition(ctx, positionReq)
	switch {
	case errors.Is(err, model.ErrCompanyIDNotFound):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	case err != nil:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, position)
}

func (r *RestServer) handlerGetPosition(c *gin.Context) {
	ctx := c.Request.Context()

	val := c.Param("id")

	id, err := strconv.Atoi(val)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	us := r.getUserSession(c)

	position, err := r.services.Position().GetPosition(ctx, us.OrgID, id)
	switch {
	case errors.Is(err, model.ErrPositionNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	case err != nil:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, position)
}

func (r *RestServer) handlerGetPositions(c *gin.Context) {
	ctx := c.Request.Context()
	us := r.getUserSession(c)

	positions, err := r.services.Position().GetPositions(ctx, us.OrgID)

	switch {
	case errors.Is(err, model.ErrPositionsNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	case err != nil:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, positions)
}

func (r *RestServer) handlerUpdatePosition(c *gin.Context) {
	ctx := c.Request.Context()
	positionReq := model.PositionUpdate{}

	val := c.Param("id")

	id, err := strconv.Atoi(val)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := c.ShouldBindJSON(&positionReq); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err = positionReq.Validation(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	us := r.getUserSession(c)
	if us.OrgID != positionReq.CompanyID {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	position, err := r.services.Position().UpdatePosition(ctx, id, positionReq)
	switch {
	case errors.Is(err, model.ErrPositionNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	case err != nil:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, position)
}

func (r *RestServer) handlerDeletePosition(c *gin.Context) {
	ctx := c.Request.Context()

	val := c.Param("id")

	id, err := strconv.Atoi(val)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	us := r.getUserSession(c)

	err = r.services.Position().DeletePosition(ctx, id, us.OrgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": "deleted"})
}
