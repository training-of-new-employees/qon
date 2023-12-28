package store

import (
	"context"

	"github.com/training-of-new-employees/qon/internal/model"
)

type RepositoryCourse interface {
	PositionCourses(ctx context.Context, userID int) ([]model.Course, error)
	CompanyCourses(ctx context.Context, companyID int) ([]model.Course, error)
}
