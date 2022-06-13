package e2e

import (
	"testing"

	"github.com/karuppiah7890/tce-e2e-test/testutils/utils"

	"github.com/karuppiah7890/tce-e2e-test/testutils/aws"
	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
)

// TODO: Make region as environment variable

// TODO: Consider making all as environment variables. Hard coded values in test code can be default.
// We can pass env vars to override stuff

// TODO: Do we really need individual test files for each provider or can we parameterise it ? thoughts

func TestAwsManagementAndWorkloadCluster(t *testing.T) {
	log.InitLogger("aws-mgmt-wkld-e2e")

	utils.RunProviderTest(aws.PROVIDER)
}
