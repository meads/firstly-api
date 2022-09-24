// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/meads/firstly-api/db (interfaces: Store)

// Package db is a generated GoMock package.
package db

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockStore is a mock of Store interface.
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore.
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance.
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// AccountExists mocks base method.
func (m *MockStore) AccountExists(arg0 context.Context, arg1 int64) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AccountExists", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AccountExists indicates an expected call of AccountExists.
func (mr *MockStoreMockRecorder) AccountExists(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AccountExists", reflect.TypeOf((*MockStore)(nil).AccountExists), arg0, arg1)
}

// CreateAccount mocks base method.
func (m *MockStore) CreateAccount(arg0 context.Context, arg1 CreateAccountParams) (Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAccount", arg0, arg1)
	ret0, _ := ret[0].(Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateAccount indicates an expected call of CreateAccount.
func (mr *MockStoreMockRecorder) CreateAccount(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAccount", reflect.TypeOf((*MockStore)(nil).CreateAccount), arg0, arg1)
}

// CreateImage mocks base method.
func (m *MockStore) CreateImage(arg0 context.Context, arg1 string) (Image, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateImage", arg0, arg1)
	ret0, _ := ret[0].(Image)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateImage indicates an expected call of CreateImage.
func (mr *MockStoreMockRecorder) CreateImage(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateImage", reflect.TypeOf((*MockStore)(nil).CreateImage), arg0, arg1)
}

// DeleteAccount mocks base method.
func (m *MockStore) DeleteAccount(arg0 context.Context, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAccount", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAccount indicates an expected call of DeleteAccount.
func (mr *MockStoreMockRecorder) DeleteAccount(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAccount", reflect.TypeOf((*MockStore)(nil).DeleteAccount), arg0, arg1)
}

// DeleteImage mocks base method.
func (m *MockStore) DeleteImage(arg0 context.Context, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteImage", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteImage indicates an expected call of DeleteImage.
func (mr *MockStoreMockRecorder) DeleteImage(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteImage", reflect.TypeOf((*MockStore)(nil).DeleteImage), arg0, arg1)
}

// GetAccount mocks base method.
func (m *MockStore) GetAccount(arg0 context.Context, arg1 int64) (Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccount", arg0, arg1)
	ret0, _ := ret[0].(Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccount indicates an expected call of GetAccount.
func (mr *MockStoreMockRecorder) GetAccount(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccount", reflect.TypeOf((*MockStore)(nil).GetAccount), arg0, arg1)
}

// GetAccountByUsername mocks base method.
func (m *MockStore) GetAccountByUsername(arg0 context.Context, arg1 string) (Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccountByUsername", arg0, arg1)
	ret0, _ := ret[0].(Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccountByUsername indicates an expected call of GetAccountByUsername.
func (mr *MockStoreMockRecorder) GetAccountByUsername(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccountByUsername", reflect.TypeOf((*MockStore)(nil).GetAccountByUsername), arg0, arg1)
}

// GetImage mocks base method.
func (m *MockStore) GetImage(arg0 context.Context, arg1 int64) (Image, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetImage", arg0, arg1)
	ret0, _ := ret[0].(Image)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetImage indicates an expected call of GetImage.
func (mr *MockStoreMockRecorder) GetImage(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetImage", reflect.TypeOf((*MockStore)(nil).GetImage), arg0, arg1)
}

// ListAccounts mocks base method.
func (m *MockStore) ListAccounts(arg0 context.Context, arg1 ListAccountsParams) ([]Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListAccounts", arg0, arg1)
	ret0, _ := ret[0].([]Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListAccounts indicates an expected call of ListAccounts.
func (mr *MockStoreMockRecorder) ListAccounts(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListAccounts", reflect.TypeOf((*MockStore)(nil).ListAccounts), arg0, arg1)
}

// ListImages mocks base method.
func (m *MockStore) ListImages(arg0 context.Context, arg1 ListImagesParams) ([]Image, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListImages", arg0, arg1)
	ret0, _ := ret[0].([]Image)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListImages indicates an expected call of ListImages.
func (mr *MockStoreMockRecorder) ListImages(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListImages", reflect.TypeOf((*MockStore)(nil).ListImages), arg0, arg1)
}

// SoftDeleteAccount mocks base method.
func (m *MockStore) SoftDeleteAccount(arg0 context.Context, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SoftDeleteAccount", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SoftDeleteAccount indicates an expected call of SoftDeleteAccount.
func (mr *MockStoreMockRecorder) SoftDeleteAccount(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SoftDeleteAccount", reflect.TypeOf((*MockStore)(nil).SoftDeleteAccount), arg0, arg1)
}

// SoftDeleteImage mocks base method.
func (m *MockStore) SoftDeleteImage(arg0 context.Context, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SoftDeleteImage", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SoftDeleteImage indicates an expected call of SoftDeleteImage.
func (mr *MockStoreMockRecorder) SoftDeleteImage(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SoftDeleteImage", reflect.TypeOf((*MockStore)(nil).SoftDeleteImage), arg0, arg1)
}

// Tx mocks base method.
func (m *MockStore) Tx(arg0 context.Context, arg1 func(*Queries, *interface{}) (interface{}, error)) (interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Tx", arg0, arg1)
	ret0, _ := ret[0].(interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Tx indicates an expected call of Tx.
func (mr *MockStoreMockRecorder) Tx(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Tx", reflect.TypeOf((*MockStore)(nil).Tx), arg0, arg1)
}

// UpdateAccount mocks base method.
func (m *MockStore) UpdateAccount(arg0 context.Context, arg1 UpdateAccountParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateAccount", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateAccount indicates an expected call of UpdateAccount.
func (mr *MockStoreMockRecorder) UpdateAccount(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAccount", reflect.TypeOf((*MockStore)(nil).UpdateAccount), arg0, arg1)
}

// UpdateImage mocks base method.
func (m *MockStore) UpdateImage(arg0 context.Context, arg1 UpdateImageParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateImage", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateImage indicates an expected call of UpdateImage.
func (mr *MockStoreMockRecorder) UpdateImage(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateImage", reflect.TypeOf((*MockStore)(nil).UpdateImage), arg0, arg1)
}
