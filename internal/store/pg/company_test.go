package pg

import (
	"context"
)

func (suite *storeTestSuite) TestCreateCompany() {
	suite.NotNil(suite.store)

	testCases := []struct {
		name        string
		companyName string
		err         error
	}{
		{
			name:        "success",
			companyName: "Test&Co",
			err:         nil,
		},
		{
			name:        "fail",
			companyName: "",
			err:         nil,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			_, err := suite.store.CompanyStorage().CreateCompanyDB(context.TODO(), tc.companyName)
			suite.Equal(tc.err, err)
		})
	}
}
