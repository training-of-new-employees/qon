package service

import (
	"context"
	"github.com/training-of-new-employees/qon/internal/model"
)

type ServicePosition interface {
	CreatePosition(ctx context.Context, position model.PositionCreate) (*model.Position, error)
	GetPosition(ctx context.Context, companyID int, positionID int) (*model.Position, error)
	GetPositions(ctx context.Context, id int) ([]*model.Position, error)
	UpdatePosition(ctx context.Context, id int, orgID int, position model.PositionUpdate) (*model.Position, error)
	DeletePosition(ctx context.Context, id int, companyID int) error
}
