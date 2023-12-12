package store

import (
	"context"

	"github.com/training-of-new-employees/qon/internal/model"
)

type RepositoryPosition interface {
	CreatePositionDB(ctx context.Context, position model.PositionCreate) (*model.Position, error)
	GetPositionDB(ctx context.Context, companyID int, positionID int) (*model.Position, error)
	GetPositionsDB(ctx context.Context, id int) ([]*model.Position, error)
	UpdatePositionDB(ctx context.Context, id int, position model.PositionUpdate) (*model.Position, error)
	DeletePositionDB(ctx context.Context, id int, companyID int) error
	AssignCourseDB(ctx context.Context, positionID int, courseID int) error
}
