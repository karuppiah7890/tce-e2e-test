package e2e

import (
	"fmt"
	"testing"
	"time"

	"github.com/karuppiah7890/tce-e2e-test/testutils/clirunner"
	"github.com/karuppiah7890/tce-e2e-test/testutils/utils"

	"github.com/karuppiah7890/tce-e2e-test/testutils/aws"
	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
	"github.com/karuppiah7890/tce-e2e-test/testutils/tanzu"
)

// TODO: Make region as environment variable

// TODO: Consider making all as environment variables. Hard coded values in test code can be default.
// We can pass env vars to override stuff

// TODO: Do we really need individual test files for each provider or can we parameterise it ? thoughts

func TestAwsManagementAndWorkloadCluster(t *testing.T) {
	log.InitLogger("aws-mgmt-wkld-e2e")
	provider := utils.AWS

	checks()

	aws.CheckRequiredAwsEnvVars()

	// Have different log levels - none/minimal, error, info, debug etc, so that we can accordingly use those in the E2E test

	// Create random names for management and workload clusters so that we can use them to name the test clusters we are going to
	// create. Ensure that these names are not already taken - check the resource group names to double check :) As Resource group name
	// is based on the cluster name
	// TODO: Create random names later, using random number or using short or long UUIDs.
	// TODO: Do we allow users to pass the cluster name for both clusters? We could. How do we take inputs? File? Env vars? Flags?
	clusterNameSuffix := time.Now().Unix()
	managementClusterName := fmt.Sprintf("test-mgmt-%d", clusterNameSuffix)
	workloadClusterName := fmt.Sprintf("test-wkld-%d", clusterNameSuffix)

	// TODO: Idea - if workload cluster and management cluster name are tied to a pipeline / workflow using
	// a unique ID, then we can use an external process to check clusters that are lying around and
	// check corresponding pipelines / workflows and if they are finished / done / cancelled, then we can
	// cleanup the cluster. We can also look at how we can add some sort of metadata (like labels, tags) to
	// the cluster or cluster resources using Tanzu to be able to do this instead of encoding the pipeline
	// metadata in the cluster name but that's a good idea too :)

	// To create Aws Cloud formation stack
	createCloudFormationStack()
	// TODO: Handle errors during deployment
	// and cleanup management cluster
	err := utils.RunCluster(managementClusterName, provider, utils.ManagementClusterType)
	if err != nil {
		log.Errorf("error while running management cluster: %v", err)

		err := tanzu.CollectManagementClusterDiagnostics(managementClusterName)
		if err != nil {
			log.Errorf("error while collecting diagnostics of management cluster: %v", err)
		}

		log.Fatal("error while running management cluster: %v", err)
	}

	// TODO: check if management cluster is running by doing something similar to
	// `tanzu management-cluster get | grep "${MANAGEMENT_CLUSTER_NAME}" | grep running`

	// TODO: Handle errors
	utils.GetClusterKubeConfig(managementClusterName, provider, utils.ManagementClusterType)

	kubeConfigPath, err := utils.GetKubeConfigPath()
	if err != nil {
		// Should we panic here and stop?
		log.Errorf("error while getting kubeconfig path: %v", err)
	}
	managementClusterKubeContext := utils.GetKubeContextForTanzuCluster(managementClusterName)

	log.Infof("Management Cluster %s Information: ", managementClusterName)
	err = utils.PrintClusterInformation(kubeConfigPath, managementClusterKubeContext)
	if err != nil {
		// Should we panic here and stop?
		log.Errorf("error while printing management cluster information: %v", err)
	}

	// TODO: Handle errors during deployment
	// and cleanup management cluster and then cleanup workload cluster
	// and cleanup management cluster and then cleanup workload cluster
	err = utils.RunCluster(workloadClusterName, provider, utils.WorkloadClusterType)
	if err != nil {
		log.Errorf("error while running workload cluster: %v", err)

		err := tanzu.CollectManagementClusterAndWorkloadClusterDiagnostics(managementClusterName, workloadClusterName, provider)
		if err != nil {
			log.Errorf("error while collecting diagnostics of management cluster and workload cluster: %v", err)
		}

		log.Fatal("error while running workload cluster: %v", err)
	}

	utils.CheckWorkloadClusterIsRunning(workloadClusterName)

	// TODO: Handle errors
	utils.GetClusterKubeConfig(workloadClusterName, provider, utils.WorkloadClusterType)

	workloadClusterKubeContext := utils.GetKubeContextForTanzuCluster(workloadClusterName)

	log.Infof("Workload Cluster %s Information: ", workloadClusterName)
	err = utils.PrintClusterInformation(kubeConfigPath, workloadClusterKubeContext)
	if err != nil {
		// Should we panic here and stop?
		log.Errorf("error while printing workload cluster information: %v", err)
	}

	// TODO: Consider testing one basic package or we can do this separately or have
	// a feature flag to test it when needed and skip it when not needed.
	// This will give us an idea of how testing packages looks like and give an example
	// to TCE package owners

	// TODO: Handle errors during cluster deletion
	// and cleanup management cluster and then cleanup workload cluster
	err = utils.DeleteCluster(workloadClusterName, provider, utils.WorkloadClusterType)
	if err != nil {
		log.Errorf("error while deleting workload cluster: %v", err)

		err := tanzu.CollectManagementClusterAndWorkloadClusterDiagnostics(managementClusterName, workloadClusterName, provider)
		if err != nil {
			log.Errorf("error while collecting diagnostics of management cluster and workload cluster: %v", err)
		}

		log.Fatal("error while deleting workload cluster: %v", err)
	}

	// TODO: Handle errors during waiting for cluster deletion.
	// We could retry in some cases, to just list the workload clusters.
	// If all retries fail, cleanup management cluster and then cleanup workload cluster
	utils.WaitForWorkloadClusterDeletion(workloadClusterName)

	// TODO: Cleanup workload cluster kube config data (cluster, user, context)
	// since tanzu cluster delete does not delete workload cluster kubeconfig entry

	// TODO: Handle errors during cluster deletion
	// and cleanup management cluster
	err = utils.DeleteCluster(managementClusterName, provider, utils.ManagementClusterType)
	if err != nil {
		log.Errorf("error while deleting management cluster: %v", err)

		err := tanzu.CollectManagementClusterDiagnostics(managementClusterName)
		if err != nil {
			log.Errorf("error while collecting diagnostics of management cluster: %v", err)
		}

		log.Fatal("error while deleting management cluster: %v", err)
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
	})

	if err != nil {
		log.Fatalf("Error occurred while creating Cloud formation stack  Exit code: %v. Error: %v", exitCode, err)
	}

}
