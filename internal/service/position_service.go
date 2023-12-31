package service

import (
	"context"

	"github.com/training-of-new-employees/qon/internal/model"
)

type ServicePosition interface {
	CreatePosition(ctx context.Context, position model.PositionSet) (*model.Position, error)
	GetPosition(ctx context.Context, companyID int, positionID int) (*model.Position, error)
	GetPositions(ctx context.Context, id int) ([]*model.Position, error)
	UpdatePosition(ctx context.Context, id int, position model.PositionSet) (*model.Position, error)
	AssignCourse(ctx context.Context, positionID int, courseID int, user_id int) error
}
