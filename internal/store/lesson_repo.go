package store

import (
	"context"

	"github.com/training-of-new-employees/qon/internal/model"
)

type RepositoryLesson interface {
	// TODO: Не совсем понятно, зачем в методах используется окончание 'DB', возможно, его следует потом убрать.
	CreateLessonDB(ctx context.Context, lesson model.LessonCreate, user_id int) (*model.Lesson, error)
	DeleteLessonDB(ctx context.Context, lessonID int) error
	GetLessonDB(ctx context.Context, lessonID int) (*model.Lesson, error)
	UpdateLessonDB(ctx context.Context, lesson model.LessonUpdate) (*model.Lesson, error)
}
