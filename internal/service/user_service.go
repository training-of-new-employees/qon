package service

import (
	"context"

	"github.com/training-of-new-employees/qon/internal/model"
)

type ServiceUser interface {
	CreateUser(ctx context.Context, user model.UserCreate) (*model.User, error)
	WriteAdminToCache(ctx context.Context, admin model.CreateAdmin) (*model.CreateAdmin, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetUserByID(ctx context.Context, id int) (*model.UserInfo, error)
	GetUsersByCompany(ctx context.Context, companyID int) ([]model.User, error)
	EditUser(ctx context.Context, val *model.UserEdit, editorCompanyID int) (*model.UserEdit, error)
	ArchiveUser(ctx context.Context, id int, editorCompanyID int) error
	GenerateTokenPair(ctx context.Context, userId int, isAdmin bool, companyID int) (*model.Tokens, error)
	CreateAdmin(ctx context.Context, val *model.CreateAdmin) (*model.User, error)
	GetAdminFromCache(context.Context, string) (*model.CreateAdmin, error)
	DeleteAdminFromCache(ctx context.Context, key string) error
	UpdatePasswordAndActivateUser(ctx context.Context, email string, password string) error
	ResetPassword(ctx context.Context, email string) error
	EditAdmin(ctx context.Context, val model.AdminEdit) (*model.AdminEdit, error)
}
