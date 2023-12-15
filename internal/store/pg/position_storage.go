package pg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jmoiron/sqlx"
	"github.com/training-of-new-employees/qon/internal/model"
	"github.com/training-of-new-employees/qon/internal/store"
)

var _ store.RepositoryPosition = (*positionStorage)(nil)

type positionStorage struct {
	db *sqlx.DB
}

func newPositionStorage(db *sqlx.DB) *positionStorage {
	return &positionStorage{db: db}
}

func (p *positionStorage) CreatePositionDB(ctx context.Context, position model.PositionCreate) (*model.Position, error) {
	var createdPosition = model.Position{}

	query := `INSERT INTO positions (company_id, name)
			  VALUES ($1, $2)
			  RETURNING id, company_id, name, active, created_at, updated_at`

	err := p.db.GetContext(ctx, &createdPosition, query, position.CompanyID, position.Name)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.ForeignKeyViolation {
			return nil, model.ErrCompanyIDNotFound
		}

		return nil, fmt.Errorf("create position: %w", err)
	}

	return &createdPosition, nil
}

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

func (p *positionStorage) GetPositionsDB(ctx context.Context, id int) ([]*model.Position, error) {
	positions := make([]*model.Position, 0)

	query := `SELECT p.id, p.company_id, p.name, p.active, p.archived, p.created_at, p.updated_at
			  FROM positions p
			  JOIN companies c ON p.company_id = c.id
			  WHERE p.company_id = $1 AND c.active = true AND p.archived = false`

	err := p.db.SelectContext(ctx, &positions, query, id)
	if err != nil {
		return []*model.Position{}, fmt.Errorf("get positions db: %w", err)
	}

	if len(positions) == 0 {
		return nil, model.ErrPositionsNotFound
	}

	return positions, nil
}

func (p *positionStorage) UpdatePositionDB(ctx context.Context, id int, orgID int, val model.PositionUpdate) (*model.Position, error) {
	position := model.Position{}

	query := `UPDATE positions SET name = $1 WHERE id = $2 AND company_id = $3 AND archived = false
              RETURNING id, name, company_id, active, archived, updated_at, created_at`

	err := p.db.QueryRowContext(ctx, query, val.Name, id, orgID).Scan(&position.ID, &position.Name,
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
