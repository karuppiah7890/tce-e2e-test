package e2e

import (
	"os"
	"testing"

	"github.com/karuppiah7890/tce-e2e-test/testutils/aws"
	"github.com/karuppiah7890/tce-e2e-test/testutils/azure"
	"github.com/karuppiah7890/tce-e2e-test/testutils/dockerprovider"
	"github.com/karuppiah7890/tce-e2e-test/testutils/github"
	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
	"github.com/karuppiah7890/tce-e2e-test/testutils/tce"
	"github.com/karuppiah7890/tce-e2e-test/testutils/utils"
	"github.com/karuppiah7890/tce-e2e-test/testutils/vsphere"
)

func TestCloneTCERepo(t *testing.T) {
	err := github.CloneRepo("https://github.com/vmware-tanzu/community-edition")
	if err != nil {
		log.Errorf("Error while cloning TCE Repo: %v", err)
	}

	packageDetails := tce.Package{}
	packageDetails.Name = os.Getenv("PACKAGE_NAME")
	packageDetails.Version = os.Getenv("PACKAGE_VERSION")
	packageDetails.ManualCreate = true
	provider := os.Getenv("PROVIDER")

	log.InitLogger(provider + "-mgmt-wkld-e2e")
	r := utils.DefaultClusterTestRunner{}

	if provider == "docker" {
		// Running package E2E test on docker
		err = utils.RunProviderTest(dockerprovider.PROVIDER, r, packageDetails)
	} else if provider == "aws" {
		// Running package E2E test on AWS
		err = utils.RunProviderTest(aws.PROVIDER, r, packageDetails)
	} else if provider == "azure" {
		// Running package E2E test on Azure
		err = utils.RunProviderTest(azure.PROVIDER, r, packageDetails)
	} else if provider == "vsphere" {
		// Running package E2E test on vSphere
		utils.RunProviderTest(vsphere.PROVIDER, r, tce.Package{})
	} else {
		// Invalid provider in PROVIDER env variable
		t.Errorf("Ivalid provider for package E2E Test")
	}
	if packageDetails.Name == "velero" {
		aws.EmptyS3Bucket()
	}
	if err != nil {
		t.Errorf("Error while running package E2E test on %v: %v", provider, err)
	}
}
