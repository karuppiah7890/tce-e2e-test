package e2e

import (
	"testing"

	"github.com/karuppiah7890/tce-e2e-test/testutils/dockerprovider"
	"github.com/karuppiah7890/tce-e2e-test/testutils/utils"

	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
)

func TestDockerManagementAndWorkloadCluster(t *testing.T) {
	log.InitLogger("docker-mgmt-wkld-e2e")

	r := utils.DefaultClusterTestRunner{}
	utils.RunProviderTest(dockerprovider.PROVIDER, r)
}
