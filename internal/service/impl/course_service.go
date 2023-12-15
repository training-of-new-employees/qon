package impl

import (
	"context"
	"fmt"
	"github.com/training-of-new-employees/qon/internal/model"
	"github.com/training-of-new-employees/qon/internal/service"
	"github.com/training-of-new-employees/qon/internal/store"
)

var _ service.ServiceCourse = (*courseService)(nil)

type courseService struct {
	db store.Storages
}

func newCourseService(db store.Storages) *courseService {
	return &courseService{db: db}
}

func (c *courseService) CreateCourse(ctx context.Context, adminID int, val model.CourseCreate) (*model.Course, error) {

	course, err := c.db.CourseStorage().CreateCourseDB(ctx, adminID, val)
	if err != nil {
		return nil, fmt.Errorf("failed to create course %w", err)
	}

	return course, nil
}

func (c *courseService) GetCourse(ctx context.Context, id int, adminID int) (*model.Course, error) {
	course, err := c.db.CourseStorage().GetCourseDB(ctx, id, adminID)
	if err != nil {
		return nil, fmt.Errorf("failed to get course %w", err)
	}

	return course, nil
}

func (c *courseService) GetCourses(ctx context.Context, userID int, companyID int) ([]*model.Course, error) {
	courses, err := c.db.CourseStorage().GetCoursesDB(ctx, userID, companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get courses %w", err)
	}

	return courses, nil
}
