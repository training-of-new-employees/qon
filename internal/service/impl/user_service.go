package impl

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"errors"
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
	host       string
}

func newUserService(
	db store.Storages, secretKey string,
	aTokenTime time.Duration, rTokenTime time.Duration,
	cache cache.Cache, jwtGen jwttoken.JWTGenerator, jwtVal jwttoken.JWTValidator,
	sender doar.EmailSender, host string,
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
		host:       host,
	}
}

func (u *uService) WriteAdminToCache(ctx context.Context, val model.CreateAdmin) (*model.CreateAdmin, error) {
	if err := val.SetPassword(); err != nil {
		return nil, fmt.Errorf("error SetPassword: %v", err)
	}

	_, err := u.GetUserByEmail(ctx, val.Email)
	if err != nil && !errors.Is(err, errs.ErrUserNotFound) && !errors.Is(err, errs.ErrNotFound) {
		return nil, err
	}
	if err == nil {
		return nil, errs.ErrEmailAlreadyExists
	}

	code := randomseq.RandomDigitNumber(4)
	key := strings.Join([]string{"register", "admin", code}, ":")

	if err := u.cache.Set(ctx, key, val); err != nil {
		return nil, err
	}

	logger.Log.Info("cache write successful", zap.String("key", key))

	if err = u.sender.SendCode(val.Email, code); err != nil {
		logger.Log.Error("service error", zap.Error(err))

		// режим мок-рассылки писем, при котором содержание письма выводится в теле пользователю
		if u.sender.Mode() == doar.TestMode {
			return nil, err
		}
	}

	return &val, nil
}

func (u *uService) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	userResp, err := u.db.UserStorage().GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return userResp, nil
}

// GetUserByID получает информацию о сотруднике, а также имя компании и должность в ней
func (u *uService) GetUserByID(ctx context.Context, id int) (*model.UserInfo, error) {
	user, err := u.db.UserStorage().GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	company, err := u.db.CompanyStorage().GetCompany(ctx, user.CompanyID)
	if err != nil {
		return nil, err
	}
	position, err := u.db.PositionStorage().GetPositionInCompany(ctx, user.CompanyID, user.PositionID)
	if err != nil {
		return nil, err
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
	// нельзя архивировать учётную запись админа
	if user.IsAdmin {
		return errs.ErrArchiveAdmin
	}
	// указанный пользователь не является сотрудником компании
	if user.CompanyID != editorCompanyID {
		return errs.ErrUserNotFound
	}

	valTrue := true
	val := &model.UserEdit{ID: id, IsArchived: &valTrue}

	_, err = u.db.UserStorage().EditUser(ctx, val)
	return err

}

func (u *uService) EditUser(ctx context.Context, val *model.UserEdit, editorCompanyID int) (*model.UserEdit, error) {
	user, err := u.db.UserStorage().GetUserByID(ctx, val.ID)
	if err != nil {
		return nil, err
	}
	// нельзя архивировать админа
	if user.IsAdmin && *val.IsArchived {
		return nil, errs.ErrArchiveAdmin
	}
	val.CompanyID = &user.CompanyID
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

	refreshToken, err := u.tokenGen.GenerateToken(userId, isAdmin, companyID, "", u.rTokenTime)
	if err != nil {
		return nil, fmt.Errorf("error failed GenerateToken %v", err)
	}

	hasher := sha1.New()
	hasher.Write([]byte(refreshToken))
	hashedRefresh := hex.EncodeToString(hasher.Sum(nil))

	accessToken, err := u.tokenGen.GenerateToken(userId, isAdmin, companyID, hashedRefresh, u.aTokenTime)
	if err != nil {
		return nil, fmt.Errorf("error failed GenerateToken %v", err)
	}

	tokens := model.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	if err := u.cache.SetRefreshToken(ctx, hashedRefresh, refreshToken); err != nil {
		return nil, fmt.Errorf("error failed SetRefreshToken %w", err)
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

	link, err := u.GenerateInvitationLinkUser(ctx, val.Email)
	if err != nil {
		logger.Log.Warn(fmt.Sprintf("Не удалось с генерировать пригласительную ссылку сотруднику с емейлом %s", val.Email))
	}

	// Отправление пригласительной ссылки сотруднику
	if err = u.sender.InviteUser(val.Email, val.Name, link); err != nil {
		logger.Log.Warn(fmt.Sprintf("Не удалось отправить пригласительную ссылку сотруднику с емейлом %s", val.Email))

		// режим мок-рассылки писем, при котором содержание письма выводится в теле пользователю
		if u.sender.Mode() == doar.TestMode {
			return nil, err
		}
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

// CreateAdmin создаёт администратора в БД
func (u *uService) CreateAdmin(ctx context.Context, val model.CreateAdmin) (*model.User, error) {

	_, err := u.GetUserByEmail(ctx, val.Email)
	if err != nil && !errors.Is(err, errs.ErrUserNotFound) {
		return nil, err
	}
	if err == nil {
		return nil, errs.ErrEmailAlreadyExists
	}

	admin := model.NewAdminCreate(val.Email, val.Password)

	createdAdmin, err := u.db.UserStorage().CreateAdmin(ctx, admin, val.Company)
	if err != nil {
		return nil, fmt.Errorf("error creating admin: %w", err)
	}

	return createdAdmin, nil
}

// UpdatePasswordAndActivateUser устанавливает пароль и активирует учётную запись пользователя.
func (u *uService) UpdatePasswordAndActivateUser(ctx context.Context, email string, password string) (*model.User, error) {
	user, err := u.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if user.IsActive {
		return nil, errs.ErrNotFirstLogin
	}

	encPassword, err := utils.EncryptPassword(password)
	if err != nil {
		return nil, err
	}

	if err = u.db.UserStorage().SetPasswordAndActivateUser(ctx, user.ID, encPassword); err != nil {
		return nil, err
	}

	return user, nil
}

// ResetPassword сбрасывает пользовательский пароль и устанавливает новый.
func (u *uService) ResetPassword(ctx context.Context, email string) error {
	user, err := u.GetUserByEmail(ctx, email)
	if err != nil {
		return err
	}

	password := randomseq.RandomPassword()

	encPassword, err := model.GenerateHash(password)
	if err != nil {
		return err
	}

	// Обновление пароля пользователя
	if err = u.db.UserStorage().UpdateUserPassword(ctx, user.ID, encPassword); err != nil {
		return err
	}

	// Отправление пароля пользователю
	if err = u.sender.SendPassword(email, user.Name, password, fmt.Sprintf("%s/login", u.host)); err != nil {
		return err
	}

	return nil
}

// EditAdmin - Валидирует полученные данные и меняет их в БД, если всё впорядке
func (u *uService) EditAdmin(ctx context.Context, val model.AdminEdit) (*model.AdminEdit, error) {
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

func (u *uService) GenerateInvitationLinkUser(ctx context.Context, email string) (string, error) {
	code := randomseq.RandomString(20)
	key := strings.Join([]string{"register", "user", email}, ":")

	if err := u.cache.SetInviteCode(ctx, key, code); err != nil {
		return "", err
	}

	logger.Log.Info("cache write successful", zap.String("key", key))

	link := fmt.Sprintf("%s/first-login?email=%s&invite=%s", u.host, email, code)

	logger.Log.Info("Generate invite link successful", zap.String("invite link", link))

	return link, nil
}

func (u *uService) GetUserInviteCodeFromCache(ctx context.Context, email string) (string, error) {
	key := strings.Join([]string{"register", "user", email}, ":")

	code, err := u.cache.GetInviteCode(ctx, key)
	if err != nil {
		logger.Log.Warn("err GetUserInviteFromCache: %v", zap.Error(err))

		return "", errs.ErrInvalidInviteCode
	}

	return code, nil
}

func (u *uService) RegenerationInvitationLinkUser(ctx context.Context, email string, companyID int) (*model.InvitationLinkResponse, error) {
	invitationLinkResponse := &model.InvitationLinkResponse{}

	employee, err := u.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if employee.IsActive {

		return nil, errs.ErrUserActivated
	}

	if employee.CompanyID != companyID {
		return nil, errs.ErrEmployeeHasAnotherCompany
	}

	link, err := u.GenerateInvitationLinkUser(ctx, email)
	if err != nil {
		return nil, errs.ErrInternal
	}

	invitationLinkResponse.Link = link
	invitationLinkResponse.Email = email

	// Отправление пригласительной ссылки сотруднику
	if err = u.sender.InviteUser(email, employee.Name, link); err != nil {
		logger.Log.Warn(fmt.Sprintf("Не удалось отправить пригласительную ссылку сотруднику с емейлом %s", email))

		// режим мок-рассылки писем, при котором содержание письма выводится в теле пользователю
		if u.sender.Mode() == doar.TestMode {
			return nil, err
		}
	}

	return invitationLinkResponse, nil
}

func (u *uService) GetInvitationLinkUser(ctx context.Context, email string, companyID int) (*model.InvitationLinkResponse, error) {
	invitationLinkResponse := &model.InvitationLinkResponse{}

	employee, err := u.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if employee.IsActive {

		return nil, errs.ErrUserActivated
	}

	if employee.CompanyID != companyID {
		return nil, errs.ErrNoAccess
	}

	code, err := u.GetUserInviteCodeFromCache(ctx, employee.Email)
	if err != nil {
		return nil, err
	}

	invitationLinkResponse.Link = fmt.Sprintf("%s/first-login?email=%s&invite=%s", u.host, email, code)
	invitationLinkResponse.Email = email

	return invitationLinkResponse, nil
}

func (u *uService) ClearSession(ctx context.Context, hashedRefreshToken string) error {
	return u.cache.DeleteRefreshToken(ctx, hashedRefreshToken)
}
