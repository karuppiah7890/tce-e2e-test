package vsphere

import (
	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
	"os"
)

const ManagementApiServerEndpoint = "VSPHERE_MANAGEMENT_CLUSTER_ENDPOINT"
const WorkloadApiServerEndpoint = "VSPHERE_WORKLOAD_CLUSTER_ENDPOINT"
const Url = "VSPHERE_SERVER"
const SshKeys = "VSPHERE_SSH_AUTHORIZED_KEY"
const Username = "VSPHERE_USERNAME"
const Password = "VSPHERE_PASSWORD"
const Datacenter = "VSPHERE_DATACENTER"
const Datastore = "VSPHERE_DATASTORE"
const VmFolder = "VSPHERE_FOLDER"
const Network = "VSPHERE_NETWORK"
const ResourcePool = "VSPHERE_RESOURCE_POOL"

type TestSecrets struct {
	ManagementApiServerEndpoint string
	WorkloadApiServerEndpoint   string
	SshKeys                     string
	Url                         string
	Username                    string
	Password                    string
	Datastore                   string
	Datacenter                  string
	VmFolder                    string
	Network                     string
	ResourcePool                string
}

func ExtractVsphereTestSecretsFromEnvVars() TestSecrets {
	CheckRequiredVsphereEnvVars()

	log.Info("Extracting AWS test secrets from environment variables")

	return TestSecrets{
		ManagementApiServerEndpoint: os.Getenv(ManagementApiServerEndpoint),
		WorkloadApiServerEndpoint:   os.Getenv(WorkloadApiServerEndpoint),
		SshKeys:                     os.Getenv(SshKeys),
		Url:                         os.Getenv(Url),
		Username:                    os.Getenv(Username),
		Password:                    os.Getenv(Password),
		Datastore:                   os.Getenv(Datastore),
		Datacenter:                  os.Getenv(Datacenter),
		VmFolder:                    os.Getenv(VmFolder),
		Network:                     os.Getenv(Network),
		ResourcePool:                os.Getenv(ResourcePool),
	}
}
