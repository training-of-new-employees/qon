package store

import (
	"context"

	"github.com/training-of-new-employees/qon/internal/model"
)

type RepositoryCourse interface {
	UserCourses(ctx context.Context, userID int) ([]model.Course, error)
	CompanyCourses(ctx context.Context, companyID int) ([]model.Course, error)
	CreateCourse(ctx context.Context, course model.CourseSet) (model.Course, error)
	EditCourse(ctx context.Context, course model.CourseSet, companyID int) (model.Course, error)
}
