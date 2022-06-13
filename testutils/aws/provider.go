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
