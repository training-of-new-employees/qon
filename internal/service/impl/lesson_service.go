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

func (l *lessonService) CreateLesson(ctx context.Context,
	lesson model.LessonCreate, user_id int) (*model.Lesson, error) {
	createdLesson, err := l.db.LessonStorage().CreateLessonDB(ctx,
		lesson, user_id)
	if err != nil {
		return nil, err
	}
	return createdLesson, nil
}

func (l *lessonService) UpdateLesson() {

}

func (l *lessonService) DeleteLesson(ctx context.Context, lessonID int) error {
	if err := l.db.LessonStorage().DeleteLessonDB(ctx, lessonID); err != nil {
		return err
	}
	return nil
}

func (l *lessonService) GetLesson(ctx context.Context, lessonID int) (*model.Lesson, error) {
	return nil, nil
}
