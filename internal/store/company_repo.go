package store

import (
	"context"

	"github.com/training-of-new-employees/qon/internal/model"
)

type RepositoryCompany interface {
	CreateCompanyDB(ctx context.Context, companyName string) (*model.Company, error)
}
