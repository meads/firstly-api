// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/meads/firstly-api/security (interfaces: Claimer)

// Package security is a generated GoMock package.
package security

import (
	reflect "reflect"
	time "time"

	jwt_go "github.com/dgrijalva/jwt-go"
	gomock "github.com/golang/mock/gomock"
)

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

// GetClaimToken mocks base method.
func (m *MockClaimer) GetClaimToken() *ClaimToken {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetClaimToken")
	ret0, _ := ret[0].(*ClaimToken)
	return ret0
}

// GetClaimToken indicates an expected call of GetClaimToken.
func (mr *MockClaimerMockRecorder) GetClaimToken() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetClaimToken", reflect.TypeOf((*MockClaimer)(nil).GetClaimToken))
}

// GetFiveMinuteExpirationToken mocks base method.
func (m *MockClaimer) GetFiveMinuteExpirationToken(arg0 string) (string, time.Time, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFiveMinuteExpirationToken", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(time.Time)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetFiveMinuteExpirationToken indicates an expected call of GetFiveMinuteExpirationToken.
func (mr *MockClaimerMockRecorder) GetFiveMinuteExpirationToken(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFiveMinuteExpirationToken", reflect.TypeOf((*MockClaimer)(nil).GetFiveMinuteExpirationToken), arg0)
}

// GetFromTokenString mocks base method.
func (m *MockClaimer) GetFromTokenString(arg0 string) (*ClaimToken, *UsernameClaims, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFromTokenString", arg0)
	ret0, _ := ret[0].(*ClaimToken)
	ret1, _ := ret[1].(*UsernameClaims)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetFromTokenString indicates an expected call of GetFromTokenString.
func (mr *MockClaimerMockRecorder) GetFromTokenString(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFromTokenString", reflect.TypeOf((*MockClaimer)(nil).GetFromTokenString), arg0)
}

// ParseWithClaims mocks base method.
func (m *MockClaimer) ParseWithClaims(arg0 string, arg1 *UsernameClaims, arg2 jwt_go.Keyfunc) (*ClaimToken, error) {
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
