package pg

import (
	"context"

	"github.com/training-of-new-employees/qon/internal/errs"
	"github.com/training-of-new-employees/qon/internal/model"
)

func (suite *storeTestSuite) prepareLessonCreation() (*model.Course,
	*model.User, error) {
	suite.NotNil(suite.store)

	company, err := suite.store.CompanyStorage().CreateCompany(context.TODO(), "test&Co")
	if err != nil {
		return nil, nil, err
	}
	suite.NotEmpty(company)

	position, err := suite.store.PositionStorage().CreatePositionDB(context.TODO(),
		model.PositionSet{CompanyID: company.ID, Name: "test-position"})
	if err != nil {
		return nil, nil, err
	}
	suite.NotEmpty(position)

	NewUser := model.UserCreate{
		CompanyID:  company.ID,
		PositionID: position.ID,
		Name:       "test",
		Email:      "test@yandex.com",
		Password:   "123456",
	}
	user, err := suite.store.UserStorage().CreateUser(context.TODO(), NewUser)
	if err != nil {
		return nil, nil, err
	}
	suite.NotEmpty(user)

	NewCourse := model.CourseCreate{
		Name:        "test",
		Description: "test",
	}
	course, err := suite.store.CourseStorage().CreateCourse(context.TODO(),
		NewCourse, user.ID)
	if err != nil {
		return nil, nil, err
	}
	return course, user, nil
}

func (suite *storeTestSuite) TestCreateLessonDB() {

	course, user, err := suite.prepareLessonCreation()

	suite.NoError(err)
	suite.NotEmpty(course)

	if err != nil || course == nil {
		return
	}

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
					CourseID:   course.ID,
					Name:       "Lesson1",
					Content:    "Content1",
					URLPicture: "http://test",
				}
				return l
			},
			user_id: user.ID,
			err:     nil,
		},
		{
			name: "empty course id",
			lesson: func() model.LessonCreate {
				l := model.LessonCreate{
					Name:       "Lesson1",
					Content:    "Content1",
					URLPicture: "http://test",
				}
				return l
			},
			user_id: user.ID,
			err:     errs.ErrCourseReference,
		},
		{
			name: "empty lesson name",
			lesson: func() model.LessonCreate {
				l := model.LessonCreate{
					CourseID:   course.ID,
					Content:    "Content1",
					URLPicture: "http://test",
				}
				return l
			},
			user_id: user.ID,
			err:     errs.ErrLessonNameNotEmpty,
		},
		{
			name: "empty user id",
			lesson: func() model.LessonCreate {
				l := model.LessonCreate{
					CourseID:   course.ID,
					Name:       "Lesson1",
					Content:    "Content",
					URLPicture: "http://test",
				}
				return l
			},
			err: errs.ErrCreaterNotFound,
		},
		{
			name: "non existing user",
			lesson: func() model.LessonCreate {
				l := model.LessonCreate{
					CourseID:   course.ID,
					Name:       "Lesson1",
					Content:    "Content1",
					URLPicture: "http://test",
				}
				return l
			},
			user_id: 34,
			err:     errs.ErrCreaterNotFound,
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

func (suite *storeTestSuite) TestDeleteLessonDB() {
	course, user, err := suite.prepareLessonCreation()
	suite.NoError(err)
	suite.NotEmpty(course)

	lesson := func() model.LessonCreate {
		l := model.LessonCreate{
			CourseID:   course.ID,
			Name:       "Lesson2",
			Content:    "Content2",
			URLPicture: "http://test",
		}
		return l
	}
	newLesson, err := suite.store.LessonStorage().CreateLessonDB(context.TODO(),
		lesson(), user.ID)
	suite.NoError(err)
	suite.NotEmpty(newLesson)

	testCases := []struct {
		name     string
		lessonId int
		err      error
	}{
		{
			name:     "success",
			lessonId: newLesson.ID,
			err:      nil,
		},
		{
			name:     "lesson id don't exist",
			lessonId: 977272727,
			err:      errs.ErrLessonNotFound,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			err := suite.store.LessonStorage().DeleteLessonDB(context.TODO(),
				tc.lessonId)
			suite.Equal(tc.err, err)
		})
	}
}

func (suite *storeTestSuite) TestUpdateLessonDB() {
	course, user, err := suite.prepareLessonCreation()
	suite.NoError(err)
	suite.NotEmpty(course)

	lesson := func() model.LessonCreate {
		l := model.LessonCreate{
			CourseID:   course.ID,
			Name:       "Lesson4",
			Content:    "Content3",
			URLPicture: "http://test",
		}
		return l
	}
	newLesson, err := suite.store.LessonStorage().CreateLessonDB(context.TODO(),
		lesson(), user.ID)
	suite.NoError(err)
	suite.NotEmpty(newLesson)

	testCases := []struct {
		name     string
		lesson   func() model.LessonUpdate
		expected model.LessonUpdate
		err      error
	}{
		{
			name: "change name",
			lesson: func() model.LessonUpdate {
				l := model.LessonUpdate{
					ID:       newLesson.ID,
					CourseID: newLesson.CourseID,
					Name:     "Lesson4",
				}
				return l
			},
			expected: model.LessonUpdate{
				Name:       lesson().Name,
				Content:    newLesson.Content,
				URLPicture: newLesson.URLPicture,
			},
			err: nil,
		},
		{
			name: "change content",
			lesson: func() model.LessonUpdate {
				l := model.LessonUpdate{
					ID:       newLesson.ID,
					CourseID: newLesson.CourseID,
					Name:     newLesson.Name,
					Content:  "NewContent",
				}
				return l
			},
			expected: model.LessonUpdate{
				Name:       newLesson.Name,
				Content:    "NewContent",
				URLPicture: newLesson.URLPicture,
			},
			err: nil,
		},
		{
			name: "change picture",
			lesson: func() model.LessonUpdate {
				l := model.LessonUpdate{
					ID:         newLesson.ID,
					CourseID:   newLesson.CourseID,
					Name:       newLesson.Name,
					URLPicture: "http://newPicture",
				}
				return l
			},
			expected: model.LessonUpdate{
				Name:       newLesson.Name,
				Content:    "NewContent",
				URLPicture: "http://newPicture",
			},
			err: nil,
		},
	}
	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			updLesson, err := suite.store.LessonStorage().UpdateLessonDB(context.TODO(),
				tc.lesson())
			suite.Equal(tc.err, err)
			suite.Equal(tc.expected.Name, updLesson.Name)
			suite.Equal(tc.expected.Content, updLesson.Content)
			suite.Equal(tc.expected.URLPicture, updLesson.URLPicture)
		})
	}
}
