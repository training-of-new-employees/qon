package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/training-of-new-employees/qon/internal/errs"
	"github.com/training-of-new-employees/qon/internal/model"
	"github.com/training-of-new-employees/qon/internal/pkg/jwttoken"
	"github.com/training-of-new-employees/qon/internal/pkg/randomseq"
	"go.uber.org/mock/gomock"
)

func (suite *handlerTestSuite) TestHandlerCreateAdminInCache() {
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

func (suite *handlerTestSuite) TestHandlerCreateUser() {
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
			suite.cache.EXPECT().GetRefreshToken(gomock.Any(), gomock.Any()).Return("", nil)

			w := httptest.NewRecorder()

			req, _ := http.NewRequest(http.MethodPost, "/api/v1/admin/employee", bytes.NewBuffer(body))
			req.Header.Set("Authorization", accessToken)

			suite.srv.ServeHTTP(w, req)
			suite.Equal(tc.expectedCode, w.Code)
		})
	}
}

func (suite *handlerTestSuite) TestHandlerGetUsers() {
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
				suite.userService.EXPECT().GetUsersByCompany(gomock.Any(), companyID).Return(model.NewTestListUsers(companyID), nil)

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
			suite.cache.EXPECT().GetRefreshToken(gomock.Any(), gomock.Any()).Return("", nil)

			w := httptest.NewRecorder()

			req, _ := http.NewRequest(http.MethodGet, "/api/v1/users", nil)
			req.Header.Set("Authorization", accessToken)

			suite.srv.ServeHTTP(w, req)
			suite.Equal(tc.expectedCode, w.Code)
		})
	}
}

func (suite *handlerTestSuite) TestHandlerGetUser() {
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
				userID := 2
				positionID := 2

				suite.userService.EXPECT().GetUserByID(gomock.Any(), userID).Return(model.NewTestUser(userID, companyID, positionID), nil)

				return fmt.Sprint(userID)
			},
		},
		{
			name:         "bad request",
			expectedCode: http.StatusBadRequest,
			prepare: func() string {
				userID := "invalid"

				return userID
			},
		},
		{
			name:         "not access",
			expectedCode: http.StatusForbidden,
			prepare: func() string {
				userID := 2
				positionID := 2

				suite.userService.EXPECT().GetUserByID(gomock.Any(), userID).Return(model.NewTestUser(userID, randomseq.RandomTestInt(), positionID), nil)

				return fmt.Sprint(userID)
			},
		},
		{
			name:         "not found",
			expectedCode: http.StatusNotFound,
			prepare: func() string {
				userID := 2

				suite.userService.EXPECT().GetUserByID(gomock.Any(), userID).Return(nil, errs.ErrUserNotFound)

				return fmt.Sprint(userID)
			},
		},
		{
			name:         "internal error",
			expectedCode: http.StatusInternalServerError,
			prepare: func() string {
				userID := 2

				suite.userService.EXPECT().GetUserByID(gomock.Any(), userID).Return(nil, errs.ErrInternal)

				return fmt.Sprint(userID)
			},
		},
	}

	// получение тестового токена для авторизации админа
	accessToken, err := jwttoken.TestAuthorizateUser(1, companyID, true)
	suite.NoError(err)

	// проверка тест-кейсов
	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			userID := tc.prepare()
			suite.cache.EXPECT().GetRefreshToken(gomock.Any(), gomock.Any()).Return("", nil)

			w := httptest.NewRecorder()

			req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/users/%s", userID), nil)
			req.Header.Set("Authorization", accessToken)

			suite.srv.ServeHTTP(w, req)
			suite.Equal(tc.expectedCode, w.Code)
		})
	}
}

func (suite *handlerTestSuite) TestHandlerEditUser() {
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
				userID := 2
				positionID := 2

				editField, expected := model.NewTestEditUser(userID, companyID, positionID)

				suite.userService.EXPECT().EditUser(gomock.Any(), &editField, companyID).Return(&expected, nil)

				body, _ := json.Marshal(editField)

				return fmt.Sprint(userID), body
			},
		},
		{
			name:         "invalid user id",
			expectedCode: http.StatusBadRequest,
			prepare: func() (string, []byte) {
				body, _ := json.Marshal(nil)

				return "invalid", body
			},
		},
		{
			name:         "invalid request body",
			expectedCode: http.StatusBadRequest,
			prepare: func() (string, []byte) {
				userID := 2

				body, _ := json.Marshal("invalid")

				return fmt.Sprint(userID), body
			},
		},
		{
			name:         "not found",
			expectedCode: http.StatusNotFound,
			prepare: func() (string, []byte) {
				userID := 2
				positionID := 2

				editField, _ := model.NewTestEditUser(userID, companyID, positionID)

				suite.userService.EXPECT().EditUser(gomock.Any(), &editField, companyID).Return(nil, errs.ErrUserNotFound)

				body, _ := json.Marshal(editField)

				return fmt.Sprint(userID), body
			},
		},
		{
			name:         "internal error",
			expectedCode: http.StatusInternalServerError,
			prepare: func() (string, []byte) {
				userID := 2
				positionID := 2

				editField, _ := model.NewTestEditUser(userID, companyID, positionID)

				suite.userService.EXPECT().EditUser(gomock.Any(), &editField, companyID).Return(nil, errs.ErrInternal)

				body, _ := json.Marshal(editField)

				return fmt.Sprint(userID), body

			},
		},
	}

	// получение тестового токена для авторизации админа
	accessToken, err := jwttoken.TestAuthorizateUser(1, companyID, true)
	suite.NoError(err)

	// проверка тест-кейсов
	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			userID, body := tc.prepare()
			suite.cache.EXPECT().GetRefreshToken(gomock.Any(), gomock.Any()).Return("", nil)

			w := httptest.NewRecorder()

			req, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/users/%s", userID), bytes.NewBuffer(body))
			req.Header.Set("Authorization", accessToken)

			suite.srv.ServeHTTP(w, req)
			suite.Equal(tc.expectedCode, w.Code)
		})
	}
}

func (suite *handlerTestSuite) TestHandlerEditAdmin() {
	userID := 1
	companyID := 1
	positionID := 1

	testCases := []struct {
		name         string
		expectedCode int
		prepare      func() []byte
	}{
		{
			name:         "success",
			expectedCode: http.StatusOK,
			prepare: func() []byte {
				editField, expected := model.NewTestAdminEdit(userID, companyID, positionID)

				suite.userService.EXPECT().EditAdmin(gomock.Any(), editField).Return(&expected, nil)

				body, _ := json.Marshal(editField)

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
			name:         "not found",
			expectedCode: http.StatusNotFound,
			prepare: func() []byte {
				editField, _ := model.NewTestAdminEdit(userID, companyID, positionID)

				suite.userService.EXPECT().EditAdmin(gomock.Any(), editField).Return(nil, errs.ErrUserNotFound)

				body, _ := json.Marshal(editField)

				return body
			},
		},
		{
			name:         "internal error",
			expectedCode: http.StatusInternalServerError,
			prepare: func() []byte {
				editField, _ := model.NewTestAdminEdit(userID, companyID, positionID)

				suite.userService.EXPECT().EditAdmin(gomock.Any(), editField).Return(nil, errs.ErrInternal)

				body, _ := json.Marshal(editField)

				return body
			},
		},
	}

	// получение тестового токена для авторизации админа
	accessToken, err := jwttoken.TestAuthorizateUser(userID, companyID, true)
	suite.NoError(err)

	// проверка тест-кейсов
	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			body := tc.prepare()
			suite.cache.EXPECT().GetRefreshToken(gomock.Any(), gomock.Any()).Return("", nil)

			w := httptest.NewRecorder()

			req, _ := http.NewRequest(http.MethodPatch, "/api/v1/admin/info", bytes.NewBuffer(body))
			req.Header.Set("Authorization", accessToken)

			suite.srv.ServeHTTP(w, req)
			suite.Equal(tc.expectedCode, w.Code)
		})
	}
}

func (suite *handlerTestSuite) TestHandlerResetPassword() {
	testCases := []struct {
		name         string
		expectedCode int
		prepare      func() []byte
	}{
		{
			name:         "success",
			expectedCode: http.StatusOK,
			prepare: func() []byte {
				u := model.NewTestResetPassword()
				suite.userService.
					EXPECT().
					ResetPassword(gomock.Any(), u.Email).
					Return(nil)

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
			name:         "user not found",
			expectedCode: http.StatusNotFound,
			prepare: func() []byte {
				u := model.NewTestResetPassword()
				suite.userService.
					EXPECT().
					ResetPassword(gomock.Any(), u.Email).
					Return(errs.ErrUserNotFound)

				body, _ := json.Marshal(u)

				return body
			},
		},
		{
			name:         "internal error",
			expectedCode: http.StatusInternalServerError,
			prepare: func() []byte {
				u := model.NewTestResetPassword()
				suite.userService.
					EXPECT().
					ResetPassword(gomock.Any(), u.Email).
					Return(errs.ErrInternal)

				body, _ := json.Marshal(u)

				return body
			},
		},
	}

	// проверка тест-кейсов
	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			body := tc.prepare()

			w := httptest.NewRecorder()

			req, _ := http.NewRequest(http.MethodPost, "/api/v1/password", bytes.NewBuffer(body))

			suite.srv.ServeHTTP(w, req)
			suite.Equal(tc.expectedCode, w.Code)
		})
	}
}

func (suite *handlerTestSuite) TestHandlerRegenerationInvitationLink() {
	userAdminID := 1
	userID := 2
	companyID := 1
	email := "user@mail.com"

	link := fmt.Sprintf("http://localhost/first-login?email=%s&invite=%s", email, randomseq.RandomString(20))

	tests := []struct {
		name         string
		userID       int
		userAdminID  int
		isAdmin      bool
		companyID    int
		prepare      func() []byte
		expectedCode int
	}{
		{
			name:         "success",
			userID:       userID,
			userAdminID:  userAdminID,
			companyID:    companyID,
			isAdmin:      true,
			expectedCode: http.StatusOK,
			prepare: func() []byte {
				response := model.NewInvitationLinkResponse(email, link)

				suite.userService.EXPECT().RegenerationInvitationLinkUser(gomock.Any(), email, companyID).Return(&response, nil)

				body, _ := json.Marshal(&response)

				return body
			},
		},
		{
			name:         "invalid request body",
			userID:       userID,
			userAdminID:  userAdminID,
			isAdmin:      true,
			companyID:    companyID,
			expectedCode: http.StatusBadRequest,
			prepare: func() []byte {

				body, _ := json.Marshal("invalid")

				return body
			},
		},
		{
			name:         "not found",
			userID:       userID,
			userAdminID:  userAdminID,
			isAdmin:      true,
			companyID:    companyID,
			expectedCode: http.StatusNotFound,
			prepare: func() []byte {
				response := model.NewInvitationLinkResponse(email, link)

				suite.userService.EXPECT().RegenerationInvitationLinkUser(gomock.Any(), email, companyID).Return(&response, errs.ErrUserNotFound)

				body, _ := json.Marshal(response)

				return body
			},
		},
		{
			name:         "internal error",
			userID:       userID,
			userAdminID:  userAdminID,
			isAdmin:      true,
			companyID:    companyID,
			expectedCode: http.StatusInternalServerError,
			prepare: func() []byte {
				response := model.NewInvitationLinkResponse(email, link)

				suite.userService.EXPECT().RegenerationInvitationLinkUser(gomock.Any(), email, companyID).Return(&response, errs.ErrInternal)

				body, _ := json.Marshal(response)

				return body
			},
		},
		{
			name:         "error the user is activated in the system",
			userID:       userID,
			userAdminID:  userAdminID,
			isAdmin:      true,
			companyID:    companyID,
			expectedCode: http.StatusConflict,
			prepare: func() []byte {
				response := model.NewInvitationLinkResponse(email, link)

				suite.userService.EXPECT().RegenerationInvitationLinkUser(gomock.Any(), email, companyID).Return(&response, errs.ErrUserActivated)

				body, _ := json.Marshal(response)

				return body
			},
		},
		{
			name:         "error the user does not have access to the handler",
			userID:       userID,
			userAdminID:  userAdminID,
			isAdmin:      false,
			companyID:    companyID,
			expectedCode: http.StatusForbidden,
			prepare: func() []byte {
				response := model.NewInvitationLinkResponse(email, link)

				body, _ := json.Marshal(&response)

				return body
			},
		},
		{
			name:         "user is not authorized",
			userID:       userID,
			isAdmin:      true,
			companyID:    companyID,
			expectedCode: http.StatusUnauthorized,
			prepare: func() []byte {
				response := model.NewInvitationLinkResponse(email, link)

				body, _ := json.Marshal(response)

				return body
			},
		},
	}

	for _, tt := range tests {
		// получение тестового токена для авторизации админа
		var accessToken string
		var err error
		if tt.userAdminID > 0 {
			accessToken, err = jwttoken.TestAuthorizateUser(tt.userAdminID, tt.companyID, tt.isAdmin)
			suite.cache.EXPECT().GetRefreshToken(gomock.Any(), gomock.Any()).Return("", nil)
			suite.NoError(err)
		}

		suite.Run(tt.name, func() {
			body := tt.prepare()

			w := httptest.NewRecorder()

			req, _ := http.NewRequest(http.MethodPatch, "/api/v1/invitation-link", bytes.NewBuffer(body))

			if len(accessToken) > 0 {
				req.Header.Set("Authorization", accessToken)
			}

			suite.srv.ServeHTTP(w, req)
			suite.Equal(tt.expectedCode, w.Code)
		})
	}
}

func (suite *handlerTestSuite) TestRestServer_handlerSetPassword() {
	userAdminID := 1
	userID := 2
	companyID := 1
	email := "user@mail.com"
	password := "user@maiL2.com"
	code := randomseq.RandomString(20)

	tests := []struct {
		name         string
		userID       int
		userAdminID  int
		isAdmin      bool
		companyID    int
		prepare      func() []byte
		expectedCode int
	}{
		{
			name:         "success",
			userID:       userID,
			userAdminID:  userAdminID,
			companyID:    companyID,
			isAdmin:      true,
			expectedCode: http.StatusOK,
			prepare: func() []byte {

				user := &model.User{
					ID:        1,
					Email:     email,
					CompanyID: companyID,
					IsAdmin:   false,
				}

				userActivate := model.UserActivation{
					Email:    email,
					Password: password,
					Invite:   code,
				}

				accessToken, err := jwttoken.TestAuthorizateUser(user.ID, user.CompanyID, user.IsAdmin)
				suite.NoError(err)
				token := &model.Tokens{
					AccessToken: accessToken,
				}

				suite.userService.EXPECT().GetUserInviteCodeFromCache(context.Background(), email).Return(code, nil)
				suite.userService.EXPECT().GenerateTokenPair(context.Background(), user.ID, user.IsAdmin, user.CompanyID).Return(token, nil)
				suite.userService.EXPECT().UpdatePasswordAndActivateUser(context.Background(), email, password).Return(user, nil)

				body, _ := json.Marshal(userActivate)

				return body
			},
		},
		{
			name:         "invalid request body",
			userID:       userID,
			userAdminID:  userAdminID,
			companyID:    companyID,
			isAdmin:      true,
			expectedCode: http.StatusBadRequest,
			prepare: func() []byte {

				body, _ := json.Marshal("invalid")

				return body
			},
		},
		{
			name:         "not found",
			userID:       userID,
			userAdminID:  userAdminID,
			companyID:    companyID,
			isAdmin:      true,
			expectedCode: http.StatusNotFound,
			prepare: func() []byte {

				userActivate := model.UserActivation{
					Email:    email,
					Password: password,
					Invite:   code,
				}

				suite.userService.EXPECT().GetUserInviteCodeFromCache(context.Background(), email).Return(code, nil)
				suite.userService.EXPECT().UpdatePasswordAndActivateUser(context.Background(), email, password).Return(nil, errs.ErrUserNotFound)

				body, _ := json.Marshal(userActivate)

				return body
			},
		},
		{
			name:         "invalid registration and authentication process could not be completed the invitation code does not match",
			userID:       userID,
			userAdminID:  userAdminID,
			companyID:    companyID,
			isAdmin:      true,
			expectedCode: http.StatusUnauthorized,
			prepare: func() []byte {

				userActivate := model.UserActivation{
					Email:    email,
					Password: password,
					Invite:   code,
				}

				suite.userService.EXPECT().GetUserInviteCodeFromCache(context.Background(), email).Return("test", nil)

				body, _ := json.Marshal(userActivate)

				return body
			},
		},

		{
			name:         "invalid registration and authentication process could not be completed user active",
			userID:       userID,
			userAdminID:  userAdminID,
			companyID:    companyID,
			isAdmin:      true,
			expectedCode: http.StatusUnauthorized,
			prepare: func() []byte {

				user := &model.User{
					ID:        1,
					Email:     email,
					CompanyID: companyID,
					IsAdmin:   false,
					IsActive:  true,
				}

				userActivate := model.UserActivation{
					Email:    email,
					Password: password,
					Invite:   code,
				}

				suite.userService.EXPECT().GetUserInviteCodeFromCache(context.Background(), email).Return(code, nil)
				suite.userService.EXPECT().UpdatePasswordAndActivateUser(context.Background(), email, password).Return(user, errs.ErrUnauthorized)

				body, _ := json.Marshal(userActivate)

				return body
			},
		},
	}
	for _, tt := range tests {

		suite.Run(tt.name, func() {
			body := tt.prepare()

			w := httptest.NewRecorder()

			req, _ := http.NewRequest(http.MethodPost, "/api/v1/users/set-password", bytes.NewBuffer(body))

			suite.srv.ServeHTTP(w, req)
			suite.Equal(tt.expectedCode, w.Code)
		})
	}
}
