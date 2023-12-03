package pg

import (
	"context"
	"database/sql"
	"github.com/training-of-new-employees/qon/internal/model"
	"github.com/training-of-new-employees/qon/internal/store"
)

var _ store.RepositoryUser = (*uStorages)(nil)

type uStorages struct {
	db *sql.DB
}

func newUStorages(db *sql.DB) *uStorages {
	return &uStorages{
		db: db,
	}
}

func (u *uStorages) CreateAdmin(ctx context.Context, admin model.CreateAdmin) (*model.User, error) {
	return nil, nil
}

func (u *uStorages) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	return nil, nil
}
