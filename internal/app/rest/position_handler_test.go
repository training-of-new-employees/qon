package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/training-of-new-employees/qon/internal/errs"
	"github.com/training-of-new-employees/qon/internal/model"
	"github.com/training-of-new-employees/qon/internal/pkg/jwttoken"
	"go.uber.org/mock/gomock"
)

func (suite *handlerTestSuite) TestCreatePosition() {
	companyID := 1

	testCases := []struct {
		name         string
		expectedCode int
		prepare      func() []byte
	}{
		{
			name:         "success",
			expectedCode: http.StatusCreated,
			prepare: func() []byte {
				positionSet := model.NewTestPositionSet()
				positionSet.CompanyID = companyID

				position := &model.Position{
					ID:         2,
					CompanyID:  positionSet.CompanyID,
					Name:       positionSet.Name,
					IsActive:   true,
					IsArchived: false,
				}

				suite.positionService.EXPECT().CreatePosition(gomock.Any(), positionSet).Return(position, nil)

				body, _ := json.Marshal(positionSet)

				return body
			},
		},
		{
			name:         "invalid request body",
			expectedCode: http.StatusBadRequest,
			prepare: func() []byte {
				body, _ := json.Marshal("invalid")

				return body
			},
		},
		{
			name:         "invalid position name",
			expectedCode: http.StatusBadRequest,
			prepare: func() []byte {
				positionSet := model.NewTestPositionSet()
				positionSet.Name = "invalid-name-#$%!"

				body, _ := json.Marshal(positionSet)

				return body
			},
		},
		{
			name:         "not found company",
			expectedCode: http.StatusNotFound,
			prepare: func() []byte {
				positionSet := model.NewTestPositionSet()
				positionSet.CompanyID = 100

				body, _ := json.Marshal(positionSet)

				return body
			},
		},
		{
			name:         "internal error",
			expectedCode: http.StatusInternalServerError,
			prepare: func() []byte {
				positionSet := model.NewTestPositionSet()
				positionSet.CompanyID = companyID

				suite.positionService.EXPECT().CreatePosition(gomock.Any(), positionSet).Return(nil, errs.ErrInternal)

				body, _ := json.Marshal(positionSet)

				return body
			},
		},
	}

	// получение тестового токена для авторизации админа
	accessToken, err := jwttoken.TestAuthorizateUser(1, companyID, true)
	suite.NoError(err)

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			body := tc.prepare()

			w := httptest.NewRecorder()

			req, _ := http.NewRequest(http.MethodPost, "/api/v1/positions", bytes.NewBuffer(body))
			req.Header.Set("Authorization", accessToken)

			suite.srv.ServeHTTP(w, req)
			suite.Equal(tc.expectedCode, w.Code)
		})
	}
}

func (suite *handlerTestSuite) TestGetPosition() {
	companyID := 1

	testCases := []struct {
		name         string
		expectedCode int
		prepare      func() int
	}{
		{
			name:         "success",
			expectedCode: http.StatusOK,
			prepare: func() int {
				positionSet := model.NewTestPositionSet()
				positionSet.CompanyID = companyID

				positionID := 2

				suite.positionService.EXPECT().GetPosition(gomock.Any(), companyID, positionID).
					Return(&model.Position{ID: 2, CompanyID: companyID, Name: "Test", IsActive: true, IsArchived: false}, nil)

				return positionID
			},
		},
		{
			name:         "not found",
			expectedCode: http.StatusNotFound,
			prepare: func() int {
				positionSet := model.NewTestPositionSet()
				positionSet.CompanyID = companyID

				positionID := 2

				suite.positionService.EXPECT().GetPosition(gomock.Any(), companyID, positionID).
					Return(nil, errs.ErrPositionNotFound)

				return positionID
			},
		},
		{
			name:         "internal error",
			expectedCode: http.StatusInternalServerError,
			prepare: func() int {
				positionSet := model.NewTestPositionSet()
				positionSet.CompanyID = companyID

				positionID := 2

				suite.positionService.EXPECT().GetPosition(gomock.Any(), companyID, positionID).
					Return(nil, errs.ErrInternal)

				return positionID
			},
		},
	}

	// получение тестового токена для авторизации админа
	accessToken, err := jwttoken.TestAuthorizateUser(1, companyID, true)
	suite.NoError(err)

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			positionID := tc.prepare()

			w := httptest.NewRecorder()

			req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/positions/%d", positionID), nil)
			req.Header.Set("Authorization", accessToken)

			suite.srv.ServeHTTP(w, req)
			suite.Equal(tc.expectedCode, w.Code)
		})
	}
}

func (suite *handlerTestSuite) TestGetPositions() {
	companyID := 1

	testCases := []struct {
		name         string
		expectedCode int
		prepare      func()
	}{
		{
			name:         "success",
			expectedCode: http.StatusOK,
			prepare: func() {
				suite.positionService.EXPECT().GetPositions(gomock.Any(), companyID).Return(model.NewTestPositions(companyID), nil)
			},
		},
		{
			name:         "not found",
			expectedCode: http.StatusNotFound,
			prepare: func() {
				suite.positionService.EXPECT().GetPositions(gomock.Any(), companyID).Return(nil, errs.ErrPositionNotFound)
			},
		},
		{
			name:         "internal error",
			expectedCode: http.StatusInternalServerError,
			prepare: func() {
				suite.positionService.EXPECT().GetPositions(gomock.Any(), companyID).Return(nil, errs.ErrInternal)
			},
		},
	}

	// получение тестового токена для авторизации админа
	accessToken, err := jwttoken.TestAuthorizateUser(1, companyID, true)
	suite.NoError(err)

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.prepare()

			w := httptest.NewRecorder()

			req, _ := http.NewRequest(http.MethodGet, "/api/v1/positions", nil)
			req.Header.Set("Authorization", accessToken)

			suite.srv.ServeHTTP(w, req)
			suite.Equal(tc.expectedCode, w.Code)
		})
	}
}
