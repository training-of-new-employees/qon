// Code generated by MockGen. DO NOT EDIT.
// Source: internal/service/lesson_service.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/training-of-new-employees/qon/internal/model"
)

// MockServiceLesson is a mock of ServiceLesson interface.
type MockServiceLesson struct {
	ctrl     *gomock.Controller
	recorder *MockServiceLessonMockRecorder
}

// MockServiceLessonMockRecorder is the mock recorder for MockServiceLesson.
type MockServiceLessonMockRecorder struct {
	mock *MockServiceLesson
}

// NewMockServiceLesson creates a new mock instance.
func NewMockServiceLesson(ctrl *gomock.Controller) *MockServiceLesson {
	mock := &MockServiceLesson{ctrl: ctrl}
	mock.recorder = &MockServiceLessonMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockServiceLesson) EXPECT() *MockServiceLessonMockRecorder {
	return m.recorder
}

// CreateLesson mocks base method.
func (m *MockServiceLesson) CreateLesson(ctx context.Context, lesson model.Lesson, userID int) (*model.Lesson, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateLesson", ctx, lesson, userID)
	ret0, _ := ret[0].(*model.Lesson)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateLesson indicates an expected call of CreateLesson.
func (mr *MockServiceLessonMockRecorder) CreateLesson(ctx, lesson, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateLesson", reflect.TypeOf((*MockServiceLesson)(nil).CreateLesson), ctx, lesson, userID)
}

// GetLesson mocks base method.
func (m *MockServiceLesson) GetLesson(ctx context.Context, lessonID int) (*model.Lesson, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLesson", ctx, lessonID)
	ret0, _ := ret[0].(*model.Lesson)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLesson indicates an expected call of GetLesson.
func (mr *MockServiceLessonMockRecorder) GetLesson(ctx, lessonID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLesson", reflect.TypeOf((*MockServiceLesson)(nil).GetLesson), ctx, lessonID)
}

// GetLessonsList mocks base method.
func (m *MockServiceLesson) GetLessonsList(ctx context.Context, courseID int) ([]model.Lesson, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLessonsList", ctx, courseID)
	ret0, _ := ret[0].([]model.Lesson)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLessonsList indicates an expected call of GetLessonsList.
func (mr *MockServiceLessonMockRecorder) GetLessonsList(ctx, courseID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLessonsList", reflect.TypeOf((*MockServiceLesson)(nil).GetLessonsList), ctx, courseID)
}

// UpdateLesson mocks base method.
func (m *MockServiceLesson) UpdateLesson(ctx context.Context, lesson model.LessonUpdate) (*model.Lesson, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateLesson", ctx, lesson)
	ret0, _ := ret[0].(*model.Lesson)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateLesson indicates an expected call of UpdateLesson.
func (mr *MockServiceLessonMockRecorder) UpdateLesson(ctx, lesson interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateLesson", reflect.TypeOf((*MockServiceLesson)(nil).UpdateLesson), ctx, lesson)
}
