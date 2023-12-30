package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"

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
	course := model.CourseSet{}
	err := c.BindJSON(&course)
	if err != nil {
		r.handleError(c, err)
		return
	}
	us := r.getUserSession(c)
	created, err := r.services.Course().CreateCourse(ctx, course, us.UserID)
	if err != nil {
		r.handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, created)

}
func (r *RestServer) handlerPutCourse(c *gin.Context) {
}
