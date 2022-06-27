package e2e

import (
	"testing"

	"github.com/karuppiah7890/tce-e2e-test/testutils/azure"
	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
	"github.com/karuppiah7890/tce-e2e-test/testutils/tce"
	"github.com/karuppiah7890/tce-e2e-test/testutils/utils"
)

func TestAzureManagementAndWorkloadCluster(t *testing.T) {
	log.InitLogger("azure-mgmt-wkld-e2e")

	r := utils.DefaultClusterTestRunner{}
	err := utils.RunProviderTest(azure.PROVIDER, r, tce.Package{})
	if err != nil {
		t.Errorf("Error while running E2E test for managed and workload cluster on Azure: %v", err)
	}
}
