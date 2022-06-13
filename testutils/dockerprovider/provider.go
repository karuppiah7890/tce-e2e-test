package dockerprovider

import (
	"context"

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

// TODO: Change name?
var PROVIDER utils.Provider = &Provider{}
