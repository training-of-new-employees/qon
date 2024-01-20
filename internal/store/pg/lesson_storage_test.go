package pg

import (
	"context"
	"fmt"

	"github.com/training-of-new-employees/qon/internal/errs"
	"github.com/training-of-new-employees/qon/internal/model"
	"github.com/training-of-new-employees/qon/internal/pkg/randomseq"
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
		Name:       randomseq.RandomName(10),
		Email:      "test@yandex.com",
		Password:   randomseq.RandomPassword(),
	}
	user, err := suite.store.UserStorage().CreateUser(context.TODO(), NewUser)
	if err != nil {
		return nil, nil, err
	}
	suite.NotEmpty(user)

	NewCourse := model.CourseSet{
		CreatedBy:   user.ID,
		Name:        randomseq.RandomName(10),
		Description: randomseq.RandomString(20),
	}
	course, err := suite.store.CourseStorage().CreateCourse(context.TODO(),
		NewCourse)
	if err != nil {
		return nil, nil, err
	}
	return course, user, nil
}

func (suite *storeTestSuite) TestGetLesson() {
	course, user, err := suite.prepareLessonCreation()

	suite.NoError(err)
	suite.NotEmpty(course)

	lesson, err := suite.store.LessonStorage().CreateLesson(
		context.TODO(),
		model.Lesson{
			CourseID:   course.ID,
			Name:       randomseq.RandomName(10),
			Content:    randomseq.RandomString(20),
			URLPicture: randomseq.RandomString(30),
		},
		user.ID,
	)
	suite.NoError(err)
	suite.NotEmpty(lesson)

	testCases := []struct {
		name           string
		lessonID       int
		expectedLesson *model.Lesson
		err            error
	}{
		{
			name:           "success",
			lessonID:       lesson.ID,
			expectedLesson: lesson,
			err:            nil,
		},
		{
			name:           "not existing lesson",
			lessonID:       randomseq.RandomTestInt(),
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

func (suite *storeTestSuite) TestCreateLesson() {

	course, user, err := suite.prepareLessonCreation()

	suite.NoError(err)
	suite.NotEmpty(course)

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
					Name:       randomseq.RandomName(10),
					Content:    randomseq.RandomString(20),
					URLPicture: randomseq.RandomString(30),
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
					Name:       randomseq.RandomName(10),
					Content:    randomseq.RandomString(20),
					URLPicture: randomseq.RandomString(30),
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
					Content:    randomseq.RandomString(20),
					URLPicture: randomseq.RandomString(30),
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
					Name:       randomseq.RandomName(10),
					Content:    randomseq.RandomString(20),
					URLPicture: randomseq.RandomString(30),
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
					Name:       randomseq.RandomName(10),
					Content:    randomseq.RandomString(20),
					URLPicture: randomseq.RandomString(30),
				}
				return l
			},
			user_id: randomseq.RandomTestInt(),
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

func (suite *storeTestSuite) TestUpdateLesson() {
	course, user, err := suite.prepareLessonCreation()
	suite.NoError(err)
	suite.NotEmpty(course)

	lesson, err := suite.store.LessonStorage().CreateLesson(
		context.TODO(),
		model.Lesson{
			CourseID:   course.ID,
			Name:       randomseq.RandomName(10),
			Content:    randomseq.RandomString(20),
			URLPicture: randomseq.RandomString(30),
		},
		user.ID,
	)
	suite.NoError(err)
	suite.NotEmpty(lesson)

	testCases := []struct {
		name    string
		prepare func() (model.LessonUpdate, model.Lesson) // возвращает данные для редактирования и ожидаемый результат для проверки тест-кейса
		err     error
	}{
		{
			name: "success",
			prepare: func() (model.LessonUpdate, model.Lesson) {
				editField := model.LessonUpdate{ID: lesson.ID}

				// ожидаемые данные урока
				expected := *lesson
				// определение случайным образом полей для редактирования:
				//
				// изменение имени урока
				if randomseq.RandomBool() {
					name := randomseq.RandomName(10)
					editField.Name = name
					expected.Name = name
				}
				// изменение содержания урока
				if randomseq.RandomBool() {
					content := randomseq.RandomString(20)
					editField.Content = content
					expected.Content = content
				}
				// изменение ссылки картинки
				if randomseq.RandomBool() {
					url := fmt.Sprintf("https://%sexample.com/%s.png", randomseq.RandomString(10), randomseq.RandomString(5))
					editField.URLPicture = url
					expected.URLPicture = url
				}

				return editField, expected
			},
			err: nil,
		},
	}
	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			edit, expected := tc.prepare()
			updLesson, err := suite.store.LessonStorage().UpdateLesson(context.TODO(), edit)

			suite.Equal(tc.err, err)

			if err != nil {
				suite.Equal(expected, *updLesson)
			}
		})
	}
}

func (suite *storeTestSuite) TestGetLessonList() {
	course, user, err := suite.prepareLessonCreation()

	suite.NoError(err)
	suite.NotEmpty(course)

	lesson1, err := suite.store.LessonStorage().CreateLesson(
		context.TODO(),
		model.Lesson{
			CourseID:   course.ID,
			Name:       randomseq.RandomName(10),
			Content:    randomseq.RandomString(20),
			URLPicture: randomseq.RandomString(30),
		},
		user.ID,
	)
	suite.NoError(err)
	suite.NotEmpty(lesson1)

	lesson2, err := suite.store.LessonStorage().CreateLesson(
		context.TODO(),
		model.Lesson{
			CourseID:   course.ID,
			Name:       randomseq.RandomName(10),
			Content:    randomseq.RandomString(20),
			URLPicture: randomseq.RandomString(30),
		},
		user.ID,
	)
	suite.NoError(err)
	suite.NotEmpty(lesson2)

	testCases := []struct {
		name     string
		courseID int
		expected []model.Lesson
		err      error
	}{
		{
			name:     "success",
			courseID: course.ID,
			expected: []model.Lesson{
				*lesson1,
				*lesson2,
			},
			err: nil,
		},
		{
			name:     "empty course",
			expected: nil,
			err:      errs.ErrCourseIDNotEmpty,
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

func (suite *storeTestSuite) TestGetUserLessonsStatus() {
	suite.NotNil(suite.store)

	ca1 := model.NewTestCreateAdmin()
	uc1 := model.NewTestUserCreate()
	admin, err := suite.store.UserStorage().CreateAdmin(context.TODO(), uc1, ca1.Company)
	suite.NoError(err)

	p := model.NewTestPositionSet()
	p.CompanyID = admin.CompanyID
	pos, err := suite.store.PositionStorage().CreatePosition(context.TODO(), p)
	suite.NoError(err)

	uc2 := model.NewTestUserCreate()
	uc2.CompanyID = admin.CompanyID
	uc2.PositionID = pos.ID
	user, err := suite.store.UserStorage().CreateUser(context.TODO(), uc2)
	suite.NoError(err)

	course, err := suite.store.CourseStorage().CreateCourse(
		context.TODO(),
		model.CourseSet{
			Name:        "Test Name",
			Description: "Test DESCRIPTION",
			CreatedBy:   user.ID,
		},
	)
	suite.NoError(err)

	lesson, err := suite.store.LessonStorage().CreateLesson(
		context.TODO(),
		model.Lesson{CourseID: course.ID, Name: "Test Lesson", Content: "Test content", URLPicture: "Test picture"},
		user.ID,
	)
	suite.NoError(err)
	err = suite.store.LessonStorage().UpdateUserLessonStatus(context.TODO(), user.ID, course.ID, lesson.ID, "done")
	suite.NoError(err)

	lesson2, err := suite.store.LessonStorage().CreateLesson(
		context.TODO(),
		model.Lesson{CourseID: course.ID, Name: "Test Lesson", Content: "Test content", URLPicture: "Test picture"},
		user.ID,
	)
	suite.NoError(err)

	testCases := []struct {
		name     string
		err      error
		ids      []int
		statuses map[int]string
	}{
		{
			name: "success",
			err:  nil,
			ids:  []int{lesson.ID},
			statuses: map[int]string{
				lesson.ID: "done",
			},
		},
		{
			name: "generated not-started status",
			err:  nil,
			ids:  []int{lesson2.ID},
			statuses: map[int]string{
				lesson2.ID: "not-started",
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			statuses, err := suite.store.LessonStorage().GetUserLessonsStatus(context.TODO(), user.ID, course.ID, tc.ids)
			suite.Equal(tc.err, err)

			if err == nil {
				for id := range statuses {
					suite.Equal(statuses[id], tc.statuses[id])
				}
			}
		})
	}
}
