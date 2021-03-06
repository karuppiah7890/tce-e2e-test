package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/karuppiah7890/tce-e2e-test/testutils"
	"github.com/karuppiah7890/tce-e2e-test/testutils/clirunner"
	"github.com/karuppiah7890/tce-e2e-test/testutils/kubeclient"
	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
	"github.com/karuppiah7890/tce-e2e-test/testutils/platforms"
	"github.com/karuppiah7890/tce-e2e-test/testutils/tanzu"
	"github.com/karuppiah7890/tce-e2e-test/testutils/tce"
)

// TODO: Further move the functions to specifics file/libs accordingly

type ClusterType struct {
	Name string
}

type WorkloadCluster struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}
type WorkloadClusters []WorkloadCluster

var ManagementClusterType = ClusterType{Name: "management-cluster"}
var WorkloadClusterType = ClusterType{Name: "cluster"}

func (clusterType ClusterType) TanzuCommand() string {
	return clusterType.Name
}

func CheckTanzuCLIInstallation() error {
	log.Info("Checking tanzu CLI installation")
	path, err := exec.LookPath("tanzu")
	if err != nil {
		log.Fatalf("tanzu CLI is not installed")
		return err
	}
	log.Infof("tanzu CLI is available at path: %s", path)
	return nil
}

func CheckKubectlCLIInstallation() {
	log.Info("Checking kubectl CLI installation")

	path, err := exec.LookPath("kubectl")
	if err != nil {
		log.Fatalf("kubectl CLI is not installed")
	}
	log.Infof("kubectl CLI is available at path: %s\n", path)
}

func CheckTanzuClusterCLIPluginInstallation(clusterType ClusterType) {
	log.Info("Checking tanzu management cluster plugin CLI installation")

	// TODO: Check for errors and return error?
	// TODO: Parse version and show warning if version is newer than what's tested by the devs while writing test
	// Refer - https://github.com/karuppiah7890/tce-e2e-test/issues/1#issuecomment-1094172278
	exitCode, err := clirunner.Run(clirunner.Cmd{
		Name: "tanzu",
		Args: []string{
			clusterType.TanzuCommand(),
			"version",
		},
		Stdout: log.InfoWriter,
		// TODO: Should we log standard errors as errors in the log? Because tanzu prints other information also
		// to standard error, which are kind of like information, apart from actual errors, so showing
		// everything as error is misleading. Gotta think what to do about this. The main problem is
		// console has only standard output and standard error, and tanzu is using standard output only for
		// giving output for things like --dry-run when it needs to print yaml content, but everything else
		// is printed to standard error
		Stderr: log.ErrorWriter,
	})

	if err != nil {
		log.Fatalf("Error occurred while checking management cluster CLI plugin installation. Exit code: %v. Error: %v", exitCode, err)
	}
}

func GetClusterNodes(kubeConfigPath string, kubeContext string) ([]string, error) {
	nodesName := []string{}
	client, err := kubeclient.GetKubeClient(kubeConfigPath, kubeContext)
	if err != nil {
		return nil, fmt.Errorf("error getting kube client: %v", err)
	}
	nodes, err := client.GetAllNodes()
	if err != nil {
		return nil, fmt.Errorf("error getting all nodes: %v", err)
	}

	for _, node := range nodes.Items {
		// TODO: There is some issue here, node.Status.Phase gives empty string I think
		log.Infof("%s\t%s", node.Name, node.Status.Phase)
		nodesName = append(nodesName, node.Name)
	}
	return nodesName, nil
}

func listWorkloadClusters() (WorkloadClusters, error) {
	var workloadClusters WorkloadClusters

	var clusterListOutput bytes.Buffer

	exitCode, err := clirunner.Run(clirunner.Cmd{
		Name: "tanzu",
		Args: []string{
			"cluster",
			"list",
			"-o",
			"json",
		},
		Env:    os.Environ(),
		Stdout: &clusterListOutput,
		// TODO: Should we log standard errors as errors in the log? Because tanzu prints other information also
		// to standard error, which are kind of like information, apart from actual errors, so showing
		// everything as error is misleading. Gotta think what to do about this. The main problem is
		// console has only standard output and standard error, and tanzu is using standard output only for
		// giving output for things like --dry-run when it needs to print yaml content, but everything else
		// is printed to standard error
		// TODO: Do we really want to output to log.ErrorWriter ? Is this
		// data necessary in the logs? This function will be called
		// a lot of times. The data in log can help development and also
		// during actual runs to check if there are any errors from the command, hmm
		Stderr: log.ErrorWriter,
	})

	if err != nil {
		return nil, fmt.Errorf("error occurred while listing workload clusters. Exit code: %v. Error: %v", exitCode, err)
	}

	// TODO: Parse JSON output from the command.
	// Check if the workload cluster name exists in the list of workload clusters.
	// Ideally, there should only be one or zero workload clusters. But let's not
	// think too much on that, for example, someone could create a separate workload
	// cluster in the meantime while the first one was being created and verified.
	// This could be done manually from their local machine to test stuff etc

	err = json.NewDecoder(&clusterListOutput).Decode(&workloadClusters)
	if err != nil {
		return nil, fmt.Errorf("error occurred while decoding JSON containing list of workload clusters. Exit code: %v. Error: %v", exitCode, err)
	}

	return workloadClusters, nil
}

func PlatformSupportCheck() {
	if runtime.GOOS == platforms.WINDOWS {
		log.Warn("Warning: This test has been tested only on Linux and Mac OS till now. Support for Windows has not been tested, so it's experimental and not guaranteed to work!")
	}
}

// TODO: Why pass provider if it's only for vSphere?
func UpdateVars(provider string, clusterType ClusterType) {
	if provider == "vsphere" {
		if clusterType == ManagementClusterType {
			os.Setenv("VSPHERE_CONTROL_PLANE_ENDPOINT", os.Getenv("VSPHERE_MANAGEMENT_CLUSTER_ENDPOINT"))
		} else if clusterType == WorkloadClusterType {
			os.Setenv("VSPHERE_CONTROL_PLANE_ENDPOINT", os.Getenv("VSPHERE_WORKLOAD_CLUSTER_ENDPOINT"))
		}
	}

}

func ManagementClusterCreationFailureTasks(ctx context.Context, r ClusterTestRunner, managementClusterName, kubeConfigPath, managementClusterKubeContext string, provider Provider) {
	err := r.CollectManagementClusterDiagnostics(managementClusterName)
	if err != nil {
		log.Errorf("error while collecting diagnostics of management cluster: %v", err)
	}

	err = r.CleanupDockerBootstrapCluster(managementClusterName)
	if err != nil {
		log.Errorf("error while cleaning up docker bootstrap cluster of the management cluster: %v", err)
	}

	err = r.DeleteContext(kubeConfigPath, managementClusterKubeContext)
	if err != nil {
		log.Errorf("error while deleting kube context %s at kubeconfig path: %v", managementClusterKubeContext, err)
	}

	err = provider.CleanupCluster(ctx, managementClusterName)
	if err != nil {
		log.Errorf("error while cleaning up the management cluster: %v", err)
	}
}

func WorkloadClusterCreationFailureTasks(ctx context.Context, r ClusterTestRunner, managementClusterName, workloadClusterName, kubeConfigPath, managementClusterKubeContext, workloadClusterKubeContext string, provider Provider) {
	err := r.CollectManagementClusterAndWorkloadClusterDiagnostics(managementClusterName, workloadClusterName, provider.Name())
	if err != nil {
		log.Errorf("error while collecting diagnostics of management cluster and workload cluster: %v", err)
	}

	err = provider.CleanupCluster(ctx, managementClusterName)
	if err != nil {
		log.Errorf("error while cleaning up the management cluster: %v", err)
	}

	err = r.DeleteContext(kubeConfigPath, managementClusterKubeContext)
	if err != nil {
		log.Errorf("error while deleting kube context %s at kubeconfig path: %v", managementClusterKubeContext, err)
	}

	err = provider.CleanupCluster(ctx, workloadClusterName)
	if err != nil {
		log.Errorf("error while cleaning up the workload cluster: %v", err)
	}

	err = r.DeleteContext(kubeConfigPath, workloadClusterKubeContext)
	if err != nil {
		log.Errorf("error while deleting kube context %s at kubeconfig path: %v", managementClusterKubeContext, err)
	}
}

func CheckRequiredEnvVars(provider Provider) error {
	requiredEnvVars := provider.RequiredEnvVars()
	errs := testutils.CheckRequiredEnvVars(requiredEnvVars)

	if len(errs) != 0 {
		fmt.Printf("abcdjjdjdksdm %v", errs)
		return fmt.Errorf("%v", errs)
	}

	return nil
}

func RunProviderTest(provider Provider, r ClusterTestRunner, packageDetails tce.Package) error {
	// Setup
	setupEnv(provider, r)
	// Setup Function complete
	managementClusterName, workloadClusterName := r.GetRandomClusterNames()

	// createManagementCluster function start
	createManagementCluster(provider, r, managementClusterName)
	// createManagementCluster Complete

	// Create Wkld Cluster Start
	createWorkloadCluster(provider, r, managementClusterName, workloadClusterName)
	// Create Wkld Cluster complete

	// package Code
	runPackageTest(r, packageDetails, workloadClusterName)
	// Package Code complete

	// TODO: Consider testing one basic package or we can do this separately or have
	// a feature flag to test it when needed and skip it when not needed.
	// This will give us an idea of how testing packages looks like and give an example
	// to TCE package owners

	// TODO: Handle errors during cluster deletion
	// and cleanup management cluster and then cleanup workload cluster

	// delete Wkld Cluster start
	deleteWorkloadCluster(provider, r, workloadClusterName, managementClusterName)
	// Delete wkld cluster complete
	// TODO: Handle errors during cluster deletion
	// and cleanup management cluster
	// Delete mgmt cluster start
	deleteManagementCluster(provider, r, managementClusterName)
	// Delete mgmt cluster complete
	return nil
}

func setupEnv(provider Provider, r ClusterTestRunner) {
	r.RunChecks()
	err := CheckRequiredEnvVars(provider)
	if err != nil {
		log.Errorf("errors while checking required environment variables: %v", err)
	}
	provider.Init()
}
func createManagementCluster(provider Provider, r ClusterTestRunner, managementClusterName string) {
	err := provider.PreClusterCreationTasks(managementClusterName, ManagementClusterType)
	if err != nil {
		log.Errorf("error while executing pre-cluster creation tasks for %v cluster: %v", managementClusterName, err)
	}

	managementClusterKubeContext := r.GetKubeContextForTanzuCluster(managementClusterName)
	kubeConfigPath, err := r.GetKubeConfigPath()
	if err != nil {
		log.Errorf("error while getting kubeconfig path: %v", err)
	}

	err = r.RunCluster(managementClusterName, provider, ManagementClusterType)
	if err != nil {
		runManagementClusterErr := err
		log.Errorf("error while running management cluster: %v", runManagementClusterErr)
		ManagementClusterCreationFailureTasks(context.TODO(), r, managementClusterName, kubeConfigPath, managementClusterKubeContext, provider)
		log.Errorf("error while running management cluster: %v", runManagementClusterErr)
	}

	// TODO: Handle errors
	r.GetClusterKubeConfig(managementClusterName, provider, ManagementClusterType)

	log.Infof("Management Cluster %s Information: ", managementClusterName)
	err = r.PrintClusterInformation(kubeConfigPath, managementClusterKubeContext)
	if err != nil {
		// Should we panic here and stop?
		log.Errorf("error while printing management cluster information: %v", err)
	}
}

func createWorkloadCluster(provider Provider, r ClusterTestRunner, managementClusterName, workloadClusterName string) {
	err := provider.PreClusterCreationTasks(workloadClusterName, WorkloadClusterType)
	if err != nil {
		log.Errorf("error while executing pre-cluster creation tasks for %v cluster: %v", err, workloadClusterName)
	}

	workloadClusterKubeContext := r.GetKubeContextForTanzuCluster(workloadClusterName)
	kubeConfigPath, err := r.GetKubeConfigPath()
	managementClusterKubeContext := r.GetKubeContextForTanzuCluster(managementClusterName)
	err = r.RunCluster(workloadClusterName, provider, WorkloadClusterType)
	if err != nil {
		runWorkloadClusterErr := err
		log.Errorf("error while running workload cluster: %v", runWorkloadClusterErr)

		WorkloadClusterCreationFailureTasks(context.TODO(), r, managementClusterName, workloadClusterName, kubeConfigPath, managementClusterKubeContext, workloadClusterKubeContext, provider)

		log.Errorf("error while running workload cluster: %v", runWorkloadClusterErr)
	}

	r.CheckWorkloadClusterIsRunning(workloadClusterName)

	// TODO: Handle errors
	r.GetClusterKubeConfig(workloadClusterName, provider, WorkloadClusterType)

	log.Infof("Workload Cluster %s Information: ", workloadClusterName)
	err = r.PrintClusterInformation(kubeConfigPath, workloadClusterKubeContext)
	if err != nil {
		// Should we panic here and stop?
		log.Errorf("error while printing workload cluster information: %v", err)
	}
}

func runPackageTest(r ClusterTestRunner, packageDetails tce.Package, workloadClusterName string) {
	workloadClusterKubeContext := r.GetKubeContextForTanzuCluster(workloadClusterName)
	if packageDetails.Name != "" {
		err := tce.PackageE2Etest(packageDetails, workloadClusterKubeContext)
		if err != nil {
			// Should we panic here and stop?
			log.Errorf("error while running e2e test for %v: %v", packageDetails.Name, err)
		}
	}
}

func deleteWorkloadCluster(provider Provider, r ClusterTestRunner, workloadClusterName, managementClusterName string) {
	err := r.DeleteCluster(workloadClusterName, provider, WorkloadClusterType)
	if err != nil {
		log.Errorf("error while deleting workload cluster: %v", err)

		err := tanzu.CollectManagementClusterAndWorkloadClusterDiagnostics(managementClusterName, workloadClusterName, provider.Name())
		if err != nil {
			log.Errorf("error while collecting diagnostics of management cluster and workload cluster: %v", err)
		}

		log.Errorf("error while deleting workload cluster: %v", err)
	}

	// TODO: Handle errors during waiting for cluster deletion.
	// We could retry in some cases, to just list the workload clusters.
	// If all retries fail, cleanup management cluster and then cleanup workload cluster
	err = r.WaitForWorkloadClusterDeletion(workloadClusterName)
	if err != nil {
		log.Errorf("error while waiting for workload cluster deletion: %v", err)
	}
	kubeConfigPath, err := r.GetKubeConfigPath()
	workloadClusterKubeContext := r.GetKubeContextForTanzuCluster(workloadClusterName)
	managementClusterKubeContext := r.GetKubeContextForTanzuCluster(managementClusterName)
	err = kubeclient.DeleteContext(kubeConfigPath, workloadClusterKubeContext)
	if err != nil {
		log.Errorf("error while deleting kube context %s at kubeconfig path: %v", managementClusterKubeContext, err)
	}

}

func deleteManagementCluster(provider Provider, r ClusterTestRunner, managementClusterName string) {
	err := r.DeleteCluster(managementClusterName, provider, ManagementClusterType)
	if err != nil {
		log.Errorf("error while deleting management cluster: %v", err)
		err := tanzu.CollectManagementClusterDiagnostics(managementClusterName)
		if err != nil {
			log.Errorf("error while collecting diagnostics of management cluster: %v", err)
		}
		log.Errorf("error while deleting management cluster: %v", err)
	}
}
