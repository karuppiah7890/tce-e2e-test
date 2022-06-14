package utils

import "context"

// Question: Move this to a package named infrastructure?
// Say something like providers or infra providers? Or it looks weird when calling utils.Provider from a caller perspective

// TODO: Change name?
type Provider interface {
	Name() string
	Init() error
	// TODO: Change CheckRequiredEnvVars to GetListOfRequiredEnvVars ? And do check in a common manner?
	CheckRequiredEnvVars() bool
	PreClusterCreationTasks(clusterName string, clusterType ClusterType) error
	CleanupCluster(ctx context.Context, clusterName string) error
}
