package impl

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/training-of-new-employees/qon/internal/model"
	"github.com/training-of-new-employees/qon/internal/pkg/jwttoken"
	"github.com/training-of-new-employees/qon/internal/service"
	"github.com/training-of-new-employees/qon/internal/store"
	"github.com/training-of-new-employees/qon/internal/store/cache"
	"time"
)

var _ service.ServiceUser = (*uService)(nil)

type uService struct {
	db         store.Storages
	cache      cache.Cache
	secretKey  string
	aTokenTime time.Duration
	rTokenTime time.Duration
	tokenGen   jwttoken.JWTGenerator
	tokenVal   jwttoken.JWTValidator
}

func newUserService(db store.Storages, secretKey string, aTokenTime time.Duration,
	rTokenTime time.Duration, cache cache.Cache, jwtGen jwttoken.JWTGenerator, jwtVal jwttoken.JWTValidator) *uService {
	return &uService{
		db:         db,
		secretKey:  secretKey,
		cache:      cache,
		aTokenTime: aTokenTime,
		rTokenTime: rTokenTime,
		tokenGen:   jwtGen,
		tokenVal:   jwtVal,
	}
}

func (u *uService) RegisterAdmin(ctx context.Context, admin model.CreateAdmin) (*model.User, error) {
	if err := admin.SetPassword(); err != nil {
		return nil, fmt.Errorf("err SetPassword: %v", err)
	}

	_, err := u.GetUserByEmail(ctx, admin.Email)
	if err != nil {
		return nil, fmt.Errorf("error failed getAdminByEmail %w", err)
	}

	key := uuid.New().String()
	if err := u.cache.Set(ctx, key, admin); err != nil {
		return nil, err
	}

	createdAdmin, err := u.db.UserStorage().CreateAdmin(ctx, admin)
	if err != nil {
		return nil, fmt.Errorf("err CreateAdmin")
	}

	return createdAdmin, nil
}

func (u *uService) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	userResp, err := u.db.UserStorage().GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("error failed getUserByEmail %w", err)
	}

	return userResp, nil
}
