package service

import (
	"context"

	"github.com/training-of-new-employees/qon/internal/model"
)

type ServiceCourse interface {
	GetCourses(ctx context.Context, u model.User) ([]model.Course, error)
	CreateCourse(ctx context.Context, c model.CourseSet) (model.Course, error)
	EditCourse(ctx context.Context, c model.CourseSet, companyID int) (model.Course, error)
}
