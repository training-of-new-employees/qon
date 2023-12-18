package pg

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jmoiron/sqlx"
	"github.com/training-of-new-employees/qon/internal/errs"
	"github.com/training-of-new-employees/qon/internal/model"
	"github.com/training-of-new-employees/qon/internal/store"
)

var _ store.RepositoryLesson = (*lessonStorage)(nil)

type lessonStorage struct {
	db *sqlx.DB
}

func newLessonStorage(db *sqlx.DB) *lessonStorage {
	return &lessonStorage{db: db}
}

func (l *lessonStorage) CreateLessonDB(ctx context.Context,
	lesson model.LessonCreate, user_id int) (*model.Lesson, error) {

	query := `INSERT INTO lessons (course_id, created_by, name, description)
	VALUES ($1, $2, $3, $4)
	RETURNING id, number, name, description, created_at`

	var createdLesson model.Lesson

	err := l.db.GetContext(ctx, &createdLesson, query,
		lesson.CourseID, user_id, lesson.Name, lesson.Description, true)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.ForeignKeyViolation {
			return nil, errs.ErrCourseNotFound
		}
		return nil, fmt.Errorf("create lesson: %w", err)
	}

	return &createdLesson, nil
}

func (l *lessonStorage) DeleteLessonDB(ctx context.Context, lessonID int) error {
	query := `UPDATE lessons SET archived = true WHERE id = $1`

	if _, err := l.db.ExecContext(ctx, query, lessonID); err != nil {
		return fmt.Errorf("delete position db: %w", err)
	}

	return nil
}
