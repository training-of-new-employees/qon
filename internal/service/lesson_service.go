package service

import (
	"context"

	"github.com/training-of-new-employees/qon/internal/model"
)

type ServiceLesson interface {
	CreateLesson(ctx context.Context, lesson model.Lesson, userID int) (*model.Lesson, error)
	GetLesson(ctx context.Context, lessonID int) (*model.Lesson, error)
	GetUserLesson(ctx context.Context, userID int, lessonID int) (*model.Lesson, error)
	UpdateLesson(ctx context.Context, lesson model.LessonUpdate) (*model.Lesson, error)
	GetLessonsList(ctx context.Context, courseID int, userID int) ([]model.Lesson, error)
	UpdateLessonStatus(ctx context.Context, userID int, lessonID int, status string) error
}
