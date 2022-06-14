package dockerprovider

import (
	"context"

	"github.com/karuppiah7890/tce-e2e-test/testutils/tanzu"
	"github.com/karuppiah7890/tce-e2e-test/testutils/utils"
)

// TODO: Change name?
type Provider struct{}

func (provider *Provider) CheckRequiredEnvVars() bool {
	return true
}

func (provider *Provider) Name() string {
	return "docker"
}

func (provider *Provider) Init() error {
	return nil
}

func (provider *Provider) PreClusterCreationTasks(clusterName string, clusterType utils.ClusterType) error {
	return nil
}

func (provider *Provider) CleanupCluster(ctx context.Context, clusterName string) error {
	// TODO: Implement using docker library
	return nil
}

//TODO: Maybe make use of https://github.com/spf13/viper to set env vars and make some values as default and parameterised.
func (provider *Provider) GetTanzuConfig(clusterName string) tanzu.TanzuConfig {
	return tanzu.TanzuConfig{
		"CLUSTER_NAME":              clusterName,
		"INFRASTRUCTURE_PROVIDER":   provider.Name(),
		"CLUSTER_PLAN":              "dev",
		"OS_ARCH":                   "",
		"OS_NAME":                   "",
		"OS_VERSION":                "",
		"CLUSTER_CIDR":              "100.96.0.0/11",
		"SERVICE_CIDR":              "100.64.0.0/13",
		"ENABLE_CEIP_PARTICIPATION": "false",
		"ENABLE_MHC":                "true",
		"IDENTITY_MANAGEMENT_TYPE":  "none",
	}
}

// TODO: Change name?
var PROVIDER utils.Provider = &Provider{}
