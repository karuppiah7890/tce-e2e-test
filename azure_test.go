package e2e

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/karuppiah7890/tce-e2e-test/testutils/azure"
	"github.com/karuppiah7890/tce-e2e-test/testutils/clirunner"
	"github.com/karuppiah7890/tce-e2e-test/testutils/docker"
	"github.com/karuppiah7890/tce-e2e-test/testutils/kubeclient"
	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
	"github.com/karuppiah7890/tce-e2e-test/testutils/tanzu"
	"github.com/karuppiah7890/tce-e2e-test/testutils/utils"

	capzv1beta1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
)

func TestAzureManagementAndWorkloadCluster(t *testing.T) {
	provider := utils.AZURE
	log.InitLogger("azure-mgmt-wkld-e2e")

	checks()

	azureTestSecrets := azure.ExtractAzureTestSecretsFromEnvVars()

	cred := azure.Login()

	clusterNameSuffix := time.Now().Unix()
	managementClusterName := fmt.Sprintf("test-mgmt-%d", clusterNameSuffix)
	workloadClusterName := fmt.Sprintf("test-wkld-%d", clusterNameSuffix)

	azureMarketplaceImageInfoForManagementCluster := getAzureMarketplaceImageInfoForCluster(managementClusterName, utils.ManagementClusterType)

	// TODO: make the below function return an error and handle the error to log and exit?
	azure.AcceptAzureImageLicenses(azureTestSecrets.SubscriptionID, cred, azureMarketplaceImageInfoForManagementCluster...)

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

		err := tanzu.CollectManagementClusterDiagnostics(managementClusterName)
		if err != nil {
			log.Errorf("error while collecting diagnostics of management cluster: %v", err)
		}

		err = utils.CleanupDockerBootstrapCluster(managementClusterName)
		if err != nil {
			log.Errorf("error while cleaning up docker bootstrap cluster of the management cluster: %v", err)
		}

		err = kubeclient.DeleteContext(kubeConfigPath, managementClusterKubeContext)
		if err != nil {
			log.Errorf("error while deleting kube context %s at kubeconfig path: %v", managementClusterKubeContext, err)
		}

		// TODO: Move this to a function named as cleanup azure cluster?
		err = azure.DeleteResourceGroup(context.TODO(), managementClusterName, azureTestSecrets.SubscriptionID, cred)
		if err != nil {
			log.Errorf("error while cleaning up azure resource group of the management cluster which has all the management cluster resources: %v", err)
		}

		log.Fatal("Summary: error while running management cluster: %v", runManagementClusterErr)
	}

	// TODO: check if management cluster is running by doing something similar to
	// `tanzu management-cluster get | grep "${MANAGEMENT_CLUSTER_NAME}" | grep running`

	// TODO: Handle errors
	utils.GetClusterKubeConfig(managementClusterName, provider, utils.ManagementClusterType)

	log.Infof("Management Cluster %s Information: ", managementClusterName)
	err = utils.PrintClusterInformation(kubeConfigPath, managementClusterKubeContext)
	if err != nil {
		// Should we panic here and stop?
		log.Errorf("error while printing management cluster information: %v", err)
	}

	azureMarketplaceImageInfoForWorkloadCluster := getAzureMarketplaceImageInfoForCluster(workloadClusterName, utils.WorkloadClusterType)

	// TODO: make the below function return an error and handle the error to log and exit?
	azure.AcceptAzureImageLicenses(azureTestSecrets.SubscriptionID, cred, azureMarketplaceImageInfoForWorkloadCluster...)

	workloadClusterKubeContext := utils.GetKubeContextForTanzuCluster(workloadClusterName)

	err = utils.RunCluster(workloadClusterName, provider, utils.WorkloadClusterType)
	if err != nil {
		runWorkloadClusterErr := err
		log.Errorf("error while running workload cluster: %v", runWorkloadClusterErr)

		err := tanzu.CollectManagementClusterAndWorkloadClusterDiagnostics(managementClusterName, workloadClusterName, provider)
		if err != nil {
			log.Errorf("error while collecting diagnostics of management cluster and workload cluster: %v", err)
		}

		// TODO: Move this to a function named as cleanup azure cluster?
		err = azure.DeleteResourceGroup(context.TODO(), managementClusterName, azureTestSecrets.SubscriptionID, cred)
		if err != nil {
			log.Errorf("error while cleaning up azure resource group of the management cluster which has all the management cluster resources: %v", err)
		}

		err = kubeclient.DeleteContext(kubeConfigPath, managementClusterKubeContext)
		if err != nil {
			log.Errorf("error while deleting kube context %s at kubeconfig path: %v", managementClusterKubeContext, err)
		}

		// TODO: Move this to a function named as cleanup azure cluster?
		err = azure.DeleteResourceGroup(context.TODO(), workloadClusterName, azureTestSecrets.SubscriptionID, cred)
		if err != nil {
			log.Errorf("error while cleaning up azure resource group of the workload cluster which has all the workload cluster resources: %v", err)
		}

		err = kubeclient.DeleteContext(kubeConfigPath, workloadClusterKubeContext)
		if err != nil {
			log.Errorf("error while deleting kube context %s at kubeconfig path: %v", managementClusterKubeContext, err)
		}

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

func getAzureMarketplaceImageInfoForCluster(clusterName string, clusterType utils.ClusterType) []*capzv1beta1.AzureMarketplaceImage {
	var clusterCreateDryRunOutputBuffer bytes.Buffer

	envVars := tanzu.TanzuConfigToEnvVars(tanzu.TanzuAzureConfig(clusterName))
	exitCode, err := clirunner.Run(clirunner.Cmd{
		Name: "tanzu",
		Args: []string{
			clusterType.TanzuCommand(),
			"create",
			clusterName,
			"--dry-run",
			// TODO: Should we add verbosity flag and value by default? or
			// let the user define the verbosity when running the tests maybe?
			// "-v",
			// "10",
		},
		Env: append(os.Environ(), envVars...),
		// TODO: Do we really want to output to log.InfoWriter ? Is this
		// data necessary in the logs? This data will contain secrets but for now we haven't masked secrets
		// in logs, also, even if we mask secrets, is this data useful and necessary?
		// The data in log can help development but that's all
		Stdout: &clusterCreateDryRunOutputBuffer,
		// TODO: Should we log standard errors as errors in the log? Because tanzu prints other information also
		// to standard error, which are kind of like information, apart from actual errors, so showing
		// everything as error is misleading. Gotta think what to do about this. The main problem is
		// console has only standard output and standard error, and tanzu is using standard output only for
		// giving output for things like --dry-run when it needs to print yaml content, but everything else
		// is printed to standard error
		// TODO: Do we really want to output to log.ErrorWriter ? Is this
		// data necessary in the logs? This data will contain secrets but for now we haven't masked secrets
		// in logs, also, even if we mask secrets, is this data useful and necessary?
		// The data in log can help development and also
		// during actual runs to check if there are any errors from the command, hmm
		Stderr: log.ErrorWriter,
	})

	if err != nil {
		log.Fatalf("Error occurred while running %v dry run. Exit code: %v. Error: %v", clusterName, exitCode, err)
	}

	clusterCreateDryRunOutput, err := io.ReadAll(&clusterCreateDryRunOutputBuffer)
	if err != nil {
		// TODO: Should we print the whole command as part of the error?
		log.Fatalf("Error occurred while reading output of %v create dry run: %v", clusterName, err)
	}

	objects := azure.ParseK8sYamlAndFetchAzureMachineTemplates(clusterCreateDryRunOutput)

	marketplaces := []*capzv1beta1.AzureMarketplaceImage{}

	for _, object := range objects {
		azureMachineTemplate, ok := object.(*capzv1beta1.AzureMachineTemplate)
		if !ok {
			log.Fatalf("Error occurred while parsing output of %v create dry run", clusterName)
		}

		marketplaces = append(marketplaces, azureMachineTemplate.Spec.Template.Spec.Image.Marketplace)
	}

	return marketplaces
}

func checks() {
	utils.CheckTanzuCLIInstallation()

	utils.CheckTanzuClusterCLIPluginInstallation(utils.ManagementClusterType)

	utils.CheckTanzuClusterCLIPluginInstallation(utils.WorkloadClusterType)

	docker.CheckDockerInstallation()

	utils.CheckKubectlCLIInstallation()

	utils.PlatformSupportCheck()
}
