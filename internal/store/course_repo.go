package store

import (
	"context"
	"github.com/training-of-new-employees/qon/internal/model"
)

type RepositoryCourse interface {
	CreateCourseDB(ctx context.Context, adminID int, course model.CourseCreate) (*model.Course, error)
	GetCourseDB(ctx context.Context, id int, adminID int) (*model.Course, error)
	GetCoursesDB(ctx context.Context, userID int, companyID int) ([]*model.Course, error)
}
