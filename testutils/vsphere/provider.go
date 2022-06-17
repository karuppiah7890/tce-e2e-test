package vsphere

import (
	"context"

	"github.com/karuppiah7890/tce-e2e-test/testutils/tanzu"
	"github.com/karuppiah7890/tce-e2e-test/testutils/utils"
)

// TODO: Change name?
type Provider struct {
	testSecrets TestSecrets
}

func (provider *Provider) RequiredEnvVars() []string {
	return []string{
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
}

func (provider *Provider) Name() string {
	return "vsphere"
}

func (provider *Provider) Init() error {
	provider.testSecrets = ExtractVsphereTestSecretsFromEnvVars()
	return nil
}

func (provider *Provider) PreClusterCreationTasks(clusterName string, clusterType utils.ClusterType) error {
	// To Update API server IP during runtime with VSPHERE_MANAGEMENT_CLUSTER_ENDPOINT to VSPHERE_CONTROL_PLANE_ENDPOINT as that is needed for cluster
	utils.UpdateVars(provider.Name(), clusterType)
	return nil
}

func (provider *Provider) CleanupCluster(ctx context.Context, clusterName string) error {
	// TODO: Implement using govmomi golang library
	return nil
}

//TODO: Maybe make use of https://github.com/spf13/viper to set env vars and make some values as default and parameterised.
func (provider *Provider) GetTanzuConfig(clusterName string) tanzu.TanzuConfig {
	return tanzu.TanzuConfig{
		"CLUSTER_NAME":                   clusterName,
		"INFRASTRUCTURE_PROVIDER":        provider.Name(),
		"CLUSTER_PLAN":                   "dev",
		"OS_ARCH":                        "amd64",
		"OS_NAME":                        "photon",
		"OS_VERSION":                     "3",
		"VSPHERE_CONTROL_PLANE_DISK_GIB": "40",
		"VSPHERE_CONTROL_PLANE_MEM_MIB":  "16384",
		"VSPHERE_CONTROL_PLANE_NUM_CPUS": "4",
		"VSPHERE_WORKER_DISK_GIB":        "40",
		"VSPHERE_WORKER_MEM_MIB":         "16384",
		"VSPHERE_WORKER_NUM_CPUS":        "4",
		"VSPHERE_INSECURE":               "true",
		"DEPLOY_TKG_ON_VSPHERE7":         "true",
		"ENABLE_TKGS_ON_VSPHERE7":        "false",
		"CLUSTER_CIDR":                   "100.96.0.0/11",
		"SERVICE_CIDR":                   "100.64.0.0/13",
		"ENABLE_CEIP_PARTICIPATION":      "false",
		"ENABLE_MHC":                     "true",
		"IDENTITY_MANAGEMENT_TYPE":       "none",
	}
}

// TODO: Change name?
var PROVIDER utils.Provider = &Provider{}
