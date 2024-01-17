package pg

import (
	"context"
	"math/rand"
	"time"

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

	position, err := suite.store.PositionStorage().CreatePosition(context.TODO(),
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

	NewCourse := model.CourseSet{
		CreatedBy:   user.ID,
		Name:        "test",
		Description: "test",
	}
	course, err := suite.store.CourseStorage().CreateCourse(context.TODO(),
		NewCourse)
	if err != nil {
		return nil, nil, err
	}
	return course, user, nil
}

func (suite *storeTestSuite) TestGetLessonDB() {
	course, user, err := suite.prepareLessonCreation()

	suite.NoError(err)
	suite.NotEmpty(course)

	lesson := func() model.Lesson {
		l := model.Lesson{
			CourseID:   course.ID,
			Name:       "Lesson2",
			Content:    "Content2",
			URLPicture: "http://test",
		}
		return l
	}
	newLesson, err := suite.store.LessonStorage().CreateLesson(context.TODO(),
		lesson(), user.ID)
	suite.NoError(err)
	suite.NotEmpty(newLesson)

	rnd := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))

	testCases := []struct {
		name           string
		lessonID       int
		expectedLesson *model.Lesson
		err            error
	}{
		{
			name:           "success",
			lessonID:       newLesson.ID,
			expectedLesson: newLesson,
			err:            nil,
		},
		{
			name:           "not existing lesson",
			lessonID:       rnd.Intn(32) + 1,
			expectedLesson: nil,
			err:            errs.ErrLessonNotFound,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			l, err := suite.store.LessonStorage().GetLesson(context.TODO(), tc.lessonID)
			suite.Equal(tc.err, err)
			suite.Equal(tc.expectedLesson, l)
		})
	}

}

func (suite *storeTestSuite) TestCreateLessonDB() {

	course, user, err := suite.prepareLessonCreation()

	suite.NoError(err)
	suite.NotEmpty(course)

	rnd := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))

	testCases := []struct {
		name    string
		lesson  func() model.Lesson
		user_id int
		err     error
	}{
		{
			name: "success",
			lesson: func() model.Lesson {
				l := model.Lesson{
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
			lesson: func() model.Lesson {
				l := model.Lesson{
					Name:       "Lesson1",
					Content:    "Content1",
					URLPicture: "http://test",
				}
				return l
			},
			user_id: user.ID,
			err:     errs.ErrCourseIDNotEmpty,
		},
		{
			name: "empty lesson name",
			lesson: func() model.Lesson {
				l := model.Lesson{
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
			lesson: func() model.Lesson {
				l := model.Lesson{
					CourseID:   course.ID,
					Name:       "Lesson1",
					Content:    "Content",
					URLPicture: "http://test",
				}
				return l
			},
			err: errs.ErrCreaterNotEmpty,
		},
		{
			name: "non existing user",
			lesson: func() model.Lesson {
				l := model.Lesson{
					CourseID:   course.ID,
					Name:       "Lesson1",
					Content:    "Content1",
					URLPicture: "http://test",
				}
				return l
			},
			user_id: rnd.Intn(32) + 1,
			err:     errs.ErrCreaterNotFound,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			_, err := suite.store.LessonStorage().CreateLesson(context.TODO(),
				tc.lesson(), tc.user_id)
			suite.Equal(tc.err, err)
		})
	}
}

func (suite *storeTestSuite) TestUpdateLessonDB() {
	course, user, err := suite.prepareLessonCreation()
	suite.NoError(err)
	suite.NotEmpty(course)

	lesson := func() model.Lesson {
		l := model.Lesson{
			CourseID:   course.ID,
			Name:       "Lesson4",
			Content:    "Content3",
			URLPicture: "http://test",
		}
		return l
	}
	newLesson, err := suite.store.LessonStorage().CreateLesson(context.TODO(),
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
					ID:   newLesson.ID,
					Name: "Lesson4",
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
					ID:      newLesson.ID,
					Name:    newLesson.Name,
					Content: "NewContent",
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
			updLesson, err := suite.store.LessonStorage().UpdateLesson(context.TODO(),
				tc.lesson())
			suite.Equal(tc.err, err)
			suite.Equal(tc.expected.Name, updLesson.Name)
			suite.Equal(tc.expected.Content, updLesson.Content)
			suite.Equal(tc.expected.URLPicture, updLesson.URLPicture)
		})
	}
}

func (suite *storeTestSuite) TestGetLessonListDB() {
	course, user, err := suite.prepareLessonCreation()

	suite.NoError(err)
	suite.NotEmpty(course)

	lesson := func() model.Lesson {
		l := model.Lesson{
			CourseID:   course.ID,
			Name:       "Lesson2",
			Content:    "Content2",
			URLPicture: "http://test",
		}
		return l
	}
	lesson1, err := suite.store.LessonStorage().CreateLesson(context.TODO(),
		lesson(), user.ID)
	suite.NoError(err)
	suite.NotEmpty(lesson1)

	lesson = func() model.Lesson {
		l := model.Lesson{
			CourseID:   course.ID,
			Name:       "Lesson3",
			Content:    "Content3",
			URLPicture: "http://test3",
		}
		return l
	}
	lesson2, err := suite.store.LessonStorage().CreateLesson(context.TODO(),
		lesson(), user.ID)
	suite.NoError(err)
	suite.NotEmpty(lesson2)

	rnd := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))

	testCases := []struct {
		name     string
		courseID int
		expected []*model.Lesson
		err      error
	}{
		{
			name:     "success",
			courseID: course.ID,
			expected: []*model.Lesson{
				lesson1,
				lesson2,
			},
			err: nil,
		},
		{
			name:     "not existing course",
			courseID: rnd.Intn(32) + 17,
			expected: nil,
			err:      errs.ErrLessonNotFound,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			lessonsList, err := suite.store.LessonStorage().GetLessonsList(context.TODO(),
				tc.courseID)
			suite.Equal(tc.err, err)
			suite.Equal(tc.expected, lessonsList)
		})
	}
}
