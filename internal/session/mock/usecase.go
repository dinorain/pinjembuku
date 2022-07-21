// Code generated by MockGen. DO NOT EDIT.
// Source: usecase.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	models "github.com/dinorain/pinjembuku/internal/models"
	gomock "github.com/golang/mock/gomock"
)

// MockSessUseCase is a mock of SessUseCase interface.
type MockSessUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockSessUseCaseMockRecorder
}

// MockSessUseCaseMockRecorder is the mock recorder for MockSessUseCase.
type MockSessUseCaseMockRecorder struct {
	mock *MockSessUseCase
}

// NewMockSessUseCase creates a new mock instance.
func NewMockSessUseCase(ctrl *gomock.Controller) *MockSessUseCase {
	mock := &MockSessUseCase{ctrl: ctrl}
	mock.recorder = &MockSessUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSessUseCase) EXPECT() *MockSessUseCaseMockRecorder {
	return m.recorder
}

// CreateSession mocks base method.
func (m *MockSessUseCase) CreateSession(ctx context.Context, session *models.Session, expire int) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSession", ctx, session, expire)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSession indicates an expected call of CreateSession.
func (mr *MockSessUseCaseMockRecorder) CreateSession(ctx, session, expire interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSession", reflect.TypeOf((*MockSessUseCase)(nil).CreateSession), ctx, session, expire)
}

// DeleteById mocks base method.
func (m *MockSessUseCase) DeleteById(ctx context.Context, sessionID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteById", ctx, sessionID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteById indicates an expected call of DeleteById.
func (mr *MockSessUseCaseMockRecorder) DeleteById(ctx, sessionID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteById", reflect.TypeOf((*MockSessUseCase)(nil).DeleteById), ctx, sessionID)
}

// GetSessionById mocks base method.
func (m *MockSessUseCase) GetSessionById(ctx context.Context, sessionID string) (*models.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSessionById", ctx, sessionID)
	ret0, _ := ret[0].(*models.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSessionById indicates an expected call of GetSessionById.
func (mr *MockSessUseCaseMockRecorder) GetSessionById(ctx, sessionID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSessionById", reflect.TypeOf((*MockSessUseCase)(nil).GetSessionById), ctx, sessionID)
}
