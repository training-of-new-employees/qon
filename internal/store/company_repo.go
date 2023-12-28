package store

import (
	"context"

	"github.com/training-of-new-employees/qon/internal/model"
)

type RepositoryCompany interface {
	CreateCompany(ctx context.Context, companyName string) (*model.Company, error)
	GetCompany(ctx context.Context, id int) (*model.Company, error)
}
