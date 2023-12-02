// Package store - пакет для работы с хранилищем.
package store

import "context"

// Store - интерфейс для хранилища.
type Store interface {
	RepositoryUser
	// TODO: Добавить интерфейсы для работы с другими объектами
}

// RepositoryUser определяет интерфейс для репозитория пользователей.
type RepositoryUser interface {
	CreateUser(ctx context.Context) error
	// TODO: здесь нужно описать все методы для RepositoryUser
}
