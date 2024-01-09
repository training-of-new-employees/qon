package pg

import (
	"context"

	"github.com/training-of-new-employees/qon/internal/model"
)

func (suite *storeTestSuite) TestCreateLessonDB() {
	suite.NotNil(suite.store)

	company, err := suite.store.CompanyStorage().CreateCompany(context.TODO(), "test&Co")
	suite.NoError(err)
	suite.NotEmpty(company)

	position, err := suite.store.PositionStorage().CreatePositionDB(context.TODO(),
		model.PositionSet{CompanyID: company.ID, Name: "test-position"})
	suite.NoError(err)
	suite.NotEmpty(position)

	NewUser := model.UserCreate{
		CompanyID:  company.ID,
		PositionID: position.ID,
		Name:       "test",
		Email:      "test@yandex.com",
		Password:   "123456",
	}
	user, err := suite.store.UserStorage().CreateUser(context.TODO(), NewUser)
	suite.NoError(err)
	suite.NotEmpty(user)

	NewCourse := model.CourseCreate{
		Name:        "test",
		Description: "test",
	}
	course, err := suite.store.CourseStorage().CreateCourse(context.TODO(),
		NewCourse, user.ID)
	suite.NoError(err)
	suite.NotEmpty(course)

	testCases := []struct {
		name    string
		lesson  func() model.LessonCreate
		user_id int
		err     error
	}{
		{
			name: "success",
			lesson: func() model.LessonCreate {
				l := model.LessonCreate{
					CourseID:    course.ID,
					Name:        "Lesson1",
					Description: "Description1",
					Path:        "http://test",
				}
				return l
			},
			user_id: user.ID,
			err:     nil,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			_, err := suite.store.LessonStorage().CreateLessonDB(context.TODO(),
				tc.lesson(), tc.user_id)
			suite.Equal(tc.err, err)
		})
	}
}
