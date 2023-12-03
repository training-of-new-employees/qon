// Package store - пакет для работы с хранилищем.
package store

import (
	"context"
	"github.com/training-of-new-employees/qon/internal/model"
)

type Storages interface {
	UserStorage() RepositoryUser
	// TODO: Добавить интерфейсы для работы с другими объектами
}

// RepositoryUser определяет интерфейс для репозитория пользователей.
type RepositoryUser interface {
	CreateAdmin(context.Context, model.CreateAdmin) (*model.User, error)
	GetUserByEmail(context.Context, string) (*model.User, error)
	// TODO: здесь нужно описать все методы для RepositoryUser
}
