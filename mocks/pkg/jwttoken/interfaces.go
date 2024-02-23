// Code generated by MockGen. DO NOT EDIT.
// Source: internal/pkg/jwttoken/interfaces.go
//
// Generated by this command:
//
//	mockgen -source=internal/pkg/jwttoken/interfaces.go -destination=mocks/pkg/jwttoken/interfaces.go
//
// Package mock_jwttoken is a generated GoMock package.
package mock_jwttoken

import (
	reflect "reflect"
	time "time"

	jwttoken "github.com/training-of-new-employees/qon/internal/pkg/jwttoken"
	gomock "go.uber.org/mock/gomock"
)

// MockJWTGenerator is a mock of JWTGenerator interface.
type MockJWTGenerator struct {
	ctrl     *gomock.Controller
	recorder *MockJWTGeneratorMockRecorder
}

// MockJWTGeneratorMockRecorder is the mock recorder for MockJWTGenerator.
type MockJWTGeneratorMockRecorder struct {
	mock *MockJWTGenerator
}

// NewMockJWTGenerator creates a new mock instance.
func NewMockJWTGenerator(ctrl *gomock.Controller) *MockJWTGenerator {
	mock := &MockJWTGenerator{ctrl: ctrl}
	mock.recorder = &MockJWTGeneratorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockJWTGenerator) EXPECT() *MockJWTGeneratorMockRecorder {
	return m.recorder
}

// GenerateToken mocks base method.
func (m *MockJWTGenerator) GenerateToken(id int, isAdmin bool, orgID int, hashedRefreshToken string, exp time.Duration) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateToken", id, isAdmin, orgID, hashedRefreshToken, exp)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateToken indicates an expected call of GenerateToken.
func (mr *MockJWTGeneratorMockRecorder) GenerateToken(id, isAdmin, orgID, hashedRefreshToken, exp any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateToken", reflect.TypeOf((*MockJWTGenerator)(nil).GenerateToken), id, isAdmin, orgID, hashedRefreshToken, exp)
}

// MockJWTValidator is a mock of JWTValidator interface.
type MockJWTValidator struct {
	ctrl     *gomock.Controller
	recorder *MockJWTValidatorMockRecorder
}

// MockJWTValidatorMockRecorder is the mock recorder for MockJWTValidator.
type MockJWTValidatorMockRecorder struct {
	mock *MockJWTValidator
}

// NewMockJWTValidator creates a new mock instance.
func NewMockJWTValidator(ctrl *gomock.Controller) *MockJWTValidator {
	mock := &MockJWTValidator{ctrl: ctrl}
	mock.recorder = &MockJWTValidatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockJWTValidator) EXPECT() *MockJWTValidatorMockRecorder {
	return m.recorder
}

// ValidateToken mocks base method.
func (m *MockJWTValidator) ValidateToken(tokenStr string) (*jwttoken.MyClaims, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateToken", tokenStr)
	ret0, _ := ret[0].(*jwttoken.MyClaims)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidateToken indicates an expected call of ValidateToken.
func (mr *MockJWTValidatorMockRecorder) ValidateToken(tokenStr any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateToken", reflect.TypeOf((*MockJWTValidator)(nil).ValidateToken), tokenStr)
}
