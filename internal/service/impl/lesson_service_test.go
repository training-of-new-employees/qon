package impl

import (
	"context"

	"github.com/training-of-new-employees/qon/internal/model"
	"go.uber.org/mock/gomock"
)

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

				suite.lessonStorage.EXPECT().GetLesson(gomock.Any(),
					lessonID).Return(&model.Lesson{}, nil)
				return lessonID
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {

			ctx := context.Background()
			lessonID := tc.prepare()
			lesson, err := suite.lessonService.GetLesson(ctx, lessonID)
			suite.Equal(tc.err, err)

			if tc.err != nil {
				suite.Equal(tc.result, lesson)
			}
		})
	}

}
