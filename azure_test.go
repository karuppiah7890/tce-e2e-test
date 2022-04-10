package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"testing"
	"time"

	"github.com/karuppiah7890/tce-e2e-test/testutils/azure"
	"github.com/karuppiah7890/tce-e2e-test/testutils/log"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/marketplaceordering/armmarketplaceordering"
)

func TestAzureManagementAndWorkloadCluster(t *testing.T) {
	log.InitLogger("azure-mgmt-wkld-e2e")

	// check if tanzu is installed
	checkTanzuCLIInstallation()

	// check if management cluster plugin is present
	checkTanzuManagementClusterCLIPluginInstallation()

	// check if workload cluster plugin is present
	checkTanzuWorkloadClusterCLIPluginInstallation()

	// check if docker is installed. This is required by tanzu CLI I think, both docker client CLI and docker daemon
	checkDockerCLIInstallation()
	// check if kubectl is installed. This is required by tanzu CLI to apply using kubectl apply to create cluster
	checkKubectlCLIInstallation()

	if runtime.GOOS == "windows" || runtime.GOOS == "linux" {
		log.Warn("Warning: This test has been tested only on Mac OS till now. Support for Linux and Windows has not been tested, so it's experimental and not guranteed to work!")
	}

	// Ensure TCE/TF is installed - check TCE installation or install it if not present. Or do it prior to the test run.

	// Ensure management and workload cluster plugins are present.

	// Ensure package plugin is present in case package tests are gonna be executed.

	azureTestSecrets := azure.ExtractAzureTestSecretsFromEnvVars()

	// Have different log levels - none/minimal, error, info, debug etc, so that we can accordingly use those in the E2E test

	cred := azure.Login()

	// TODO: make the below function return an error and handle the error to log and exit?
	acceptImageLicense(azureTestSecrets.SubscriptionID, cred)

	// Create random names for management and workload clusters so that we can use them to name the test clusters we are going to
	// create. Ensure that these names are not already taken - check the resource group names to double check :) As Resource group name
	// is based on the cluster name
	// TODO: Create random names later, using random number or using short or long UUIDs.
	// TODO: Do we allow users to pass the cluster name for both clusters? We could. How do we take inputs? File? Env vars? Flags?
	clusterNameSuffix := time.Now().Unix()
	managementClusterName := fmt.Sprintf("test-mgmt-%d", clusterNameSuffix)
	workloadClusterName := fmt.Sprintf("test-wkld-%d", clusterNameSuffix)

	runManagementClusterDryRun(managementClusterName)

	// TODO: Handle errors during deployment
	runManagementCluster(managementClusterName)

	runWorkloadClusterDryRun(workloadClusterName)

	// TODO: Handle errors during deployment
	runWorkloadCluster(workloadClusterName)

	// TODO: Handle errors during cluster deletion
	deleteWorkloadCluster(workloadClusterName)

	// TODO: Handle errors during waiting for cluster deletion.
	// We could retry in some cases, to just list the workload clusters
	waitForWorkloadClusterDeletion(workloadClusterName)

	// TODO: Handle errors during cluster deletion
	deleteManagementCluster(managementClusterName)
}

func acceptImageLicense(subscriptionID string, cred *azidentity.ClientSecretCredential) {
	// We have hardcoded the value of the inputs required for accepting Azure VM image license terms.
	// TODO: Use management-cluster / workload cluster dry run (--dry-run) to get Azure VM image names / skus, offering, publisher
	azureVmImagePublisher := "vmware-inc"
	// The value k8s-1dot21dot5-ubuntu-2004 comes from latest TKG BOM file based on OS arch, OS name and OS version
	// provided in test/azure/cluster-config.yaml in TCE repo. This value needs to be changed manually whenever there's going to
	// be a change in the underlying Tanzu Framework CLI version (management-cluster and cluster plugins) causing new
	// TKr BOMs to be used with new Azure VM images which have different image billing plan SKU
	azureVmImageBillingPlanSku := "k8s-1dot22dot8-ubuntu-2004"
	azureVmImageOffer := "tkg-capi"

	ctx := context.Background()
	client := armmarketplaceordering.NewMarketplaceAgreementsClient(subscriptionID, cred, nil)

	log.Info("Getting marketplace terms for Azure VM image")
	res, err := client.Get(ctx,
		armmarketplaceordering.OfferType(armmarketplaceordering.OfferTypeVirtualmachine),
		azureVmImagePublisher,
		azureVmImageOffer,
		azureVmImageBillingPlanSku,
		nil)
	if err != nil {
		log.Fatalf("Error while getting marketplace terms for Azure VM image: %+v", err)
	}

	agreementTerms := res.MarketplaceAgreementsClientGetResult.AgreementTerms

	if agreementTerms.Properties == nil {
		log.Fatalf("Error: Azure VM image agreement terms Properties field is not available")
	}

	if agreementTerms.Properties.Accepted == nil {
		log.Fatalf("Error: Azure VM image agreement terms Properties Accepted field is not available")
	}

	if isTermsAccepted := *agreementTerms.Properties.Accepted; isTermsAccepted {
		log.Info("Azure VM image agreement terms are already accepted")
	} else {
		log.Info("Azure VM image agreement terms is not already accepted. Accepting the Azure VM image agreement terms now")

		*agreementTerms.Properties.Accepted = true
		// Note: We sign using a PUT request to change the `accepted` property in the agreement. This is how Azure CLI does it too.
		// This is because the sign API does not work as of this comment. Reference - https://docs.microsoft.com/en-us/answers/questions/52637/cannot-sign-azure-marketplace-vm-image-licence-thr.html
		createResponse, err := client.Create(ctx, armmarketplaceordering.OfferTypeVirtualmachine, azureVmImagePublisher, azureVmImageOffer, azureVmImageBillingPlanSku, agreementTerms, nil)
		if err != nil {
			log.Fatalf("Error while signing and accepting the agreement terms for Azure VM image: %+v", err)
		}

		signedAgreementTerms := createResponse.AgreementTerms

		if signedAgreementTerms.Properties == nil {
			log.Fatalf("Error while signing and accepting the agreement terms for Azure VM image: Azure VM image agreement terms Properties field is not available")
		}

		if signedAgreementTerms.Properties.Accepted == nil {
			log.Fatalf("Error while signing and accepting the agreement terms for Azure VM image: Azure VM image agreement terms Properties Accepted field is not available")
		}

		if isTermsSignedAndAccepted := *signedAgreementTerms.Properties.Accepted; !isTermsSignedAndAccepted {
			log.Fatalf("Error while signing and accepting the agreement terms for Azure VM image: Azure VM image agreement terms was not signed and accepted")
		} else {
			log.Info("Accepted the Azure VM image agreement terms!")
		}
	}
}

func runManagementClusterDryRun(managementClusterName string) {
	envVars := tanzuConfigToEnvVars(tanzuAzureConfig(managementClusterName))
	exitCode, err := cliRunner(Cmd{
		Name: "tanzu",
		Args: []string{
			"management-cluster",
			"create",
			managementClusterName,
			"--dry-run",
			// TODO: Should we add verbosity flag and value by default? or
			// let the user define the verbosity when running the tests maybe?
			// "-v",
			// "10",
		},
		Env: append(os.Environ(), envVars...),
		// TODO: Output to log files in the future and if needed, to console also
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})

	if err != nil {
		log.Fatalf("Error occurred while running management cluster dry run. Exit code: %v. Error: %v", exitCode, err)
	}
}

func runManagementCluster(managementClusterName string) {
	envVars := tanzuConfigToEnvVars(tanzuAzureConfig(managementClusterName))
	exitCode, err := cliRunner(Cmd{
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
		Env: append(os.Environ(), envVars...),
		// TODO: Output to log files in the future and if needed, to console also
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})

	if err != nil {
		log.Fatalf("Error occurred while deploying management cluster. Exit code: %v. Error: %v", exitCode, err)
	}
}

func deleteManagementCluster(managementClusterName string) {
	envVars := tanzuConfigToEnvVars(tanzuAzureConfig(managementClusterName))
	exitCode, err := cliRunner(Cmd{
		Name: "tanzu",
		Args: []string{
			"management-cluster",
			"delete",
			managementClusterName,
			// TODO: Should we add verbosity flag and value by default? or
			// let the user define the verbosity when running the tests maybe?
			// "-v",
			// "10",
		},
		Env: append(os.Environ(), envVars...),
		// TODO: Output to log files in the future and if needed, to console also
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})

	if err != nil {
		log.Fatalf("Error occurred while deleting management cluster. Exit code: %v. Error: %v", exitCode, err)
	}
}

func runWorkloadClusterDryRun(workloadClusterName string) {
	envVars := tanzuConfigToEnvVars(tanzuAzureConfig(workloadClusterName))
	exitCode, err := cliRunner(Cmd{
		Name: "tanzu",
		Args: []string{
			"cluster",
			"create",
			workloadClusterName,
			"--dry-run",
			// TODO: Should we add verbosity flag and value by default? or
			// let the user define the verbosity when running the tests maybe?
			// "-v",
			// "10",
		},
		Env: append(os.Environ(), envVars...),
		// TODO: Output to log files in the future and if needed, to console also
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})

	if err != nil {
		log.Fatalf("Error occurred while running workload cluster dry run. Exit code: %v. Error: %v", exitCode, err)
	}
}

func runWorkloadCluster(workloadClusterName string) {
	envVars := tanzuConfigToEnvVars(tanzuAzureConfig(workloadClusterName))
	exitCode, err := cliRunner(Cmd{
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
		Env: append(os.Environ(), envVars...),
		// TODO: Output to log files in the future and if needed, to console also
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})

	if err != nil {
		log.Fatalf("Error occurred while deploying workload cluster. Exit code: %v. Error: %v", exitCode, err)
	}
}

func deleteWorkloadCluster(workloadClusterName string) {
	envVars := tanzuConfigToEnvVars(tanzuAzureConfig(workloadClusterName))
	exitCode, err := cliRunner(Cmd{
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
		Env: append(os.Environ(), envVars...),
		// TODO: Output to log files in the future and if needed, to console also
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})

	if err != nil {
		log.Fatalf("Error occurred while deleting workload cluster. Exit code: %v. Error: %v", exitCode, err)
	}
}

type WorkloadCluster struct {
	Name string `json:"name"`
}

type WorkloadClusters []WorkloadCluster

func waitForWorkloadClusterDeletion(workloadClusterName string) {
	for i := 0; i < 100; i++ {
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
	}

	// TODO: maybe return error instead of fatal stop?
	log.Fatalf("Timed out waiting for workload cluster %s to get deleted", workloadClusterName)
}

func listWorkloadClusters() WorkloadClusters {
	var workloadClusters WorkloadClusters

	var clusterListOutput bytes.Buffer

	multiWriter := io.MultiWriter(&clusterListOutput, os.Stdout)

	exitCode, err := cliRunner(Cmd{
		Name: "tanzu",
		Args: []string{
			"cluster",
			"list",
			"-o",
			"json",
		},
		Env: os.Environ(),
		// TODO: Output to log files in the future and if needed, to console also
		Stdout: multiWriter,
		Stderr: os.Stderr,
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

func checkTanzuCLIInstallation() {
	log.Info("Checking tanzu CLI installation")
	path, err := exec.LookPath("tanzu")
	if err != nil {
		log.Fatalf("tanzu CLI is not installed")
	}
	log.Infof("tanzu CLI is available at path: %s\n", path)
}

func checkTanzuManagementClusterCLIPluginInstallation() {
	log.Info("Checking tanzu management cluster plugin CLI installation")

	// TODO: Check for errors and return error?
	// TODO: Parse version and show warning if version is newer than what's tested by the devs while writing test
	// Refer - https://github.com/karuppiah7890/tce-e2e-test/issues/1#issuecomment-1094172278
	exitCode, err := cliRunner(Cmd{
		Name: "tanzu",
		Args: []string{
			"management-cluster",
			"version",
		},
		// TODO: Output to log files in the future and if needed, to console also
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})

	if err != nil {
		log.Fatalf("Error occurred while checking management cluster CLI plugin installation. Exit code: %v. Error: %v", exitCode, err)
	}
}

func checkTanzuWorkloadClusterCLIPluginInstallation() {
	log.Info("Checking tanzu workload cluster plugin CLI installation")

	// TODO: Check for errors and return error?
	// TODO: Parse version and show warning if version is newer than what's tested by the devs while writing test
	// Refer - https://github.com/karuppiah7890/tce-e2e-test/issues/1#issuecomment-1094172278
	exitCode, err := cliRunner(Cmd{
		Name: "tanzu",
		Args: []string{
			"cluster",
			"version",
		},
		// TODO: Output to log files in the future and if needed, to console also
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})

	if err != nil {
		log.Fatalf("Error occurred while checking workload cluster CLI plugin installation. Exit code: %v. Error: %v", exitCode, err)
	}
}

func checkDockerCLIInstallation() {
	log.Info("Checking docker CLI installation")

	path, err := exec.LookPath("docker")
	if err != nil {
		log.Fatalf("docker CLI is not installed")
	}
	log.Infof("docker CLI is available at path: %s\n", path)
}

func checkKubectlCLIInstallation() {
	log.Info("Checking kubectl CLI installation")

	path, err := exec.LookPath("kubectl")
	if err != nil {
		log.Fatalf("kubectl CLI is not installed")
	}
	log.Infof("kubectl CLI is available at path: %s\n", path)
}

type Cmd struct {
	// Name is the Name of the command to run.
	//
	// This is the only field that must be set to a non-zero
	// value.
	Name string

	// Args holds command line arguments, including the command as Args[0].
	// If the Args field is empty or nil, Run uses {Name}.
	Args []string

	// Env specifies the environment of the process.
	// Each entry is of the form "key=value".
	// If Env is nil, the new process uses the current process's
	// environment.
	// If Env contains duplicate environment keys, only the last
	// value in the slice for each duplicate key is used.
	// As a special case on Windows, SYSTEMROOT is always added if
	// missing and not explicitly set to the empty string.
	Env []string

	// Stdout and Stderr specify the process's standard output and error.
	//
	// If either is nil, Run connects the corresponding file descriptor
	// to the null device (os.DevNull).
	//
	// If either is an *os.File, the corresponding output from the process
	// is connected directly to that file.
	//
	// Otherwise, during the execution of the command a separate goroutine
	// reads from the process over a pipe and delivers that data to the
	// corresponding Writer. In this case, Wait does not complete until the
	// goroutine reaches EOF or encounters an error.
	//
	// If Stdout and Stderr are the same writer, and have a type that can
	// be compared with ==, at most one goroutine at a time will call Write.
	Stdout io.Writer
	Stderr io.Writer
}

// TODO: Maybe create a wrapper function called Tanzu() around cliRunner?
func cliRunner(command Cmd) (int, error) {
	cmd := exec.Command(command.Name, command.Args...)
	cmd.Stdout = command.Stdout
	cmd.Stderr = command.Stderr
	cmd.Env = command.Env

	// TODO: Maybe set cmd.Env explicitly to a narrow set of env vars to just inject the secrets
	// that we want to inject and nothing else. But system level env vars maybe needed for the CLI.
	// Think about how to inject the env vars. Use single struct as function argument?
	// Use a struct to store the data including env vars and use that as data/context in
	// it's methods where method runs the command by injecting the env vars
	// Or something like
	// Tanzu({ env: []string{"key=value", "key2=value2"}, command: "management-cluster version" })
	// But the above is not exactly readable, hmm

	log.Infof("Running the command `%v`", cmd.String())

	err := cmd.Run()
	if err != nil {
		// TODO: Handle the error by returning it?
		log.Infof("Error occurred while running the command `%v`: %v", cmd.String(), err)
		return cmd.ProcessState.ExitCode(), err
	}

	return cmd.ProcessState.ExitCode(), nil
}

type TanzuConfig map[string]string

type EnvVars []string

func tanzuAzureConfig(clusterName string) TanzuConfig {
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
		"CLUSTER_NAME":                     clusterName,
		"INFRASTRUCTURE_PROVIDER":          "azure",
		"CLUSTER_PLAN":                     "dev",
		"AZURE_LOCATION":                   "australiacentral",
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
}

func tanzuConfigToEnvVars(tanzuConfig TanzuConfig) EnvVars {
	envVars := make(EnvVars, 0, len(tanzuConfig))

	for key, value := range tanzuConfig {
		envVars = append(envVars, fmt.Sprintf("%s=%s", key, value))
	}

	return envVars
}
