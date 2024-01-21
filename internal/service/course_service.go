package service

import (
	"context"

	"github.com/training-of-new-employees/qon/internal/model"
)

type ServiceCourse interface {
	GetUserCourses(ctx context.Context, userID int) ([]model.Course, error)
	GetUserCourseLessons(ctx context.Context, userID int, courseID int) ([]model.Lesson, error)
	GetCompanyCourses(ctx context.Context, companyID int) ([]model.Course, error)
	CreateCourse(ctx context.Context, c model.CourseSet) (*model.Course, error)
	EditCourse(ctx context.Context, c model.CourseSet, companyID int) (*model.Course, error)
}
