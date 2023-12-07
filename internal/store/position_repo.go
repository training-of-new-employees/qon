package store

import (
	"context"
	"github.com/training-of-new-employees/qon/internal/model"
)

type RepositoryPosition interface {
	CreatePosition(ctx context.Context, position model.CreatePosition) (*model.Position, error)
}
