package cache

import (
	"context"

	"github.com/training-of-new-employees/qon/internal/model"
)

type Cache interface {
	Get(ctx context.Context, key string) (*model.CreateAdmin, error)
	Set(ctx context.Context, uuid string, admin model.CreateAdmin) error
	Delete(ctx context.Context, key string) error
}
