package e2e

import (
	"os"
	"testing"

	"github.com/karuppiah7890/tce-e2e-test/testutils/dockerprovider"
	"github.com/karuppiah7890/tce-e2e-test/testutils/github"
	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
	"github.com/karuppiah7890/tce-e2e-test/testutils/utils"
)

func TestCloneTCERepo(t *testing.T) {
	err := github.CloneRepo("https://github.com/vmware-tanzu/community-edition")
	if err != nil {
		log.Errorf("Error while cloning TCE Repo: %v", err)
	}

	log.InitLogger("docker-mgmt-wkld-e2e")

	packageDetails := utils.Package{}
	packageDetails.Name = os.Getenv("PACKAGE_NAME")
	packageDetails.Version = os.Getenv("PACKAGE_VERSION")

	r := utils.DefaultClusterTestRunner{}
	utils.RunProviderTest(dockerprovider.PROVIDER, r, packageDetails)
}
