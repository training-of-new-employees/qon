package rest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"go.uber.org/mock/gomock"

	"github.com/training-of-new-employees/qon/internal/errs"
	"github.com/training-of-new-employees/qon/internal/model"
	"github.com/training-of-new-employees/qon/internal/pkg/jwttoken"
)

func (suite *handlerTestSuite) TestCreateCourse() {
	creatorID := 1
	courseID := 2

	testCases := []struct {
		name         string
		expectedCode int
		prepare      func() []byte
	}{
		{
			name:         "success",
			expectedCode: http.StatusCreated,
			prepare: func() []byte {
				courseSet := model.NewTestCreateCourse()
				courseSet.CreatedBy = creatorID

				course := &model.Course{
					ID:          courseID,
					Name:        courseSet.Name,
					Description: courseSet.Description,
					CreatedBy:   creatorID,
					IsActive:    true,
					IsArchived:  false,
				}

				suite.courseService.EXPECT().CreateCourse(gomock.Any(), courseSet).Return(course, nil)

				body, _ := json.Marshal(courseSet)

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
			name:         "invalid course name",
			expectedCode: http.StatusBadRequest,
			prepare: func() []byte {
				courseSet := model.NewTestCreateCourse()
				courseSet.Name = "invalid-name-#$%!"

				body, _ := json.Marshal(courseSet)

				return body
			},
		},
		{
			name:         "internal error",
			expectedCode: http.StatusInternalServerError,
			prepare: func() []byte {
				courseSet := model.NewTestCreateCourse()
				courseSet.CreatedBy = creatorID

				suite.courseService.EXPECT().CreateCourse(gomock.Any(), courseSet).Return(nil, errs.ErrInternal)

				body, _ := json.Marshal(courseSet)

				return body
			},
		},
	}

	// получение тестового токена для авторизации админа
	accessToken, err := jwttoken.TestAuthorizateUser(creatorID, 1, true)
	suite.NoError(err)

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			body := tc.prepare()

			w := httptest.NewRecorder()

			req, _ := http.NewRequest(http.MethodPost, "/api/v1/admin/courses", bytes.NewBuffer(body))
			req.Header.Set("Authorization", accessToken)

			suite.srv.ServeHTTP(w, req)
			suite.Equal(tc.expectedCode, w.Code)
		})
	}
}
