// Code generated by MockGen. DO NOT EDIT.
// Source: internal/service/position_service.go
//
// Generated by this command:
//
//	mockgen -source=internal/service/position_service.go -destination=mocks//service/position_service.go
//
// Package mock_service is a generated GoMock package.
package mock_service

import (
	context "context"
	reflect "reflect"

	model "github.com/training-of-new-employees/qon/internal/model"
	gomock "go.uber.org/mock/gomock"
)

// MockServicePosition is a mock of ServicePosition interface.
type MockServicePosition struct {
	ctrl     *gomock.Controller
	recorder *MockServicePositionMockRecorder
}

// MockServicePositionMockRecorder is the mock recorder for MockServicePosition.
type MockServicePositionMockRecorder struct {
	mock *MockServicePosition
}

// NewMockServicePosition creates a new mock instance.
func NewMockServicePosition(ctrl *gomock.Controller) *MockServicePosition {
	mock := &MockServicePosition{ctrl: ctrl}
	mock.recorder = &MockServicePositionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockServicePosition) EXPECT() *MockServicePositionMockRecorder {
	return m.recorder
}

// AssignCourse mocks base method.
func (m *MockServicePosition) AssignCourse(ctx context.Context, positionID, courseID, user_id int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AssignCourse", ctx, positionID, courseID, user_id)
	ret0, _ := ret[0].(error)
	return ret0
}

// AssignCourse indicates an expected call of AssignCourse.
func (mr *MockServicePositionMockRecorder) AssignCourse(ctx, positionID, courseID, user_id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AssignCourse", reflect.TypeOf((*MockServicePosition)(nil).AssignCourse), ctx, positionID, courseID, user_id)
}

// CreatePosition mocks base method.
func (m *MockServicePosition) CreatePosition(ctx context.Context, position model.PositionCreate) (*model.Position, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePosition", ctx, position)
	ret0, _ := ret[0].(*model.Position)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreatePosition indicates an expected call of CreatePosition.
func (mr *MockServicePositionMockRecorder) CreatePosition(ctx, position any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePosition", reflect.TypeOf((*MockServicePosition)(nil).CreatePosition), ctx, position)
}

// DeletePosition mocks base method.
func (m *MockServicePosition) DeletePosition(ctx context.Context, id, companyID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeletePosition", ctx, id, companyID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeletePosition indicates an expected call of DeletePosition.
func (mr *MockServicePositionMockRecorder) DeletePosition(ctx, id, companyID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePosition", reflect.TypeOf((*MockServicePosition)(nil).DeletePosition), ctx, id, companyID)
}

// GetPosition mocks base method.
func (m *MockServicePosition) GetPosition(ctx context.Context, companyID, positionID int) (*model.Position, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPosition", ctx, companyID, positionID)
	ret0, _ := ret[0].(*model.Position)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPosition indicates an expected call of GetPosition.
func (mr *MockServicePositionMockRecorder) GetPosition(ctx, companyID, positionID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPosition", reflect.TypeOf((*MockServicePosition)(nil).GetPosition), ctx, companyID, positionID)
}

// GetPositions mocks base method.
func (m *MockServicePosition) GetPositions(ctx context.Context, id int) ([]*model.Position, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPositions", ctx, id)
	ret0, _ := ret[0].([]*model.Position)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPositions indicates an expected call of GetPositions.
func (mr *MockServicePositionMockRecorder) GetPositions(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPositions", reflect.TypeOf((*MockServicePosition)(nil).GetPositions), ctx, id)
}

// UpdatePosition mocks base method.
func (m *MockServicePosition) UpdatePosition(ctx context.Context, id int, position model.PositionUpdate) (*model.Position, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePosition", ctx, id, position)
	ret0, _ := ret[0].(*model.Position)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdatePosition indicates an expected call of UpdatePosition.
func (mr *MockServicePositionMockRecorder) UpdatePosition(ctx, id, position any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePosition", reflect.TypeOf((*MockServicePosition)(nil).UpdatePosition), ctx, id, position)
}
