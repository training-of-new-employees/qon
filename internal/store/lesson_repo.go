package store

import (
	"context"

	"github.com/training-of-new-employees/qon/internal/model"
)

type RepositoryLesson interface {
	CreateLessonDB(ctx context.Context, lesson model.LessonCreate,
		user_id int) (*model.Lesson, error)
	DeleteLessonDB(ctx context.Context, lessonID int) error
}
