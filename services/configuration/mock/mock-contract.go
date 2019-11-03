// Code generated by MockGen. DO NOT EDIT.
// Source: ../contract.go

// Package mock_configuration is a generated GoMock package.
package mock_configuration

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockConfigurationContract is a mock of ConfigurationContract interface
type MockConfigurationContract struct {
	ctrl     *gomock.Controller
	recorder *MockConfigurationContractMockRecorder
}

// MockConfigurationContractMockRecorder is the mock recorder for MockConfigurationContract
type MockConfigurationContractMockRecorder struct {
	mock *MockConfigurationContract
}

// NewMockConfigurationContract creates a new mock instance
func NewMockConfigurationContract(ctrl *gomock.Controller) *MockConfigurationContract {
	mock := &MockConfigurationContract{ctrl: ctrl}
	mock.recorder = &MockConfigurationContractMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockConfigurationContract) EXPECT() *MockConfigurationContractMockRecorder {
	return m.recorder
}

// GetGrpcHost mocks base method
func (m *MockConfigurationContract) GetGrpcHost() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGrpcHost")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGrpcHost indicates an expected call of GetGrpcHost
func (mr *MockConfigurationContractMockRecorder) GetGrpcHost() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGrpcHost", reflect.TypeOf((*MockConfigurationContract)(nil).GetGrpcHost))
}

// GetGrpcPort mocks base method
func (m *MockConfigurationContract) GetGrpcPort() (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGrpcPort")
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGrpcPort indicates an expected call of GetGrpcPort
func (mr *MockConfigurationContractMockRecorder) GetGrpcPort() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGrpcPort", reflect.TypeOf((*MockConfigurationContract)(nil).GetGrpcPort))
}

// GetHttpsHost mocks base method
func (m *MockConfigurationContract) GetHttpsHost() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetHttpsHost")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetHttpsHost indicates an expected call of GetHttpsHost
func (mr *MockConfigurationContractMockRecorder) GetHttpsHost() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHttpsHost", reflect.TypeOf((*MockConfigurationContract)(nil).GetHttpsHost))
}

// GetHttpsPort mocks base method
func (m *MockConfigurationContract) GetHttpsPort() (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetHttpsPort")
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetHttpsPort indicates an expected call of GetHttpsPort
func (mr *MockConfigurationContractMockRecorder) GetHttpsPort() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHttpsPort", reflect.TypeOf((*MockConfigurationContract)(nil).GetHttpsPort))
}

// GetDatabaseConnectionString mocks base method
func (m *MockConfigurationContract) GetDatabaseConnectionString() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDatabaseConnectionString")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDatabaseConnectionString indicates an expected call of GetDatabaseConnectionString
func (mr *MockConfigurationContractMockRecorder) GetDatabaseConnectionString() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDatabaseConnectionString", reflect.TypeOf((*MockConfigurationContract)(nil).GetDatabaseConnectionString))
}

// GetDatabaseName mocks base method
func (m *MockConfigurationContract) GetDatabaseName() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDatabaseName")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDatabaseName indicates an expected call of GetDatabaseName
func (mr *MockConfigurationContractMockRecorder) GetDatabaseName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDatabaseName", reflect.TypeOf((*MockConfigurationContract)(nil).GetDatabaseName))
}
