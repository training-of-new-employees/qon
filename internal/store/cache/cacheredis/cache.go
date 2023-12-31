package cacheredis

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/redis/go-redis/v9"
	"github.com/training-of-new-employees/qon/internal/model"
	"github.com/training-of-new-employees/qon/internal/store/cache"
)

var _ cache.Cache = (*Redis)(nil)

var ErrWritingCache = errors.New("error writing to Redis")

var ErrDeletingCache = errors.New("error deleting from Redis")

type Redis struct {
	client *redis.Client
}

func NewRedis(client *redis.Client) *Redis {
	return &Redis{client: client}
}

func (r *Redis) Get(ctx context.Context, key string) (*model.CreateAdmin, error) {
	val, err := r.client.Get(ctx, key).Result()

	if err == redis.Nil {
		return &model.CreateAdmin{}, err
	}

	var admin model.CreateAdmin
	if err = json.Unmarshal([]byte(val), &admin); err != nil {
		return nil, err
	}

	return &admin, nil
}

func (r *Redis) Set(ctx context.Context, uuid string, admin model.CreateAdmin) error {
	adminJSON, err := json.Marshal(admin)
	if err != nil {
		return err
	}

	err = r.client.Set(ctx, uuid, adminJSON, 0).Err()
	if err != nil {
		return ErrWritingCache
	}

	return nil
}

func (r *Redis) Delete(ctx context.Context, key string) error {
	_, err := r.client.Del(ctx, key).Result()
	if err != nil {
		return ErrDeletingCache
	}

	return nil
}
