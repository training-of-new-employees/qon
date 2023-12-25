package pg

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/training-of-new-employees/qon/internal/logger"
	"github.com/training-of-new-employees/qon/internal/model"
	"github.com/training-of-new-employees/qon/internal/store"
	"go.uber.org/zap"
)

var _ store.RepositoryCompany = (*companyStorage)(nil)

type companyStorage struct {
	db    *sqlx.DB
	store *Store
}

func newCompanyStorage(db *sqlx.DB, s *Store) *companyStorage {
	return &companyStorage{db: db, store: s}
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

// createCompanyTx - создание компании в транзакции.
// ВAЖНО: использовать только внутри транзакции.
func (c *companyStorage) createCompanyTx(ctx context.Context, tx *sqlx.Tx, companyName string) (*model.Company, error) {
	company := model.Company{}

	query := `INSERT INTO companies (name) VALUES ($1) RETURNING id, name`

	if err := tx.GetContext(ctx, &company, query, companyName); err != nil {
		return nil, err
	}

	return &company, nil
}

// tx - обёртка для простого использования транзакций без дублирования кода.
func (c *companyStorage) tx(f func(*sqlx.Tx) error) error {
	// открываем транзакцию
	tx, err := c.db.Beginx()
	if err != nil {
		return fmt.Errorf("beginning tx: %w", err)
	}
	// отмена транзакции
	defer func() {
		if err := tx.Rollback(); err != nil {
			logger.Log.Warn("err during tx rollback %v", zap.Error(err))
		}
	}()

	if err = f(tx); err != nil {
		return err
	}

	// фиксация транзакции
	return tx.Commit()
}
