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

func (suite *storeTestSuite) TestGetPositionDB() {
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
			_, err := suite.store.PositionStorage().GetPositionInComp(context.TODO(), companyID, positionID)
			suite.Equal(tc.err, err)
		})
	}
}

func (suite *storeTestSuite) TestGetPositionsDB() {
	suite.NotNil(suite.store)

	// поиск в пустой базе должностей по id несуществующей компании
	positions, err := suite.store.PositionStorage().ListPositions(context.TODO(), randomseq.RandomTestInt())
	suite.Error(err)
	suite.Empty(positions)

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

func (suite *storeTestSuite) TestUpdatePositionDB() {
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

func (suite *storeTestSuite) TestAssignCourseDB() {
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
