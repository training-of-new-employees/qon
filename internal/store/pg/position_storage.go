package pg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/training-of-new-employees/qon/internal/logger"
	"github.com/training-of-new-employees/qon/internal/model"
	"github.com/training-of-new-employees/qon/internal/store"
	"go.uber.org/zap"
)

var _ store.RepositoryPosition = (*positionStorage)(nil)

type positionStorage struct {
	db    *sqlx.DB
	store *Store
}

func newPositionStorage(db *sqlx.DB, s *Store) *positionStorage {
	return &positionStorage{db: db, store: s}
}

// CreatePositionDB создание должности в рамках компании.
func (p *positionStorage) CreatePositionDB(ctx context.Context, position model.PositionCreate) (*model.Position, error) {
	// открываем транзакцию
	tx, err := p.db.Beginx()
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

	// создание должности
	createdPosition, err := p.createPositionTx(ctx, tx, position.CompanyID, position.Name)
	if err != nil {
		return nil, handleError(err)
	}

	// фиксация транзакции
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("committing tx: %w", err)
	}

	return createdPosition, nil
}

// GetPositionDB - получить данные должности, привязанной к компании.
func (p *positionStorage) GetPositionDB(ctx context.Context, companyID int, positionID int) (*model.Position, error) {
	position := model.Position{}

	query := `SELECT p.id, p.company_id, p.name, p.active, p.archived, p.created_at, p.updated_at
              FROM positions p 
              JOIN companies c ON p.company_id = c.id
              WHERE p.company_id = $1 AND p.id = $2 AND p.archived = false`

	err := p.db.GetContext(ctx, &position, query, companyID, positionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrPositionNotFound
		}

		return &model.Position{}, fmt.Errorf("get position db: %w", err)
	}

	return &position, nil
}

// GetPositionsDB - получить список должностей компании.
func (p *positionStorage) GetPositionsDB(ctx context.Context, companyID int) ([]*model.Position, error) {
	positions := make([]*model.Position, 0)

	query := `SELECT p.id, p.company_id, p.name, p.active, p.archived, p.created_at, p.updated_at
			  FROM positions p
			  JOIN companies c ON p.company_id = c.id
			  WHERE p.company_id = $1 AND c.active = true AND p.archived = false`

	err := p.db.SelectContext(ctx, &positions, query, companyID)
	if err != nil {
		return []*model.Position{}, fmt.Errorf("get positions db: %w", err)
	}

	if len(positions) == 0 {
		return nil, model.ErrPositionsNotFound
	}

	return positions, nil
}

// UpdatePositionDB - обновить компанию.
func (p *positionStorage) UpdatePositionDB(ctx context.Context, positionID int, val model.PositionUpdate) (*model.Position, error) {
	position := model.Position{}

	query := `UPDATE positions SET name = $1 WHERE id = $2 AND company_id = $3
              RETURNING id, name, company_id, active, archived, updated_at, created_at`

	err := p.db.QueryRowContext(ctx, query, val.Name, positionID, val.CompanyID).Scan(&position.ID, &position.Name,
		&position.CompanyID, &position.IsActive, &position.IsArchived, &position.UpdatedAt, &position.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {

			return &model.Position{}, model.ErrPositionNotFound
		}

		return &model.Position{}, fmt.Errorf("update position db: %w", err)
	}

	return &position, nil
}

func (p *positionStorage) DeletePositionDB(ctx context.Context, id int, companyID int) error {
	query := `UPDATE positions SET archived = true WHERE id = $1 AND company_id = $2`

	if _, err := p.db.ExecContext(ctx, query, id, companyID); err != nil {
		return fmt.Errorf("delete position db: %w", err)
	}

	return nil
}

// GetPositionByID - получить данные должности по идентификатору.
func (p *positionStorage) GetPositionByID(ctx context.Context, positionID int) (*model.Position, error) {
	position := model.Position{}

	query := `SELECT id, company_id, name, active, archived, created_at, updated_at
              FROM positions 
              WHERE id = $1 AND archived = false`

	err := p.db.GetContext(ctx, &position, query, positionID)
	if err != nil {
		return &model.Position{}, handleError(err)
	}

	return &position, nil
}

// AssignCourseDB - назначить курс на должность.
func (p *positionStorage) AssignCourseDB(ctx context.Context, positionID int,
	courseID int, user_id int) error {

	query := `INSERT INTO position_course (position_id, course_id)
			  VALUES ($1, $2)`
	if _, err := p.db.ExecContext(ctx, query, positionID, courseID); err != nil {
		return handleError(err)
	}
	return nil
}

// createPositionTx - создание должности в рамках компании.
// ВAЖНО: использовать только внутри транзакции.
func (p *positionStorage) createPositionTx(ctx context.Context, tx *sqlx.Tx, companyID int, positionName string) (*model.Position, error) {
	position := model.Position{}

	query := `INSERT INTO
				positions (company_id, name) VALUES ($1, $2)
			  RETURNING id, company_id, active, name, created_at, updated_at`

	if err := tx.GetContext(ctx, &position, query, companyID, positionName); err != nil {
		return nil, err
	}

	return &position, nil
}
