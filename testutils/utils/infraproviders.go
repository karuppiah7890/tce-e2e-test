package utils

import "context"

// Question: Move this to a package named infrastructure?
// Say something like providers or infra providers? Or it looks when calling utils.AWS from a caller perspective

const AWS = "aws"
const VSPHERE = "vsphere"
const Docker = "docker"

// TODO: Change name?
type Provider interface {
	Name() string
	Init() error
	// TODO: Change CheckRequiredEnvVars to GetListOfRequiredEnvVars ? And do check in a common manner?
	CheckRequiredEnvVars() bool
	PreClusterCreationTasks(clusterName string, clusterType ClusterType) error
	CleanupCluster(ctx context.Context, clusterName string) error
}
