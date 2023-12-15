package rest

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/training-of-new-employees/qon/internal/model"
	"net/http"
	"strconv"
)

// @Summary Создание курса
// @Produce json
// @Param object body model.CourseCreate true "Course"
// @Success 201 {object} model.Course
// @Failure 400 {object} Error
// @Failure 401 {object} Error
// @Failure 403 {object} Error
// @Failure 500 {object} Error
// @Router /api/v1/course [post]
func (r *RestServer) handlerCreateCourse(c *gin.Context) {
	ctx := c.Request.Context()
	courseReq := model.CourseCreate{}

	if err := c.ShouldBindJSON(&courseReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := courseReq.Validation(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	us := r.getUserSession(c)

	course, err := r.services.Course().CreateCourse(ctx, us.UserID, courseReq)
	switch {
	case errors.Is(err, model.ErrAdminIDNotFound):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	case err != nil:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, course)
}

// @Summary Получение курса
// @Produce json
// @Param id path int true "Course ID"
// @Success 200 {object} model.Course
// @Failure 400 {object} Error
// @Failure 401 {object} Error
// @Failure 403 {object} Error
// @Failure 404 {object} Error
// @Failure 500 {object} Error
// @Router /api/v1/course/{id} [get]
func (r *RestServer) handlerGetCourse(c *gin.Context) {
	ctx := c.Request.Context()

	val := c.Param("id")

	id, err := strconv.Atoi(val)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	us := r.getUserSession(c)

	course, err := r.services.Course().GetCourse(ctx, id, us.UserID)
	switch {
	case errors.Is(err, model.ErrCourseNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	case err != nil:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, course)
}

// @Summary Получение списка курсов
// @Produce json
// @Success 200 {array} model.Course
// @Failure 401 {object} Error
// @Failure 403 {object} Error
// @Failure 404 {object} Error
// @Failure 500 {object} Error
// @Router /api/v1/course/ [get]
func (r *RestServer) handlerGetCourses(c *gin.Context) {
	ctx := c.Request.Context()

	us := r.getUserSession(c)

	courses, err := r.services.Course().GetCourses(ctx, us.UserID, us.OrgID)
	switch {
	case errors.Is(err, model.ErrCoursesNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	case err != nil:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, courses)
}
