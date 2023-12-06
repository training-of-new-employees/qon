package service

import (
	"context"
	"github.com/training-of-new-employees/qon/internal/model"
)

type ServicePosition interface {
	CreatePosition(ctx context.Context, position model.CreatePosition) (*model.Position, error)
}
