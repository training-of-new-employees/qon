package impl

import (
	"context"

	"github.com/training-of-new-employees/qon/internal/model"
	"github.com/training-of-new-employees/qon/internal/service"
	"github.com/training-of-new-employees/qon/internal/store"
)

var _ service.ServiceLesson = (*lessonService)(nil)

type lessonService struct {
	db store.Storages
}

func newLessonService(db store.Storages) *lessonService {
	return &lessonService{db: db}
}

func (l *lessonService) CreateLesson(ctx context.Context, lesson model.Lesson, userID int) (*model.Lesson, error) {
	createdLesson, err := l.db.LessonStorage().CreateLesson(ctx,
		lesson, userID)
	if err != nil {
		return nil, err
	}
	return createdLesson, nil
}

func (l *lessonService) GetLesson(ctx context.Context, lessonID int) (*model.Lesson, error) {
	lesson, err := l.db.LessonStorage().GetLesson(ctx, lessonID)
	if err != nil {
		return nil, err
	}
	return lesson, nil
}

func (l *lessonService) UpdateLesson(ctx context.Context, lesson model.LessonUpdate) (*model.Lesson, error) {
	updatedLesson, err := l.db.LessonStorage().UpdateLesson(ctx, lesson)
	if err != nil {
		return nil, err
	}
	return updatedLesson, nil
}

func (l *lessonService) GetLessonsList(ctx context.Context, courseID int) ([]model.Lesson, error) {
	lessonsList, err := l.db.LessonStorage().GetLessonsList(ctx, courseID)
	if err != nil {
		return nil, err
	}
	return lessonsList, nil
}
