// Code generated by MockGen. DO NOT EDIT.
// Source: internal/store/position_repo.go

// Package mock_store is a generated GoMock package.
package mock_store

import (
	context "context"
	reflect "reflect"

	model "github.com/training-of-new-employees/qon/internal/model"
	gomock "go.uber.org/mock/gomock"
)

// MockRepositoryPosition is a mock of RepositoryPosition interface.
type MockRepositoryPosition struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryPositionMockRecorder
}

// MockRepositoryPositionMockRecorder is the mock recorder for MockRepositoryPosition.
type MockRepositoryPositionMockRecorder struct {
	mock *MockRepositoryPosition
}

// NewMockRepositoryPosition creates a new mock instance.
func NewMockRepositoryPosition(ctrl *gomock.Controller) *MockRepositoryPosition {
	mock := &MockRepositoryPosition{ctrl: ctrl}
	mock.recorder = &MockRepositoryPositionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepositoryPosition) EXPECT() *MockRepositoryPositionMockRecorder {
	return m.recorder
}

// AssignCourseDB mocks base method.
func (m *MockRepositoryPosition) AssignCourseDB(ctx context.Context, positionID, courseID, user_id int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AssignCourseDB", ctx, positionID, courseID, user_id)
	ret0, _ := ret[0].(error)
	return ret0
}

// AssignCourseDB indicates an expected call of AssignCourseDB.
func (mr *MockRepositoryPositionMockRecorder) AssignCourseDB(ctx, positionID, courseID, user_id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AssignCourseDB", reflect.TypeOf((*MockRepositoryPosition)(nil).AssignCourseDB), ctx, positionID, courseID, user_id)
}

// CreatePositionDB mocks base method.
func (m *MockRepositoryPosition) CreatePositionDB(ctx context.Context, position model.PositionSet) (*model.Position, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePositionDB", ctx, position)
	ret0, _ := ret[0].(*model.Position)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreatePositionDB indicates an expected call of CreatePositionDB.
func (mr *MockRepositoryPositionMockRecorder) CreatePositionDB(ctx, position interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePositionDB", reflect.TypeOf((*MockRepositoryPosition)(nil).CreatePositionDB), ctx, position)
}

// GetPositionByID mocks base method.
func (m *MockRepositoryPosition) GetPositionByID(ctx context.Context, positionID int) (*model.Position, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPositionByID", ctx, positionID)
	ret0, _ := ret[0].(*model.Position)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPositionByID indicates an expected call of GetPositionByID.
func (mr *MockRepositoryPositionMockRecorder) GetPositionByID(ctx, positionID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPositionByID", reflect.TypeOf((*MockRepositoryPosition)(nil).GetPositionByID), ctx, positionID)
}

// GetPositionDB mocks base method.
func (m *MockRepositoryPosition) GetPositionDB(ctx context.Context, companyID, positionID int) (*model.Position, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPositionDB", ctx, companyID, positionID)
	ret0, _ := ret[0].(*model.Position)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPositionDB indicates an expected call of GetPositionDB.
func (mr *MockRepositoryPositionMockRecorder) GetPositionDB(ctx, companyID, positionID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPositionDB", reflect.TypeOf((*MockRepositoryPosition)(nil).GetPositionDB), ctx, companyID, positionID)
}

// GetPositionsDB mocks base method.
func (m *MockRepositoryPosition) GetPositionsDB(ctx context.Context, companyID int) ([]*model.Position, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPositionsDB", ctx, companyID)
	ret0, _ := ret[0].([]*model.Position)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPositionsDB indicates an expected call of GetPositionsDB.
func (mr *MockRepositoryPositionMockRecorder) GetPositionsDB(ctx, companyID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPositionsDB", reflect.TypeOf((*MockRepositoryPosition)(nil).GetPositionsDB), ctx, companyID)
}

// UpdatePositionDB mocks base method.
func (m *MockRepositoryPosition) UpdatePositionDB(ctx context.Context, id int, position model.PositionSet) (*model.Position, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePositionDB", ctx, id, position)
	ret0, _ := ret[0].(*model.Position)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdatePositionDB indicates an expected call of UpdatePositionDB.
func (mr *MockRepositoryPositionMockRecorder) UpdatePositionDB(ctx, id, position interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePositionDB", reflect.TypeOf((*MockRepositoryPosition)(nil).UpdatePositionDB), ctx, id, position)
}
