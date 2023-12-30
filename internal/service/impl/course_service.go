package impl

import (
	"context"

	"github.com/training-of-new-employees/qon/internal/model"
	"github.com/training-of-new-employees/qon/internal/service"
	"github.com/training-of-new-employees/qon/internal/store"
)

var _ service.ServiceCourse = (*courseService)(nil)

type courseService struct {
	db store.Storages
}

func newCourseService(db store.Storages) *courseService {
	return &courseService{
		db: db,
	}
}

func (cs *courseService) GetCourses(ctx context.Context, u model.User) ([]model.Course, error) {
	if !u.IsAdmin {
		return cs.db.CourseStorage().PositionCourses(ctx, u.ID)
	}
	return cs.db.CourseStorage().CompanyCourses(ctx, u.CompanyID)
}

func (cs *courseService) CreateCourse(ctx context.Context, c model.CourseSet, creatorID int) (*model.Course, error) {
	err := c.Validation()
	if err != nil {
		return nil, err
	}
	return cs.db.CourseStorage().CreateCourse(ctx, c, creatorID)
}
