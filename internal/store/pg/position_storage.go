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

var _ store.RepositoryPosition = (*positionStorage)(nil)

// positionStorage - репозиторий должности.
type positionStorage struct {
	db *sqlx.DB
	transaction
}

// newPositionStorage - конструктор репозитория должности.
func newPositionStorage(db *sqlx.DB) *positionStorage {
	return &positionStorage{
		db:          db,
		transaction: transaction{db: db},
	}
}

// CreatePosition - создание должности в рамках компании.
func (p *positionStorage) CreatePosition(ctx context.Context, position model.PositionSet) (*model.Position, error) {
	var createdPosition *model.Position

	// открываем транзакцию
	err := p.tx(func(tx *sqlx.Tx) error {
		var err error

		// создание должности
		createdPosition, err = p.createPositionTx(ctx, tx, position.CompanyID, position.Name)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, handleError(err)
	}

	return createdPosition, nil
}

// GetPositionByID - получение данных должности по идентификатору.
func (p *positionStorage) GetPositionByID(ctx context.Context, positionID int) (*model.Position, error) {
	position := model.Position{}

	query := `
		SELECT
			id, company_id, name, active, archived, created_at, updated_at
        FROM positions 
        WHERE id = $1 AND archived = false
	`
	err := p.db.GetContext(ctx, &position, query, positionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrPositionNotFound
		}
		return nil, handleError(err)
	}

	return &position, nil
}

// GetPositionInCompany - получение данных должности, привязанной к компании.
func (p *positionStorage) GetPositionInCompany(ctx context.Context, companyID int, positionID int) (*model.Position, error) {
	position := &model.Position{}

	// открываем транзакцию
	err := p.tx(func(tx *sqlx.Tx) error {
		var err error

		// создание должности
		position, err = p.getPositionInCompanyTx(ctx, tx, companyID, positionID)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrPositionNotFound
		}
		return nil, handleError(err)
	}

	return position, nil
}

// ListPositions - получение списка должностей компании.
func (p *positionStorage) ListPositions(ctx context.Context, companyID int) ([]*model.Position, error) {
	positions := make([]*model.Position, 0)

	query := `
		SELECT
			p.id, p.company_id, p.name, p.active, p.archived, p.created_at, p.updated_at
		FROM positions p
		JOIN companies c ON p.company_id = c.id
		WHERE p.company_id = $1 AND c.active = true AND p.archived = false AND p.name != 'admin'
	`

	err := p.db.SelectContext(ctx, &positions, query, companyID)
	if err != nil {
		return nil, handleError(err)
	}

	return positions, nil
}

// UpdatePosition - обновление данных должности.
func (p *positionStorage) UpdatePosition(ctx context.Context, positionID int, val model.PositionSet) (*model.Position, error) {
	position := model.Position{}

	query := `
		UPDATE positions
		SET name = COALESCE($1, name), archived = COALESCE($2, archived)
		WHERE id = $3 AND company_id = $4
        RETURNING id, name, company_id, active, archived, updated_at, created_at
	`

	err := p.db.QueryRowContext(ctx, query, val.Name, val.IsArchived, positionID, val.CompanyID).Scan(
		&position.ID, &position.Name, &position.CompanyID,
		&position.IsActive, &position.IsArchived,
		&position.UpdatedAt, &position.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrPositionNotFound
		}
		return nil, handleError(err)
	}

	return &position, nil
}

// AssignCourse - назначение курса на должность.
func (p *positionStorage) AssignCourse(ctx context.Context, positionID int, courseID int) error {
	// открываем транзакцию
	err := p.tx(func(tx *sqlx.Tx) error {
		// назначение курса на должность
		return p.assignCourseTx(ctx, tx, positionID, courseID)
	})

	if err != nil {
		return handleError(err)
	}

	return nil
}

// AssignCourses - назначение нескольких курсов на должность.
func (p *positionStorage) AssignCourses(ctx context.Context, positionID int, courseIDs []int) error {
	err := p.tx(func(tx *sqlx.Tx) error {
		query := `DELETE FROM position_course WHERE position_id = $1`
		_, err := tx.ExecContext(ctx, query, positionID)
		if err != nil {
			return err
		}

		if len(courseIDs) == 0 {
			return nil
		}

		return p.assignCoursesTx(ctx, tx, positionID, courseIDs)
	})

	if err != nil {
		return handleError(err)
	}

	return nil
}

func (p *positionStorage) GetCourseForPosition(ctx context.Context, positionID int) ([]int, error) {
	query := `
		SELECT course_id
		FROM position_course
		WHERE position_id = $1
	`

	rows, err := p.db.QueryContext(ctx, query, positionID)
	if err != nil {
		return nil, handleError(err)
	}

	if rows.Err() != nil {
		return nil, handleError(rows.Err())
	}

	ids := make([]int, 0)

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, handleError(err)
		}

		ids = append(ids, id)
	}

	return ids, nil
}
