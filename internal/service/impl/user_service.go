package impl

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/training-of-new-employees/qon/internal/errs"
	"github.com/training-of-new-employees/qon/internal/logger"
	"github.com/training-of-new-employees/qon/internal/model"
	"github.com/training-of-new-employees/qon/internal/pkg/doar"
	"github.com/training-of-new-employees/qon/internal/pkg/jwttoken"
	"github.com/training-of-new-employees/qon/internal/pkg/randomseq"
	"github.com/training-of-new-employees/qon/internal/service"
	"github.com/training-of-new-employees/qon/internal/store"
	"github.com/training-of-new-employees/qon/internal/store/cache"
	"github.com/training-of-new-employees/qon/internal/utils"
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

	code := randomseq.RandomDigitNumber(4)
	key := strings.Join([]string{"register", "admin", code}, ":")

	if err := u.cache.Set(ctx, key, val); err != nil {
		return nil, err
	}

	logger.Log.Info("cache write successful", zap.String("key", key))

	if err = u.sender.SendCode(val.Email, code); err != nil {
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

// GetUserByID получает информацию о сотруднике, а также имя компании и должность в ней
func (u *uService) GetUserByID(ctx context.Context, id int) (*model.UserInfo, error) {
	user, err := u.db.UserStorage().GetUserByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("can't get user by service: %w", err)
	}
	company, err := u.db.UserStorage().GetCompany(ctx, user.CompanyID)
	if err != nil {
		return nil, fmt.Errorf("can't get user by service: %w", err)
	}
	position, err := u.db.PositionStorage().GetPositionDB(ctx, user.CompanyID, user.PositionID)
	if err != nil {
		return nil, fmt.Errorf("can't get user by service: %w", err)
	}
	info := &model.UserInfo{
		User:         *user,
		CompanyName:  company.Name,
		PositionName: position.Name,
	}
	return info, nil
}

func (u *uService) ArchiveUser(ctx context.Context, id int, editorCompanyID int) error {
	user, err := u.db.UserStorage().GetUserByID(ctx, id)
	if err != nil {
		return err
	}
	if user.CompanyID != editorCompanyID {
		return errs.ErrUserNotFound
	}
	val := &model.UserEdit{
		ID:         id,
		IsArchived: true,
	}
	_, err = u.db.UserStorage().EditUser(ctx, val)
	return err

}

func (u *uService) EditUser(ctx context.Context, val *model.UserEdit, editorCompanyID int) (*model.UserEdit, error) {
	user, err := u.db.UserStorage().GetUserByID(ctx, val.ID)
	if err != nil {
		return nil, err
	}
	if user.CompanyID != editorCompanyID {
		return nil, errs.ErrUserNotFound
	}
	return u.db.UserStorage().EditUser(ctx, val)
}

// GetUsersByCompany - получает данные о пользователях в компании
func (u *uService) GetUsersByCompany(ctx context.Context, companyID int) ([]model.User, error) {
	users, err := u.db.UserStorage().GetUsersByCompany(ctx, companyID)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (u *uService) GenerateTokenPair(ctx context.Context, userId int, isAdmin bool, companyID int) (*model.Tokens, error) {

	accessToken, err := u.tokenGen.GenerateToken(userId, isAdmin, companyID, u.aTokenTime)
	if err != nil {
		return nil, fmt.Errorf("error failed GenerateToken %v", err)
	}

	refreshToken, err := u.tokenGen.GenerateToken(userId, isAdmin, companyID, u.rTokenTime)
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
		return nil, err
	}

	// TODO: генерирация пригласительной ссылки
	link := fmt.Sprintf("https://sample?email=%s", val.Email)

	// Отправление пригласительной ссылки сотруднику
	if err = u.sender.InviteUser(val.Email, link); err != nil {
		logger.Log.Warn(fmt.Sprintf("Не удалось отправить пригласительную ссылку сотруднику с емейлом %s", val.Email))
	}

	return createdUser, nil
}

func (u *uService) GetAdminFromCache(ctx context.Context, code string) (*model.CreateAdmin, error) {
	key := strings.Join([]string{"register", "admin", code}, ":")

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

// UpdatePasswordAndActivateUser устанавливает пароль и активирует учётную запись пользователя.
func (u *uService) UpdatePasswordAndActivateUser(ctx context.Context, email string, password string) error {
	user, err := u.GetUserByEmail(ctx, email)
	if err != nil {
		return err
	}
	if user.Email == "" {
		return model.ErrUserNotFound
	}

	encPassword, err := utils.EncryptPassword(password)
	if err != nil {
		return err
	}

	if err = u.db.UserStorage().SetPasswordAndActivateUser(ctx, user.ID, encPassword); err != nil {
		return err
	}

	return nil
}

// ResetPassword сбрасывает пользовательский пароль и устанавливает новый.
func (u *uService) ResetPassword(ctx context.Context, email string) error {
	user, err := u.GetUserByEmail(ctx, email)
	if err != nil {
		return err
	}
	if user.Email == "" {
		return model.ErrUserNotFound
	}

	password := model.GeneratePassword()

	encPassword, err := model.GenerateHash(password)
	if err != nil {
		return err
	}

	// Обновление пароля пользователя
	if err = u.db.UserStorage().UpdateUserPassword(ctx, user.ID, encPassword); err != nil {
		return err
	}

	// Отправление пароля пользователю
	if err = u.sender.SendPassword(email, password); err != nil {
		return err
	}

	return nil
}

// EditAdmin - Валидирует полученные данные и меняет их в БД, если всё впорядке
func (u *uService) EditAdmin(ctx context.Context, val *model.AdminEdit) (*model.AdminEdit, error) {

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
