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
		lesson.CourseID, user_id, lesson.Name, lesson.Description)
	if err != nil {
		return nil, handleError(err)
	}

	return &createdLesson, nil
}

func (l *lessonStorage) DeleteLessonDB(ctx context.Context, lessonID int) error {
	query := `UPDATE lessons SET archived = true WHERE id = $1`

	if _, err := l.db.ExecContext(ctx, query, lessonID); err != nil {
		return handleError(err)
	}

	return nil
}

func (l *lessonStorage) GetLessonDB(ctx context.Context,
	lessonID int) (*model.Lesson, error) {
	query := `SELECT id, course_id, created_by, number, name, 
			         description, created_at, updated_at
			  FROM lessons
		      WHERE id = $1 AND archived = false`
	var lesson model.Lesson
	err := l.db.GetContext(ctx, &lesson, query, lessonID)
	if err != nil {
		return nil, handleError(err)
	}
	return &lesson, nil
}

func (l *lessonStorage) UpdateLessonDB(ctx context.Context,
	lesson model.LessonUpdate) (*model.Lesson, error) {
	var updatedLesson model.Lesson
	var err error

	tx, err := l.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("beginning tx: %w", err)
	}

	defer func() {
		if err != nil {
			if err := tx.Rollback(); err != nil {
				logger.Log.Warn("err during tx rollback %v", zap.Error(err))
			}
		}
	}()

	query := `SELECT id
				FROM lessons
				WHERE id = $1 AND course_id = $2 AND archived = false`
	_, err = tx.ExecContext(ctx, query, lesson.ID, lesson.CourseID)
	if err != nil {
		return nil, handleError(err)
	}

	if lesson.Name != "" {
		query := `UPDATE lessons
			  	  SET name = $1 WHERE id = $2 AND course_id = $3 `
		_, err := tx.ExecContext(ctx, query, lesson.Name, lesson.ID, lesson.CourseID)
		if err != nil {
			return nil, handleError(err)
		}
	}

	if lesson.Description != "" {
		query := `UPDATE lessons
			  	  SET description = $1 WHERE id = $2 AND course_id = $3 `
		_, err := tx.ExecContext(ctx, query, lesson.Description, lesson.ID, lesson.CourseID)
		if err != nil {
			return nil, handleError(err)
		}
	}

	if lesson.Path != "" {
		query := `UPDATE lessons
			  	  SET path = $1 WHERE id = $2 AND course_id = $3 `
		_, err := tx.ExecContext(ctx, query, lesson.Path, lesson.ID, lesson.CourseID)
		if err != nil {
			return nil, handleError(err)
		}
	}

	query = `SELECT id, course_id, created_by, number, name, 
			         description, created_at, updated_at
			  FROM lessons
		      WHERE id = $1 AND course_id = $2`
	err = tx.GetContext(ctx, &updatedLesson, query, lesson.ID, lesson.CourseID)
	if err != nil {
		return nil, handleError(err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("committing tx: %w", err)
	}
	return &updatedLesson, nil
}
