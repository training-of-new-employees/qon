// Package service - пакет, который содержит бизнес-логику приложения.
package service

import "context"

// Service - общий сервис.
type Service interface {
	ServiceUser
	// NOTE: здесь определяются интерфейсы для других сущностей
	// ...
}

type ServiceUser interface {
	RegisterUser(ctx context.Context)
	// TODO: здесь нужно описать все методы для ServiceUser
	// ...
}
