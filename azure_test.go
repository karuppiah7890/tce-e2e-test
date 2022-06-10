package e2e

import (
	"testing"

	"github.com/karuppiah7890/tce-e2e-test/testutils/azure"
	"github.com/karuppiah7890/tce-e2e-test/testutils/kubeclient"
	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
	"github.com/karuppiah7890/tce-e2e-test/testutils/tanzu"
	"github.com/karuppiah7890/tce-e2e-test/testutils/utils"
)

func TestAzureManagementAndWorkloadCluster(t *testing.T) {
	provider := utils.AZURE
	log.InitLogger("azure-mgmt-wkld-e2e")

	utils.RunChecks()

	azure.PROVIDER.CheckRequiredEnvVars()

	azure.PROVIDER.Init()

	managementClusterName, workloadClusterName := utils.GetRandomClusterNames()

	azure.PROVIDER.PreClusterCreationTasks(managementClusterName, utils.ManagementClusterType)

	managementClusterKubeContext := utils.GetKubeContextForTanzuCluster(managementClusterName)
	kubeConfigPath, err := utils.GetKubeConfigPath()
	if err != nil {
		// TODO: Should we continue here for any reason without stopping? As kubeconfig path is not available
		log.Fatalf("error while getting kubeconfig path: %v", err)
	}

	err = utils.RunCluster(managementClusterName, provider, utils.ManagementClusterType)
	if err != nil {
		runManagementClusterErr := err
		log.Errorf("error while running management cluster: %v", runManagementClusterErr)
		utils.ManagementClusterCreationFailureTasks(managementClusterName, kubeConfigPath, managementClusterKubeContext, azure.PROVIDER)
		log.Fatal("Summary: error while running management cluster: %v", runManagementClusterErr)
	}

	// TODO: Handle errors
	utils.GetClusterKubeConfig(managementClusterName, provider, utils.ManagementClusterType)

	log.Infof("Management Cluster %s Information: ", managementClusterName)
	err = utils.PrintClusterInformation(kubeConfigPath, managementClusterKubeContext)
	if err != nil {
		// Should we panic here and stop?
		log.Errorf("error while printing management cluster information: %v", err)
	}

	azure.PROVIDER.PreClusterCreationTasks(workloadClusterName, utils.WorkloadClusterType)

	workloadClusterKubeContext := utils.GetKubeContextForTanzuCluster(workloadClusterName)

	err = utils.RunCluster(workloadClusterName, provider, utils.WorkloadClusterType)
	if err != nil {
		runWorkloadClusterErr := err
		log.Errorf("error while running workload cluster: %v", runWorkloadClusterErr)

		utils.WorkloadClusterFailureTasks(managementClusterName, workloadClusterName, kubeConfigPath, managementClusterKubeContext, workloadClusterKubeContext, azure.PROVIDER)

		log.Fatal("error while running workload cluster: %v", runWorkloadClusterErr)
	}

	utils.CheckWorkloadClusterIsRunning(workloadClusterName)

	// TODO: Handle errors
	utils.GetClusterKubeConfig(workloadClusterName, provider, utils.WorkloadClusterType)

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

	err = kubeclient.DeleteContext(kubeConfigPath, workloadClusterKubeContext)
	if err != nil {
		log.Errorf("error while deleting kube context %s at kubeconfig path: %v", managementClusterKubeContext, err)
	}

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
