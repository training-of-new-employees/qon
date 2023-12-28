package pg

import (
	"context"
	"math/rand"
	"time"

	"github.com/training-of-new-employees/qon/internal/errs"
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
			name:        "empty name",
			companyName: "",
			err:         errs.ErrCompanyNameNotEmpty,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			_, err := suite.store.CompanyStorage().CreateCompany(context.TODO(), tc.companyName)
			suite.Equal(tc.err, err)
		})
	}
}

func (suite *storeTestSuite) TestGetCompany() {
	suite.NotNil(suite.store)

	rnd := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))

	company, err := suite.store.CompanyStorage().CreateCompany(context.TODO(), "test-company")
	suite.NoError(err)
	suite.NotEmpty(company)

	testCases := []struct {
		name      string
		companyID int
		err       error
	}{
		{
			name:      "success",
			companyID: company.ID,
			err:       nil,
		},
		{
			name:      "not found",
			companyID: rnd.Intn(32) + 2,
			err:       errs.ErrCompanyNotFound,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			_, err := suite.store.CompanyStorage().GetCompany(context.TODO(), tc.companyID)
			suite.Equal(tc.err, err)
		})
	}
}
