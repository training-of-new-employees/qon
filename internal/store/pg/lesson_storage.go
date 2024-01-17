package pg

import (
	"context"
	"database/sql"
	"errors"
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

func (l *lessonStorage) CreateLesson(ctx context.Context,
	lesson model.Lesson, user_id int) (*model.Lesson, error) {
	var createdLesson *model.Lesson
	var err error

	err = l.tx(func(tx *sqlx.Tx) error {
		createdLesson, err = l.createLessonTx(ctx, tx, lesson, user_id)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, handleError(err)
	}

	return createdLesson, nil
}

// updateLessonTx - обновление урока.
// ВAЖНО: использовать только внутри транзакции.
func (l *lessonStorage) createLessonTx(ctx context.Context,
	tx *sqlx.Tx, lesson model.Lesson, userId int) (*model.Lesson, error) {

	query := `INSERT INTO lessons (course_id, created_by, name )
			  VALUES ($1, $2, $3)
		      RETURNING id, course_id, name`

	var createdLesson model.Lesson

	err := tx.GetContext(ctx, &createdLesson, query,
		lesson.CourseID, userId, lesson.Name)
	if err != nil {
		return nil, err
	}

	createdLesson.Content, err = l.insertTextsTx(ctx, tx, createdLesson.ID,
		lesson.Content, userId)
	if err != nil {
		return nil, err
	}

	createdLesson.URLPicture, err = l.insertPicturesTx(ctx, tx, createdLesson.ID,
		lesson.URLPicture, userId)
	if err != nil {
		return nil, err
	}
	return &createdLesson, nil
}

func (l *lessonStorage) GetLesson(ctx context.Context,
	lessonID int) (*model.Lesson, error) {
	query := `SELECT l.id, l.course_id, l.name, t.content,
					 p.url_picture, l.archived
			  FROM lessons l
			  JOIN texts t ON l.id = t.lesson_id
			  JOIN pictures p ON  l.id = p.lesson_id
		      WHERE l.id = $1`
	var lesson model.Lesson
	err := l.db.GetContext(ctx, &lesson, query, lessonID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrLessonNotFound
		}
		return nil, handleError(err)
	}
	return &lesson, nil
}

func (l *lessonStorage) UpdateLesson(ctx context.Context,
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
				WHERE id = $1 AND archived = false`
	_, err := tx.ExecContext(ctx, query, lesson.ID)
	if err != nil {
		return nil, err
	}

	content, err := l.updateTextsTx(ctx, tx, lesson.ID, lesson.Content)
	if err != nil {
		return nil, err
	}
	updatedLesson.Content = content

	urlPicture, err := l.updatePicturesTx(ctx, tx, lesson.ID, lesson.URLPicture)
	if err != nil {
		return nil, err
	}
	updatedLesson.URLPicture = urlPicture

	query = `UPDATE lessons
			  	SET name = COALESCE(NULLIF($1, ''), name)
				WHERE id = $2
				RETURNING name, archived`
	err = tx.GetContext(ctx, &updatedLesson,
		query, lesson.Name, lesson.ID)
	if err != nil {
		return nil, err
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
func (l *lessonStorage) GetLessonsList(ctx context.Context, courseID int) ([]*model.Lesson, error) {
	lessonsList := []*model.Lesson{}

	query := `SELECT l.id, l.course_id, l.name, t.content,
					 p.url_picture, l.archived
			  FROM lessons l
			  JOIN texts t ON l.id = t.lesson_id
			  JOIN pictures p ON  l.id = p.lesson_id
			  WHERE l.course_id = $1`

	err := l.db.SelectContext(ctx, &lessonsList, query, courseID)
	if err != nil {
		return nil, handleError(err)
	}

	if len(lessonsList) == 0 {
		return nil, errs.ErrLessonNotFound
	}

	return lessonsList, nil
}
