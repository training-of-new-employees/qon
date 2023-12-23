// Code generated by MockGen. DO NOT EDIT.
// Source: internal/pkg/doar/doar.go
//
// Generated by this command:
//
//	mockgen -source=internal/pkg/doar/doar.go -destination=mocks//pkg/doar/doar.go
//
// Package mock_doar is a generated GoMock package.
package mock_doar

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockEmailSender is a mock of EmailSender interface.
type MockEmailSender struct {
	ctrl     *gomock.Controller
	recorder *MockEmailSenderMockRecorder
}

// MockEmailSenderMockRecorder is the mock recorder for MockEmailSender.
type MockEmailSenderMockRecorder struct {
	mock *MockEmailSender
}

// NewMockEmailSender creates a new mock instance.
func NewMockEmailSender(ctrl *gomock.Controller) *MockEmailSender {
	mock := &MockEmailSender{ctrl: ctrl}
	mock.recorder = &MockEmailSenderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEmailSender) EXPECT() *MockEmailSenderMockRecorder {
	return m.recorder
}

// InviteUser mocks base method.
func (m *MockEmailSender) InviteUser(email, invitationLink string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InviteUser", email, invitationLink)
	ret0, _ := ret[0].(error)
	return ret0
}

// InviteUser indicates an expected call of InviteUser.
func (mr *MockEmailSenderMockRecorder) InviteUser(email, invitationLink any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InviteUser", reflect.TypeOf((*MockEmailSender)(nil).InviteUser), email, invitationLink)
}

// SendCode mocks base method.
func (m *MockEmailSender) SendCode(email, code string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendCode", email, code)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendCode indicates an expected call of SendCode.
func (mr *MockEmailSenderMockRecorder) SendCode(email, code any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendCode", reflect.TypeOf((*MockEmailSender)(nil).SendCode), email, code)
}

// SendPassword mocks base method.
func (m *MockEmailSender) SendPassword(email, password string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendPassword", email, password)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendPassword indicates an expected call of SendPassword.
func (mr *MockEmailSenderMockRecorder) SendPassword(email, password any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendPassword", reflect.TypeOf((*MockEmailSender)(nil).SendPassword), email, password)
}

// MockMailer is a mock of Mailer interface.
type MockMailer struct {
	ctrl     *gomock.Controller
	recorder *MockMailerMockRecorder
}

// MockMailerMockRecorder is the mock recorder for MockMailer.
type MockMailerMockRecorder struct {
	mock *MockMailer
}

// NewMockMailer creates a new mock instance.
func NewMockMailer(ctrl *gomock.Controller) *MockMailer {
	mock := &MockMailer{ctrl: ctrl}
	mock.recorder = &MockMailerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMailer) EXPECT() *MockMailerMockRecorder {
	return m.recorder
}

// SendEmail mocks base method.
func (m *MockMailer) SendEmail(email, title, body string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendEmail", email, title, body)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendEmail indicates an expected call of SendEmail.
func (mr *MockMailerMockRecorder) SendEmail(email, title, body any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendEmail", reflect.TypeOf((*MockMailer)(nil).SendEmail), email, title, body)
}