package store

import (
	"context"

	"github.com/training-of-new-employees/qon/internal/model"
)

type RepositoryPosition interface {
	CreatePositionDB(ctx context.Context, position model.PositionSet) (*model.Position, error)
	GetPositionDB(ctx context.Context, companyID int, positionID int) (*model.Position, error)
	GetPositionsDB(ctx context.Context, companyID int) ([]*model.Position, error)
	GetPositionByID(ctx context.Context, positionID int) (*model.Position, error)
	UpdatePositionDB(ctx context.Context, id int, position model.PositionSet) (*model.Position, error)
	AssignCourseDB(ctx context.Context, positionID int, courseID int, user_id int) error
}
