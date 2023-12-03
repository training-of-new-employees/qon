// Package service - пакет, который содержит бизнес-логику приложения.
package service

import (
	"context"
	"github.com/training-of-new-employees/qon/internal/model"
)

// Service - общий сервис.
type Service interface {
	User() ServiceUser
	// NOTE: здесь определяются интерфейсы для других сущностей
	// ...
}

type ServiceUser interface {
	//RegisterUser(ctx context.Context)
	RegisterAdmin(ctx context.Context, admin model.CreateAdmin) (*model.User, error)
	// TODO: здесь нужно описать все методы для ServiceUser
	// ...
}
