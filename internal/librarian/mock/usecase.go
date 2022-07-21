// Code generated by MockGen. DO NOT EDIT.
// Source: usecase.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	models "github.com/dinorain/pinjembuku/internal/models"
	utils "github.com/dinorain/pinjembuku/pkg/utils"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockLibrarianUseCase is a mock of LibrarianUseCase interface.
type MockLibrarianUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockLibrarianUseCaseMockRecorder
}

// MockLibrarianUseCaseMockRecorder is the mock recorder for MockLibrarianUseCase.
type MockLibrarianUseCaseMockRecorder struct {
	mock *MockLibrarianUseCase
}

// NewMockLibrarianUseCase creates a new mock instance.
func NewMockLibrarianUseCase(ctrl *gomock.Controller) *MockLibrarianUseCase {
	mock := &MockLibrarianUseCase{ctrl: ctrl}
	mock.recorder = &MockLibrarianUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLibrarianUseCase) EXPECT() *MockLibrarianUseCaseMockRecorder {
	return m.recorder
}

// CachedFindById mocks base method.
func (m *MockLibrarianUseCase) CachedFindById(ctx context.Context, librarianID uuid.UUID) (*models.Librarian, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CachedFindById", ctx, librarianID)
	ret0, _ := ret[0].(*models.Librarian)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CachedFindById indicates an expected call of CachedFindById.
func (mr *MockLibrarianUseCaseMockRecorder) CachedFindById(ctx, librarianID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CachedFindById", reflect.TypeOf((*MockLibrarianUseCase)(nil).CachedFindById), ctx, librarianID)
}

// DeleteById mocks base method.
func (m *MockLibrarianUseCase) DeleteById(ctx context.Context, librarianID uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteById", ctx, librarianID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteById indicates an expected call of DeleteById.
func (mr *MockLibrarianUseCaseMockRecorder) DeleteById(ctx, librarianID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteById", reflect.TypeOf((*MockLibrarianUseCase)(nil).DeleteById), ctx, librarianID)
}

// FindAll mocks base method.
func (m *MockLibrarianUseCase) FindAll(ctx context.Context, pagination *utils.Pagination) ([]models.Librarian, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAll", ctx, pagination)
	ret0, _ := ret[0].([]models.Librarian)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAll indicates an expected call of FindAll.
func (mr *MockLibrarianUseCaseMockRecorder) FindAll(ctx, pagination interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAll", reflect.TypeOf((*MockLibrarianUseCase)(nil).FindAll), ctx, pagination)
}

// FindByEmail mocks base method.
func (m *MockLibrarianUseCase) FindByEmail(ctx context.Context, email string) (*models.Librarian, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByEmail", ctx, email)
	ret0, _ := ret[0].(*models.Librarian)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByEmail indicates an expected call of FindByEmail.
func (mr *MockLibrarianUseCaseMockRecorder) FindByEmail(ctx, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByEmail", reflect.TypeOf((*MockLibrarianUseCase)(nil).FindByEmail), ctx, email)
}

// FindById mocks base method.
func (m *MockLibrarianUseCase) FindById(ctx context.Context, librarianID uuid.UUID) (*models.Librarian, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindById", ctx, librarianID)
	ret0, _ := ret[0].(*models.Librarian)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindById indicates an expected call of FindById.
func (mr *MockLibrarianUseCaseMockRecorder) FindById(ctx, librarianID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindById", reflect.TypeOf((*MockLibrarianUseCase)(nil).FindById), ctx, librarianID)
}

// GenerateTokenPair mocks base method.
func (m *MockLibrarianUseCase) GenerateTokenPair(librarian *models.Librarian, sessionID string) (string, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateTokenPair", librarian, sessionID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GenerateTokenPair indicates an expected call of GenerateTokenPair.
func (mr *MockLibrarianUseCaseMockRecorder) GenerateTokenPair(librarian, sessionID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateTokenPair", reflect.TypeOf((*MockLibrarianUseCase)(nil).GenerateTokenPair), librarian, sessionID)
}

// Login mocks base method.
func (m *MockLibrarianUseCase) Login(ctx context.Context, email, password string) (*models.Librarian, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", ctx, email, password)
	ret0, _ := ret[0].(*models.Librarian)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Login indicates an expected call of Login.
func (mr *MockLibrarianUseCaseMockRecorder) Login(ctx, email, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockLibrarianUseCase)(nil).Login), ctx, email, password)
}

// Register mocks base method.
func (m *MockLibrarianUseCase) Register(ctx context.Context, librarian *models.Librarian) (*models.Librarian, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", ctx, librarian)
	ret0, _ := ret[0].(*models.Librarian)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Register indicates an expected call of Register.
func (mr *MockLibrarianUseCaseMockRecorder) Register(ctx, librarian interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockLibrarianUseCase)(nil).Register), ctx, librarian)
}

// UpdateById mocks base method.
func (m *MockLibrarianUseCase) UpdateById(ctx context.Context, librarian *models.Librarian) (*models.Librarian, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateById", ctx, librarian)
	ret0, _ := ret[0].(*models.Librarian)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateById indicates an expected call of UpdateById.
func (mr *MockLibrarianUseCaseMockRecorder) UpdateById(ctx, librarian interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateById", reflect.TypeOf((*MockLibrarianUseCase)(nil).UpdateById), ctx, librarian)
}
