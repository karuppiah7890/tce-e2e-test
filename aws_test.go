package e2e

import (
	"testing"

	"github.com/karuppiah7890/tce-e2e-test/testutils/tce"
	"github.com/karuppiah7890/tce-e2e-test/testutils/utils"

	"github.com/karuppiah7890/tce-e2e-test/testutils/aws"
	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
)

func TestAwsManagementAndWorkloadCluster(t *testing.T) {
	log.InitLogger("aws-mgmt-wkld-e2e")

	r := utils.DefaultClusterTestRunner{}
	err := utils.RunProviderTest(aws.PROVIDER, r, tce.Package{})
	if err != nil {
		t.Errorf("Error while running E2E test for managed and workload cluster on AWS: %v", err)
	}
}
