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
	lesson model.Lesson, userID int) (*model.Lesson, error) {
	var createdLesson *model.Lesson
	var err error

	err = l.tx(func(tx *sqlx.Tx) error {
		createdLesson, err = l.createLessonTx(ctx, tx, lesson, userID)
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
func (l *lessonStorage) createLessonTx(ctx context.Context, tx *sqlx.Tx, lesson model.Lesson, userId int) (*model.Lesson, error) {
	query := `
		INSERT INTO lessons (course_id, created_by, name )
		VALUES ($1, $2, $3)
		RETURNING id, course_id, name
	`

	var createdLesson model.Lesson

	err := tx.GetContext(ctx, &createdLesson, query, lesson.CourseID, userId, lesson.Name)
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
	query := `
		SELECT
			l.id, l.course_id, l.name, t.content, p.url_picture, l.archived
		FROM lessons l
		JOIN texts t ON l.id = t.lesson_id
		JOIN pictures p ON  l.id = p.lesson_id
		WHERE l.id = $1
	`
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

func (l *lessonStorage) UpdateLesson(ctx context.Context, lesson model.LessonUpdate) (*model.Lesson, error) {
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

	query := `
		SELECT id
		FROM lessons
		WHERE id = $1 AND archived = false
	`
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

	query = `
		UPDATE lessons
		SET name = COALESCE(NULLIF($1, ''), name)
		WHERE id = $2
		RETURNING id, course_id, name, archived
	`
	err = tx.GetContext(ctx, &updatedLesson,
		query, lesson.Name, lesson.ID)
	if err != nil {
		return nil, err
	}

	return &updatedLesson, nil
}

// GetLessonsListDB - получить список уроков курса.
func (l *lessonStorage) GetLessonsList(ctx context.Context, courseID int) ([]model.Lesson, error) {
	lessonsList := []model.Lesson{}

	if courseID == 0 {
		return nil, errs.ErrCourseIDNotEmpty
	}

	query := `
		SELECT
			l.id, l.course_id, l.name, t.content, p.url_picture, l.archived
		FROM lessons l
		JOIN texts t ON l.id = t.lesson_id
		JOIN pictures p ON  l.id = p.lesson_id
		WHERE l.course_id = $1
	`
	err := l.db.SelectContext(ctx, &lessonsList, query, courseID)
	if err != nil {
		return nil, handleError(err)
	}

	return lessonsList, nil
}

func (l *lessonStorage) GetUserLessonsStatus(ctx context.Context, userID int, courseID int, lessonsIds []int) (map[int]string, error) {
	statuses := make(map[int]string)

	err := l.tx(func(tx *sqlx.Tx) error {
		var errInTx error
		statuses, errInTx = l.transaction.getUserLessonsStatusTx(ctx, tx, userID, courseID, lessonsIds)
		if errInTx != nil {
			return errInTx
		}

		return nil
	})
	if err != nil {
		return nil, handleError(err)
	}

	return statuses, nil
}

func (l *lessonStorage) UpdateUserLessonStatus(ctx context.Context, userID, courseID, lessonID int, status string) error {
	err := l.tx(func(tx *sqlx.Tx) error {
		updateStatusQuery := `
			INSERT INTO lesson_results (user_id, lesson_id, course_id, status)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (course_id, lesson_id, user_id) DO UPDATE SET status = EXCLUDED.status
		`
		_, err := tx.ExecContext(ctx, updateStatusQuery, userID, lessonID, courseID, status)
		if err != nil {
			return err
		}

		return l.transaction.syncUserCourseProgressTx(ctx, tx, userID, courseID)
	})

	if err != nil {
		return handleError(err)
	}

	return nil
}
