package impl

import (
	"context"
	"errors"

	"github.com/training-of-new-employees/qon/internal/errs"
	"github.com/training-of-new-employees/qon/internal/model"
	"github.com/training-of-new-employees/qon/internal/service"
	"github.com/training-of-new-employees/qon/internal/store"
)

var _ service.ServiceCourse = (*courseService)(nil)

type courseService struct {
	db store.Storages
}

func newCourseService(db store.Storages) *courseService {
	return &courseService{
		db: db,
	}
}

func (cs *courseService) GetUserCourses(ctx context.Context, userID int) ([]model.Course, error) {
	courses, err := cs.db.CourseStorage().UserCourses(ctx, userID)
	if err != nil {
		return nil, err
	}

	coursesIds := make([]int, 0, len(courses))
	for _, course := range courses {
		coursesIds = append(coursesIds, course.ID)
	}

	statuses, err := cs.db.CourseStorage().GetUserCoursesStatus(ctx, userID, coursesIds)
	if err != nil {
		return nil, err
	}

	for idx, course := range courses {
		courses[idx].Status = statuses[course.ID]
	}

	return courses, nil
}

func (cs *courseService) GetUserCourse(ctx context.Context, courseID, userID int) (*model.Course, error) {
	return cs.db.CourseStorage().GetUserCourse(ctx, courseID, userID)
}

func (cs *courseService) GetUserCourseLessons(ctx context.Context, userID int, courseID int) ([]model.Lesson, error) {
	course, err := cs.db.CourseStorage().GetUserCourse(ctx, userID, courseID)
	if err != nil {
		return nil, err
	}

	lessons, err := cs.db.LessonStorage().GetLessonsList(ctx, course.ID)
	if err != nil {
		return nil, err
	}

	lessonsIds := make([]int, 0, len(lessons))
	for _, lesson := range lessons {
		lessonsIds = append(lessonsIds, lesson.ID)
	}

	statuses, err := cs.db.LessonStorage().GetUserLessonsStatus(ctx, userID, course.ID, lessonsIds)
	if err != nil {
		return nil, err
	}

	for idx, lesson := range lessons {
		lessons[idx].Status = statuses[lesson.ID]
	}

	return lessons, nil
}

func (cs *courseService) GetCompanyCourses(ctx context.Context, companyID int) ([]model.Course, error) {
	return cs.db.CourseStorage().CompanyCourses(ctx, companyID)
}

func (cs *courseService) GetCompanyCourse(ctx context.Context, courseID, companyID int) (*model.Course, error) {
	return cs.db.CourseStorage().CompanyCourse(ctx, courseID, companyID)
}

func (cs *courseService) CreateCourse(ctx context.Context, c model.CourseSet) (*model.Course, error) {
	return cs.db.CourseStorage().CreateCourse(ctx, c)
}

func (cs *courseService) EditCourse(ctx context.Context, c model.CourseSet, companyID int) (*model.Course, error) {
	err := c.Validation()
	if !errors.Is(err, errs.ErrCourseNameIsEmpty) && err != nil {
		return nil, err
	}
	return cs.db.CourseStorage().EditCourse(ctx, c, companyID)
}
