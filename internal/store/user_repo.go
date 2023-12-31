package store

import (
	"context"

	"github.com/training-of-new-employees/qon/internal/model"
)

type RepositoryUser interface {
	CreateAdmin(context.Context, model.UserCreate, string) (*model.User, error)
	CreateUser(context.Context, model.UserCreate) (*model.User, error)
	EditAdmin(context.Context, model.AdminEdit) (*model.AdminEdit, error)
	GetUserByEmail(context.Context, string) (*model.User, error)
	GetUsersByCompany(ctx context.Context, companyID int) ([]model.User, error)
	EditUser(ctx context.Context, edit *model.UserEdit) (*model.UserEdit, error)
	SetPasswordAndActivateUser(context.Context, int, string) error
	UpdateUserPassword(context.Context, int, string) error
	GetUserByID(context.Context, int) (*model.User, error)
}
