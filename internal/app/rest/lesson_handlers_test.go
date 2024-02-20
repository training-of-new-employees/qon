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

func (suite *handlerTestSuite) TestUpdateLessonStatus() {
	userID := 1

	testCases := []struct {
		name         string
		expectedCode int
		prepare      func() ([]byte, string)
	}{
		{
			name:         "success",
			expectedCode: http.StatusOK,
			prepare: func() ([]byte, string) {
				dto := model.LessonStatusUpdate{
					Status: "in-process",
				}
				lessonID := "1"

				suite.lessonService.
					EXPECT().
					UpdateLessonStatus(gomock.Any(), userID, 1, dto.Status).
					Return(nil)

				body, _ := json.Marshal(dto)
				return body, lessonID
			},
		},
		{
			name:         "invalid lesson id",
			expectedCode: http.StatusBadRequest,
			prepare: func() ([]byte, string) {
				dto := model.LessonStatusUpdate{
					Status: "in-process",
				}
				lessonID := "hello"

				body, _ := json.Marshal(dto)
				return body, lessonID
			},
		},
		{
			name:         "invalid body",
			expectedCode: http.StatusBadRequest,
			prepare: func() ([]byte, string) {
				lessonID := "8"
				return []byte("hello"), lessonID
			},
		},
		{
			name:         "lesson not found",
			expectedCode: http.StatusNotFound,
			prepare: func() ([]byte, string) {
				dto := model.LessonStatusUpdate{
					Status: "in-process",
				}
				lessonID := "1"

				suite.lessonService.
					EXPECT().
					UpdateLessonStatus(gomock.Any(), userID, 1, dto.Status).
					Return(errs.ErrLessonNotFound)

				body, _ := json.Marshal(dto)
				return body, lessonID
			},
		},
		{
			name:         "internal error",
			expectedCode: http.StatusInternalServerError,
			prepare: func() ([]byte, string) {
				dto := model.LessonStatusUpdate{
					Status: "in-process",
				}
				lessonID := "1"

				suite.lessonService.
					EXPECT().
					UpdateLessonStatus(gomock.Any(), userID, 1, dto.Status).
					Return(errs.ErrInternal)

				body, _ := json.Marshal(dto)
				return body, lessonID
			},
		},
	}

	accessToken, err := jwttoken.TestAuthorizateUser(userID, 1, false)
	suite.NoError(err)
	suite.cache.EXPECT().GetRefreshToken(gomock.Any(), gomock.Any()).Return("", nil).AnyTimes()

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			body, lessonID := tc.prepare()

			w := httptest.NewRecorder()

			req, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/users/lessons/%s", lessonID), bytes.NewBuffer(body))
			req.Header.Set("Authorization", accessToken)

			suite.srv.ServeHTTP(w, req)
			suite.Equal(tc.expectedCode, w.Code)
		})
	}
}

func (suite *handlerTestSuite) TestGetUserLesson() {
	userID := 1

	testCases := []struct {
		name         string
		expectedCode int
		prepare      func() string
	}{
		{
			name:         "success",
			expectedCode: http.StatusOK,
			prepare: func() string {
				lessonID := "1"

				suite.lessonService.
					EXPECT().
					GetUserLesson(gomock.Any(), userID, 1).
					Return(&model.Lesson{}, nil)

				return lessonID
			},
		},
		{
			name:         "invalid lesson id",
			expectedCode: http.StatusBadRequest,
			prepare: func() string {
				lessonID := "hello"

				return lessonID
			},
		},
		{
			name:         "lesson not found",
			expectedCode: http.StatusNotFound,
			prepare: func() string {
				lessonID := "1"

				suite.lessonService.
					EXPECT().
					GetUserLesson(gomock.Any(), userID, 1).
					Return(nil, errs.ErrLessonNotFound)

				return lessonID
			},
		},
		{
			name:         "course not found",
			expectedCode: http.StatusNotFound,
			prepare: func() string {
				lessonID := "1"

				suite.lessonService.
					EXPECT().
					GetUserLesson(gomock.Any(), userID, 1).
					Return(nil, errs.ErrCourseNotFound)

				return lessonID
			},
		},
		{
			name:         "internal error",
			expectedCode: http.StatusInternalServerError,
			prepare: func() string {
				lessonID := "1"

				suite.lessonService.
					EXPECT().
					GetUserLesson(gomock.Any(), userID, 1).
					Return(nil, errs.ErrInternal)

				return lessonID
			},
		},
	}

	accessToken, err := jwttoken.TestAuthorizateUser(userID, 1, false)
	suite.NoError(err)
	suite.cache.EXPECT().GetRefreshToken(gomock.Any(), gomock.Any()).Return("", nil).AnyTimes()

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			lessonID := tc.prepare()

			w := httptest.NewRecorder()

			req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/users/lessons/%s", lessonID), nil)
			req.Header.Set("Authorization", accessToken)

			suite.srv.ServeHTTP(w, req)
			suite.Equal(tc.expectedCode, w.Code)
		})
	}
}
