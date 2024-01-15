package pg

import (
	"context"

	"github.com/training-of-new-employees/qon/internal/errs"
	"github.com/training-of-new-employees/qon/internal/model"
)

func (suite *storeTestSuite) TestCreateCourse() {
	suite.NotNil(suite.store)

	ca := model.NewTestCreateAdmin()
	uc := model.NewTestUserCreate()
	admin, _ := suite.store.UserStorage().CreateAdmin(context.TODO(), uc, ca.Company)

	testCases := []struct {
		name    string
		prepare func() model.CourseSet
		err     error
	}{
		{
			name: "success",
			prepare: func() model.CourseSet {
				course := model.NewTestCourseSet()
				course.CreatedBy = admin.ID
				return course
			},
			err: nil,
		},
		{
			name: "empty name",
			prepare: func() model.CourseSet {
				course := model.NewTestCourseSet()
				course.Name = ""
				return course

			},
			err: errs.ErrCourseNameIsEmpty,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			_, err := suite.store.CourseStorage().CreateCourse(context.TODO(), tc.prepare())
			suite.Equal(tc.err, err)
		})
	}
}
