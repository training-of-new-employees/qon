package impl

import (
	"context"

	"go.uber.org/mock/gomock"

	"github.com/training-of-new-employees/qon/internal/errs"
	"github.com/training-of-new-employees/qon/internal/model"
)

func (suite *serviceTestSuite) TestGetPositionCourses() {
	testCases := []struct {
		name      string
		err       error
		resultLen int
		prepare   func() (int, int)
	}{
		{
			name:      "success",
			err:       nil,
			resultLen: 2,
			prepare: func() (int, int) {
				companyID := 2
				positionID := 1

				suite.positionStorage.
					EXPECT().
					GetPositionInCompany(gomock.Any(), companyID, positionID).
					Return(&model.Position{}, nil)

				suite.positionStorage.
					EXPECT().
					GetCourseForPosition(gomock.Any(), positionID).
					Return([]int{1, 2}, nil)

				return companyID, positionID
			},
		},
		{
			name:      "success (empty)",
			err:       nil,
			resultLen: 0,
			prepare: func() (int, int) {
				companyID := 2
				positionID := 1

				suite.positionStorage.
					EXPECT().
					GetPositionInCompany(gomock.Any(), companyID, positionID).
					Return(&model.Position{}, nil)

				suite.positionStorage.
					EXPECT().
					GetCourseForPosition(gomock.Any(), positionID).
					Return([]int{}, nil)

				return companyID, positionID
			},
		},
		{
			name:      "company not found",
			err:       errs.ErrCompanyNotFound,
			resultLen: 0,
			prepare: func() (int, int) {
				companyID := 2
				positionID := 1

				suite.positionStorage.
					EXPECT().
					GetPositionInCompany(gomock.Any(), companyID, positionID).
					Return(nil, errs.ErrCompanyNotFound)

				return companyID, positionID
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			ctx := context.Background()

			companyID, positionID := tc.prepare()
			courses, err := suite.positionService.GetPositionCourses(ctx, companyID, positionID)
			suite.Equal(tc.err, err)

			if tc.err == nil {
				suite.Len(courses, tc.resultLen)
			}
		})
	}
}

func (suite *serviceTestSuite) TestAssignCourses() {
	testCases := []struct {
		name    string
		err     error
		prepare func() (int, int, []int)
	}{
		{
			name: "success",
			err:  nil,
			prepare: func() (int, int, []int) {
				userID := 1
				courses := []int{1, 2}
				positionID := 1

				suite.userStorage.
					EXPECT().
					GetUserByID(gomock.Any(), userID).
					Return(&model.User{}, nil)

				suite.positionStorage.
					EXPECT().
					GetPositionByID(gomock.Any(), positionID).
					Return(&model.Position{}, nil)

				suite.positionStorage.
					EXPECT().
					AssignCourses(gomock.Any(), positionID, courses).
					Return(nil)

				return userID, positionID, courses
			},
		},
		{
			name: "user not found",
			err:  errs.ErrUserNotFound,
			prepare: func() (int, int, []int) {
				userID := 1
				courses := []int{1, 2}
				positionID := 1

				suite.userStorage.
					EXPECT().
					GetUserByID(gomock.Any(), userID).
					Return(nil, errs.ErrUserNotFound)

				return userID, positionID, courses
			},
		},
		{
			name: "position not found",
			err:  errs.ErrPositionNotFound,
			prepare: func() (int, int, []int) {
				userID := 1
				courses := []int{1, 2}
				positionID := 1

				suite.userStorage.
					EXPECT().
					GetUserByID(gomock.Any(), userID).
					Return(&model.User{}, nil)

				suite.positionStorage.
					EXPECT().
					GetPositionByID(gomock.Any(), positionID).
					Return(nil, errs.ErrPositionNotFound)

				return userID, positionID, courses
			},
		},
		{
			name: "course already assigned",
			err:  errs.ErrPositionCourseUsed,
			prepare: func() (int, int, []int) {
				userID := 1
				courses := []int{1, 2}
				positionID := 1

				suite.userStorage.
					EXPECT().
					GetUserByID(gomock.Any(), userID).
					Return(&model.User{}, nil)

				suite.positionStorage.
					EXPECT().
					GetPositionByID(gomock.Any(), positionID).
					Return(&model.Position{}, nil)

				suite.positionStorage.
					EXPECT().
					AssignCourses(gomock.Any(), positionID, courses).
					Return(errs.ErrPositionCourseUsed)

				return userID, positionID, courses
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			ctx := context.Background()

			userID, positionID, courses := tc.prepare()
			err := suite.positionService.AssignCourses(ctx, positionID, courses, userID)
			suite.Equal(tc.err, err)
		})
	}
}
