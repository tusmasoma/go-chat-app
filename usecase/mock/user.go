// Code generated by MockGen. DO NOT EDIT.
// Source: user.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockUserUseCase is a mock of UserUseCase interface.
type MockUserUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockUserUseCaseMockRecorder
}

// MockUserUseCaseMockRecorder is the mock recorder for MockUserUseCase.
type MockUserUseCaseMockRecorder struct {
	mock *MockUserUseCase
}

// NewMockUserUseCase creates a new mock instance.
func NewMockUserUseCase(ctrl *gomock.Controller) *MockUserUseCase {
	mock := &MockUserUseCase{ctrl: ctrl}
	mock.recorder = &MockUserUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserUseCase) EXPECT() *MockUserUseCaseMockRecorder {
	return m.recorder
}

// LoginAndGenerateToken mocks base method.
func (m *MockUserUseCase) LoginAndGenerateToken(ctx context.Context, email, password string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoginAndGenerateToken", ctx, email, password)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoginAndGenerateToken indicates an expected call of LoginAndGenerateToken.
func (mr *MockUserUseCaseMockRecorder) LoginAndGenerateToken(ctx, email, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoginAndGenerateToken", reflect.TypeOf((*MockUserUseCase)(nil).LoginAndGenerateToken), ctx, email, password)
}

// SignUpAndGenerateToken mocks base method.
func (m *MockUserUseCase) SignUpAndGenerateToken(ctx context.Context, email, passward string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignUpAndGenerateToken", ctx, email, passward)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignUpAndGenerateToken indicates an expected call of SignUpAndGenerateToken.
func (mr *MockUserUseCaseMockRecorder) SignUpAndGenerateToken(ctx, email, passward interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignUpAndGenerateToken", reflect.TypeOf((*MockUserUseCase)(nil).SignUpAndGenerateToken), ctx, email, passward)
}
