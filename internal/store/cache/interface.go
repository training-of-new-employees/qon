package cache

import (
	"context"

	"github.com/training-of-new-employees/qon/internal/model"
)

type Cache interface {
	// TODO: нужно переименовать методы Get и Set, чтобы было понятнее их назначение (для примера SetAdmin, GetAdmin).
	Get(ctx context.Context, key string) (*model.CreateAdmin, error)
	Set(ctx context.Context, uuid string, admin model.CreateAdmin) error

	SetInviteCode(ctx context.Context, key string, code string) error
	GetInviteCode(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
}
