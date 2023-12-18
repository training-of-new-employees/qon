package service

import (
	"context"

	"github.com/training-of-new-employees/qon/internal/model"
)

type ServiceLesson interface {
	CreateLesson(ctx context.Context, lesson model.LessonCreate,
		user_id int) (*model.Lesson, error)
	UpdateLesson()
	DeleteLesson(ctx context.Context, lessonID int) error
	GetLesson(ctx context.Context, lessonID int) (*model.Lesson, error)
}
