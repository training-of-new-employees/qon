package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"go.uber.org/mock/gomock"

	"github.com/training-of-new-employees/qon/internal/errs"
	"github.com/training-of-new-employees/qon/internal/model"
	"github.com/training-of-new-employees/qon/internal/pkg/jwttoken"
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
			suite.cache.EXPECT().GetRefreshToken(gomock.Any(), gomock.Any()).Return("", nil)

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
		prepare      func() string
	}{
		{
			name:         "success",
			expectedCode: http.StatusOK,
			prepare: func() string {
				positionSet := model.NewTestPositionSet()
				positionSet.CompanyID = companyID

				positionID := 2

				suite.positionService.EXPECT().GetPosition(gomock.Any(), companyID, positionID).
					Return(&model.Position{ID: 2, CompanyID: companyID, Name: "Test", IsActive: true, IsArchived: false}, nil)

				return fmt.Sprint(positionID)
			},
		},
		{
			name:         "not found",
			expectedCode: http.StatusNotFound,
			prepare: func() string {
				positionSet := model.NewTestPositionSet()
				positionSet.CompanyID = companyID

				positionID := 2

				suite.positionService.EXPECT().GetPosition(gomock.Any(), companyID, positionID).
					Return(nil, errs.ErrPositionNotFound)

				return fmt.Sprint(positionID)
			},
		},
		{
			name:         "invalid position id",
			expectedCode: http.StatusBadRequest,
			prepare: func() string {
				positionID := "invalid"

				return fmt.Sprint(positionID)
			},
		},
		{
			name:         "internal error",
			expectedCode: http.StatusInternalServerError,
			prepare: func() string {
				positionSet := model.NewTestPositionSet()
				positionSet.CompanyID = companyID

				positionID := 2

				suite.positionService.EXPECT().GetPosition(gomock.Any(), companyID, positionID).
					Return(nil, errs.ErrInternal)

				return fmt.Sprint(positionID)
			},
		},
	}

	// получение тестового токена для авторизации админа
	accessToken, err := jwttoken.TestAuthorizateUser(1, companyID, true)
	suite.NoError(err)

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			positionID := tc.prepare()
			suite.cache.EXPECT().GetRefreshToken(gomock.Any(), gomock.Any()).Return("", nil)

			w := httptest.NewRecorder()

			req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/positions/%s", positionID), nil)
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
			suite.cache.EXPECT().GetRefreshToken(gomock.Any(), gomock.Any()).Return("", nil)

			w := httptest.NewRecorder()

			req, _ := http.NewRequest(http.MethodGet, "/api/v1/positions", nil)
			req.Header.Set("Authorization", accessToken)

			suite.srv.ServeHTTP(w, req)
			suite.Equal(tc.expectedCode, w.Code)
		})
	}
}

func (suite *handlerTestSuite) TestUpdatePosition() {
	companyID := 1

	testCases := []struct {
		name         string
		expectedCode int
		prepare      func() (string, []byte)
	}{
		{
			name:         "success",
			expectedCode: http.StatusOK,
			prepare: func() (string, []byte) {
				positionID := 2

				positionSet := model.NewTestPositionSet()
				positionSet.CompanyID = companyID

				position := &model.Position{
					ID:         positionID,
					CompanyID:  positionSet.CompanyID,
					Name:       positionSet.Name,
					IsActive:   true,
					IsArchived: false,
				}

				suite.positionService.EXPECT().UpdatePosition(gomock.Any(), positionID, positionSet).Return(position, nil)

				body, _ := json.Marshal(positionSet)

				return fmt.Sprint(positionID), body
			},
		},
		{
			name:         "not found",
			expectedCode: http.StatusNotFound,
			prepare: func() (string, []byte) {
				positionID := 2

				positionSet := model.NewTestPositionSet()
				positionSet.CompanyID = companyID

				suite.positionService.EXPECT().UpdatePosition(gomock.Any(), positionID, positionSet).Return(nil, errs.ErrPositionNotFound)

				body, _ := json.Marshal(positionSet)

				return fmt.Sprint(positionID), body
			},
		},
		{
			name:         "bad request",
			expectedCode: http.StatusBadRequest,
			prepare: func() (string, []byte) {
				positionID := 2

				body, _ := json.Marshal("invalid")

				return fmt.Sprint(positionID), body
			},
		},
		{
			name:         "invalid company id",
			expectedCode: http.StatusBadRequest,
			prepare: func() (string, []byte) {
				positionID := "invalid"

				return fmt.Sprint(positionID), nil
			},
		},
		{
			name:         "invalid company name",
			expectedCode: http.StatusBadRequest,
			prepare: func() (string, []byte) {
				positionID := 2

				positionSet := model.NewTestPositionSet()
				positionSet.CompanyID = companyID
				positionSet.Name = "invalid-name-#$%!"

				body, _ := json.Marshal(positionSet)

				return fmt.Sprint(positionID), body
			},
		},
		{
			name:         "internal error",
			expectedCode: http.StatusInternalServerError,
			prepare: func() (string, []byte) {
				positionID := 2

				positionSet := model.NewTestPositionSet()
				positionSet.CompanyID = companyID

				suite.positionService.EXPECT().UpdatePosition(gomock.Any(), positionID, positionSet).Return(nil, errs.ErrInternal)

				body, _ := json.Marshal(positionSet)

				return fmt.Sprint(positionID), body
			},
		},
	}

	// получение тестового токена для авторизации админа
	accessToken, err := jwttoken.TestAuthorizateUser(1, companyID, true)
	suite.NoError(err)

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			positionID, body := tc.prepare()
			suite.cache.EXPECT().GetRefreshToken(gomock.Any(), gomock.Any()).Return("", nil)

			w := httptest.NewRecorder()

			req, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/positions/update/%s", positionID), bytes.NewBuffer(body))
			req.Header.Set("Authorization", accessToken)

			suite.srv.ServeHTTP(w, req)
			suite.Equal(tc.expectedCode, w.Code)
		})
	}
}
