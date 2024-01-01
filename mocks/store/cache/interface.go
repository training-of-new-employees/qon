// Code generated by MockGen. DO NOT EDIT.
// Source: internal/store/cache/interface.go
//
// Generated by this command:
//
//	mockgen -source=internal/store/cache/interface.go -destination=mocks/store/cache/interface.go
//

// Package mock_cache is a generated GoMock package.
package mock_cache

import (
	context "context"
	reflect "reflect"

	model "github.com/training-of-new-employees/qon/internal/model"
	gomock "go.uber.org/mock/gomock"
)

// MockCache is a mock of Cache interface.
type MockCache struct {
	ctrl     *gomock.Controller
	recorder *MockCacheMockRecorder
}

// MockCacheMockRecorder is the mock recorder for MockCache.
type MockCacheMockRecorder struct {
	mock *MockCache
}

// NewMockCache creates a new mock instance.
func NewMockCache(ctrl *gomock.Controller) *MockCache {
	mock := &MockCache{ctrl: ctrl}
	mock.recorder = &MockCacheMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCache) EXPECT() *MockCacheMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockCache) Delete(ctx context.Context, key string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, key)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockCacheMockRecorder) Delete(ctx, key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockCache)(nil).Delete), ctx, key)
}

// Get mocks base method.
func (m *MockCache) Get(ctx context.Context, key string) (*model.CreateAdmin, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, key)
	ret0, _ := ret[0].(*model.CreateAdmin)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockCacheMockRecorder) Get(ctx, key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockCache)(nil).Get), ctx, key)
}

// Set mocks base method.
func (m *MockCache) Set(ctx context.Context, uuid string, admin model.CreateAdmin) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", ctx, uuid, admin)
	ret0, _ := ret[0].(error)
	return ret0
}

// Set indicates an expected call of Set.
func (mr *MockCacheMockRecorder) Set(ctx, uuid, admin any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockCache)(nil).Set), ctx, uuid, admin)
}
