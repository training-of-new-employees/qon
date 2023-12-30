package impl

import (
	"context"
	"errors"

	"github.com/training-of-new-employees/qon/internal/errs"
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

func (cs *courseService) CreateCourse(ctx context.Context, c model.CourseSet) (*model.Course, error) {
	err := c.Validation()
	if err != nil {
		return nil, err
	}
	return cs.db.CourseStorage().CreateCourse(ctx, c)
}

func (cs *courseService) EditCourse(ctx context.Context, c model.CourseSet) (*model.Course, error) {
	err := c.Validation()
	if !errors.Is(err, errs.ErrCourseNameNotEmpty) && err != nil {
		return nil, err
	}
	return cs.db.CourseStorage().EditCourse(ctx, c)
}
