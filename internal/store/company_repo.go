package store

import (
	"context"

	"github.com/training-of-new-employees/qon/internal/model"
)

// RepositoryCompany - интерфейс репозитория компании/организации.
type RepositoryCompany interface {
	// CreatePositionDB - создание компании
	CreateCompany(ctx context.Context, companyName string) (*model.Company, error)
	// CreatePositionDB - получение данных компании по ИД
	GetCompany(ctx context.Context, id int) (*model.Company, error)
}
