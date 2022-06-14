package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/karuppiah7890/tce-e2e-test/testutils/clirunner"
	"github.com/karuppiah7890/tce-e2e-test/testutils/docker"
	"github.com/karuppiah7890/tce-e2e-test/testutils/kubeclient"
	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
	"github.com/karuppiah7890/tce-e2e-test/testutils/platforms"
	"github.com/karuppiah7890/tce-e2e-test/testutils/tanzu"
	"k8s.io/client-go/util/homedir"
)

// TODO: Further move the functions to specifics file/libs accordingly

type TanzuConfig map[string]string

type ClusterType struct {
	Name string
}

type WorkloadCluster struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}
type WorkloadClusters []WorkloadCluster

type EnvVars []string

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

func GetKubeConfigPath() (string, error) {
	home := homedir.HomeDir()

	if home == "" {
		return "", fmt.Errorf("could not find home directory to get absolute path of kubeconfig")
	}

	return filepath.Join(home, ".kube", "config"), nil
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

func RunCluster(clusterName string, provider Provider, clusterType ClusterType) error {
	envVars := tanzuConfigToEnvVars(provider.GetTanzuConfig(clusterName))
	exitCode, err := clirunner.Run(clirunner.Cmd{
		Name: "tanzu",
		Args: []string{
			clusterType.TanzuCommand(),
			"create",
			clusterName,
			// TODO: Should we add verbosity flag and value by default? or
			// let the user define the verbosity when running the tests maybe?
			// "-v",
			// "10",
		},
		Env:    append(os.Environ(), envVars...),
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
		return fmt.Errorf("error occurred while deploying %v. exit code: %v. error: %v", clusterName, exitCode, err)
	}
	return nil
}

func GetClusterKubeConfig(clusterName string, provider Provider, clusterType ClusterType) {
	// TODO: Do we really need the secrets here?
	envVars := tanzuConfigToEnvVars(provider.GetTanzuConfig(clusterName))
	exitCode, err := clirunner.Run(clirunner.Cmd{
		// TODO: Replace magic strings like "tanzu", "management-cluster" etc
		Name: "tanzu",
		Args: []string{
			clusterType.TanzuCommand(),
			"kubeconfig",
			"get",
			clusterName,
			"--admin",
			// TODO: Should we add verbosity flag and value by default? or
			// let the user define the verbosity (eg 0-9) when running the tests maybe?
			// "-v",
			// "9",
		},
		Env:    append(os.Environ(), envVars...),
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
		log.Fatalf("Error occurred while getting %v kubeconfig. Exit code: %v. Error: %v", clusterName, exitCode, err)
	}
}

func GetKubeContextForTanzuCluster(clusterName string) string {
	return fmt.Sprintf("%s-admin@%s", clusterName, clusterName)
}

func PrintClusterInformation(kubeConfigPath string, kubeContext string) error {
	client, err := kubeclient.GetKubeClient(kubeConfigPath, kubeContext)
	if err != nil {
		return fmt.Errorf("error getting kube client: %v", err)
	}

	versionInfo, err := client.Discovery().ServerVersion()
	if err != nil {
		return fmt.Errorf("error getting kubernetes api server version: %v", err)
	}

	log.Infof("Kubernetes API server version is %s", versionInfo.String())

	// TODO: Should we get exact details as `kubectl get pod -A`? Showing age, restart count, how many containers are ready,
	// pod's phase (running) etc

	pods, err := client.GetAllPodsFromAllNamespaces()
	if err != nil {
		return fmt.Errorf("error getting all pods from all namespaces: %v", err)
	}

	// TODO: Should we check pods.RemainingItemCount value to see if it is 0 to ensure we have got all the pods?

	log.Info("Pod Name\tPod Namespace\tPod Phase")
	for _, pod := range pods.Items {
		// TODO: Use some library to print / format in some sort of table format? With proper spacing
		log.Infof("%s\t%s\t%s", pod.Name, pod.Namespace, pod.Status.Phase)
	}

	nodes, err := client.GetAllNodes()
	if err != nil {
		return fmt.Errorf("error getting all nodes: %v", err)
	}

	log.Info("\n\nNode Name\tNode Phase")
	for _, node := range nodes.Items {
		// TODO: There is some issue here, node.Status.Phase gives empty string I think
		log.Infof("%s\t%s", node.Name, node.Status.Phase)
	}

	return nil
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

func CheckWorkloadClusterIsRunning(workloadClusterName string) {
	// TODO: Should this function use a loop and wait (with timeout) for workload cluster to show up in the list
	// of workload clusters and be in running state? Or will Tanzu exit workload cluster creation
	// command only when workload cluster shows up in the list and is in running state? Gotta check
	workloadClusters := listWorkloadClusters()

	isClusterPresent := false
	clusterStatus := ""

	for _, workloadCluster := range workloadClusters {
		if workloadCluster.Name == workloadClusterName {
			isClusterPresent = true
			clusterStatus = workloadCluster.Status
		}
	}

	if !isClusterPresent {
		// Return errors for caller to handle maybe? Instead of abrupt stop?
		log.Fatalf("Workload cluster %s is not present in the list of workload clusters", workloadClusterName)
	}

	if clusterStatus != "running" {
		// Return errors for caller to handle maybe? Instead of abrupt stop?
		log.Fatalf("Workload cluster %s is not in running status, it is in %s status", workloadClusterName, clusterStatus)
	}

	log.Infof("Workload cluster %s is running successfully\n", workloadClusterName)
}

func WaitForWorkloadClusterDeletion(workloadClusterName string) {
	// TODO: Use timer for timeout and ticker for polling every few seconds
	// instead of using sleep
	for i := 0; i < 60; i++ {
		workloadClusters := listWorkloadClusters()

		isClusterPresent := false

		for _, workloadCluster := range workloadClusters {
			if workloadCluster.Name == workloadClusterName {
				isClusterPresent = true
			}
		}

		if isClusterPresent {
			log.Info("Waiting for workload cluster to get deleted")
		} else {
			log.Infof("Workload cluster %s successfully deleted\n", workloadClusterName)
			return
		}

		time.Sleep(10 * time.Second)
	}

	// TODO: maybe return error instead of fatal stop?
	log.Fatalf("Timed out waiting for workload cluster %s to get deleted", workloadClusterName)
}

func listWorkloadClusters() WorkloadClusters {
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
		// TODO: return error instead of fatal? So that the caller can retry if they want to or stop execution
		log.Fatalf("Error occurred while listing workload clusters. Exit code: %v. Error: %v", exitCode, err)
	}

	// TODO: Parse JSON output from the command.
	// Check if the workload cluster name exists in the list of workload clusters.
	// Ideally, there should only be one or zero workload clusters. But let's not
	// think too much on that, for example, someone could create a separate workload
	// cluster in the meantime while the first one was being created and verified.
	// This could be done manually from their local machine to test stuff etc

	json.NewDecoder(&clusterListOutput).Decode(&workloadClusters)

	return workloadClusters
}

func CleanupDockerBootstrapCluster(managementClusterName string) error {
	bootstrapClusterDockerContainerName, err := tanzu.GetBootstrapClusterDockerContainerNameForManagementCluster(managementClusterName)
	if err != nil {
		return fmt.Errorf("error getting bootstrap cluster docker container name for the management cluster %s: %v", managementClusterName, err)
	}

	err = docker.ForceRemoveRunningContainer(bootstrapClusterDockerContainerName)
	if err != nil {
		return fmt.Errorf("error force stopping and removing bootstrap cluster docker container name for the management cluster %s: %v", managementClusterName, err)
	}

	return nil
}

func DeleteCluster(clusterName string, provider Provider, clusterType ClusterType) error {
	// TODO: Do we really need the  secrets here?
	envVars := tanzuConfigToEnvVars(provider.GetTanzuConfig(clusterName))
	exitCode, err := clirunner.Run(clirunner.Cmd{
		Name: "tanzu",
		Args: []string{
			clusterType.TanzuCommand(),
			"delete",
			clusterName,
			"--yes",
			// TODO: Should we add verbosity flag and value by default? or
			// let the user define the verbosity when running the tests maybe?
			// "-v",
			// "10",
		},
		Env:    append(os.Environ(), envVars...),
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
		return fmt.Errorf("error occurred while deleting %v. exit code: %v. error: %v", clusterName, exitCode, err)
	}

	return nil
}

//TODO: Should we stick to env vars for cluster config or can we use yaml like tanzu cli consumes
func tanzuConfigToEnvVars(tanzuConfig TanzuConfig) EnvVars {
	envVars := make(EnvVars, 0, len(tanzuConfig))

	for key, value := range tanzuConfig {
		envVars = append(envVars, fmt.Sprintf("%s=%s", key, value))
	}

	return envVars
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

func RunChecks() {
	_ = CheckTanzuCLIInstallation()

	CheckTanzuClusterCLIPluginInstallation(ManagementClusterType)

	CheckTanzuClusterCLIPluginInstallation(WorkloadClusterType)

	docker.CheckDockerInstallation()

	CheckKubectlCLIInstallation()

	PlatformSupportCheck()
}

func GetRandomClusterNames() (string, string) {
	clusterNameSuffix := time.Now().Unix()
	managementClusterName := fmt.Sprintf("test-mgmt-%d", clusterNameSuffix)
	workloadClusterName := fmt.Sprintf("test-wkld-%d", clusterNameSuffix)
	log.Infof("Management Cluster Name : %s", managementClusterName)
	log.Infof("Workload Cluster Name : %s", workloadClusterName)
	return managementClusterName, workloadClusterName
}

func ManagementClusterCreationFailureTasks(ctx context.Context, managementClusterName, kubeConfigPath, managementClusterKubeContext string, provider Provider) {
	err := tanzu.CollectManagementClusterDiagnostics(managementClusterName)
	if err != nil {
		log.Errorf("error while collecting diagnostics of management cluster: %v", err)
	}

	err = CleanupDockerBootstrapCluster(managementClusterName)
	if err != nil {
		log.Errorf("error while cleaning up docker bootstrap cluster of the management cluster: %v", err)
	}

	err = kubeclient.DeleteContext(kubeConfigPath, managementClusterKubeContext)
	if err != nil {
		log.Errorf("error while deleting kube context %s at kubeconfig path: %v", managementClusterKubeContext, err)
	}

	err = provider.CleanupCluster(ctx, managementClusterName)
	if err != nil {
		log.Errorf("error while cleaning up the management cluster: %v", err)
	}
}

func WorkloadClusterCreationFailureTasks(ctx context.Context, managementClusterName, workloadClusterName, kubeConfigPath, managementClusterKubeContext, workloadClusterKubeContext string, provider Provider) {
	err := tanzu.CollectManagementClusterAndWorkloadClusterDiagnostics(managementClusterName, workloadClusterName, provider.Name())
	if err != nil {
		log.Errorf("error while collecting diagnostics of management cluster and workload cluster: %v", err)
	}

	err = provider.CleanupCluster(ctx, managementClusterName)
	if err != nil {
		log.Errorf("error while cleaning up the management cluster: %v", err)
	}

	err = kubeclient.DeleteContext(kubeConfigPath, managementClusterKubeContext)
	if err != nil {
		log.Errorf("error while deleting kube context %s at kubeconfig path: %v", managementClusterKubeContext, err)
	}

	err = provider.CleanupCluster(ctx, workloadClusterName)
	if err != nil {
		log.Errorf("error while cleaning up the workload cluster: %v", err)
	}

	err = kubeclient.DeleteContext(kubeConfigPath, workloadClusterKubeContext)
	if err != nil {
		log.Errorf("error while deleting kube context %s at kubeconfig path: %v", managementClusterKubeContext, err)
	}
}

func RunProviderTest(provider Provider) {
	RunChecks()

	provider.CheckRequiredEnvVars()

	provider.Init()

	managementClusterName, workloadClusterName := GetRandomClusterNames()

	provider.PreClusterCreationTasks(managementClusterName, ManagementClusterType)

	managementClusterKubeContext := GetKubeContextForTanzuCluster(managementClusterName)
	kubeConfigPath, err := GetKubeConfigPath()
	if err != nil {
		// TODO: Should we continue here for any reason without stopping? As kubeconfig path is not available
		log.Fatalf("error while getting kubeconfig path: %v", err)
	}

	err = RunCluster(managementClusterName, provider, ManagementClusterType)
	if err != nil {
		runManagementClusterErr := err
		log.Errorf("error while running management cluster: %v", runManagementClusterErr)
		ManagementClusterCreationFailureTasks(context.TODO(), managementClusterName, kubeConfigPath, managementClusterKubeContext, provider)
		log.Fatal("Summary: error while running management cluster: %v", runManagementClusterErr)
	}

	// TODO: Handle errors
	GetClusterKubeConfig(managementClusterName, provider, ManagementClusterType)

	log.Infof("Management Cluster %s Information: ", managementClusterName)
	err = PrintClusterInformation(kubeConfigPath, managementClusterKubeContext)
	if err != nil {
		// Should we panic here and stop?
		log.Errorf("error while printing management cluster information: %v", err)
	}

	provider.PreClusterCreationTasks(workloadClusterName, WorkloadClusterType)

	workloadClusterKubeContext := GetKubeContextForTanzuCluster(workloadClusterName)

	err = RunCluster(workloadClusterName, provider, WorkloadClusterType)
	if err != nil {
		runWorkloadClusterErr := err
		log.Errorf("error while running workload cluster: %v", runWorkloadClusterErr)

		WorkloadClusterCreationFailureTasks(context.TODO(), managementClusterName, workloadClusterName, kubeConfigPath, managementClusterKubeContext, workloadClusterKubeContext, provider)

		log.Fatal("error while running workload cluster: %v", runWorkloadClusterErr)
	}

	CheckWorkloadClusterIsRunning(workloadClusterName)

	// TODO: Handle errors
	GetClusterKubeConfig(workloadClusterName, provider, WorkloadClusterType)

	log.Infof("Workload Cluster %s Information: ", workloadClusterName)
	err = PrintClusterInformation(kubeConfigPath, workloadClusterKubeContext)
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
	err = DeleteCluster(workloadClusterName, provider, WorkloadClusterType)
	if err != nil {
		log.Errorf("error while deleting workload cluster: %v", err)

		err := tanzu.CollectManagementClusterAndWorkloadClusterDiagnostics(managementClusterName, workloadClusterName, provider.Name())
		if err != nil {
			log.Errorf("error while collecting diagnostics of management cluster and workload cluster: %v", err)
		}

		log.Fatal("error while deleting workload cluster: %v", err)
	}

	// TODO: Handle errors during waiting for cluster deletion.
	// We could retry in some cases, to just list the workload clusters.
	// If all retries fail, cleanup management cluster and then cleanup workload cluster
	WaitForWorkloadClusterDeletion(workloadClusterName)

	err = kubeclient.DeleteContext(kubeConfigPath, workloadClusterKubeContext)
	if err != nil {
		log.Errorf("error while deleting kube context %s at kubeconfig path: %v", managementClusterKubeContext, err)
	}

	// TODO: Handle errors during cluster deletion
	// and cleanup management cluster
	err = DeleteCluster(managementClusterName, provider, ManagementClusterType)
	if err != nil {
		log.Errorf("error while deleting management cluster: %v", err)

		err := tanzu.CollectManagementClusterDiagnostics(managementClusterName)
		if err != nil {
			log.Errorf("error while collecting diagnostics of management cluster: %v", err)
		}

		log.Fatal("error while deleting management cluster: %v", err)
	}
}
