package service

import (
	"context"
	"github.com/training-of-new-employees/qon/internal/model"
)

type ServiceCourse interface {
	CreateCourse(ctx context.Context, adminID int, course model.CourseCreate) (*model.Course, error)
	GetCourse(ctx context.Context, id int, adminID int) (*model.Course, error)
	GetCourses(ctx context.Context, userID int, companyID int) ([]*model.Course, error)
}
