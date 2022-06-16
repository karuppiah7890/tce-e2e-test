// Code generated by MockGen. DO NOT EDIT.
// Source: testutils/utils/cluster_test_runner.go

// Package mock_utils is a generated GoMock package.
package mock_utils

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	utils "github.com/karuppiah7890/tce-e2e-test/testutils/utils"
)

// MockClusterTestRunner is a mock of ClusterTestRunner interface.
type MockClusterTestRunner struct {
	ctrl     *gomock.Controller
	recorder *MockClusterTestRunnerMockRecorder
}

// MockClusterTestRunnerMockRecorder is the mock recorder for MockClusterTestRunner.
type MockClusterTestRunnerMockRecorder struct {
	mock *MockClusterTestRunner
}

// NewMockClusterTestRunner creates a new mock instance.
func NewMockClusterTestRunner(ctrl *gomock.Controller) *MockClusterTestRunner {
	mock := &MockClusterTestRunner{ctrl: ctrl}
	mock.recorder = &MockClusterTestRunnerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClusterTestRunner) EXPECT() *MockClusterTestRunnerMockRecorder {
	return m.recorder
}

// CheckWorkloadClusterIsRunning mocks base method.
func (m *MockClusterTestRunner) CheckWorkloadClusterIsRunning(workloadClusterName string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckWorkloadClusterIsRunning", workloadClusterName)
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckWorkloadClusterIsRunning indicates an expected call of CheckWorkloadClusterIsRunning.
func (mr *MockClusterTestRunnerMockRecorder) CheckWorkloadClusterIsRunning(workloadClusterName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckWorkloadClusterIsRunning", reflect.TypeOf((*MockClusterTestRunner)(nil).CheckWorkloadClusterIsRunning), workloadClusterName)
}

// CleanupDockerBootstrapCluster mocks base method.
func (m *MockClusterTestRunner) CleanupDockerBootstrapCluster(managementClusterName string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CleanupDockerBootstrapCluster", managementClusterName)
	ret0, _ := ret[0].(error)
	return ret0
}

// CleanupDockerBootstrapCluster indicates an expected call of CleanupDockerBootstrapCluster.
func (mr *MockClusterTestRunnerMockRecorder) CleanupDockerBootstrapCluster(managementClusterName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CleanupDockerBootstrapCluster", reflect.TypeOf((*MockClusterTestRunner)(nil).CleanupDockerBootstrapCluster), managementClusterName)
}

// CollectManagementClusterAndWorkloadClusterDiagnostics mocks base method.
func (m *MockClusterTestRunner) CollectManagementClusterAndWorkloadClusterDiagnostics(managementClusterName, workloadClusterName, workloadClusterInfra string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CollectManagementClusterAndWorkloadClusterDiagnostics", managementClusterName, workloadClusterName, workloadClusterInfra)
	ret0, _ := ret[0].(error)
	return ret0
}

// CollectManagementClusterAndWorkloadClusterDiagnostics indicates an expected call of CollectManagementClusterAndWorkloadClusterDiagnostics.
func (mr *MockClusterTestRunnerMockRecorder) CollectManagementClusterAndWorkloadClusterDiagnostics(managementClusterName, workloadClusterName, workloadClusterInfra interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CollectManagementClusterAndWorkloadClusterDiagnostics", reflect.TypeOf((*MockClusterTestRunner)(nil).CollectManagementClusterAndWorkloadClusterDiagnostics), managementClusterName, workloadClusterName, workloadClusterInfra)
}

// CollectManagementClusterDiagnostics mocks base method.
func (m *MockClusterTestRunner) CollectManagementClusterDiagnostics(managementClusterName string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CollectManagementClusterDiagnostics", managementClusterName)
	ret0, _ := ret[0].(error)
	return ret0
}

// CollectManagementClusterDiagnostics indicates an expected call of CollectManagementClusterDiagnostics.
func (mr *MockClusterTestRunnerMockRecorder) CollectManagementClusterDiagnostics(managementClusterName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CollectManagementClusterDiagnostics", reflect.TypeOf((*MockClusterTestRunner)(nil).CollectManagementClusterDiagnostics), managementClusterName)
}

// DeleteCluster mocks base method.
func (m *MockClusterTestRunner) DeleteCluster(clusterName string, provider utils.Provider, clusterType utils.ClusterType) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCluster", clusterName, provider, clusterType)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCluster indicates an expected call of DeleteCluster.
func (mr *MockClusterTestRunnerMockRecorder) DeleteCluster(clusterName, provider, clusterType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCluster", reflect.TypeOf((*MockClusterTestRunner)(nil).DeleteCluster), clusterName, provider, clusterType)
}

// DeleteContext mocks base method.
func (m *MockClusterTestRunner) DeleteContext(kubeConfigPath, contextName string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteContext", kubeConfigPath, contextName)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteContext indicates an expected call of DeleteContext.
func (mr *MockClusterTestRunnerMockRecorder) DeleteContext(kubeConfigPath, contextName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteContext", reflect.TypeOf((*MockClusterTestRunner)(nil).DeleteContext), kubeConfigPath, contextName)
}

// GetClusterKubeConfig mocks base method.
func (m *MockClusterTestRunner) GetClusterKubeConfig(clusterName string, provider utils.Provider, clusterType utils.ClusterType) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "GetClusterKubeConfig", clusterName, provider, clusterType)
}

// GetClusterKubeConfig indicates an expected call of GetClusterKubeConfig.
func (mr *MockClusterTestRunnerMockRecorder) GetClusterKubeConfig(clusterName, provider, clusterType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetClusterKubeConfig", reflect.TypeOf((*MockClusterTestRunner)(nil).GetClusterKubeConfig), clusterName, provider, clusterType)
}

// GetKubeConfigPath mocks base method.
func (m *MockClusterTestRunner) GetKubeConfigPath() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetKubeConfigPath")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetKubeConfigPath indicates an expected call of GetKubeConfigPath.
func (mr *MockClusterTestRunnerMockRecorder) GetKubeConfigPath() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetKubeConfigPath", reflect.TypeOf((*MockClusterTestRunner)(nil).GetKubeConfigPath))
}

// GetKubeContextForTanzuCluster mocks base method.
func (m *MockClusterTestRunner) GetKubeContextForTanzuCluster(clusterName string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetKubeContextForTanzuCluster", clusterName)
	ret0, _ := ret[0].(string)
	return ret0
}

// GetKubeContextForTanzuCluster indicates an expected call of GetKubeContextForTanzuCluster.
func (mr *MockClusterTestRunnerMockRecorder) GetKubeContextForTanzuCluster(clusterName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetKubeContextForTanzuCluster", reflect.TypeOf((*MockClusterTestRunner)(nil).GetKubeContextForTanzuCluster), clusterName)
}

// GetRandomClusterNames mocks base method.
func (m *MockClusterTestRunner) GetRandomClusterNames() (string, string) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRandomClusterNames")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(string)
	return ret0, ret1
}

// GetRandomClusterNames indicates an expected call of GetRandomClusterNames.
func (mr *MockClusterTestRunnerMockRecorder) GetRandomClusterNames() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRandomClusterNames", reflect.TypeOf((*MockClusterTestRunner)(nil).GetRandomClusterNames))
}

// PrintClusterInformation mocks base method.
func (m *MockClusterTestRunner) PrintClusterInformation(kubeConfigPath, kubeContext string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PrintClusterInformation", kubeConfigPath, kubeContext)
	ret0, _ := ret[0].(error)
	return ret0
}

// PrintClusterInformation indicates an expected call of PrintClusterInformation.
func (mr *MockClusterTestRunnerMockRecorder) PrintClusterInformation(kubeConfigPath, kubeContext interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PrintClusterInformation", reflect.TypeOf((*MockClusterTestRunner)(nil).PrintClusterInformation), kubeConfigPath, kubeContext)
}

// RunChecks mocks base method.
func (m *MockClusterTestRunner) RunChecks() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RunChecks")
}

// RunChecks indicates an expected call of RunChecks.
func (mr *MockClusterTestRunnerMockRecorder) RunChecks() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RunChecks", reflect.TypeOf((*MockClusterTestRunner)(nil).RunChecks))
}

// RunCluster mocks base method.
func (m *MockClusterTestRunner) RunCluster(clusterName string, provider utils.Provider, clusterType utils.ClusterType) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RunCluster", clusterName, provider, clusterType)
	ret0, _ := ret[0].(error)
	return ret0
}

// RunCluster indicates an expected call of RunCluster.
func (mr *MockClusterTestRunnerMockRecorder) RunCluster(clusterName, provider, clusterType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RunCluster", reflect.TypeOf((*MockClusterTestRunner)(nil).RunCluster), clusterName, provider, clusterType)
}

// WaitForWorkloadClusterDeletion mocks base method.
func (m *MockClusterTestRunner) WaitForWorkloadClusterDeletion(workloadClusterName string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WaitForWorkloadClusterDeletion", workloadClusterName)
	ret0, _ := ret[0].(error)
	return ret0
}

// WaitForWorkloadClusterDeletion indicates an expected call of WaitForWorkloadClusterDeletion.
func (mr *MockClusterTestRunnerMockRecorder) WaitForWorkloadClusterDeletion(workloadClusterName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WaitForWorkloadClusterDeletion", reflect.TypeOf((*MockClusterTestRunner)(nil).WaitForWorkloadClusterDeletion), workloadClusterName)
}