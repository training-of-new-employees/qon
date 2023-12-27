package service

import (
	"context"

	"github.com/training-of-new-employees/qon/internal/model"
)

type ServiceCourse interface {
	GetCourses(ctx context.Context, u model.User) ([]model.Course, error)
}
