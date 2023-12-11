package impl

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/training-of-new-employees/qon/internal/logger"
	"github.com/training-of-new-employees/qon/internal/model"
	"github.com/training-of-new-employees/qon/internal/pkg/doar"
	"github.com/training-of-new-employees/qon/internal/pkg/jwttoken"
	"github.com/training-of-new-employees/qon/internal/service"
	"github.com/training-of-new-employees/qon/internal/store"
	"github.com/training-of-new-employees/qon/internal/store/cache"
)

var _ service.ServiceUser = (*uService)(nil)

// uService - структура, которая реализует интерфейс ServiceUser.
type uService struct {
	db         store.Storages
	cache      cache.Cache
	secretKey  string
	aTokenTime time.Duration
	rTokenTime time.Duration
	tokenGen   jwttoken.JWTGenerator
	tokenVal   jwttoken.JWTValidator
	sender     doar.EmailSender
}

func newUserService(
	db store.Storages,
	secretKey string,
	aTokenTime time.Duration,
	rTokenTime time.Duration,
	cache cache.Cache,
	jwtGen jwttoken.JWTGenerator,
	jwtVal jwttoken.JWTValidator,
	sender doar.EmailSender,
) *uService {
	return &uService{
		db:         db,
		secretKey:  secretKey,
		cache:      cache,
		aTokenTime: aTokenTime,
		rTokenTime: rTokenTime,
		tokenGen:   jwtGen,
		tokenVal:   jwtVal,
		sender:     sender,
	}
}

func (u *uService) WriteAdminToCache(
	ctx context.Context,
	val model.CreateAdmin,
) (*model.CreateAdmin, error) {

	if err := val.SetPassword(); err != nil {
		return nil, fmt.Errorf("error SetPassword: %v", err)
	}

	user, err := u.GetUserByEmail(ctx, val.Email)
	if err != nil {
		return nil, err
	}

	if user.Email == val.Email {
		return nil, model.ErrEmailAlreadyExists
	}

	key := uuid.New().String()
	if err := u.cache.Set(ctx, key, val); err != nil {
		return nil, err
	}

	logger.Log.Info("cache write successful", zap.String("key", key))

	if err = u.sender.SendEmail(val.Email, key); err != nil {
		return nil, err
	}

	return &val, nil
}

func (u *uService) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {

	userResp, err := u.db.UserStorage().GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("error failed GetUserByEmail %w", err)
	}

	return userResp, nil
}

func (u *uService) GenerateTokenPair(ctx context.Context, userId int,	isAdmin bool,	companyId int) (*model.Tokens, error) {

	accessToken, err := u.tokenGen.GenerateToken(userId, isAdmin, companyId, u.aTokenTime)
	if err != nil {
		return nil, fmt.Errorf("error failed GenerateToken %v", err)
	}

	refreshToken, err := u.tokenGen.GenerateToken(userId, isAdmin, companyId, u.rTokenTime)
	if err != nil {
		return nil, fmt.Errorf("error failed GenerateToken %v", err)
	}

	tokens := model.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	//TODO save token

	return &tokens, nil
}

func (u *uService) CreateUser(ctx context.Context, val model.UserCreate) (*model.User, error) {

	if err := val.SetPassword(); err != nil {
		return nil, fmt.Errorf("err SetPassword: %w", err)
	}

	createdUser, err := u.db.UserStorage().CreateUser(ctx, val)
	if err != nil {
		return nil, fmt.Errorf("err CreateUser")
	}

	return createdUser, nil
}

func (u *uService) GetAdminFromCache(ctx context.Context, key string) (*model.CreateAdmin, error) {

	admin, err := u.cache.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("err GetAdminFromCache: %v", err)
	}

	return admin, nil
}

func (u *uService) DeleteAdminFromCache(ctx context.Context, key string) error {

	if err := u.cache.Delete(ctx, key); err != nil {
		return fmt.Errorf("err DeleteAdminFromCache: %v", err)
	}

	return nil
}

func (u *uService) CreateAdmin(ctx context.Context, val *model.CreateAdmin) (*model.User, error) {

	user, err := u.GetUserByEmail(ctx, val.Email)
	if err != nil {
		return nil, err
	}

	if user.Email == val.Email {
		return nil, model.ErrEmailAlreadyExists
	}

	admin := model.NewAdminCreate(val.Email, val.Password)

	createdAdmin, err := u.db.UserStorage().CreateAdmin(ctx, admin, val.Company)
	if err != nil {
		return nil, fmt.Errorf("error creating admin: %w", err)
	}

	return createdAdmin, nil
}

// EditAdmin - Валидирует полученные данные и меняет их в БД, если всё впорядке
func (u *uService) EditAdmin(
	ctx context.Context,
	val *model.AdminEdit,
) (*model.AdminEdit, error) {

	err := val.Validation()
	if err != nil {
		return nil, err
	}

	edited, err := u.db.UserStorage().EditAdmin(ctx, val)
	if err != nil {
		return nil, fmt.Errorf("can't edit user: %w", err)
	}
	return edited, nil

}
