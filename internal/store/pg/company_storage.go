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
	// открываем транзакцию
	tx, err := c.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("beginning tx: %w", err)
	}
	// отмена транзакции
	defer func() {
		if err != nil {
			if err := tx.Rollback(); err != nil {
				logger.Log.Warn("err during tx rollback %v", zap.Error(err))
			}
		}
	}()

	// создание компании
	createdCompany, err := c.createCompanyTx(ctx, tx, companyName)
	if err != nil {
		return nil, handleError(err)
	}

	// фиксация транзакции
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("committing tx: %w", err)
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
