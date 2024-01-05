package rest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/training-of-new-employees/qon/internal/errs"
	"github.com/training-of-new-employees/qon/internal/model"
	"github.com/training-of-new-employees/qon/internal/pkg/jwttoken"

	"go.uber.org/mock/gomock"
)

func (suite *handlerTestSuite) TestCreateAdminInCache() {
	testCases := []struct {
		name         string
		expectedCode int
		prepareMock  func() []byte
	}{
		{
			name:         "success",
			expectedCode: http.StatusCreated,
			prepareMock: func() []byte {
				// подготовка моков для выполнения тест-кейса
				admin := model.NewTestCreateAdmin()

				suite.userService.EXPECT().WriteAdminToCache(gomock.Any(), admin).Return(&admin, nil)

				body, _ := json.Marshal(admin)

				return body
			},
		},
		{
			name:         "already exist email",
			expectedCode: http.StatusConflict,
			prepareMock: func() []byte {
				// подготовка моков для выполнения тест-кейса
				admin := model.NewTestCreateAdmin()

				suite.userService.EXPECT().WriteAdminToCache(gomock.Any(), admin).Return(nil, errs.ErrEmailAlreadyExists)

				body, _ := json.Marshal(admin)

				return body
			},
		},
		{
			name:         "cannot send email",
			expectedCode: http.StatusInternalServerError,
			prepareMock: func() []byte {
				// подготовка моков для выполнения тест-кейса
				admin := model.NewTestCreateAdmin()

				suite.userService.EXPECT().WriteAdminToCache(gomock.Any(), admin).Return(nil, errs.ErrNotSendEmail)

				body, _ := json.Marshal(admin)

				return body
			},
		},
		{
			name:         "invalid data",
			expectedCode: http.StatusBadRequest,
			prepareMock: func() []byte {
				// подготовка моков для выполнения тест-кейса
				admin := model.NewTestCreateAdmin()
				admin.Email = "invalid"

				// TODO: если валидация входящих данных будет перенесена в сервис, раскомментировать
				// suite.userService.EXPECT().WriteAdminToCache(gomock.Any(), admin).Return(nil, errs.ErrInvalidEmail)

				body, _ := json.Marshal(admin)

				return body
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			body := tc.prepareMock()

			w := httptest.NewRecorder()

			req, _ := http.NewRequest(http.MethodPost, "/api/v1/admin/register", bytes.NewBuffer(body))

			suite.srv.ServeHTTP(w, req)
			suite.Equal(tc.expectedCode, w.Code)
		})
	}
}

func (suite *handlerTestSuite) TestCreateUser() {
	testCases := []struct {
		name         string
		expectedCode int
		prepare      func() []byte
	}{
		{
			name:         "success",
			expectedCode: http.StatusCreated,
			prepare: func() []byte {
				u := model.NewTestUserCreate()

				u.CompanyID = 1
				u.PositionID = 2

				user := &model.User{
					ID:         2,
					Email:      u.Email,
					CompanyID:  u.CompanyID,
					PositionID: u.PositionID,
					Name:       u.Name,
					Patronymic: u.Patronymic,
					Surname:    u.Surname,
				}
				suite.userService.EXPECT().CreateUser(gomock.Any(), u).Return(user, nil)

				body, _ := json.Marshal(u)

				return body
			},
		},
		{
			name:         "invalid email",
			expectedCode: http.StatusBadRequest,
			prepare: func() []byte {
				u := model.NewTestUserCreate()

				u.Email = "invalid"
				u.CompanyID = 1
				u.PositionID = 2

				// TODO: если валидация входящих данных будет перенесена в сервис, раскомментировать
				//
				// user := &model.User{
				//	ID: 2,
				//	Email: u.Email,
				//	CompanyID: u.CompanyID,
				//	PositionID: u.PositionID,
				//	Name: u.Name,
				//	Patronymic: u.Patronymic,
				//	Surname: u.Surname,
				// }
				// suite.userService.EXPECT().CreateUser(gomock.Any(), u).Return(user, nil)

				body, _ := json.Marshal(u)

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
			name:         "already exist email",
			expectedCode: http.StatusConflict,
			prepare: func() []byte {
				// подготовка моков для выполнения тест-кейса
				u := model.NewTestUserCreate()

				u.CompanyID = 1
				u.PositionID = 2

				suite.userService.EXPECT().CreateUser(gomock.Any(), u).Return(nil, errs.ErrEmailAlreadyExists)

				body, _ := json.Marshal(u)

				return body
			},
		},
		{
			name:         "internal error",
			expectedCode: http.StatusInternalServerError,
			prepare: func() []byte {
				// подготовка моков для выполнения тест-кейса
				u := model.NewTestUserCreate()

				u.CompanyID = 1
				u.PositionID = 2

				suite.userService.EXPECT().CreateUser(gomock.Any(), u).Return(nil, errs.ErrInternal)

				body, _ := json.Marshal(u)

				return body
			},
		},
	}

	// получение тестового токена для авторизации админа
	accessToken, err := jwttoken.TestAuthorizateUser(1, 1, true)
	suite.NoError(err)

	// проверка тест-кейсов
	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			body := tc.prepare()

			w := httptest.NewRecorder()

			req, _ := http.NewRequest(http.MethodPost, "/api/v1/admin/employee", bytes.NewBuffer(body))
			req.Header.Set("Authorization", accessToken)

			suite.srv.ServeHTTP(w, req)
			suite.Equal(tc.expectedCode, w.Code)
		})
	}
}

func (suite *handlerTestSuite) TestGetUsers() {
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
				suite.userService.EXPECT().GetUsersByCompany(gomock.Any(), companyID).Return(model.NewTestUsers(companyID), nil)

			},
		},
		{
			name:         "not found",
			expectedCode: http.StatusNotFound,
			prepare: func() {
				suite.userService.EXPECT().GetUsersByCompany(gomock.Any(), companyID).Return(nil, errs.ErrUserNotFound)
			},
		},
		{
			name:         "internal error",
			expectedCode: http.StatusInternalServerError,
			prepare: func() {
				suite.userService.EXPECT().GetUsersByCompany(gomock.Any(), companyID).Return(nil, errs.ErrInternal)
			},
		},
	}

	// получение тестового токена для авторизации админа
	accessToken, err := jwttoken.TestAuthorizateUser(1, companyID, true)
	suite.NoError(err)

	// проверка тест-кейсов
	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.prepare()

			w := httptest.NewRecorder()

			req, _ := http.NewRequest(http.MethodGet, "/api/v1/users", nil)
			req.Header.Set("Authorization", accessToken)

			suite.srv.ServeHTTP(w, req)
			suite.Equal(tc.expectedCode, w.Code)
		})
	}
}
