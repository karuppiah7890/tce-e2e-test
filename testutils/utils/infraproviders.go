package utils

import (
	"context"

	"github.com/karuppiah7890/tce-e2e-test/testutils/tanzu"
)

// Question: Move this to a package named infrastructure?
// Say something like providers or infra providers? Or it looks weird when calling utils.Provider from a caller perspective

// TODO: Change name?
type Provider interface {
	Name() string
	Init() error
	RequiredEnvVars() []string
	PreClusterCreationTasks(clusterName string, clusterType ClusterType) error
	CleanupCluster(ctx context.Context, clusterName string) error
	GetTanzuConfig(clusterName string) tanzu.TanzuConfig
}
