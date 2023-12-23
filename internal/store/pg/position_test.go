package pg

import (
	"context"

	"github.com/training-of-new-employees/qon/internal/model"
)

func (suite *storeTestSuite) TestCreatePosition() {
	suite.NotNil(suite.store)

	admin := model.UserCreate{
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
