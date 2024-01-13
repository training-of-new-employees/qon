package cache

import (
	"context"

	"github.com/training-of-new-employees/qon/internal/model"
)

type Cache interface {
	Get(ctx context.Context, key string) (*model.CreateAdmin, error)
	Set(ctx context.Context, uuid string, admin model.CreateAdmin) error
	SetInviteCode(ctx context.Context, key string, code string) error
	GetInviteCode(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
}
