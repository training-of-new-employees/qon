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
				courseSet := model.NewTestCourseSet()
				courseSet.CreatedBy = creatorID
				courseSet.ID = 0

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
			name:         "success punctuation name",
			expectedCode: http.StatusCreated,
			prepare: func() []byte {
				courseSet := model.NewTestCourseSet()
				courseSet.CreatedBy = creatorID
				courseSet.ID = 0
				courseSet.Name = courseSet.Name + ",!№:;"

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
				courseSet := model.NewTestCourseSet()
				courseSet.Name = "invalid-name-*"

				body, _ := json.Marshal(courseSet)

				return body
			},
		},
		{
			name:         "invalid course name",
			expectedCode: http.StatusBadRequest,
			prepare: func() []byte {
				courseSet := model.NewTestCourseSet()
				courseSet.Name = "invalid-name-#"

				body, _ := json.Marshal(courseSet)

				return body
			},
		},
		{
			name:         "internal error",
			expectedCode: http.StatusInternalServerError,
			prepare: func() []byte {
				courseSet := model.NewTestCourseSet()
				courseSet.ID = 0
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
			suite.cache.EXPECT().GetRefreshToken(gomock.Any(), gomock.Any()).Return("", nil)

			w := httptest.NewRecorder()

			req, _ := http.NewRequest(http.MethodPost, "/api/v1/admin/courses", bytes.NewBuffer(body))
			req.Header.Set("Authorization", accessToken)

			suite.srv.ServeHTTP(w, req)
			suite.Equal(tc.expectedCode, w.Code)
		})
	}
}

func (suite *handlerTestSuite) TestGetAdminCourses() {
	companyID := 1

	testCases := []struct {
		name         string
		expectedCode int
		prepare      func() []byte
	}{
		{
			name:         "success",
			expectedCode: http.StatusOK,
			prepare: func() []byte {
				courseSet := model.NewTestCourseSet()

				courses := []model.Course{
					{
						ID:          1,
						Name:        courseSet.Name,
						Description: courseSet.Description,
						CreatedBy:   1,
						IsActive:    true,
						IsArchived:  false,
					},
				}

				suite.courseService.EXPECT().GetCompanyCourses(gomock.Any(), companyID).Return(courses, nil)
				return nil
			},
		},
		{
			name:         "not found courses",
			expectedCode: http.StatusNotFound,
			prepare: func() []byte {
				suite.courseService.EXPECT().GetCompanyCourses(gomock.Any(), companyID).Return(nil, errs.ErrCourseNotFound)

				return nil
			},
		},
		{
			name:         "internal error",
			expectedCode: http.StatusInternalServerError,
			prepare: func() []byte {
				suite.courseService.EXPECT().GetCompanyCourses(gomock.Any(), companyID).Return(nil, errs.ErrInternal)

				return nil
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

			req, _ := http.NewRequest(http.MethodGet, "/api/v1/admin/courses", bytes.NewBuffer(body))
			req.Header.Set("Authorization", accessToken)

			suite.srv.ServeHTTP(w, req)
			suite.Equal(tc.expectedCode, w.Code)
		})
	}
}

func (suite *handlerTestSuite) TestGetUserCourses() {
	userID := 1

	testCases := []struct {
		name         string
		expectedCode int
		prepare      func() []byte
	}{
		{
			name:         "success",
			expectedCode: http.StatusOK,
			prepare: func() []byte {
				courseSet := model.NewTestCourseSet()

				courses := []model.Course{
					{
						ID:          1,
						Name:        courseSet.Name,
						Description: courseSet.Description,
						CreatedBy:   1,
						IsActive:    true,
						IsArchived:  false,
					},
				}

				suite.courseService.EXPECT().GetUserCourses(gomock.Any(), userID).Return(courses, nil)
				return nil
			},
		},
		{
			name:         "not found courses",
			expectedCode: http.StatusNotFound,
			prepare: func() []byte {
				suite.courseService.EXPECT().GetUserCourses(gomock.Any(), userID).Return(nil, errs.ErrCourseNotFound)

				return nil
			},
		},
		{
			name:         "internal error",
			expectedCode: http.StatusInternalServerError,
			prepare: func() []byte {
				suite.courseService.EXPECT().GetUserCourses(gomock.Any(), userID).Return(nil, errs.ErrInternal)

				return nil
			},
		},
	}

	// получение тестового токена для авторизации пользователя
	accessToken, err := jwttoken.TestAuthorizateUser(userID, 1, false)
	suite.NoError(err)

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			body := tc.prepare()
			suite.cache.EXPECT().GetRefreshToken(gomock.Any(), gomock.Any()).Return("", nil).AnyTimes()

			w := httptest.NewRecorder()

			req, _ := http.NewRequest(http.MethodGet, "/api/v1/users/courses", bytes.NewBuffer(body))
			req.Header.Set("Authorization", accessToken)

			suite.srv.ServeHTTP(w, req)
			suite.Equal(tc.expectedCode, w.Code)
		})
	}
}

func (suite *handlerTestSuite) TestEditCourse() {
	userID := 1
	companyID := 2

	testCases := []struct {
		name         string
		expectedCode int
		prepare      func() (string, []byte)
	}{
		{
			name:         "success",
			expectedCode: http.StatusOK,
			prepare: func() (string, []byte) {
				courseSet := model.NewTestCourseSet()
				courseSet.CreatedBy = userID

				course := &model.Course{
					ID:          courseSet.ID,
					Name:        courseSet.Name,
					Description: courseSet.Description,
					CreatedBy:   userID,
					IsActive:    true,
					IsArchived:  false,
				}

				suite.courseService.EXPECT().EditCourse(gomock.Any(), courseSet, companyID).Return(course, nil)

				body, _ := json.Marshal(courseSet)

				return fmt.Sprintf("%d", courseSet.ID), body
			},
		},
		{
			name:         "invalid course id",
			expectedCode: http.StatusBadRequest,
			prepare: func() (string, []byte) {
				id := "invalid"
				body, _ := json.Marshal("invalid")

				return id, body
			},
		},
		{
			name:         "invalid body",
			expectedCode: http.StatusBadRequest,
			prepare: func() (string, []byte) {
				id := "1"
				body, _ := json.Marshal("invalid")

				return id, body
			},
		},
		{
			name:         "invalid course name",
			expectedCode: http.StatusBadRequest,
			prepare: func() (string, []byte) {
				courseSet := model.NewTestCourseSet()
				courseSet.Name = "invalid-name-#$%!"

				body, _ := json.Marshal(courseSet)

				return fmt.Sprintf("%d", courseSet.ID), body
			},
		},
		{
			name:         "empty course name",
			expectedCode: http.StatusOK,
			prepare: func() (string, []byte) {
				courseSet := model.NewTestCourseSet()
				courseSet.Name = ""
				courseSet.CreatedBy = userID

				course := &model.Course{
					ID:          courseSet.ID,
					Name:        courseSet.Name,
					Description: courseSet.Description,
					CreatedBy:   userID,
					IsActive:    true,
					IsArchived:  false,
				}

				suite.courseService.EXPECT().EditCourse(gomock.Any(), courseSet, companyID).Return(course, nil)

				body, _ := json.Marshal(courseSet)

				return fmt.Sprintf("%d", courseSet.ID), body
			},
		},
		{
			name:         "internal error",
			expectedCode: http.StatusInternalServerError,
			prepare: func() (string, []byte) {
				courseSet := model.NewTestCourseSet()
				courseSet.CreatedBy = userID

				suite.courseService.EXPECT().EditCourse(gomock.Any(), courseSet, companyID).Return(nil, errs.ErrInternal)

				body, _ := json.Marshal(courseSet)

				return fmt.Sprintf("%d", courseSet.ID), body
			},
		},
	}

	// получение тестового токена для авторизации админа
	accessToken, err := jwttoken.TestAuthorizateUser(userID, companyID, true)
	suite.NoError(err)

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			id, body := tc.prepare()
			suite.cache.EXPECT().GetRefreshToken(gomock.Any(), gomock.Any()).Return("", nil)

			w := httptest.NewRecorder()

			req, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/admin/courses/%s", id), bytes.NewBuffer(body))
			req.Header.Set("Authorization", accessToken)

			suite.srv.ServeHTTP(w, req)
			suite.Equal(tc.expectedCode, w.Code)
		})
	}
}

func (suite *handlerTestSuite) GetUserCourseLessons() {
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
				courseID := "8"

				suite.courseService.
					EXPECT().
					GetUserCourseLessons(gomock.Any(), userID, 8).
					Return([]model.Lesson{}, nil)

				return courseID
			},
		},
		{
			name:         "invalid course id",
			expectedCode: http.StatusBadRequest,
			prepare: func() string {
				return "invalid"
			},
		},
		{
			name:         "course not found",
			expectedCode: http.StatusNotFound,
			prepare: func() string {
				courseID := "8"

				suite.courseService.
					EXPECT().
					GetUserCourseLessons(gomock.Any(), userID, 8).
					Return(nil, errs.ErrCourseNotFound)

				return courseID
			},
		},
		{
			name:         "internal error",
			expectedCode: http.StatusNotFound,
			prepare: func() string {
				courseID := "8"

				suite.courseService.
					EXPECT().
					GetUserCourseLessons(gomock.Any(), userID, 8).
					Return(nil, errs.ErrInternal)

				return courseID
			},
		},
	}

	accessToken, err := jwttoken.TestAuthorizateUser(userID, 1, true)
	suite.NoError(err)
	suite.cache.EXPECT().GetRefreshToken(gomock.Any(), gomock.Any()).Return("", nil).AnyTimes()

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			courseID := tc.prepare()

			w := httptest.NewRecorder()

			req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/users/courses/%s/lessons", courseID), nil)
			req.Header.Set("Authorization", accessToken)

			suite.srv.ServeHTTP(w, req)
			suite.Equal(tc.expectedCode, w.Code)
		})
	}
}
