package impl

import (
	"context"

	"go.uber.org/mock/gomock"

	"github.com/training-of-new-employees/qon/internal/errs"
	"github.com/training-of-new-employees/qon/internal/model"
	"github.com/training-of-new-employees/qon/internal/pkg/randomseq"
)

func (suite *serviceTestSuite) TestCreateLesson() {
	testCases := []struct {
		name    string
		err     error
		result  *model.Lesson
		prepare func() model.Lesson
	}{
		{
			name:   "success",
			err:    nil,
			result: &model.Lesson{},
			prepare: func() model.Lesson {
				courseID := 1

				lesson := model.NewTestLesson(courseID)

				suite.lessonStorage.EXPECT().CreateLesson(gomock.Any(), lesson, 1).Return(&lesson, nil)

				return lesson
			},
		},
		{
			name:   "incorrect course reference",
			err:    errs.ErrCourseReference,
			result: &model.Lesson{},
			prepare: func() model.Lesson {
				courseID := randomseq.RandomTestInt()

				lesson := model.NewTestLesson(courseID)

				suite.lessonStorage.EXPECT().CreateLesson(gomock.Any(), lesson, 1).Return(nil, errs.ErrCourseReference)

				return lesson
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			_, err := suite.lessonService.CreateLesson(context.TODO(), tc.prepare(), 1)
			suite.Equal(tc.err, err)
		})
	}
}

func (suite *serviceTestSuite) TestGetLesson() {
	testCases := []struct {
		name    string
		err     error
		result  *model.Lesson
		prepare func() (int, int)
	}{
		{
			name:   "success",
			err:    nil,
			result: &model.Lesson{},
			prepare: func() (int, int) {
				lessonID := 1
				companyID := 1

				suite.lessonStorage.EXPECT().GetLesson(gomock.Any(), lessonID).Return(&model.Lesson{CourseID: 1}, nil)
				suite.courseStorage.EXPECT().CompanyCourse(gomock.Any(), 1, companyID).Return(&model.Course{}, nil)
				return lessonID, companyID
			},
		},
		{
			name:   "not exist",
			err:    errs.ErrLessonNotFound,
			result: nil,
			prepare: func() (int, int) {
				lessonID := 1
				companyID := 1

				suite.lessonStorage.EXPECT().GetLesson(gomock.Any(), lessonID).Return(nil, errs.ErrLessonNotFound)
				return lessonID, companyID
			},
		},
		{
			name:   "company course not found",
			err:    errs.ErrLessonNotFound,
			result: nil,
			prepare: func() (int, int) {
				lessonID := 1
				companyID := 1

				suite.lessonStorage.EXPECT().GetLesson(gomock.Any(), lessonID).Return(&model.Lesson{CourseID: 1}, nil)
				suite.courseStorage.EXPECT().CompanyCourse(gomock.Any(), 1, companyID).Return(nil, errs.ErrCourseNotFound)
				return lessonID, companyID
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			lessonID, companyID := tc.prepare()
			lesson, err := suite.lessonService.GetLesson(context.Background(), lessonID, companyID)
			suite.Equal(tc.err, err)

			if tc.err != nil {
				suite.Equal(tc.result, lesson)
			}
		})
	}
}

func (suite *serviceTestSuite) TestUpdateLesson() {
	testCases := []struct {
		name    string
		err     error
		prepare func() model.LessonUpdate
	}{
		{
			name: "success",
			err:  nil,
			prepare: func() model.LessonUpdate {
				edit := model.NewTestEditLesson(1)

				suite.lessonStorage.EXPECT().UpdateLesson(gomock.Any(), edit).Return(&model.Lesson{}, nil)

				return edit
			},
		},
		{
			name: "fail",
			err:  errs.ErrCourseNotFound,
			prepare: func() model.LessonUpdate {
				edit := model.NewTestEditLesson(randomseq.RandomTestInt())
				suite.lessonStorage.EXPECT().UpdateLesson(gomock.Any(), edit).Return(nil, errs.ErrCourseNotFound)
				return edit
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			_, err := suite.lessonService.UpdateLesson(context.TODO(), tc.prepare())
			suite.Equal(tc.err, err)
		})
	}
}

func (suite *serviceTestSuite) TestGetUserLesson() {
	testCases := []struct {
		name    string
		err     error
		prepare func() (int, int)
	}{
		{
			name: "success",
			err:  nil,
			prepare: func() (int, int) {
				userID := 1
				lessonID := 1
				courseID := 1

				suite.lessonStorage.
					EXPECT().
					GetLesson(gomock.Any(), lessonID).
					Return(&model.Lesson{ID: lessonID, CourseID: courseID}, nil)

				suite.courseStorage.
					EXPECT().
					UserCourses(gomock.Any(), userID).
					Return([]model.Course{{ID: courseID}}, nil)

				suite.lessonStorage.
					EXPECT().
					GetUserLessonsStatus(gomock.Any(), userID, courseID, []int{lessonID}).
					Return(map[int]string{lessonID: "not-started"}, nil)

				return userID, lessonID
			},
		},
		{
			name: "lesson not found",
			err:  errs.ErrLessonNotFound,
			prepare: func() (int, int) {
				userID := 1
				lessonID := 1

				suite.lessonStorage.
					EXPECT().
					GetLesson(gomock.Any(), lessonID).
					Return(nil, errs.ErrLessonNotFound)

				return userID, lessonID
			},
		},
		{
			name: "course not found",
			err:  errs.ErrCourseNotFound,
			prepare: func() (int, int) {
				userID := 1
				lessonID := 1
				courseID := 1

				suite.lessonStorage.
					EXPECT().
					GetLesson(gomock.Any(), lessonID).
					Return(&model.Lesson{ID: lessonID, CourseID: courseID}, nil)

				suite.courseStorage.
					EXPECT().
					UserCourses(gomock.Any(), userID).
					Return([]model.Course{{ID: 2}}, nil)

				return userID, lessonID
			},
		},
		{
			name: "error receive status",
			err:  errs.ErrInternal,
			prepare: func() (int, int) {
				userID := 1
				lessonID := 1
				courseID := 1

				suite.lessonStorage.
					EXPECT().
					GetLesson(gomock.Any(), lessonID).
					Return(&model.Lesson{ID: lessonID, CourseID: courseID}, nil)

				suite.courseStorage.
					EXPECT().
					UserCourses(gomock.Any(), userID).
					Return([]model.Course{{ID: courseID}}, nil)

				suite.lessonStorage.
					EXPECT().
					GetUserLessonsStatus(gomock.Any(), userID, courseID, []int{lessonID}).
					Return(nil, errs.ErrInternal)

				return userID, lessonID
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			userID, lessonID := tc.prepare()
			_, err := suite.lessonService.GetUserLesson(context.TODO(), userID, lessonID)
			suite.Equal(tc.err, err)
		})
	}
}

func (suite *serviceTestSuite) TestGetLessonsList() {
	testCases := []struct {
		name    string
		err     error
		prepare func() int
	}{
		{
			name: "success",
			err:  nil,
			prepare: func() int {
				courseID := 1
				suite.lessonStorage.EXPECT().GetLessonsList(gomock.Any(), courseID).Return(model.NewTestListLessons(courseID), nil)

				return courseID
			},
		},
		{
			name: "success (empty)",
			err:  nil,
			prepare: func() int {
				courseID := randomseq.RandomTestInt()
				suite.lessonStorage.EXPECT().GetLessonsList(gomock.Any(), courseID).Return(nil, nil)

				return courseID
			},
		},
		{
			name: "internal error",
			err:  errs.ErrInternal,
			prepare: func() int {
				courseID := randomseq.RandomTestInt()
				suite.lessonStorage.EXPECT().GetLessonsList(gomock.Any(), courseID).Return(nil, errs.ErrInternal)

				return courseID
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			_, err := suite.lessonService.GetLessonsList(context.TODO(), tc.prepare())
			suite.Equal(tc.err, err)
		})
	}
}
