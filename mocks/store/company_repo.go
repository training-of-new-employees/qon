// Code generated by MockGen. DO NOT EDIT.
// Source: internal/store/company_repo.go

// Package mock_store is a generated GoMock package.
package mock_store

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/training-of-new-employees/qon/internal/model"
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

// CreateCompanyDB mocks base method.
func (m *MockRepositoryCompany) CreateCompanyDB(ctx context.Context, companyName string) (*model.Company, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCompanyDB", ctx, companyName)
	ret0, _ := ret[0].(*model.Company)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateCompanyDB indicates an expected call of CreateCompanyDB.
func (mr *MockRepositoryCompanyMockRecorder) CreateCompanyDB(ctx, companyName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCompanyDB", reflect.TypeOf((*MockRepositoryCompany)(nil).CreateCompanyDB), ctx, companyName)
}