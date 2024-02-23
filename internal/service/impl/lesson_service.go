package impl

import (
	"context"

	"github.com/training-of-new-employees/qon/internal/errs"
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

func (l *lessonService) GetLesson(ctx context.Context, lessonID int, companyID int) (*model.Lesson, error) {
	lesson, err := l.db.LessonStorage().GetLesson(ctx, lessonID)
	if err != nil {
		return nil, err
	}

	_, err = l.db.CourseStorage().CompanyCourse(ctx, lesson.CourseID, companyID)
	if err != nil {
		return nil, errs.ErrLessonNotFound
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

func (l *lessonService) GetUserLesson(ctx context.Context, userID int, lessonID int) (*model.Lesson, error) {
	lesson, err := l.db.LessonStorage().GetLesson(ctx, lessonID)
	if err != nil {
		return nil, err
	}

	courses, err := l.db.CourseStorage().UserCourses(ctx, userID)
	if err != nil {
		return nil, err
	}

	isUserHasCourse := false
	courseID := 0

	for _, course := range courses {
		if course.ID == lesson.CourseID {
			isUserHasCourse = true
			courseID = course.ID
			break
		}
	}

	if !isUserHasCourse {
		return nil, errs.ErrCourseNotFound
	}

	statuses, err := l.db.LessonStorage().GetUserLessonsStatus(ctx, userID, courseID, []int{lessonID})
	if err != nil {
		return nil, err
	}

	lesson.Status = statuses[lessonID]
	return lesson, nil
}

func (l *lessonService) UpdateLessonStatus(ctx context.Context, userID int, lessonID int, status string) error {
	lesson, err := l.db.LessonStorage().GetLesson(ctx, lessonID)
	if err != nil {
		return err
	}

	courses, err := l.db.CourseStorage().UserCourses(ctx, userID)
	if err != nil {
		return err
	}

	isUserHasCourse := false

	for _, course := range courses {
		if course.ID == lesson.CourseID {
			isUserHasCourse = true
			break
		}
	}

	if !isUserHasCourse {
		return errs.ErrCourseNotFound
	}

	return l.db.LessonStorage().UpdateUserLessonStatus(ctx, userID, lesson.CourseID, lesson.ID, status)
}
