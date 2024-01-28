// Code generated by MockGen. DO NOT EDIT.
// Source: internal/store/course_repo.go
//
// Generated by this command:
//
//	mockgen -source=internal/store/course_repo.go -destination=mocks/store/course_repo.go
//
// Package mock_store is a generated GoMock package.
package mock_store

import (
	context "context"
	reflect "reflect"

	model "github.com/training-of-new-employees/qon/internal/model"
	gomock "go.uber.org/mock/gomock"
)

// MockRepositoryCourse is a mock of RepositoryCourse interface.
type MockRepositoryCourse struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryCourseMockRecorder
}

// MockRepositoryCourseMockRecorder is the mock recorder for MockRepositoryCourse.
type MockRepositoryCourseMockRecorder struct {
	mock *MockRepositoryCourse
}

// NewMockRepositoryCourse creates a new mock instance.
func NewMockRepositoryCourse(ctrl *gomock.Controller) *MockRepositoryCourse {
	mock := &MockRepositoryCourse{ctrl: ctrl}
	mock.recorder = &MockRepositoryCourseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepositoryCourse) EXPECT() *MockRepositoryCourseMockRecorder {
	return m.recorder
}

// CompanyCourse mocks base method.
func (m *MockRepositoryCourse) CompanyCourse(ctx context.Context, courseID, companyID int) (*model.Course, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CompanyCourse", ctx, courseID, companyID)
	ret0, _ := ret[0].(*model.Course)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CompanyCourse indicates an expected call of CompanyCourse.
func (mr *MockRepositoryCourseMockRecorder) CompanyCourse(ctx, courseID, companyID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CompanyCourse", reflect.TypeOf((*MockRepositoryCourse)(nil).CompanyCourse), ctx, courseID, companyID)
}

// CompanyCourses mocks base method.
func (m *MockRepositoryCourse) CompanyCourses(ctx context.Context, companyID int) ([]model.Course, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CompanyCourses", ctx, companyID)
	ret0, _ := ret[0].([]model.Course)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CompanyCourses indicates an expected call of CompanyCourses.
func (mr *MockRepositoryCourseMockRecorder) CompanyCourses(ctx, companyID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CompanyCourses", reflect.TypeOf((*MockRepositoryCourse)(nil).CompanyCourses), ctx, companyID)
}

// CreateCourse mocks base method.
func (m *MockRepositoryCourse) CreateCourse(ctx context.Context, course model.CourseSet) (*model.Course, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCourse", ctx, course)
	ret0, _ := ret[0].(*model.Course)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateCourse indicates an expected call of CreateCourse.
func (mr *MockRepositoryCourseMockRecorder) CreateCourse(ctx, course any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCourse", reflect.TypeOf((*MockRepositoryCourse)(nil).CreateCourse), ctx, course)
}

// EditCourse mocks base method.
func (m *MockRepositoryCourse) EditCourse(ctx context.Context, course model.CourseSet, companyID int) (*model.Course, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EditCourse", ctx, course, companyID)
	ret0, _ := ret[0].(*model.Course)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EditCourse indicates an expected call of EditCourse.
func (mr *MockRepositoryCourseMockRecorder) EditCourse(ctx, course, companyID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EditCourse", reflect.TypeOf((*MockRepositoryCourse)(nil).EditCourse), ctx, course, companyID)
}

// GetUserCourse mocks base method.
func (m *MockRepositoryCourse) GetUserCourse(ctx context.Context, courseID, userID int) (*model.Course, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserCourse", ctx, courseID, userID)
	ret0, _ := ret[0].(*model.Course)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserCourse indicates an expected call of GetUserCourse.
func (mr *MockRepositoryCourseMockRecorder) GetUserCourse(ctx, courseID, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserCourse", reflect.TypeOf((*MockRepositoryCourse)(nil).GetUserCourse), ctx, courseID, userID)
}

// GetUserCoursesStatus mocks base method.
func (m *MockRepositoryCourse) GetUserCoursesStatus(ctx context.Context, userID int, coursesIds []int) (map[int]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserCoursesStatus", ctx, userID, coursesIds)
	ret0, _ := ret[0].(map[int]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserCoursesStatus indicates an expected call of GetUserCoursesStatus.
func (mr *MockRepositoryCourseMockRecorder) GetUserCoursesStatus(ctx, userID, coursesIds any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserCoursesStatus", reflect.TypeOf((*MockRepositoryCourse)(nil).GetUserCoursesStatus), ctx, userID, coursesIds)
}

// UserCourses mocks base method.
func (m *MockRepositoryCourse) UserCourses(ctx context.Context, userID int) ([]model.Course, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserCourses", ctx, userID)
	ret0, _ := ret[0].([]model.Course)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UserCourses indicates an expected call of UserCourses.
func (mr *MockRepositoryCourseMockRecorder) UserCourses(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserCourses", reflect.TypeOf((*MockRepositoryCourse)(nil).UserCourses), ctx, userID)
}
