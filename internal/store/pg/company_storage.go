package pg

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/training-of-new-employees/qon/internal/model"
	"github.com/training-of-new-employees/qon/internal/store"
)

var _ store.RepositoryCompany = (*companyStorage)(nil)

type companyStorage struct {
	db *sqlx.DB
	transaction
}

func newCompanyStorage(db *sqlx.DB) *companyStorage {
	return &companyStorage{
		db:          db,
		transaction: transaction{db: db},
	}
}

// createCompanyDB - создание компании.
func (c *companyStorage) CreateCompanyDB(ctx context.Context, companyName string) (*model.Company, error) {
	var createdCompany *model.Company

	// открываем транзакцию
	err := c.tx(func(tx *sqlx.Tx) error {
		var err error

		// создание должности
		createdCompany, err = c.createCompanyTx(ctx, tx, companyName)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, handleError(err)
	}

	return createdCompany, nil

}
