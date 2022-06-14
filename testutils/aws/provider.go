package aws

import (
	"context"
	"os"

	"github.com/karuppiah7890/tce-e2e-test/testutils/clirunner"
	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
	"github.com/karuppiah7890/tce-e2e-test/testutils/utils"
)

// TODO: Change name?
type Provider struct {
	testSecrets TestSecrets
}

func (provider *Provider) CheckRequiredEnvVars() bool {
	CheckRequiredAwsEnvVars()
	return true
}

func (provider *Provider) Name() string {
	return "aws"
}

func (provider *Provider) Init() error {
	provider.testSecrets = ExtractAwsTestSecretsFromEnvVars()
	return nil
}

func (provider *Provider) PreClusterCreationTasks(clusterName string, clusterType utils.ClusterType) error {
	createCloudFormationStack()
	return nil
}

func (provider *Provider) CleanupCluster(ctx context.Context, clusterName string) error {
	// TODO: Implement using aws-nuke library
	return nil
}

//TODO: Maybe make use of https://github.com/spf13/viper to set env vars and make some values as default and parameterised.
func (provider *Provider) GetTanzuConfig(clusterName string) utils.TanzuConfig {
	return utils.TanzuConfig{
		"CLUSTER_NAME":               clusterName,
		"INFRASTRUCTURE_PROVIDER":    provider.Name(),
		"CLUSTER_PLAN":               "dev",
		"AWS_NODE_AZ":                "us-east-1a",
		"AWS_REGION":                 "us-east-1",
		"OS_ARCH":                    "amd64",
		"OS_NAME":                    "amazon",
		"OS_VERSION":                 "2",
		"CONTROL_PLANE_MACHINE_TYPE": "m5.xlarge",
		"NODE_MACHINE_TYPE":          "m5.xlarge",
		"AWS_PRIVATE_NODE_CIDR":      "10.0.16.0/20",
		"AWS_PUBLIC_NODE_CIDR":       "10.0.0.0/20",
		"AWS_VPC_CIDR":               "10.0.0.0/16",
		"CLUSTER_CIDR":               "100.96.0.0/11",
		"SERVICE_CIDR":               "100.64.0.0/13",
		"ENABLE_CEIP_PARTICIPATION":  "false",
		"ENABLE_MHC":                 "true",
		"BASTION_HOST_ENABLED":       "true",
		"IDENTITY_MANAGEMENT_TYPE":   "none",
	}
}

func createCloudFormationStack() {
	log.Info("Creating Cloud formation stack ")
	exitCode, err := clirunner.Run(clirunner.Cmd{
		Name: "tanzu",
		Args: []string{
			"management-cluster",
			"permissions",
			"aws",
			"set",
		},
		Stdout: log.InfoWriter,
		Stderr: log.ErrorWriter,
		Env:    os.Environ(),
	})

	if err != nil {
		log.Fatalf("Error occurred while creating Cloud formation stack  Exit code: %v. Error: %v", exitCode, err)
	}

}

// TODO: Change name?
var PROVIDER utils.Provider = &Provider{}
