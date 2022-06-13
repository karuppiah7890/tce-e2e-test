package vsphere

import (
	"context"

	"github.com/karuppiah7890/tce-e2e-test/testutils/utils"
)

// TODO: Change name?
type Provider struct {
	testSecrets TestSecrets
}

func (provider *Provider) CheckRequiredEnvVars() bool {
	CheckRequiredVsphereEnvVars()
	return true
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

// TODO: Change name?
var PROVIDER utils.Provider = &Provider{}
