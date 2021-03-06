// Code generated by MockGen. DO NOT EDIT.
// Source: testutils/utils/infraproviders.go

// Package mock_utils is a generated GoMock package.
package mock_utils

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	tanzu "github.com/karuppiah7890/tce-e2e-test/testutils/tanzu"
	utils "github.com/karuppiah7890/tce-e2e-test/testutils/utils"
)

// MockProvider is a mock of Provider interface.
type MockProvider struct {
	ctrl     *gomock.Controller
	recorder *MockProviderMockRecorder
}

// MockProviderMockRecorder is the mock recorder for MockProvider.
type MockProviderMockRecorder struct {
	mock *MockProvider
}

// NewMockProvider creates a new mock instance.
func NewMockProvider(ctrl *gomock.Controller) *MockProvider {
	mock := &MockProvider{ctrl: ctrl}
	mock.recorder = &MockProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProvider) EXPECT() *MockProviderMockRecorder {
	return m.recorder
}

// CleanupCluster mocks base method.
func (m *MockProvider) CleanupCluster(ctx context.Context, clusterName string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CleanupCluster", ctx, clusterName)
	ret0, _ := ret[0].(error)
	return ret0
}

// CleanupCluster indicates an expected call of CleanupCluster.
func (mr *MockProviderMockRecorder) CleanupCluster(ctx, clusterName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CleanupCluster", reflect.TypeOf((*MockProvider)(nil).CleanupCluster), ctx, clusterName)
}

// GetTanzuConfig mocks base method.
func (m *MockProvider) GetTanzuConfig(clusterName string) tanzu.TanzuConfig {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTanzuConfig", clusterName)
	ret0, _ := ret[0].(tanzu.TanzuConfig)
	return ret0
}

// GetTanzuConfig indicates an expected call of GetTanzuConfig.
func (mr *MockProviderMockRecorder) GetTanzuConfig(clusterName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTanzuConfig", reflect.TypeOf((*MockProvider)(nil).GetTanzuConfig), clusterName)
}

// Init mocks base method.
func (m *MockProvider) Init() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Init")
	ret0, _ := ret[0].(error)
	return ret0
}

// Init indicates an expected call of Init.
func (mr *MockProviderMockRecorder) Init() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Init", reflect.TypeOf((*MockProvider)(nil).Init))
}

// Name mocks base method.
func (m *MockProvider) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name.
func (mr *MockProviderMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockProvider)(nil).Name))
}

// PreClusterCreationTasks mocks base method.
func (m *MockProvider) PreClusterCreationTasks(clusterName string, clusterType utils.ClusterType) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PreClusterCreationTasks", clusterName, clusterType)
	ret0, _ := ret[0].(error)
	return ret0
}

// PreClusterCreationTasks indicates an expected call of PreClusterCreationTasks.
func (mr *MockProviderMockRecorder) PreClusterCreationTasks(clusterName, clusterType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PreClusterCreationTasks", reflect.TypeOf((*MockProvider)(nil).PreClusterCreationTasks), clusterName, clusterType)
}

// RequiredEnvVars mocks base method.
func (m *MockProvider) RequiredEnvVars() []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RequiredEnvVars")
	ret0, _ := ret[0].([]string)
	return ret0
}

// RequiredEnvVars indicates an expected call of RequiredEnvVars.
func (mr *MockProviderMockRecorder) RequiredEnvVars() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RequiredEnvVars", reflect.TypeOf((*MockProvider)(nil).RequiredEnvVars))
}
