// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/meads/firstly-api/security (interfaces: MethodSigner,ClaimValidator,Claimer)

// Package security is a generated GoMock package.
package security

import (
	reflect "reflect"

	jwt_go "github.com/dgrijalva/jwt-go"
	gomock "github.com/golang/mock/gomock"
)

// MockMethodSigner is a mock of MethodSigner interface.
type MockMethodSigner struct {
	ctrl     *gomock.Controller
	recorder *MockMethodSignerMockRecorder
}

// MockMethodSignerMockRecorder is the mock recorder for MockMethodSigner.
type MockMethodSignerMockRecorder struct {
	mock *MockMethodSigner
}

// NewMockMethodSigner creates a new mock instance.
func NewMockMethodSigner(ctrl *gomock.Controller) *MockMethodSigner {
	mock := &MockMethodSigner{ctrl: ctrl}
	mock.recorder = &MockMethodSignerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMethodSigner) EXPECT() *MockMethodSignerMockRecorder {
	return m.recorder
}

// Alg mocks base method.
func (m *MockMethodSigner) Alg() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Alg")
	ret0, _ := ret[0].(string)
	return ret0
}

// Alg indicates an expected call of Alg.
func (mr *MockMethodSignerMockRecorder) Alg() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Alg", reflect.TypeOf((*MockMethodSigner)(nil).Alg))
}

// Sign mocks base method.
func (m *MockMethodSigner) Sign(arg0 string, arg1 interface{}) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Sign", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Sign indicates an expected call of Sign.
func (mr *MockMethodSignerMockRecorder) Sign(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Sign", reflect.TypeOf((*MockMethodSigner)(nil).Sign), arg0, arg1)
}

// Verify mocks base method.
func (m *MockMethodSigner) Verify(arg0, arg1 string, arg2 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Verify", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Verify indicates an expected call of Verify.
func (mr *MockMethodSignerMockRecorder) Verify(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Verify", reflect.TypeOf((*MockMethodSigner)(nil).Verify), arg0, arg1, arg2)
}

// MockClaimValidator is a mock of ClaimValidator interface.
type MockClaimValidator struct {
	ctrl     *gomock.Controller
	recorder *MockClaimValidatorMockRecorder
}

// MockClaimValidatorMockRecorder is the mock recorder for MockClaimValidator.
type MockClaimValidatorMockRecorder struct {
	mock *MockClaimValidator
}

// NewMockClaimValidator creates a new mock instance.
func NewMockClaimValidator(ctrl *gomock.Controller) *MockClaimValidator {
	mock := &MockClaimValidator{ctrl: ctrl}
	mock.recorder = &MockClaimValidatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClaimValidator) EXPECT() *MockClaimValidatorMockRecorder {
	return m.recorder
}

// Valid mocks base method.
func (m *MockClaimValidator) Valid() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Valid")
	ret0, _ := ret[0].(error)
	return ret0
}

// Valid indicates an expected call of Valid.
func (mr *MockClaimValidatorMockRecorder) Valid() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Valid", reflect.TypeOf((*MockClaimValidator)(nil).Valid))
}

// MockClaimer is a mock of Claimer interface.
type MockClaimer struct {
	ctrl     *gomock.Controller
	recorder *MockClaimerMockRecorder
}

// MockClaimerMockRecorder is the mock recorder for MockClaimer.
type MockClaimerMockRecorder struct {
	mock *MockClaimer
}

// NewMockClaimer creates a new mock instance.
func NewMockClaimer(ctrl *gomock.Controller) *MockClaimer {
	mock := &MockClaimer{ctrl: ctrl}
	mock.recorder = &MockClaimerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClaimer) EXPECT() *MockClaimerMockRecorder {
	return m.recorder
}

// NewWithClaims mocks base method.
func (m *MockClaimer) NewWithClaims(arg0 MethodSigner, arg1 ClaimValidator) *ClaimToken {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewWithClaims", arg0, arg1)
	ret0, _ := ret[0].(*ClaimToken)
	return ret0
}

// NewWithClaims indicates an expected call of NewWithClaims.
func (mr *MockClaimerMockRecorder) NewWithClaims(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewWithClaims", reflect.TypeOf((*MockClaimer)(nil).NewWithClaims), arg0, arg1)
}

// ParseWithClaims mocks base method.
func (m *MockClaimer) ParseWithClaims(arg0 string, arg1 ClaimsValidator, arg2 jwt_go.Keyfunc) (*ClaimToken, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseWithClaims", arg0, arg1, arg2)
	ret0, _ := ret[0].(*ClaimToken)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ParseWithClaims indicates an expected call of ParseWithClaims.
func (mr *MockClaimerMockRecorder) ParseWithClaims(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseWithClaims", reflect.TypeOf((*MockClaimer)(nil).ParseWithClaims), arg0, arg1, arg2)
}

// SignedString mocks base method.
func (m *MockClaimer) SignedString(arg0 interface{}) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignedString", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignedString indicates an expected call of SignedString.
func (mr *MockClaimerMockRecorder) SignedString(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignedString", reflect.TypeOf((*MockClaimer)(nil).SignedString), arg0)
}
