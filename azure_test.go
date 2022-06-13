package e2e

import (
	"testing"

	"github.com/karuppiah7890/tce-e2e-test/testutils/azure"
	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
	"github.com/karuppiah7890/tce-e2e-test/testutils/utils"
)

func TestAzureManagementAndWorkloadCluster(t *testing.T) {
	log.InitLogger("azure-mgmt-wkld-e2e")

	utils.RunProviderTest(azure.PROVIDER)
}
