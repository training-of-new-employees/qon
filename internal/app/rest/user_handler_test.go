package rest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/training-of-new-employees/qon/internal/errs"
	"github.com/training-of-new-employees/qon/internal/model"

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
