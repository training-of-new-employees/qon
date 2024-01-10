package pg

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"github.com/training-of-new-employees/qon/internal/errs"
	"github.com/training-of-new-employees/qon/internal/model"
	"github.com/training-of-new-employees/qon/internal/store"
)

var _ store.RepositoryCompany = (*companyStorage)(nil)

// companyStorage - репозиторий для компании/организации.
type companyStorage struct {
	db *sqlx.DB
	transaction
}

// newCompanyStorage - конструктор репозитория компании/организации.
func newCompanyStorage(db *sqlx.DB) *companyStorage {
	return &companyStorage{
		db:          db,
		transaction: transaction{db: db},
	}
}

// createCompany - создание компании/организации.
func (c *companyStorage) CreateCompany(ctx context.Context, companyName string) (*model.Company, error) {
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

// GetCompany - получение информации о компании/организации по id.
func (c *companyStorage) GetCompany(ctx context.Context, id int) (*model.Company, error) {
	var comp *model.Company

	// открываем транзакцию
	err := c.tx(func(tx *sqlx.Tx) error {
		var err error
		comp, err = c.getCompanyTx(ctx, tx, id)
		return err
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrCompanyNotFound
		}
		return nil, handleError(err)
	}

	return comp, nil
}
