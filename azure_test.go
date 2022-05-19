package e2e

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/karuppiah7890/tce-e2e-test/testutils/azure"
	"github.com/karuppiah7890/tce-e2e-test/testutils/clirunner"
	"github.com/karuppiah7890/tce-e2e-test/testutils/docker"
	"github.com/karuppiah7890/tce-e2e-test/testutils/kubeclient"
	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
	"github.com/karuppiah7890/tce-e2e-test/testutils/platforms"
	"github.com/karuppiah7890/tce-e2e-test/testutils/tanzu"
	"github.com/karuppiah7890/tce-e2e-test/testutils/utils"

	capzv1beta1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
)

// TODO: Make region as environment variable

// TODO: Consider making all as environment variables. Hard coded values in test code can be default
// We can pass env vars to override stuff

// TODO: Use the utils package in testutils package.
// Check vsphere_test.go for usage of utils package

func TestAzureManagementAndWorkloadCluster(t *testing.T) {
	// TODO: Convert magic strings like "azure" to constants
	provider := "azure"
	log.InitLogger("azure-mgmt-wkld-e2e")

	// Ensure TCE/TF is installed - check TCE installation or install it if not present. Or do it prior to the test run.
	// check if tanzu is installed
	utils.CheckTanzuCLIInstallation()

	// Ensure management and workload cluster plugins are present.
	// check if management cluster plugin is present
	ManagementClusterType := utils.ClusterType{Name: "management-cluster"}
	utils.CheckTanzuClusterCLIPluginInstallation(ManagementClusterType)

	// check if workload cluster plugin is present
	WorkloadClusterType := utils.ClusterType{Name: "cluster"}
	utils.CheckTanzuClusterCLIPluginInstallation(WorkloadClusterType)

	// check if docker is installed. This is required by tanzu CLI I think, both docker client CLI and docker daemon
	docker.CheckDockerInstallation()
	// check if kubectl is installed. This is required by tanzu CLI to apply using kubectl apply to create cluster
	utils.CheckKubectlCLIInstallation()

	if runtime.GOOS == platforms.WINDOWS {
		log.Warn("Warning: This test has been tested only on Linux and Mac OS till now. Support for Windows has not been tested, so it's experimental and not guranteed to work!")
	}

	// TODO: Ensure package plugin is present in case package tests are gonna be executed.
	azureTestSecrets := azure.ExtractAzureTestSecretsFromEnvVars()

	// Have different log levels - none/minimal, error, info, debug etc, so that we can accordingly use those in the E2E test

	cred := azure.Login()

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

	azureMarketplaceImageInfoForManagementCluster := getAzureMarketplaceImageInfoForClusters(managementClusterName, ManagementClusterType)

	// TODO: make the below function return an error and handle the error to log and exit?
	azure.AcceptAzureImageLicenses(azureTestSecrets.SubscriptionID, cred, azureMarketplaceImageInfoForManagementCluster...)

	managementClusterKubeContext := utils.GetKubeContextForTanzuCluster(managementClusterName)
	kubeConfigPath, err := utils.GetKubeConfigPath()
	if err != nil {
		// TODO: Should we continue here for any reason without stopping? As kubeconfig path is not available
		log.Fatalf("error while getting kubeconfig path: %v", err)
	}

	// TODO: Handle errors during deployment
	// and cleanup management cluster
	err = utils.RunCluster(managementClusterName, provider, ManagementClusterType)
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
	utils.GetClusterKubeConfig(managementClusterName, provider, ManagementClusterType)

	log.Infof("Management Cluster %s Information: ", managementClusterName)
	err = utils.PrintClusterInformation(kubeConfigPath, managementClusterKubeContext)
	if err != nil {
		// Should we panic here and stop?
		log.Errorf("error while printing management cluster information: %v", err)
	}

	azureMarketplaceImageInfoForWorkloadCluster := getAzureMarketplaceImageInfoForClusters(workloadClusterName, WorkloadClusterType)

	// TODO: make the below function return an error and handle the error to log and exit?
	azure.AcceptAzureImageLicenses(azureTestSecrets.SubscriptionID, cred, azureMarketplaceImageInfoForWorkloadCluster...)

	workloadClusterKubeContext := utils.GetKubeContextForTanzuCluster(workloadClusterName)

	// TODO: Handle errors during deployment
	// and cleanup management cluster and then cleanup workload cluster
	err = utils.RunCluster(workloadClusterName, provider, WorkloadClusterType)
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
	utils.GetClusterKubeConfig(workloadClusterName, provider, WorkloadClusterType)

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
	err = utils.DeleteCluster(workloadClusterName, provider, WorkloadClusterType)
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
	err = utils.DeleteCluster(managementClusterName, provider, ManagementClusterType)
	if err != nil {
		log.Errorf("error while deleting management cluster: %v", err)

		err := tanzu.CollectManagementClusterDiagnostics(managementClusterName)
		if err != nil {
			log.Errorf("error while collecting diagnostics of management cluster: %v", err)
		}

		log.Fatal("error while deleting management cluster: %v", err)
	}
}

func getAzureMarketplaceImageInfoForClusters(clusterName string, clusterType utils.ClusterType) []*capzv1beta1.AzureMarketplaceImage {
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
