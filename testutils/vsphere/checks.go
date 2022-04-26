package vsphere

import (
	"github.com/karuppiah7890/tce-e2e-test/testutils"
	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
)

func CheckRequiredVsphereEnvVars() {
	log.Info("Checking required vSphere environment variables")

	requiredEnvVars := []string{
		ManagementApiServerEndpoint,
		WorkloadApiServerEndpoint,
		SshKeys,
		Url,
		Username,
		Password,
		Datastore,
		Datacenter,
		VmFolder,
		Network,
		ResourcePool,
	}
	errs := testutils.CheckRequiredEnvVars(requiredEnvVars)

	if len(errs) != 0 {
		log.Fatalf("Errors while checking required environment variables: %v\n", errs)
	}
}
