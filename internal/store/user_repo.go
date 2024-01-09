package store

import (
	"context"

	"github.com/training-of-new-employees/qon/internal/model"
)

// RepositoryUser - интерфейс репозитория пользователя.
type RepositoryUser interface {
	CreateUser(ctx context.Context, user model.UserCreate) (*model.User, error)

	// CreateAdmin - создание администратора по данным пользователя и названию компании
	CreateAdmin(ctx context.Context, admin model.UserCreate, companyName string) (*model.User, error)

	GetUserByID(context.Context, int) (*model.User, error)

	GetUserByEmail(ctx context.Context, email string) (*model.User, error)

	GetUsersByCompany(ctx context.Context, companyID int) ([]model.User, error)

	EditUser(ctx context.Context, edit *model.UserEdit) (*model.UserEdit, error)

	EditAdmin(ctx context.Context, edit model.AdminEdit) (*model.AdminEdit, error)

	UpdateUserPassword(context.Context, int, string) error

	// SetPasswordAndActivateUser - установить пароль и активировать пользователя
	SetPasswordAndActivateUser(context.Context, int, string) error
}
