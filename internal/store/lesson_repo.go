package store

import (
	"context"

	"github.com/training-of-new-employees/qon/internal/model"
)

type RepositoryLesson interface {
	CreateLesson(ctx context.Context, lesson model.Lesson, userID int) (*model.Lesson, error)
	GetLesson(ctx context.Context, lessonID int) (*model.Lesson, error)
	UpdateLesson(ctx context.Context, lesson model.LessonUpdate) (*model.Lesson, error)
	GetLessonsList(ctx context.Context, courseID int) ([]model.Lesson, error)
}
