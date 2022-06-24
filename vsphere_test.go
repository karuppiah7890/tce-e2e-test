package e2e

import (
	"testing"

	"github.com/karuppiah7890/tce-e2e-test/testutils/tce"
	"github.com/karuppiah7890/tce-e2e-test/testutils/utils"
	"github.com/karuppiah7890/tce-e2e-test/testutils/vsphere"

	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
)

func TestManagementAndWorkloadCluster(t *testing.T) {
	log.InitLogger("vsphere-mgmt-wkld-e2e")

	r := utils.DefaultClusterTestRunner{}
	utils.RunProviderTest(vsphere.PROVIDER, r, tce.Package{})
}
