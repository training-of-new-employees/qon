package service

import (
	"context"

	"github.com/training-of-new-employees/qon/internal/model"
)

type ServiceUser interface {
	CreateUser(ctx context.Context, user model.UserCreate) (*model.User, error)
	WriteAdminToCache(ctx context.Context, admin model.CreateAdmin) (*model.CreateAdmin, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GenerateTokenPair(
		ctx context.Context,
		userId int,
		isAdmin bool,
		companyId int,
	) (*model.Tokens, error)
	CreateAdmin(ctx context.Context, val *model.CreateAdmin) (*model.User, error)
	GetAdminFromCache(context.Context, string) (*model.CreateAdmin, error)
	DeleteAdminFromCache(ctx context.Context, key string) error
	EditAdmin(ctx context.Context, val *model.AdminEdit) (*model.AdminEdit, error)
}
