package pg

import (
	"context"

	"github.com/training-of-new-employees/qon/internal/model"
)

func (suite *storeTestSuite) TestCreateAdmin() {
	suite.NotNil(suite.store)

	admin := model.AdminCreate{
		Email:      "some@example.org",
		Password:   "Qwerty12345",
		Name:       "Admin",
		Patronymic: "Admin",
		Surname:    "Admin",
		IsAdmin:    true,
		IsActive:   true,
	}

	companyName := "test&Co"

	_, err := suite.store.UserStorage().CreateAdmin(context.TODO(), admin, companyName)

	suite.NoError(err)
}

func (suite *storeTestSuite) TestGetUserByEmail() {
	suite.NotNil(suite.store)

	email := "some@example.org"

	_, err := suite.store.UserStorage().GetUserByEmail(context.TODO(), email)

	suite.Error(err)
}
