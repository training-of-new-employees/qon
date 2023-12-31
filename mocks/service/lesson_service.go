// Code generated by MockGen. DO NOT EDIT.
// Source: internal/service/lesson_service.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	context "context"
	reflect "reflect"

	model "github.com/training-of-new-employees/qon/internal/model"
	gomock "go.uber.org/mock/gomock"
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
func (m *MockServiceLesson) CreateLesson(ctx context.Context, lesson model.LessonCreate, user_id int) (*model.Lesson, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateLesson", ctx, lesson, user_id)
	ret0, _ := ret[0].(*model.Lesson)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateLesson indicates an expected call of CreateLesson.
func (mr *MockServiceLessonMockRecorder) CreateLesson(ctx, lesson, user_id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateLesson", reflect.TypeOf((*MockServiceLesson)(nil).CreateLesson), ctx, lesson, user_id)
}

// DeleteLesson mocks base method.
func (m *MockServiceLesson) DeleteLesson(ctx context.Context, lessonID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteLesson", ctx, lessonID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteLesson indicates an expected call of DeleteLesson.
func (mr *MockServiceLessonMockRecorder) DeleteLesson(ctx, lessonID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteLesson", reflect.TypeOf((*MockServiceLesson)(nil).DeleteLesson), ctx, lessonID)
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
