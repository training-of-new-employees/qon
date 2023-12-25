package pg

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/training-of-new-employees/qon/internal/errs"
	"github.com/training-of-new-employees/qon/internal/model"
)

func (suite *storeTestSuite) TestCreatePositionDB() {
	suite.NotNil(suite.store)

	company, err := suite.store.CompanyStorage().CreateCompanyDB(context.TODO(), "test&Co")
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
			_, err := suite.store.PositionStorage().CreatePositionDB(context.TODO(), tc.position())
			suite.Equal(tc.err, err)
		})
	}

	suite.NoError(err)
}

func (suite *storeTestSuite) TestGetPositionDB() {
	suite.NotNil(suite.store)

	rnd := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))

	company, err := suite.store.CompanyStorage().CreateCompanyDB(context.TODO(), "test&Co")
	suite.NoError(err)
	suite.NotEmpty(company)

	position, err := suite.store.PositionStorage().CreatePositionDB(
		context.TODO(),
		model.PositionSet{CompanyID: company.ID, Name: "test-position"},
	)
	suite.NoError(err)
	suite.NotEmpty(position)

	testCases := []struct {
		name    string
		payload func() (int, int)
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
				companyID := rnd.Intn(32) + 1
				return companyID, position.ID
			},
			err: errs.ErrPositionNotFound,
		},
		{
			name: "random position",
			payload: func() (int, int) {
				positionID := rnd.Intn(32)
				return company.ID, positionID
			},
			err: errs.ErrPositionNotFound,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			companyID, positionID := tc.payload()
			_, err := suite.store.PositionStorage().GetPositionDB(context.TODO(), companyID, positionID)
			suite.Equal(tc.err, err)
		})
	}
}

func (suite *storeTestSuite) TestGetPositionsDB() {
	suite.NotNil(suite.store)

	// генерация случайного id
	rnd := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))
	companyID := rnd.Intn(32)

	// поиск в пустой базе должностей по id несуществующей компании
	positions, err := suite.store.PositionStorage().GetPositionsDB(context.TODO(), companyID)
	suite.Error(err)
	suite.Empty(positions)

	// создание компании
	company, err := suite.store.CompanyStorage().CreateCompanyDB(context.TODO(), "test&Co")
	suite.NoError(err)
	suite.NotEmpty(company)

	// добавление новых должностей к созданной компании (кол-во должностей, которое необходимо добавить, генерируется случайно)
	countPositions := rnd.Intn(32)

	expectedIDs := []int{}
	for i := 0; i < countPositions; i++ {
		position, err := suite.store.PositionStorage().CreatePositionDB(
			context.TODO(),
			model.PositionSet{CompanyID: company.ID, Name: fmt.Sprintf("test%d-position", i)},
		)
		suite.NoError(err)

		// добавляем в массив идентификаторы добавленных должностей
		expectedIDs = append(expectedIDs, position.ID)
	}

	// Получаем добавленные должности
	positions, err = suite.store.PositionStorage().GetPositionsDB(context.TODO(), company.ID)
	suite.NotEmpty(positions)
	suite.NoError(err)

	// добавляем в массив идентификаторы полученных должностей
	actualIDs := []int{}
	for _, p := range positions {
		actualIDs = append(actualIDs, p.ID)
	}

	suite.EqualValues(expectedIDs, actualIDs)
}

func (suite *storeTestSuite) TestGetPositionByID() {
	suite.NotNil(suite.store)

	rnd := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))

	company, err := suite.store.CompanyStorage().CreateCompanyDB(context.TODO(), "test&Co")
	suite.NoError(err)
	suite.NotEmpty(company)

	position, err := suite.store.PositionStorage().CreatePositionDB(
		context.TODO(),
		model.PositionSet{CompanyID: company.ID, Name: "test-position"},
	)
	suite.NoError(err)
	suite.NotEmpty(position)

	testCases := []struct {
		name    string
		payload func() int
		err     error
	}{
		{
			name: "success",
			payload: func() int {
				return position.ID
			},
			err: nil,
		},
		{
			name: "random position",
			payload: func() int {
				positionID := rnd.Intn(32)
				return positionID
			},
			err: errs.ErrPositionNotFound,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			positionID := tc.payload()
			_, err := suite.store.PositionStorage().GetPositionByID(context.TODO(), positionID)
			suite.Equal(tc.err, err)
		})
	}
}

func (suite *storeTestSuite) TestUpdatePositionDB() {
	suite.NotNil(suite.store)

	rnd := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))

	company, err := suite.store.CompanyStorage().CreateCompanyDB(context.TODO(), "test&Co")
	suite.NoError(err)
	suite.NotEmpty(company)

	position, err := suite.store.PositionStorage().CreatePositionDB(
		context.TODO(),
		model.PositionSet{CompanyID: company.ID, Name: "test-position"},
	)
	suite.NoError(err)
	suite.NotEmpty(position)

	testCases := []struct {
		name    string
		payload func() (int, model.PositionSet)
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
				companyID := rnd.Intn(32)
				return position.ID, model.PositionSet{CompanyID: companyID, Name: "new-position-name"}
			},
			err: errs.ErrPositionNotFound,
		},
		{
			name: "position not found",
			payload: func() (int, model.PositionSet) {
				positionID := rnd.Intn(32)
				return positionID, model.PositionSet{CompanyID: company.ID, Name: "new-position-name"}
			},
			err: errs.ErrPositionNotFound,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			positionID, positionUpdate := tc.payload()

			_, err := suite.store.PositionStorage().UpdatePositionDB(
				context.TODO(),
				positionID,
				positionUpdate,
			)

			suite.Equal(tc.err, err)
		})
	}
}
