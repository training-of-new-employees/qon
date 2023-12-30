package rest

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/training-of-new-employees/qon/internal/errs"
	"github.com/training-of-new-employees/qon/internal/model"
)

func (r *RestServer) handlerGetCourses(c *gin.Context) {
	ctx := c.Request.Context()
	us := r.getUserSession(c)
	u := model.User{
		ID:        us.UserID,
		IsAdmin:   us.IsAdmin,
		CompanyID: us.OrgID,
	}
	courses, err := r.services.Course().GetCourses(ctx, u)
	if err != nil {
		r.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, courses)

}
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
func (r *RestServer) handlerEditCourse(c *gin.Context) {
	ctx := c.Request.Context()
	sID := c.Param("id")
	id, err := strconv.Atoi(sID)
	if err != nil {
		r.handleError(c, errs.ErrBadRequest)
		return
	}
	course := model.NewCourseSet(id, r.getUserSession(c).UserID)
	err = c.BindJSON(&course)
	if err != nil {
		r.handleError(c, errs.ErrBadRequest)
		return
	}
	_, err = r.services.Course().EditCourse(ctx, course)
	if err != nil {
		r.handleError(c, err)
		return
	}
	c.Status(http.StatusOK)

}
