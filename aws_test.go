package e2e

import (
	"fmt"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/karuppiah7890/tce-e2e-test/testutils/aws"
	"github.com/karuppiah7890/tce-e2e-test/testutils/clirunner"
	"github.com/karuppiah7890/tce-e2e-test/testutils/docker"
	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
	"github.com/karuppiah7890/tce-e2e-test/testutils/platforms"
)

func TestAwsManagementAndWorkloadCluster(t *testing.T) {
	log.InitLogger("aws-mgmt-wkld-e2e")

	// TODO: Think about installing TCE / TF from tar ball and from source
	// make release based on OS? Windows has make? Hmm
	// release-dir
	// tar ball, zip based on OS
	// install.sh and install.bat based on OS
	// TODO: use tce.Install("<version>")?

	// Ensure TCE/TF is installed - check TCE installation or install it if not present. Or do it prior to the test run.
	// check if tanzu is installed
	checkTanzuCLIInstallation()

	// Ensure management and workload cluster plugins are present.
	// check if management cluster plugin is present
	checkTanzuManagementClusterCLIPluginInstallation()

	// check if workload cluster plugin is present
	checkTanzuWorkloadClusterCLIPluginInstallation()

	// check if docker is installed. This is required by tanzu CLI I think, both docker client CLI and docker daemon
	docker.CheckDockerInstallation()
	// check if kubectl is installed. This is required by tanzu CLI to apply using kubectl apply to create cluster
	checkKubectlCLIInstallation()

	if runtime.GOOS == platforms.WINDOWS {
		log.Warn("Warning: This test has been tested only on Linux and Mac OS till now. Support for Windows has not been tested, so it's experimental and not guranteed to work!")
	}

	// TODO: Ensure package plugin is present in case package tests are gonna be executed.

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

	// TODO: Handle errors during deployment
	// and cleanup management cluster
	runAwsManagementCluster(managementClusterName)

	// TODO: check if management cluster is running by doing something similar to
	// `tanzu management-cluster get | grep "${MANAGEMENT_CLUSTER_NAME}" | grep running`

	// TODO: Handle errors
	getAwsManagementClusterKubeConfig(managementClusterName)

	kubeConfigPath, err := getKubeConfigPath()
	if err != nil {
		// Should we panic here and stop?
		log.Errorf("error while getting kubeconfig path: %v", err)
	}
	managementClusterKubeContext := getKubeContextForTanzuCluster(managementClusterName)

	log.Infof("Management Cluster %s Information: ", managementClusterName)
	err = printClusterInformation(kubeConfigPath, managementClusterKubeContext)
	if err != nil {
		// Should we panic here and stop?
		log.Errorf("error while printing management cluster information: %v", err)
	}

	// TODO: Handle errors during deployment
	// and cleanup management cluster and then cleanup workload cluster
	// and cleanup management cluster and then cleanup workload cluster
	runAwsWorkloadCluster(workloadClusterName)

	checkWorkloadClusterIsRunning(workloadClusterName)

	// TODO: Handle errors
	getAwsWorkloadClusterKubeConfig(workloadClusterName)

	workloadClusterKubeContext := getKubeContextForTanzuCluster(workloadClusterName)

	log.Infof("Workload Cluster %s Information: ", workloadClusterName)
	err = printClusterInformation(kubeConfigPath, workloadClusterKubeContext)
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
	deleteAwsWorkloadCluster(workloadClusterName)

	// TODO: Handle errors during waiting for cluster deletion.
	// We could retry in some cases, to just list the workload clusters.
	// If all retries fail, cleanup management cluster and then cleanup workload cluster
	waitForWorkloadClusterDeletion(workloadClusterName)

	// TODO: Cleanup workload cluster kube config data (cluster, user, context)
	// since tanzu cluster delete does not delete workload cluster kubeconfig entry

	// TODO: Handle errors during cluster deletion
	// and cleanup management cluster
	deleteAwsManagementCluster(managementClusterName)
}

// TODO: Duplicate of runManagementCluster in azure_test.go , just config is different
func runAwsManagementCluster(managementClusterName string) {
	envVars := tanzuConfigToEnvVars(tanzuAwsConfig(managementClusterName))
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
		log.Fatalf("Error occurred while deploying management cluster. Exit code: %v. Error: %v", exitCode, err)
	}
}

func getAwsManagementClusterKubeConfig(managementClusterName string) {
	// TODO: Do we really need the AWS secrets here?
	envVars := tanzuConfigToEnvVars(tanzuAwsConfig(managementClusterName))
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

func deleteAwsManagementCluster(managementClusterName string) {
	// TODO: Do we really need the AWS secrets here?
	envVars := tanzuConfigToEnvVars(tanzuAwsConfig(managementClusterName))
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
		log.Fatalf("Error occurred while deleting management cluster. Exit code: %v. Error: %v", exitCode, err)
	}
}

func runAwsWorkloadCluster(workloadClusterName string) {
	// TODO: Do we really need the AWS secrets here?
	envVars := tanzuConfigToEnvVars(tanzuAwsConfig(workloadClusterName))
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
		log.Fatalf("Error occurred while deploying workload cluster. Exit code: %v. Error: %v", exitCode, err)
	}
}

func getAwsWorkloadClusterKubeConfig(workloadClusterName string) {
	// TODO: Do we really need the AWS secrets here?
	envVars := tanzuConfigToEnvVars(tanzuAwsConfig(workloadClusterName))
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

func deleteAwsWorkloadCluster(workloadClusterName string) {
	// TODO: Do we really need the AWS secrets here?
	envVars := tanzuConfigToEnvVars(tanzuAwsConfig(workloadClusterName))
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
		log.Fatalf("Error occurred while deleting workload cluster. Exit code: %v. Error: %v", exitCode, err)
	}
}

// TODO: Consider using Tanzu golang client library instead of running tanzu as a CLI.
// We could invoke plugins using their names and have tight integration. It comes with it's
// own pros and cons. Pro - tight and smooth integration with Tanzu Framework.
// Con - same as Pro - tight and smooth integration with Tanzu Framework - why? because Tanzu
// Framework does not provide any guarantee for API support, also it's in 0.x.y series which
// means they can break a lot of things which can break already generally fragile E2E tests
// more easily and more often. Also, if we import Tanzu Framework as a library, to test different
// versions of Tanzu Framework, we have import different versions of it, unlike CLI where we can
// just install the appropriate CLI version before testing it. For example, test v0.11.4 TF that
// TCE currently uses and also test v0.20.0 TF which is the latest version of TF. Of course it's not
// easy to concurrently / simultaneously test both versions, at least not in CLI, and with library, idk,
// it might be possible and easy? not sure for now, gotta experiment. We can also consider dynamically linked
// libraries and similar concept, we currently instead have tanzu CLI tool which is dynamically invoked and linked
// to this test program

// TODO: Maybe create a wrapper function called Tanzu() around clirunner.Run()?

func tanzuAwsConfig(clusterName string) TanzuConfig {
	// TODO: Ideas:
	// We could also represent this config in a test data yaml file,
	// but config as code is more powerful - we can do more over here
	// for example, run tests for multiple plans - dev and prod very
	// easily instead of duplicating the whole config yaml file just to
	// run same test with different plan. We can then easily run many tests
	// with different set of config values by defining the range / possible set
	// of test values for each config.
	// Some configs that can be changed -
	// 1. Infra provider can change if it's AWS, vSphere etc
	// but this function is named as tanzuAzureConfig so it's okay.
	// 2. Cluster plan - dev and prod
	// 3. Azure location - the whole big list of azure locations
	// 4. Azure control plane machine type - the whole big list of azure VM machine types. Note: https://github.com/vmware-tanzu/community-edition/issues/1749. Also note, we might need VMs of some minimum size for cluster creation to work
	// 5. Azure worker node machine type - the whole big list of azure VM machine types. Note: https://github.com/vmware-tanzu/community-edition/issues/1749. Also note, we might need VMs of some minimum size for cluster creation to work
	// 6. OS_ARCH - amd64, arm64 . There's talks around ARM support at different levels now
	// 7. OS_VERSION - 20.04 or 18.04 as of now
	// 8. AZURE_VNET_CIDR - any CIDR range
	// 9. AZURE_CONTROL_PLANE_SUBNET_CIDR - any CIDR range
	// 10. AZURE_NODE_SUBNET_CIDR - any CIDR range
	// 11. CLUSTER_CIDR - any CIDR range
	// 12. CLUSTER_CIDR - any CIDR range
	// 13. ENABLE_CEIP_PARTICIPATION - true or false
	// 14. ENABLE_MHC - true or false. In one issue, someone had to set this to false or else their cluster creation was failing
	// 15. IDENTITY_MANAGEMENT_TYPE - none or some particular set of identity management types
	return TanzuConfig{
		"CLUSTER_NAME":               clusterName,
		"AWS_AMI_ID":                 "ami-0bcd9ed3ef40fad77",
		"INFRASTRUCTURE_PROVIDER":    "aws",
		"CLUSTER_PLAN":               "dev",
		"AWS_NODE_AZ":                "us-east-1a",
		"AWS_REGION":                 "us-east-1",
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
