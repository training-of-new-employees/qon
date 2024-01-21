package store

import (
	"context"

	"github.com/training-of-new-employees/qon/internal/model"
)

type RepositoryCourse interface {
	UserCourses(ctx context.Context, userID int) ([]model.Course, error)
	GetUserCourse(ctx context.Context, userID int, courseID int) (*model.Course, error)
	GetUserCoursesStatus(ctx context.Context, userID int, coursesIds []int) (map[int]string, error)
	CompanyCourses(ctx context.Context, companyID int) ([]model.Course, error)
	CreateCourse(ctx context.Context, course model.CourseSet) (*model.Course, error)
	EditCourse(ctx context.Context, course model.CourseSet, companyID int) (*model.Course, error)
}
