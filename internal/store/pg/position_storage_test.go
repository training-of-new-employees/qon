package pg

import (
	"context"
	"fmt"

	"github.com/training-of-new-employees/qon/internal/errs"
	"github.com/training-of-new-employees/qon/internal/model"
	"github.com/training-of-new-employees/qon/internal/pkg/randomseq"
)

func (suite *storeTestSuite) TestCreatePositionDB() {
	suite.NotNil(suite.store)

	// добавление компании
	company, err := suite.store.CompanyStorage().CreateCompany(context.TODO(), "test&Co")
	suite.NoError(err)
	suite.NotEmpty(company)

	testCases := []struct {
		name     string
		position func() model.PositionSet
		err      error
	}{
		{
			name: "success",
			position: func() model.PositionSet {
				position := model.PositionSet{
					CompanyID: company.ID,
					Name:      "test-position",
				}

				return position
			},
			err: nil,
		},
		{
			name: "empty company id",
			position: func() model.PositionSet {
				position := model.PositionSet{
					Name: "test-position",
				}

				return position
			},
			err: errs.ErrCompanyReference,
		},
		{
			name: "empty position name",
			position: func() model.PositionSet {
				position := model.PositionSet{
					CompanyID: company.ID,
					Name:      "",
				}

				return position
			},
			err: errs.ErrPositionNameNotEmpty,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			_, err := suite.store.PositionStorage().CreatePosition(context.TODO(), tc.position())
			suite.Equal(tc.err, err)
		})
	}

	suite.NoError(err)
}

func (suite *storeTestSuite) TestGetPositionByID() {
	suite.NotNil(suite.store)

	// добавление компании
	company, err := suite.store.CompanyStorage().CreateCompany(context.TODO(), "test&Co")
	suite.NoError(err)
	suite.NotEmpty(company)

	// добавление должности
	position, err := suite.store.PositionStorage().CreatePosition(
		context.TODO(),
		model.PositionSet{CompanyID: company.ID, Name: "test-position"},
	)
	suite.NoError(err)
	suite.NotEmpty(position)

	testCases := []struct {
		name       string
		positionID int
		err        error
	}{
		{
			name:       "success",
			positionID: position.ID,
			err:        nil,
		},
		{
			name:       "random position",
			positionID: randomseq.RandomTestInt(),
			err:        errs.ErrPositionNotFound,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			_, err := suite.store.PositionStorage().GetPositionByID(context.TODO(), tc.positionID)
			suite.Equal(tc.err, err)
		})
	}
}

func (suite *storeTestSuite) TestGetPositionInCompany() {
	suite.NotNil(suite.store)

	// добавление компании
	company, err := suite.store.CompanyStorage().CreateCompany(context.TODO(), "test&Co")
	suite.NoError(err)
	suite.NotEmpty(company)

	// добавление должности
	position, err := suite.store.PositionStorage().CreatePosition(
		context.TODO(),
		model.PositionSet{CompanyID: company.ID, Name: "test-position"},
	)
	suite.NoError(err)
	suite.NotEmpty(position)

	testCases := []struct {
		name    string
		payload func() (int, int) // получение ИД компании и ИД должности
		err     error
	}{
		{
			name: "success",
			payload: func() (int, int) {
				return company.ID, position.ID
			},
			err: nil,
		},
		{
			name: "random company",
			payload: func() (int, int) {
				companyID := randomseq.RandomTestInt()
				return companyID, position.ID
			},
			err: errs.ErrPositionNotFound,
		},
		{
			name: "random position",
			payload: func() (int, int) {
				positionID := randomseq.RandomTestInt()
				return company.ID, positionID
			},
			err: errs.ErrPositionNotFound,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			companyID, positionID := tc.payload()
			_, err := suite.store.PositionStorage().GetPositionInCompany(context.TODO(), companyID, positionID)
			suite.Equal(tc.err, err)
		})
	}
}

func (suite *storeTestSuite) TestListPositions() {
	suite.NotNil(suite.store)

	// поиск в пустой базе должностей по id несуществующей компании
	positions, err := suite.store.PositionStorage().ListPositions(context.TODO(), randomseq.RandomTestInt())
	suite.NoError(err)
	suite.Equal([]*model.Position{}, positions)

	// создание компании
	company, err := suite.store.CompanyStorage().CreateCompany(context.TODO(), "test&Co")
	suite.NoError(err)
	suite.NotEmpty(company)

	// генерация случайного кол-ва должностей (от 100 до 356)
	countPositions := randomseq.RandomTestInt()

	// массив с ожидаемыми идентификаторами должностей
	expectedIDs := []int{}

	// добавление случайного кол-ва должностей
	for i := 0; i < countPositions; i++ {
		position, err := suite.store.PositionStorage().CreatePosition(
			context.TODO(),
			model.PositionSet{CompanyID: company.ID, Name: fmt.Sprintf("test-position-%d", i)},
		)
		suite.NoError(err)

		// добавление в массив идентификаторов добавленных должностей
		expectedIDs = append(expectedIDs, position.ID)
	}

	// получение добавленных должностей
	positions, err = suite.store.PositionStorage().ListPositions(context.TODO(), company.ID)
	suite.NotEmpty(positions)
	suite.NoError(err)

	// добавление идентификаторов должностей в массив
	actualIDs := []int{}
	for _, p := range positions {
		actualIDs = append(actualIDs, p.ID)
	}

	suite.EqualValues(expectedIDs, actualIDs)
}

func (suite *storeTestSuite) TestUpdatePosition() {
	suite.NotNil(suite.store)

	// добавление компании
	company, err := suite.store.CompanyStorage().CreateCompany(context.TODO(), "test&Co")
	suite.NoError(err)
	suite.NotEmpty(company)

	// добавление должности
	position, err := suite.store.PositionStorage().CreatePosition(
		context.TODO(),
		model.PositionSet{CompanyID: company.ID, Name: "test-position"},
	)
	suite.NoError(err)
	suite.NotEmpty(position)

	testCases := []struct {
		name    string
		payload func() (int, model.PositionSet) // возвращает ИД должности и данные позиции для обновления
		err     error
	}{
		{
			name: "success",
			payload: func() (int, model.PositionSet) {
				return position.ID, model.PositionSet{CompanyID: company.ID, Name: "new-position-name"}
			},
			err: nil,
		},
		{
			name: "empty company name",
			payload: func() (int, model.PositionSet) {
				return position.ID, model.PositionSet{CompanyID: company.ID, Name: ""}
			},
			err: errs.ErrPositionNameNotEmpty,
		},
		{
			name: "company not found",
			payload: func() (int, model.PositionSet) {
				return position.ID, model.PositionSet{CompanyID: randomseq.RandomTestInt(), Name: "new-position-name"}
			},
			err: errs.ErrPositionNotFound,
		},
		{
			name: "position not found",
			payload: func() (int, model.PositionSet) {
				return randomseq.RandomTestInt(), model.PositionSet{CompanyID: company.ID, Name: "new-position-name"}
			},
			err: errs.ErrPositionNotFound,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			positionID, positionUpdate := tc.payload()

			_, err := suite.store.PositionStorage().UpdatePosition(
				context.TODO(),
				positionID,
				positionUpdate,
			)

			suite.Equal(tc.err, err)
		})
	}
}

func (suite *storeTestSuite) TestGetCourseForPosition() {
	suite.NotNil(suite.store)

	company, err := suite.store.CompanyStorage().CreateCompany(context.TODO(), "test&Co")
	suite.NoError(err)
	suite.NotEmpty(company)

	position, err := suite.store.PositionStorage().CreatePosition(
		context.TODO(),
		model.PositionSet{CompanyID: company.ID, Name: "test-position"},
	)
	suite.NoError(err)
	suite.NotEmpty(position)

	position2, err := suite.store.PositionStorage().CreatePosition(
		context.TODO(),
		model.PositionSet{CompanyID: company.ID, Name: "test-position2"},
	)
	suite.NoError(err)
	suite.NotEmpty(position2)

	u := model.NewTestUserCreate()
	u.CompanyID = company.ID
	u.PositionID = position.ID

	// добавление пользователя
	user, err := suite.store.UserStorage().CreateUser(context.TODO(), u)
	suite.NoError(err)
	suite.NotEmpty(user)

	course, err := suite.store.CourseStorage().CreateCourse(context.TODO(), model.CourseSet{
		Name:        "Test Course",
		Description: "Test Description Test Description",
		IsArchived:  false,
		CreatedBy:   user.ID,
	})
	suite.NoError(err)
	suite.NotEmpty(course)

	err = suite.store.PositionStorage().AssignCourses(context.TODO(), position.ID, []int{course.ID})
	suite.NoError(err)

	testCases := []struct {
		name      string
		payload   func() int // возвращает ИД должности
		err       error
		resultLen int
	}{
		{
			name: "success",
			payload: func() int {
				return position.ID
			},
			resultLen: 1,
			err:       nil,
		},
		{
			name: "success (empty)",
			payload: func() int {
				return position2.ID
			},
			resultLen: 0,
			err:       nil,
		},
		{
			name: "position not found",
			payload: func() int {
				return randomseq.RandomTestInt()
			},
			resultLen: 0,
			err:       nil,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			positionID := tc.payload()

			courses, err := suite.store.PositionStorage().GetCourseForPosition(
				context.TODO(),
				positionID,
			)

			suite.Equal(tc.err, err)
			if tc.err == nil {
				suite.Len(courses, tc.resultLen)
			}
		})
	}
}

func (suite *storeTestSuite) TestAssignCourses() {
	suite.NotNil(suite.store)

	company, err := suite.store.CompanyStorage().CreateCompany(context.TODO(), "test&Co")
	suite.NoError(err)
	suite.NotEmpty(company)

	position, err := suite.store.PositionStorage().CreatePosition(
		context.TODO(),
		model.PositionSet{CompanyID: company.ID, Name: "test-position"},
	)
	suite.NoError(err)
	suite.NotEmpty(position)

	u := model.NewTestUserCreate()
	u.CompanyID = company.ID
	u.PositionID = position.ID

	// добавление пользователя
	user, err := suite.store.UserStorage().CreateUser(context.TODO(), u)
	suite.NoError(err)
	suite.NotEmpty(user)

	course, err := suite.store.CourseStorage().CreateCourse(context.TODO(), model.CourseSet{
		Name:        "Test Course",
		Description: "Test Description Test Description",
		IsArchived:  false,
		CreatedBy:   user.ID,
	})
	suite.NoError(err)
	suite.NotEmpty(course)

	testCases := []struct {
		name    string
		payload func() (int, []int)
		err     error
	}{
		{
			name: "success",
			payload: func() (int, []int) {
				return position.ID, []int{course.ID}
			},
			err: nil,
		},
		{
			name: "position not found",
			payload: func() (int, []int) {
				return randomseq.RandomTestInt(), []int{course.ID}
			},
			err: errs.ErrPositionReference,
		},
		{
			name: "course not found",
			payload: func() (int, []int) {
				return position.ID, []int{randomseq.RandomTestInt()}
			},
			err: errs.ErrCourseReference,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			positionID, courses := tc.payload()

			err := suite.store.PositionStorage().AssignCourses(
				context.TODO(),
				positionID,
				courses,
			)

			suite.Equal(tc.err, err)
		})
	}
}

func (suite *storeTestSuite) TestAssignCourse() {
	suite.NotNil(suite.store)

	// добавление компании
	company, err := suite.store.CompanyStorage().CreateCompany(context.TODO(), "test&Co")
	suite.NoError(err)
	suite.NotEmpty(company)

	// добавление должности
	position, err := suite.store.PositionStorage().CreatePosition(
		context.TODO(),
		model.PositionSet{CompanyID: company.ID, Name: "test-position"},
	)
	suite.NoError(err)
	suite.NotEmpty(position)

	// TODO: после добавления репозитория курсов, нужно дописать тест
}
