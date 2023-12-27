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
	return nil, nil
}
