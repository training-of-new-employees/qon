package pg

import (
	"context"

	"github.com/training-of-new-employees/qon/internal/errs"
	"github.com/training-of-new-employees/qon/internal/model"
)

func (suite *storeTestSuite) TestCreateUser() {
	suite.NotNil(suite.store)

	company, err := suite.store.CompanyStorage().CreateCompanyDB(context.TODO(), "test-company")
	suite.NoError(err)
	suite.NotEmpty(company)

	position, err := suite.store.PositionStorage().CreatePositionDB(context.TODO(), model.PositionSet{CompanyID: company.ID, Name: "test-position"})
	suite.NoError(err)
	suite.NotEmpty(position)

	uniqueEmail := "some@example.org"

	testCases := []struct {
		name string
		data func() model.UserCreate
		err  error
	}{
		{
			name: "success",
			data: func() model.UserCreate {
				u := model.NewTestUserCreate()
				u.CompanyID = company.ID
				u.PositionID = position.ID
				u.Email = uniqueEmail

				return u
			},
			err: nil,
		},
		{
			name: "unique email",
			data: func() model.UserCreate {
				u := model.NewTestUserCreate()
				u.CompanyID = company.ID
				u.PositionID = position.ID
				u.Email = uniqueEmail

				return u
			},
			err: errs.ErrEmailAlreadyExists,
		},
		{
			name: "not reference company id",
			data: func() model.UserCreate {
				u := model.NewTestUserCreate()
				u.CompanyID = 0

				return u
			},
			err: errs.ErrCompanyReference,
		},
		{
			name: "not reference position id",
			data: func() model.UserCreate {
				u := model.NewTestUserCreate()
				u.CompanyID = company.ID
				u.PositionID = 0

				return u
			},
			err: errs.ErrPositionReference,
		},
		{
			name: "empty email",
			data: func() model.UserCreate {
				u := model.NewTestUserCreate()
				u.CompanyID = company.ID
				u.PositionID = position.ID
				u.Email = ""

				return u
			},
			err: errs.ErrEmailNotEmpty,
		},

		{
			name: "empty password",
			data: func() model.UserCreate {
				u := model.NewTestUserCreate()
				u.CompanyID = company.ID
				u.PositionID = position.ID
				u.Password = ""

				return u
			},
			err: errs.ErrPasswordNotEmpty,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			user := tc.data()
			_, err := suite.store.UserStorage().CreateUser(context.TODO(), user)

			suite.Equal(tc.err, err)
		})
	}
}

func (suite *storeTestSuite) TestCreateAdmin() {
	suite.NotNil(suite.store)

	uniqueEmail := "some@example.org"

	testCases := []struct {
		name string
		data func() (string, model.UserCreate)
		err  error
	}{
		{
			name: "success",
			data: func() (string, model.UserCreate) {
				u := model.NewTestUserCreate()
				u.Email = uniqueEmail

				return "test-company", u
			},
			err: nil,
		},
		{
			name: "success",
			data: func() (string, model.UserCreate) {
				u := model.NewTestUserCreate()
				u.Email = uniqueEmail

				return "test-company", u
			},
			err: errs.ErrEmailAlreadyExists,
		},
		{
			name: "empty company name",
			data: func() (string, model.UserCreate) {
				u := model.NewTestUserCreate()

				return "", u
			},
			err: errs.ErrCompanyNameNotEmpty,
		},
		{
			name: "empty email",
			data: func() (string, model.UserCreate) {
				u := model.NewTestUserCreate()
				u.Email = ""

				return "test-company", u
			},
			err: errs.ErrEmailNotEmpty,
		},
		{
			name: "empty password",
			data: func() (string, model.UserCreate) {
				u := model.NewTestUserCreate()
				u.Password = ""

				return "test-company", u
			},
			err: errs.ErrPasswordNotEmpty,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			companyName, user := tc.data()
			_, err := suite.store.UserStorage().CreateAdmin(context.TODO(), user, companyName)

			suite.Equal(tc.err, err)
		})
	}
}

func (suite *storeTestSuite) TestGetUserByEmail() {
	suite.NotNil(suite.store)

	email := "some@example.org"

	_, err := suite.store.UserStorage().GetUserByEmail(context.TODO(), email)

	suite.Error(err)
}
