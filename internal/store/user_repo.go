package store

import (
	"context"

	"github.com/training-of-new-employees/qon/internal/model"
)

// RepositoryUser - интерфейс репозитория пользователя.
type RepositoryUser interface {
	// CreateAdmin - создание администратора по данным пользователя и названию компании
	CreateAdmin(context.Context, model.UserCreate, string) (*model.User, error)
	// CreateUser - создание пользователя
	CreateUser(context.Context, model.UserCreate) (*model.User, error)
	// EditAdmin - редактирование данных администратора
	EditAdmin(context.Context, model.AdminEdit) (*model.AdminEdit, error)
	// GetUserByEmail - получение данных пользователя по емейлу
	GetUserByEmail(context.Context, string) (*model.User, error)
	// GetUsersByCompany - получение данных пользователей определённой компании
	GetUsersByCompany(ctx context.Context, companyID int) ([]model.User, error)
	// EditUser - редактирование данных пользователя
	EditUser(ctx context.Context, edit *model.UserEdit) (*model.UserEdit, error)
	// SetPasswordAndActivateUser - установить пароль и активировать пользователя
	SetPasswordAndActivateUser(context.Context, int, string) error
	// UpdateUserPassword - обновить пароль пользователя
	UpdateUserPassword(context.Context, int, string) error
	// GetUserByID - получение данных пользователя по ИД
	GetUserByID(context.Context, int) (*model.User, error)
}
