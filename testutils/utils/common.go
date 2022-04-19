package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/karuppiah7890/tce-e2e-test/testutils/clirunner"
	"github.com/karuppiah7890/tce-e2e-test/testutils/kubeclient"
	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
	"k8s.io/client-go/util/homedir"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// TODO: Further move the functions to specifics file/libs accordingly

// TODO: Move this to a tanzu specific lib
type TanzuConfig map[string]string

type WorkloadCluster struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}
type WorkloadClusters []WorkloadCluster

// TODO: Move this to a common util / tanzu specific lib
type EnvVars []string

func CheckTanzuCLIInstallation() {
	log.Info("Checking tanzu CLI installation")
	path, err := exec.LookPath("tanzu")
	if err != nil {
		log.Fatalf("tanzu CLI is not installed")
	}
	log.Infof("tanzu CLI is available at path: %s\n", path)
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

func CheckTanzuManagementClusterCLIPluginInstallation() {
	log.Info("Checking tanzu management cluster plugin CLI installation")

	// TODO: Check for errors and return error?
	// TODO: Parse version and show warning if version is newer than what's tested by the devs while writing test
	// Refer - https://github.com/karuppiah7890/tce-e2e-test/issues/1#issuecomment-1094172278
	exitCode, err := clirunner.Run(clirunner.Cmd{
		Name: "tanzu",
		Args: []string{
			"management-cluster",
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

func CheckTanzuWorkloadClusterCLIPluginInstallation() {
	log.Info("Checking tanzu workload cluster plugin CLI installation")

	// TODO: Check for errors and return error?
	// TODO: Parse version and show warning if version is newer than what's tested by the devs while writing test
	// Refer - https://github.com/karuppiah7890/tce-e2e-test/issues/1#issuecomment-1094172278
	exitCode, err := clirunner.Run(clirunner.Cmd{
		Name: "tanzu",
		Args: []string{
			"cluster",
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
		log.Fatalf("Error occurred while checking workload cluster CLI plugin installation. Exit code: %v. Error: %v", exitCode, err)
	}
}

func RunManagementCluster(managementClusterName string, provider string) error {
	envVars := tanzuConfigToEnvVars(tanzuConfig(managementClusterName, provider))
	exitCode, err := clirunner.Run(clirunner.Cmd{
		Name: "tanzu",
		Args: []string{
			"management-cluster",
			"create",
			managementClusterName,
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
		return fmt.Errorf("error occurred while deploying management cluster. exit code: %v. error: %v", exitCode, err)
	}
	return nil
}

func GetManagementClusterKubeConfig(managementClusterName string, provider string) {
	// TODO: Do we really need the AWS secrets here?
	envVars := tanzuConfigToEnvVars(tanzuConfig(managementClusterName, provider))
	exitCode, err := clirunner.Run(clirunner.Cmd{
		// TODO: Replace magic strings like "tanzu", "management-cluster" etc
		Name: "tanzu",
		Args: []string{
			"management-cluster",
			"kubeconfig",
			"get",
			managementClusterName,
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
		log.Fatalf("Error occurred while getting management cluster kubeconfig. Exit code: %v. Error: %v", exitCode, err)
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

func RunWorkloadCluster(workloadClusterName string, provider string) error {
	// TODO: Do we really need the AWS secrets here?
	envVars := tanzuConfigToEnvVars(tanzuConfig(workloadClusterName, provider))
	exitCode, err := clirunner.Run(clirunner.Cmd{
		Name: "tanzu",
		Args: []string{
			"cluster",
			"create",
			workloadClusterName,
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
		return fmt.Errorf("error occurred while deploying workload cluster. exit code: %v. error: %v", exitCode, err)
	}

	return nil
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

func GetWorkloadClusterKubeConfig(workloadClusterName string, provider string) {
	// TODO: Do we really need the AWS secrets here?
	envVars := tanzuConfigToEnvVars(tanzuConfig(workloadClusterName, provider))
	exitCode, err := clirunner.Run(clirunner.Cmd{
		Name: "tanzu",
		Args: []string{
			"cluster",
			"kubeconfig",
			"get",
			workloadClusterName,
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
		log.Fatalf("Error occurred while getting workload cluster kubeconfig. Exit code: %v. Error: %v", exitCode, err)
	}
}

func DeleteWorkloadCluster(workloadClusterName string, provider string) error {
	// TODO: Do we really need the AWS secrets here?
	envVars := tanzuConfigToEnvVars(tanzuConfig(workloadClusterName, provider))
	exitCode, err := clirunner.Run(clirunner.Cmd{
		Name: "tanzu",
		Args: []string{
			"cluster",
			"delete",
			workloadClusterName,
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
		return fmt.Errorf("error occurred while deleting workload cluster. exit code: %v. error: %v", exitCode, err)
	}

	return nil
}

func DeleteManagementCluster(managementClusterName string, provider string) error {
	// TODO: Do we really need the AWS secrets here?
	envVars := tanzuConfigToEnvVars(tanzuConfig(managementClusterName, provider))
	exitCode, err := clirunner.Run(clirunner.Cmd{
		Name: "tanzu",
		Args: []string{
			"management-cluster",
			"delete",
			managementClusterName,
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
		return fmt.Errorf("error occurred while deleting management cluster. exit code: %v. error: %v", exitCode, err)
	}

	return nil
}

func tanzuConfigToEnvVars(tanzuConfig TanzuConfig) EnvVars {
	envVars := make(EnvVars, 0, len(tanzuConfig))

	for key, value := range tanzuConfig {
		envVars = append(envVars, fmt.Sprintf("%s=%s", key, value))
	}

	return envVars
}

func tanzuConfig(clusterName string, infraProvider string) TanzuConfig {
	switch infraProvider {
	case "aws":
		return TanzuConfig{
			"CLUSTER_NAME":               clusterName,
			"AWS_AMI_ID":                 "ami-0bcd9ed3ef40fad77",
			"INFRASTRUCTURE_PROVIDER":    "aws",
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
	case "azure":
		return TanzuConfig{
			"CLUSTER_NAME":                     clusterName,
			"INFRASTRUCTURE_PROVIDER":          "azure",
			"CLUSTER_PLAN":                     "dev",
			"AZURE_LOCATION":                   "australiaeast",
			"AZURE_CONTROL_PLANE_MACHINE_TYPE": "Standard_D4s_v3",
			"AZURE_NODE_MACHINE_TYPE":          "Standard_D4s_v3",
			"OS_ARCH":                          "amd64",
			"OS_NAME":                          "ubuntu",
			"OS_VERSION":                       "20.04",
			"AZURE_VNET_CIDR":                  "10.0.0.0/16",
			"AZURE_CONTROL_PLANE_SUBNET_CIDR":  "10.0.0.0/24",
			"AZURE_NODE_SUBNET_CIDR":           "10.0.1.0/24",
			"CLUSTER_CIDR":                     "100.96.0.0/11",
			"SERVICE_CIDR":                     "100.64.0.0/13",
			"ENABLE_CEIP_PARTICIPATION":        "false",
			"ENABLE_MHC":                       "true",
			"IDENTITY_MANAGEMENT_TYPE":         "none",
		}
	case "vsphere":
		return TanzuConfig{
			"CLUSTER_NAME":                   clusterName,
			"INFRASTRUCTURE_PROVIDER":        "vsphere",
			"CLUSTER_PLAN":                   "dev",
			"OS_ARCH":                        "amd64",
			"OS_NAME":                        "photon",
			"OS_VERSION":                     "3",
			"VSPHERE_CONTROL_PLANE_DISK_GIB": "40",
			"VSPHERE_CONTROL_PLANE_MEM_MIB":  "16384",
			"VSPHERE_CONTROL_PLANE_NUM_CPUS": "4",
			"VSPHERE_WORKER_DISK_GIB":        "40",
			"VSPHERE_WORKER_MEM_MIB":         "16384",
			"VSPHERE_WORKER_NUM_CPUS":        "4",
			"VSPHERE_INSECURE":               "true",
			"DEPLOY_TKG_ON_VSPHERE7":         "true",
			"ENABLE_TKGS_ON_VSPHERE7":        "false",
			"CLUSTER_CIDR":                   "100.96.0.0/11",
			"SERVICE_CIDR":                   "100.64.0.0/13",
			"ENABLE_CEIP_PARTICIPATION":      "false",
			"ENABLE_MHC":                     "true",
			"IDENTITY_MANAGEMENT_TYPE":       "none",
		}
	}
	return TanzuConfig{
		"CLUSTER_NAME": clusterName,
	}
}
