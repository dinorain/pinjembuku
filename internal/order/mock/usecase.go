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

// MockOrderUseCase is a mock of OrderUseCase interface.
type MockOrderUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockOrderUseCaseMockRecorder
}

// MockOrderUseCaseMockRecorder is the mock recorder for MockOrderUseCase.
type MockOrderUseCaseMockRecorder struct {
	mock *MockOrderUseCase
}

// NewMockOrderUseCase creates a new mock instance.
func NewMockOrderUseCase(ctrl *gomock.Controller) *MockOrderUseCase {
	mock := &MockOrderUseCase{ctrl: ctrl}
	mock.recorder = &MockOrderUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOrderUseCase) EXPECT() *MockOrderUseCaseMockRecorder {
	return m.recorder
}

// CachedFindById mocks base method.
func (m *MockOrderUseCase) CachedFindById(ctx context.Context, orderID uuid.UUID) (*models.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CachedFindById", ctx, orderID)
	ret0, _ := ret[0].(*models.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CachedFindById indicates an expected call of CachedFindById.
func (mr *MockOrderUseCaseMockRecorder) CachedFindById(ctx, orderID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CachedFindById", reflect.TypeOf((*MockOrderUseCase)(nil).CachedFindById), ctx, orderID)
}

// Create mocks base method.
func (m *MockOrderUseCase) Create(ctx context.Context, order *models.Order) (*models.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, order)
	ret0, _ := ret[0].(*models.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockOrderUseCaseMockRecorder) Create(ctx, order interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockOrderUseCase)(nil).Create), ctx, order)
}

// DeleteById mocks base method.
func (m *MockOrderUseCase) DeleteById(ctx context.Context, orderID uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteById", ctx, orderID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteById indicates an expected call of DeleteById.
func (mr *MockOrderUseCaseMockRecorder) DeleteById(ctx, orderID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteById", reflect.TypeOf((*MockOrderUseCase)(nil).DeleteById), ctx, orderID)
}

// FindAll mocks base method.
func (m *MockOrderUseCase) FindAll(ctx context.Context, pagination *utils.Pagination) ([]models.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAll", ctx, pagination)
	ret0, _ := ret[0].([]models.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAll indicates an expected call of FindAll.
func (mr *MockOrderUseCaseMockRecorder) FindAll(ctx, pagination interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAll", reflect.TypeOf((*MockOrderUseCase)(nil).FindAll), ctx, pagination)
}

// FindAllByLibrarianId mocks base method.
func (m *MockOrderUseCase) FindAllByLibrarianId(ctx context.Context, librarianID uuid.UUID, pagination *utils.Pagination) ([]models.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAllByLibrarianId", ctx, librarianID, pagination)
	ret0, _ := ret[0].([]models.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAllByLibrarianId indicates an expected call of FindAllByLibrarianId.
func (mr *MockOrderUseCaseMockRecorder) FindAllByLibrarianId(ctx, librarianID, pagination interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAllByLibrarianId", reflect.TypeOf((*MockOrderUseCase)(nil).FindAllByLibrarianId), ctx, librarianID, pagination)
}

// FindAllByUserId mocks base method.
func (m *MockOrderUseCase) FindAllByUserId(ctx context.Context, userID uuid.UUID, pagination *utils.Pagination) ([]models.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAllByUserId", ctx, userID, pagination)
	ret0, _ := ret[0].([]models.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAllByUserId indicates an expected call of FindAllByUserId.
func (mr *MockOrderUseCaseMockRecorder) FindAllByUserId(ctx, userID, pagination interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAllByUserId", reflect.TypeOf((*MockOrderUseCase)(nil).FindAllByUserId), ctx, userID, pagination)
}

// FindAllByUserIdLibrarianId mocks base method.
func (m *MockOrderUseCase) FindAllByUserIdLibrarianId(ctx context.Context, userID, librarianID uuid.UUID, pagination *utils.Pagination) ([]models.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAllByUserIdLibrarianId", ctx, userID, librarianID, pagination)
	ret0, _ := ret[0].([]models.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAllByUserIdLibrarianId indicates an expected call of FindAllByUserIdLibrarianId.
func (mr *MockOrderUseCaseMockRecorder) FindAllByUserIdLibrarianId(ctx, userID, librarianID, pagination interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAllByUserIdLibrarianId", reflect.TypeOf((*MockOrderUseCase)(nil).FindAllByUserIdLibrarianId), ctx, userID, librarianID, pagination)
}

// FindById mocks base method.
func (m *MockOrderUseCase) FindById(ctx context.Context, orderID uuid.UUID) (*models.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindById", ctx, orderID)
	ret0, _ := ret[0].(*models.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindById indicates an expected call of FindById.
func (mr *MockOrderUseCaseMockRecorder) FindById(ctx, orderID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindById", reflect.TypeOf((*MockOrderUseCase)(nil).FindById), ctx, orderID)
}

// UpdateById mocks base method.
func (m *MockOrderUseCase) UpdateById(ctx context.Context, order *models.Order) (*models.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateById", ctx, order)
	ret0, _ := ret[0].(*models.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateById indicates an expected call of UpdateById.
func (mr *MockOrderUseCaseMockRecorder) UpdateById(ctx, order interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateById", reflect.TypeOf((*MockOrderUseCase)(nil).UpdateById), ctx, order)
}