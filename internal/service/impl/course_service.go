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

func (cs *courseService) GetUserCourses(ctx context.Context, userID int) ([]model.Course, error) {
	return cs.db.CourseStorage().UserCourses(ctx, userID)

}
func (cs *courseService) GetCompanyCourses(ctx context.Context, companyID int) ([]model.Course, error) {
	return cs.db.CourseStorage().CompanyCourses(ctx, companyID)
}

func (cs *courseService) CreateCourse(ctx context.Context, c model.CourseSet) (*model.Course, error) {
	return cs.db.CourseStorage().CreateCourse(ctx, c)
}

func (cs *courseService) EditCourse(ctx context.Context, c model.CourseSet, companyID int) (*model.Course, error) {
	err := c.Validation()
	if !errors.Is(err, errs.ErrCourseNameIsEmpty) && err != nil {
		return nil, err
	}
	return cs.db.CourseStorage().EditCourse(ctx, c, companyID)
}
