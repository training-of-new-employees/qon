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

var _ store.RepositoryCourse = (*courseStorage)(nil)

type courseStorage struct {
	db *sqlx.DB
}

func newCourseStorage(db *sqlx.DB) *courseStorage {
	return &courseStorage{db: db}
}

func (c *courseStorage) CreateCourseDB(ctx context.Context, adminID int, val model.CourseCreate) (*model.Course, error) {
	course := model.Course{}

	query := `INSERT INTO courses (name, description, created_by)
			  VALUES ($1, $2, $3)
			  RETURNING id, name, description, active, created_by, archived, created_at, updated_at`

	err := c.db.GetContext(ctx, &course, query, val.Name, val.Description, adminID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.ForeignKeyViolation {
			return nil, model.ErrAdminIDNotFound
		}

		return &model.Course{}, fmt.Errorf("create course: %w", err)
	}

	return &course, nil
}

func (c *courseStorage) GetCourseDB(ctx context.Context, id int, adminID int) (*model.Course, error) {
	course := model.Course{}

	query := `SELECT c.id, c.name, c.description, c.created_by, c.active, c.archived, c.created_at, c.updated_at
              FROM courses c
			  JOIN users u ON c.created_by = u.id
			  WHERE c.created_by = $1 AND u.archived = false AND c.id = $2 AND c.archived = false`

	err := c.db.GetContext(ctx, &course, query, adminID, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &model.Course{}, model.ErrCourseNotFound
		}

		return &model.Course{}, fmt.Errorf("get course: %w", err)

	}

	return &course, nil
}

func (c *courseStorage) GetCoursesDB(ctx context.Context, userId int, companyID int) ([]*model.Course, error) {
	courses := make([]*model.Course, 0)

	query := `SELECT c.id, c.name, c.description, c.created_by, c.active, c.archived, c.created_at, c.updated_at
			  FROM courses c
			  JOIN users u ON c.created_by = u.id
			  WHERE c.created_by = $1 AND u.company_id = $2 AND u.archived = false AND c.archived = false`

	err := c.db.SelectContext(ctx, &courses, query, userId, companyID)
	if err != nil {
		return []*model.Course{}, fmt.Errorf("get courses: %w", err)
	}

	if len(courses) == 0 {
		return nil, model.ErrCoursesNotFound
	}

	return courses, nil
}
