package e2e

import (
	"fmt"
	"testing"

	"github.com/karuppiah7890/tce-e2e-test/testutils/dockerprovider"
	"github.com/karuppiah7890/tce-e2e-test/testutils/utils"

	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
)

func TestDockerManagementAndWorkloadCluster(t *testing.T) {
	log.InitLogger(fmt.Sprintf("docker-mgmt-wkld-e2e"))

	utils.RunProviderTest(dockerprovider.PROVIDER)
}
