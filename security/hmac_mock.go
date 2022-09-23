// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/meads/firstly-api/security (interfaces: Hasher)

// Package security is a generated GoMock package.
package security

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockHasher is a mock of Hasher interface.
type MockHasher struct {
	ctrl     *gomock.Controller
	recorder *MockHasherMockRecorder
}

// MockHasherMockRecorder is the mock recorder for MockHasher.
type MockHasherMockRecorder struct {
	mock *MockHasher
}

// NewMockHasher creates a new mock instance.
func NewMockHasher(ctrl *gomock.Controller) *MockHasher {
	mock := &MockHasher{ctrl: ctrl}
	mock.recorder = &MockHasherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHasher) EXPECT() *MockHasherMockRecorder {
	return m.recorder
}

// GeneratePasswordHash mocks base method.
func (m *MockHasher) GeneratePasswordHash(arg0 []byte, arg1 string) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GeneratePasswordHash", arg0, arg1)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GeneratePasswordHash indicates an expected call of GeneratePasswordHash.
func (mr *MockHasherMockRecorder) GeneratePasswordHash(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GeneratePasswordHash", reflect.TypeOf((*MockHasher)(nil).GeneratePasswordHash), arg0, arg1)
}

// GenerateSalt mocks base method.
func (m *MockHasher) GenerateSalt() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateSalt")
	ret0, _ := ret[0].(string)
	return ret0
}

// GenerateSalt indicates an expected call of GenerateSalt.
func (mr *MockHasherMockRecorder) GenerateSalt() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateSalt", reflect.TypeOf((*MockHasher)(nil).GenerateSalt))
}

// IsValidPasswordHash mocks base method.
func (m *MockHasher) IsValidPasswordHash(arg0, arg1, arg2 []byte) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsValidPasswordHash", arg0, arg1, arg2)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsValidPasswordHash indicates an expected call of IsValidPasswordHash.
func (mr *MockHasherMockRecorder) IsValidPasswordHash(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsValidPasswordHash", reflect.TypeOf((*MockHasher)(nil).IsValidPasswordHash), arg0, arg1, arg2)
}
