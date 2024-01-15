// Code generated by MockGen. DO NOT EDIT.
// Source: internal/service/user_service.go
//
// Generated by this command:
//
//	mockgen -source=internal/service/user_service.go -destination=mocks/service/user_service.go
//

// Package mock_service is a generated GoMock package.
package mock_service

import (
	context "context"
	reflect "reflect"

	model "github.com/training-of-new-employees/qon/internal/model"
	gomock "go.uber.org/mock/gomock"
)

// MockServiceUser is a mock of ServiceUser interface.
type MockServiceUser struct {
	ctrl     *gomock.Controller
	recorder *MockServiceUserMockRecorder
}

// MockServiceUserMockRecorder is the mock recorder for MockServiceUser.
type MockServiceUserMockRecorder struct {
	mock *MockServiceUser
}

// NewMockServiceUser creates a new mock instance.
func NewMockServiceUser(ctrl *gomock.Controller) *MockServiceUser {
	mock := &MockServiceUser{ctrl: ctrl}
	mock.recorder = &MockServiceUserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockServiceUser) EXPECT() *MockServiceUserMockRecorder {
	return m.recorder
}

// ArchiveUser mocks base method.
func (m *MockServiceUser) ArchiveUser(ctx context.Context, id, editorCompanyID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ArchiveUser", ctx, id, editorCompanyID)
	ret0, _ := ret[0].(error)
	return ret0
}

// ArchiveUser indicates an expected call of ArchiveUser.
func (mr *MockServiceUserMockRecorder) ArchiveUser(ctx, id, editorCompanyID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ArchiveUser", reflect.TypeOf((*MockServiceUser)(nil).ArchiveUser), ctx, id, editorCompanyID)
}

// CreateAdmin mocks base method.
func (m *MockServiceUser) CreateAdmin(ctx context.Context, val model.CreateAdmin) (*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAdmin", ctx, val)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateAdmin indicates an expected call of CreateAdmin.
func (mr *MockServiceUserMockRecorder) CreateAdmin(ctx, val interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAdmin", reflect.TypeOf((*MockServiceUser)(nil).CreateAdmin), ctx, val)
}

// CreateUser mocks base method.
func (m *MockServiceUser) CreateUser(ctx context.Context, user model.UserCreate) (*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", ctx, user)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockServiceUserMockRecorder) CreateUser(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockServiceUser)(nil).CreateUser), ctx, user)
}

// DeleteAdminFromCache mocks base method.
func (m *MockServiceUser) DeleteAdminFromCache(ctx context.Context, key string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAdminFromCache", ctx, key)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAdminFromCache indicates an expected call of DeleteAdminFromCache.
func (mr *MockServiceUserMockRecorder) DeleteAdminFromCache(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAdminFromCache", reflect.TypeOf((*MockServiceUser)(nil).DeleteAdminFromCache), ctx, key)
}

// EditAdmin mocks base method.
func (m *MockServiceUser) EditAdmin(ctx context.Context, val model.AdminEdit) (*model.AdminEdit, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EditAdmin", ctx, val)
	ret0, _ := ret[0].(*model.AdminEdit)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EditAdmin indicates an expected call of EditAdmin.
func (mr *MockServiceUserMockRecorder) EditAdmin(ctx, val interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EditAdmin", reflect.TypeOf((*MockServiceUser)(nil).EditAdmin), ctx, val)
}

// EditUser mocks base method.
func (m *MockServiceUser) EditUser(ctx context.Context, val *model.UserEdit, editorCompanyID int) (*model.UserEdit, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EditUser", ctx, val, editorCompanyID)
	ret0, _ := ret[0].(*model.UserEdit)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EditUser indicates an expected call of EditUser.
func (mr *MockServiceUserMockRecorder) EditUser(ctx, val, editorCompanyID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EditUser", reflect.TypeOf((*MockServiceUser)(nil).EditUser), ctx, val, editorCompanyID)
}

// GenerateInvitationLinkUser mocks base method.
func (m *MockServiceUser) GenerateInvitationLinkUser(ctx context.Context, email string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateInvitationLinkUser", ctx, email)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateInvitationLinkUser indicates an expected call of GenerateInvitationLinkUser.
func (mr *MockServiceUserMockRecorder) GenerateInvitationLinkUser(ctx, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateInvitationLinkUser", reflect.TypeOf((*MockServiceUser)(nil).GenerateInvitationLinkUser), ctx, email)
}

// GenerateTokenPair mocks base method.
func (m *MockServiceUser) GenerateTokenPair(ctx context.Context, userId int, isAdmin bool, companyID int) (*model.Tokens, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateTokenPair", ctx, userId, isAdmin, companyID)
	ret0, _ := ret[0].(*model.Tokens)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateTokenPair indicates an expected call of GenerateTokenPair.
func (mr *MockServiceUserMockRecorder) GenerateTokenPair(ctx, userId, isAdmin, companyID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateTokenPair", reflect.TypeOf((*MockServiceUser)(nil).GenerateTokenPair), ctx, userId, isAdmin, companyID)
}

// GetAdminFromCache mocks base method.
func (m *MockServiceUser) GetAdminFromCache(arg0 context.Context, arg1 string) (*model.CreateAdmin, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAdminFromCache", arg0, arg1)
	ret0, _ := ret[0].(*model.CreateAdmin)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAdminFromCache indicates an expected call of GetAdminFromCache.
func (mr *MockServiceUserMockRecorder) GetAdminFromCache(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAdminFromCache", reflect.TypeOf((*MockServiceUser)(nil).GetAdminFromCache), arg0, arg1)
}

// GetUserByEmail mocks base method.
func (m *MockServiceUser) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByEmail", ctx, email)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByEmail indicates an expected call of GetUserByEmail.
func (mr *MockServiceUserMockRecorder) GetUserByEmail(ctx, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByEmail", reflect.TypeOf((*MockServiceUser)(nil).GetUserByEmail), ctx, email)
}

// GetUserByID mocks base method.
func (m *MockServiceUser) GetUserByID(ctx context.Context, id int) (*model.UserInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByID", ctx, id)
	ret0, _ := ret[0].(*model.UserInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByID indicates an expected call of GetUserByID.
func (mr *MockServiceUserMockRecorder) GetUserByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByID", reflect.TypeOf((*MockServiceUser)(nil).GetUserByID), ctx, id)
}

// GetUserInviteCodeFromCache mocks base method.
func (m *MockServiceUser) GetUserInviteCodeFromCache(ctx context.Context, email string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserInviteCodeFromCache", ctx, email)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserInviteCodeFromCache indicates an expected call of GetUserInviteCodeFromCache.
func (mr *MockServiceUserMockRecorder) GetUserInviteCodeFromCache(ctx, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserInviteCodeFromCache", reflect.TypeOf((*MockServiceUser)(nil).GetUserInviteCodeFromCache), ctx, email)
}

// GetUsersByCompany mocks base method.
func (m *MockServiceUser) GetUsersByCompany(ctx context.Context, companyID int) ([]model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUsersByCompany", ctx, companyID)
	ret0, _ := ret[0].([]model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUsersByCompany indicates an expected call of GetUsersByCompany.
func (mr *MockServiceUserMockRecorder) GetUsersByCompany(ctx, companyID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsersByCompany", reflect.TypeOf((*MockServiceUser)(nil).GetUsersByCompany), ctx, companyID)
}

// RegenerationInvitationLinkUser mocks base method.
func (m *MockServiceUser) RegenerationInvitationLinkUser(ctx context.Context, email string, orgID int) (*model.InvitationLinkResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegenerationInvitationLinkUser", ctx, email, orgID)
	ret0, _ := ret[0].(*model.InvitationLinkResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RegenerationInvitationLinkUser indicates an expected call of RegenerationInvitationLinkUser.
func (mr *MockServiceUserMockRecorder) RegenerationInvitationLinkUser(ctx, email, orgID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegenerationInvitationLinkUser", reflect.TypeOf((*MockServiceUser)(nil).RegenerationInvitationLinkUser), ctx, email, orgID)
}

// ResetPassword mocks base method.
func (m *MockServiceUser) ResetPassword(ctx context.Context, email string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResetPassword", ctx, email)
	ret0, _ := ret[0].(error)
	return ret0
}

// ResetPassword indicates an expected call of ResetPassword.
func (mr *MockServiceUserMockRecorder) ResetPassword(ctx, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResetPassword", reflect.TypeOf((*MockServiceUser)(nil).ResetPassword), ctx, email)
}

// UpdatePasswordAndActivateUser mocks base method.
func (m *MockServiceUser) UpdatePasswordAndActivateUser(ctx context.Context, email, password string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePasswordAndActivateUser", ctx, email, password)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdatePasswordAndActivateUser indicates an expected call of UpdatePasswordAndActivateUser.
func (mr *MockServiceUserMockRecorder) UpdatePasswordAndActivateUser(ctx, email, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePasswordAndActivateUser", reflect.TypeOf((*MockServiceUser)(nil).UpdatePasswordAndActivateUser), ctx, email, password)
}

// WriteAdminToCache mocks base method.
func (m *MockServiceUser) WriteAdminToCache(ctx context.Context, admin model.CreateAdmin) (*model.CreateAdmin, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteAdminToCache", ctx, admin)
	ret0, _ := ret[0].(*model.CreateAdmin)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WriteAdminToCache indicates an expected call of WriteAdminToCache.
func (mr *MockServiceUserMockRecorder) WriteAdminToCache(ctx, admin interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteAdminToCache", reflect.TypeOf((*MockServiceUser)(nil).WriteAdminToCache), ctx, admin)
}
