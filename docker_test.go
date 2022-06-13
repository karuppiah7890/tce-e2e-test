package e2e

import (
	"fmt"
	"testing"

	"github.com/karuppiah7890/tce-e2e-test/testutils/dockerprovider"
	"github.com/karuppiah7890/tce-e2e-test/testutils/utils"

	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
)

// TODO: Make region as environment variable

// TODO: Consider making all as environment variables. Hard coded values in test code can be default.
// We can pass env vars to override stuff
func TestDockerManagementAndWorkloadCluster(t *testing.T) {
	provider := utils.Docker
	log.InitLogger(fmt.Sprintf("%s-mgmt-wkld-e2e", provider))

	utils.RunProviderTest(dockerprovider.PROVIDER)
}
