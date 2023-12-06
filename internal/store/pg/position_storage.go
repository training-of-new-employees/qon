package pg

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/training-of-new-employees/qon/internal/model"
	"github.com/training-of-new-employees/qon/internal/store"
)

var _ store.RepositoryPosition = (*positionStorage)(nil)

type positionStorage struct {
	db *sqlx.DB
}

func newPositionStorage(db *sqlx.DB) *positionStorage {
	return &positionStorage{db: db}
}

func (s *positionStorage) CreatePosition(ctx context.Context, position model.CreatePosition) (*model.Position, error) {

	return nil, nil
}
