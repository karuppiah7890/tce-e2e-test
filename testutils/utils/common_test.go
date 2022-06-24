package utils_test

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/karuppiah7890/tce-e2e-test/testutils/tce"
	"github.com/karuppiah7890/tce-e2e-test/testutils/utils"
	"github.com/karuppiah7890/tce-e2e-test/testutils/utils/mock_utils"
)

func TestRunProviderTest(t *testing.T) {

	t.Run("when management cluster creation fails it should collect diagnostics", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		defer ctrl.Finish()

		provider := mock_utils.NewMockProvider(ctrl)
		r := mock_utils.NewMockClusterTestRunner(ctrl)

		gomock.InOrder(
			r.EXPECT().RunChecks(),

			provider.EXPECT().RequiredEnvVars(),

			provider.EXPECT().Init(),

			r.EXPECT().
				GetRandomClusterNames().Return("test-mgmt", "test-wkld"),

			provider.EXPECT().PreClusterCreationTasks("test-mgmt", utils.ManagementClusterType),

			r.EXPECT().GetKubeContextForTanzuCluster("test-mgmt").Return("mock-context"),

			r.EXPECT().GetKubeConfigPath().Return("mock-config-path", nil),

			r.EXPECT().
				RunCluster("test-mgmt", provider, utils.ManagementClusterType).
				Return(fmt.Errorf("some error in management cluster creation")),

			r.EXPECT().CollectManagementClusterDiagnostics("test-mgmt"),
			r.EXPECT().CleanupDockerBootstrapCluster("test-mgmt"),
			r.EXPECT().DeleteContext("mock-config-path", "mock-context"),
			provider.EXPECT().CleanupCluster(gomock.Any(), "test-mgmt"),
		)

		err := utils.RunProviderTest(provider, r, tce.Package{})
		expectedError := "error while running management cluster: some error in management cluster creation"
		if err.Error() != expectedError {
			t.Logf("expected error to be: %v. But got: %v", expectedError, err)
			t.Fail()
		}
	})

	t.Run("when workload cluster creation fails it should collect diagnostics", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		defer ctrl.Finish()

		provider := mock_utils.NewMockProvider(ctrl)
		r := mock_utils.NewMockClusterTestRunner(ctrl)

		provider.EXPECT().Name().Return("mock-infra")

		gomock.InOrder(
			r.EXPECT().RunChecks(),

			provider.EXPECT().RequiredEnvVars(),

			provider.EXPECT().Init(),

			r.EXPECT().
				GetRandomClusterNames().Return("test-mgmt", "test-wkld"),

			provider.EXPECT().PreClusterCreationTasks("test-mgmt", utils.ManagementClusterType),

			r.EXPECT().GetKubeContextForTanzuCluster("test-mgmt").Return("mock-context-1"),

			r.EXPECT().GetKubeConfigPath().Return("mock-config-path", nil),

			r.EXPECT().
				RunCluster("test-mgmt", provider, utils.ManagementClusterType),

			r.EXPECT().GetClusterKubeConfig("test-mgmt", provider, utils.ManagementClusterType),

			r.EXPECT().PrintClusterInformation("mock-config-path", "mock-context-1"),

			provider.EXPECT().PreClusterCreationTasks("test-wkld", utils.WorkloadClusterType),

			r.EXPECT().GetKubeContextForTanzuCluster("test-wkld").Return("mock-context-2"),

			r.EXPECT().
				RunCluster("test-wkld", provider, utils.WorkloadClusterType).
				Return(fmt.Errorf("some error in workload cluster creation")),

			r.EXPECT().CollectManagementClusterAndWorkloadClusterDiagnostics("test-mgmt", "test-wkld", "mock-infra"),
			provider.EXPECT().CleanupCluster(gomock.Any(), "test-mgmt"),
			r.EXPECT().DeleteContext("mock-config-path", "mock-context-1"),
			provider.EXPECT().CleanupCluster(gomock.Any(), "test-wkld"),
			r.EXPECT().DeleteContext("mock-config-path", "mock-context-2"),
		)

		err := utils.RunProviderTest(provider, r, tce.Package{})
		expectedError := "error while running workload cluster: some error in workload cluster creation"
		if err.Error() != expectedError {
			t.Logf("expected error to be: %v. But got: %v", expectedError, err)
			t.Fail()
		}
	})
}
