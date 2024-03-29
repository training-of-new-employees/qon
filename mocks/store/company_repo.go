// Code generated by MockGen. DO NOT EDIT.
// Source: internal/store/company_repo.go
//
// Generated by this command:
//
//	mockgen -source=internal/store/company_repo.go -destination=mocks/store/company_repo.go
//

// Package mock_store is a generated GoMock package.
package mock_store

import (
	context "context"
	reflect "reflect"

	model "github.com/training-of-new-employees/qon/internal/model"
	gomock "go.uber.org/mock/gomock"
)

// MockRepositoryCompany is a mock of RepositoryCompany interface.
type MockRepositoryCompany struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryCompanyMockRecorder
}

// MockRepositoryCompanyMockRecorder is the mock recorder for MockRepositoryCompany.
type MockRepositoryCompanyMockRecorder struct {
	mock *MockRepositoryCompany
}

// NewMockRepositoryCompany creates a new mock instance.
func NewMockRepositoryCompany(ctrl *gomock.Controller) *MockRepositoryCompany {
	mock := &MockRepositoryCompany{ctrl: ctrl}
	mock.recorder = &MockRepositoryCompanyMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepositoryCompany) EXPECT() *MockRepositoryCompanyMockRecorder {
	return m.recorder
}

// CreateCompany mocks base method.
func (m *MockRepositoryCompany) CreateCompany(ctx context.Context, companyName string) (*model.Company, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCompany", ctx, companyName)
	ret0, _ := ret[0].(*model.Company)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateCompany indicates an expected call of CreateCompany.
func (mr *MockRepositoryCompanyMockRecorder) CreateCompany(ctx, companyName any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCompany", reflect.TypeOf((*MockRepositoryCompany)(nil).CreateCompany), ctx, companyName)
}

// GetCompany mocks base method.
func (m *MockRepositoryCompany) GetCompany(ctx context.Context, id int) (*model.Company, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCompany", ctx, id)
	ret0, _ := ret[0].(*model.Company)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCompany indicates an expected call of GetCompany.
func (mr *MockRepositoryCompanyMockRecorder) GetCompany(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCompany", reflect.TypeOf((*MockRepositoryCompany)(nil).GetCompany), ctx, id)
}
