// Code generated by MockGen. DO NOT EDIT.
// Source: internal/store/store.go
//
// Generated by this command:
//
//	mockgen -source=internal/store/store.go -destination=mocks/store/store.go
//
// Package mock_store is a generated GoMock package.
package mock_store

import (
	reflect "reflect"

	store "github.com/training-of-new-employees/qon/internal/store"
	gomock "go.uber.org/mock/gomock"
)

// MockStorages is a mock of Storages interface.
type MockStorages struct {
	ctrl     *gomock.Controller
	recorder *MockStoragesMockRecorder
}

// MockStoragesMockRecorder is the mock recorder for MockStorages.
type MockStoragesMockRecorder struct {
	mock *MockStorages
}

// NewMockStorages creates a new mock instance.
func NewMockStorages(ctrl *gomock.Controller) *MockStorages {
	mock := &MockStorages{ctrl: ctrl}
	mock.recorder = &MockStoragesMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorages) EXPECT() *MockStoragesMockRecorder {
	return m.recorder
}

// CompanyStorage mocks base method.
func (m *MockStorages) CompanyStorage() store.RepositoryCompany {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CompanyStorage")
	ret0, _ := ret[0].(store.RepositoryCompany)
	return ret0
}

// CompanyStorage indicates an expected call of CompanyStorage.
func (mr *MockStoragesMockRecorder) CompanyStorage() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CompanyStorage", reflect.TypeOf((*MockStorages)(nil).CompanyStorage))
}

// CourseStorage mocks base method.
func (m *MockStorages) CourseStorage() store.RepositoryCourse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CourseStorage")
	ret0, _ := ret[0].(store.RepositoryCourse)
	return ret0
}

// CourseStorage indicates an expected call of CourseStorage.
func (mr *MockStoragesMockRecorder) CourseStorage() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CourseStorage", reflect.TypeOf((*MockStorages)(nil).CourseStorage))
}

// PositionStorage mocks base method.
func (m *MockStorages) PositionStorage() store.RepositoryPosition {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PositionStorage")
	ret0, _ := ret[0].(store.RepositoryPosition)
	return ret0
}

// PositionStorage indicates an expected call of PositionStorage.
func (mr *MockStoragesMockRecorder) PositionStorage() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PositionStorage", reflect.TypeOf((*MockStorages)(nil).PositionStorage))
}

// UserStorage mocks base method.
func (m *MockStorages) UserStorage() store.RepositoryUser {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserStorage")
	ret0, _ := ret[0].(store.RepositoryUser)
	return ret0
}

// UserStorage indicates an expected call of UserStorage.
func (mr *MockStoragesMockRecorder) UserStorage() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserStorage", reflect.TypeOf((*MockStorages)(nil).UserStorage))
}
