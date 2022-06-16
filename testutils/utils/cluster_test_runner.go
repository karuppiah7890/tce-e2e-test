package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/karuppiah7890/tce-e2e-test/testutils/clirunner"
	"github.com/karuppiah7890/tce-e2e-test/testutils/docker"
	"github.com/karuppiah7890/tce-e2e-test/testutils/kubeclient"
	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
	"github.com/karuppiah7890/tce-e2e-test/testutils/tanzu"
	"k8s.io/client-go/util/homedir"
)

type ClusterTestRunner interface {
	RunChecks()
	GetRandomClusterNames() (string, string)
	GetKubeContextForTanzuCluster(clusterName string) string
	GetKubeConfigPath() (string, error)
	RunCluster(clusterName string, provider Provider, clusterType ClusterType) error
	GetClusterKubeConfig(clusterName string, provider Provider, clusterType ClusterType)
	PrintClusterInformation(kubeConfigPath string, kubeContext string) error
	CheckWorkloadClusterIsRunning(workloadClusterName string)
	DeleteCluster(clusterName string, provider Provider, clusterType ClusterType) error
	WaitForWorkloadClusterDeletion(workloadClusterName string)
	CollectManagementClusterDiagnostics(managementClusterName string) error
	CollectManagementClusterAndWorkloadClusterDiagnostics(managementClusterName string, workloadClusterName string, workloadClusterInfra string) error
	DeleteContext(kubeConfigPath string, contextName string) error
}

type DefaultClusterTestRunner struct{}

// This is to ensure that DefaultClusterTestRunner implements ClusterTestRunner interface
// and if not, compiler level errors will be thrown
var _ ClusterTestRunner = DefaultClusterTestRunner{}

func (r DefaultClusterTestRunner) RunChecks() {
	_ = CheckTanzuCLIInstallation()

	CheckTanzuClusterCLIPluginInstallation(ManagementClusterType)

	CheckTanzuClusterCLIPluginInstallation(WorkloadClusterType)

	docker.CheckDockerInstallation()

	CheckKubectlCLIInstallation()

	PlatformSupportCheck()
}

func (r DefaultClusterTestRunner) GetRandomClusterNames() (string, string) {
	clusterNameSuffix := time.Now().Unix()
	managementClusterName := fmt.Sprintf("test-mgmt-%d", clusterNameSuffix)
	workloadClusterName := fmt.Sprintf("test-wkld-%d", clusterNameSuffix)
	log.Infof("Management Cluster Name : %s", managementClusterName)
	log.Infof("Workload Cluster Name : %s", workloadClusterName)
	return managementClusterName, workloadClusterName
}

func (r DefaultClusterTestRunner) GetKubeContextForTanzuCluster(clusterName string) string {
	return fmt.Sprintf("%s-admin@%s", clusterName, clusterName)
}

func (r DefaultClusterTestRunner) GetKubeConfigPath() (string, error) {
	home := homedir.HomeDir()

	if home == "" {
		return "", fmt.Errorf("could not find home directory to get absolute path of kubeconfig")
	}

	return filepath.Join(home, ".kube", "config"), nil
}

func (r DefaultClusterTestRunner) RunCluster(clusterName string, provider Provider, clusterType ClusterType) error {
	envVars := tanzu.TanzuConfigToEnvVars(provider.GetTanzuConfig(clusterName))
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

func (r DefaultClusterTestRunner) GetClusterKubeConfig(clusterName string, provider Provider, clusterType ClusterType) {
	// TODO: Do we really need the secrets here?
	envVars := tanzu.TanzuConfigToEnvVars(provider.GetTanzuConfig(clusterName))
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

func (r DefaultClusterTestRunner) PrintClusterInformation(kubeConfigPath string, kubeContext string) error {
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

func (r DefaultClusterTestRunner) CheckWorkloadClusterIsRunning(workloadClusterName string) {
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

func (r DefaultClusterTestRunner) DeleteCluster(clusterName string, provider Provider, clusterType ClusterType) error {
	// TODO: Do we really need the  secrets here?
	envVars := tanzu.TanzuConfigToEnvVars(provider.GetTanzuConfig(clusterName))
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

func (r DefaultClusterTestRunner) WaitForWorkloadClusterDeletion(workloadClusterName string) {
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

func (r DefaultClusterTestRunner) CollectManagementClusterDiagnostics(managementClusterName string) error {
	return nil
}

func (r DefaultClusterTestRunner) CollectManagementClusterAndWorkloadClusterDiagnostics(managementClusterName string, workloadClusterName string, workloadClusterInfra string) error {
	return tanzu.CollectManagementClusterAndWorkloadClusterDiagnostics(managementClusterName, workloadClusterName, workloadClusterInfra)
}

func (r DefaultClusterTestRunner) DeleteContext(kubeConfigPath string, contextName string) error {
	return kubeclient.DeleteContext(kubeConfigPath, contextName)
}
