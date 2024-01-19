package impl

import (
	"context"

	"github.com/training-of-new-employees/qon/internal/errs"
	"github.com/training-of-new-employees/qon/internal/model"
	"github.com/training-of-new-employees/qon/internal/pkg/randomseq"
	"go.uber.org/mock/gomock"
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
		prepare func() int
	}{
		{
			name:   "success",
			err:    nil,
			result: &model.Lesson{},
			prepare: func() int {
				lessonID := 1

				suite.lessonStorage.EXPECT().GetLesson(gomock.Any(), lessonID).Return(&model.Lesson{}, nil)
				return lessonID
			},
		},
		{
			name:   "not exist",
			err:    errs.ErrLessonNotFound,
			result: nil,
			prepare: func() int {
				lessonID := 1

				suite.lessonStorage.EXPECT().GetLesson(gomock.Any(), lessonID).Return(nil, errs.ErrLessonNotFound)
				return lessonID
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			lesson, err := suite.lessonService.GetLesson(context.Background(), tc.prepare())
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
