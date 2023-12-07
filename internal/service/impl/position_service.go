package impl

import (
	"context"
	"github.com/training-of-new-employees/qon/internal/model"
	"github.com/training-of-new-employees/qon/internal/store"
)

type positionService struct {
	db store.Storages
}

func newPositionService(db store.Storages) *positionService {
	return &positionService{db: db}
}

func (p *positionService) CreatePosition(ctx context.Context, position model.CreatePosition) (*model.Position, error) {
	return nil, nil
}
