// Code generated by MockGen. DO NOT EDIT.
// Source: ../contracts/TenantRepositoryServiceContract.go

// Package mock_contracts is a generated GoMock package.
package mock_contracts

import (
	context "context"
	contracts "github.com/decentralized-cloud/Tenant/repository/contracts"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockTenantRepositoryServiceContract is a mock of TenantRepositoryServiceContract interface
type MockTenantRepositoryServiceContract struct {
	ctrl     *gomock.Controller
	recorder *MockTenantRepositoryServiceContractMockRecorder
}

// MockTenantRepositoryServiceContractMockRecorder is the mock recorder for MockTenantRepositoryServiceContract
type MockTenantRepositoryServiceContractMockRecorder struct {
	mock *MockTenantRepositoryServiceContract
}

// NewMockTenantRepositoryServiceContract creates a new mock instance
func NewMockTenantRepositoryServiceContract(ctrl *gomock.Controller) *MockTenantRepositoryServiceContract {
	mock := &MockTenantRepositoryServiceContract{ctrl: ctrl}
	mock.recorder = &MockTenantRepositoryServiceContractMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockTenantRepositoryServiceContract) EXPECT() *MockTenantRepositoryServiceContractMockRecorder {
	return m.recorder
}

// CreateTenant mocks base method
func (m *MockTenantRepositoryServiceContract) CreateTenant(ctx context.Context, request *contracts.CreateTenantRequest) (*contracts.CreateTenantResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTenant", ctx, request)
	ret0, _ := ret[0].(*contracts.CreateTenantResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTenant indicates an expected call of CreateTenant
func (mr *MockTenantRepositoryServiceContractMockRecorder) CreateTenant(ctx, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTenant", reflect.TypeOf((*MockTenantRepositoryServiceContract)(nil).CreateTenant), ctx, request)
}

// ReadTenant mocks base method
func (m *MockTenantRepositoryServiceContract) ReadTenant(ctx context.Context, request *contracts.ReadTenantRequest) (*contracts.ReadTenantResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadTenant", ctx, request)
	ret0, _ := ret[0].(*contracts.ReadTenantResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadTenant indicates an expected call of ReadTenant
func (mr *MockTenantRepositoryServiceContractMockRecorder) ReadTenant(ctx, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadTenant", reflect.TypeOf((*MockTenantRepositoryServiceContract)(nil).ReadTenant), ctx, request)
}

// UpdateTenant mocks base method
func (m *MockTenantRepositoryServiceContract) UpdateTenant(ctx context.Context, request *contracts.UpdateTenantRequest) (*contracts.UpdateTenantResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTenant", ctx, request)
	ret0, _ := ret[0].(*contracts.UpdateTenantResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateTenant indicates an expected call of UpdateTenant
func (mr *MockTenantRepositoryServiceContractMockRecorder) UpdateTenant(ctx, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTenant", reflect.TypeOf((*MockTenantRepositoryServiceContract)(nil).UpdateTenant), ctx, request)
}

// DeleteTenant mocks base method
func (m *MockTenantRepositoryServiceContract) DeleteTenant(ctx context.Context, request *contracts.DeleteTenantRequest) (*contracts.DeleteTenantResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTenant", ctx, request)
	ret0, _ := ret[0].(*contracts.DeleteTenantResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteTenant indicates an expected call of DeleteTenant
func (mr *MockTenantRepositoryServiceContractMockRecorder) DeleteTenant(ctx, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTenant", reflect.TypeOf((*MockTenantRepositoryServiceContract)(nil).DeleteTenant), ctx, request)
}
