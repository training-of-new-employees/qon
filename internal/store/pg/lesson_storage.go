package pg

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/training-of-new-employees/qon/internal/errs"
	"github.com/training-of-new-employees/qon/internal/logger"
	"github.com/training-of-new-employees/qon/internal/model"
	"github.com/training-of-new-employees/qon/internal/store"
	"go.uber.org/zap"
)

var _ store.RepositoryLesson = (*lessonStorage)(nil)

type lessonStorage struct {
	db *sqlx.DB
	transaction
}

func newLessonStorage(db *sqlx.DB) *lessonStorage {
	return &lessonStorage{
		db:          db,
		transaction: transaction{db: db},
	}
}

func (l *lessonStorage) CreateLessonDB(ctx context.Context,
	lesson model.LessonCreate, user_id int) (*model.Lesson, error) {

	query := `INSERT INTO lessons (course_id, created_by, name, content, url_picture)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, course_id, created_by, name, content, url_picture, active, archived,
	created_at, updated_at`

	var createdLesson model.Lesson

	err := l.db.GetContext(ctx, &createdLesson, query,
		lesson.CourseID, user_id, lesson.Name, lesson.Content, lesson.URLPicture)
	if err != nil {
		return nil, handleError(err)
	}

	return &createdLesson, nil
}

func (l *lessonStorage) DeleteLessonDB(ctx context.Context, lessonID int) error {

	query := `UPDATE lessons SET archived = true WHERE id = $1`

	var result sql.Result
	var err error

	if result, err = l.db.ExecContext(ctx, query, lessonID); err != nil {
		fmt.Println(err.Error())
		return handleError(err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errs.ErrLessonNotFound
	}

	return nil
}

func (l *lessonStorage) GetLessonDB(ctx context.Context,
	lessonID int) (*model.Lesson, error) {
	query := `SELECT id, course_id, created_by, active, archived, name, 
			         content, url_picture, created_at, updated_at
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
	var updatedLesson *model.Lesson
	var err error

	err = l.tx(func(tx *sqlx.Tx) error {
		updatedLesson, err = l.updateLessonTx(ctx, tx, lesson)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, handleError(err)
	}

	return updatedLesson, nil
}

// updateLessonTx - обновление урока.
// ВAЖНО: использовать только внутри транзакции.
func (l *lessonStorage) updateLessonTx(ctx context.Context,
	tx *sqlx.Tx, lesson model.LessonUpdate) (*model.Lesson, error) {
	updatedLesson := model.Lesson{}

	query := `SELECT id
				FROM lessons
				WHERE id = $1 AND course_id = $2 AND archived = false`
	_, err := tx.ExecContext(ctx, query, lesson.ID, lesson.CourseID)
	if err != nil {
		return nil, handleError(err)
	}

	query = `UPDATE lessons
			  	SET 
					name 		= COALESCE(NULLIF($1, ''), name),
					content     = COALESCE(NULLIF($2, ''), content),
					url_picture = COALESCE(NULLIF($3, ''), url_picture)
				WHERE id = $4 AND course_id = $5 `
	_, err = tx.ExecContext(ctx, query, lesson.Name, lesson.Content,
		lesson.URLPicture, lesson.ID, lesson.CourseID)
	if err != nil {
		return nil, handleError(err)
	}

	query = `SELECT id, course_id, created_by, name, 
			        content, url_picture, active, archived,
					created_at, updated_at
			  FROM lessons
		      WHERE id = $1 AND course_id = $2`
	err = tx.GetContext(ctx, &updatedLesson, query, lesson.ID, lesson.CourseID)
	if err != nil {
		return nil, handleError(err)
	}
	return &updatedLesson, nil
}

// tx - обёртка для простого использования транзакций без дублирования кода.
func (l *lessonStorage) tx(f func(*sqlx.Tx) error) error {
	// открываем транзакцию
	tx, err := l.db.Beginx()
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

// GetLessonsListDB - получить список уроков курса.
func (l *lessonStorage) GetLessonsListDB(ctx context.Context, courseID int) ([]*model.Lesson, error) {
	lessonsList := []*model.Lesson{}

	query := `
	SELECT id, course_id, created_by, active, archived, name, 
	content, url_picture, created_at, updated_at
	FROM lessons
	WHERE course_id = $1 AND archived = false`

	err := l.db.SelectContext(ctx, &lessonsList, query, courseID)
	if err != nil {
		return nil, handleError(err)
	}

	if len(lessonsList) == 0 {
		return nil, errs.ErrLessonNotFound
	}

	return lessonsList, nil
}
