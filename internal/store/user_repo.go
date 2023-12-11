package store

import (
	"context"

	"github.com/training-of-new-employees/qon/internal/model"
)

type RepositoryUser interface {
	CreateAdmin(context.Context, model.AdminCreate, string) (*model.User, error)
	CreateUser(context.Context, model.UserCreate) (*model.User, error)
	EditAdmin(context.Context, *model.AdminEdit) (*model.AdminEdit, error)
	GetUserByEmail(context.Context, string) (*model.User, error)
	GetUserByID(context.Context, int) (*model.User, error)
}
