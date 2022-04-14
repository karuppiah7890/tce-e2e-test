package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"testing"
	"time"

	"gopkg.in/yaml.v3"
	serializer "k8s.io/apimachinery/pkg/runtime/serializer"

	// TODO: Rename imports in a better manner, like, use capzABC and capiDEF etc for naming
	// CAPZ and CAPI stuff

	infrav1alpha3 "sigs.k8s.io/cluster-api-provider-azure/api/v1alpha3"
	infrav1alpha4 "sigs.k8s.io/cluster-api-provider-azure/api/v1alpha4"
	infrav1alpha3exp "sigs.k8s.io/cluster-api-provider-azure/exp/api/v1alpha3"
	infrav1alpha4exp "sigs.k8s.io/cluster-api-provider-azure/exp/api/v1alpha4"
	infrav1beta1exp "sigs.k8s.io/cluster-api-provider-azure/exp/api/v1beta1"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	clusterv1alpha3 "sigs.k8s.io/cluster-api/api/v1alpha3"
	clusterv1alpha4 "sigs.k8s.io/cluster-api/api/v1alpha4"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"

	capiBootstrapKubeadmv1beta "sigs.k8s.io/cluster-api/bootstrap/kubeadm/api/v1beta1"
	capiControlplaneKubeadmv1beta "sigs.k8s.io/cluster-api/controlplane/kubeadm/api/v1beta1"

	expv1alpha3 "sigs.k8s.io/cluster-api/exp/api/v1alpha3"
	expv1alpha4 "sigs.k8s.io/cluster-api/exp/api/v1alpha4"
	expv1 "sigs.k8s.io/cluster-api/exp/api/v1beta1"

	addonsv1alpha3 "sigs.k8s.io/cluster-api/exp/addons/api/v1alpha3"
	addonsv1alpha4 "sigs.k8s.io/cluster-api/exp/addons/api/v1alpha4"
	addonsv1 "sigs.k8s.io/cluster-api/exp/addons/api/v1beta1"

	"github.com/karuppiah7890/tce-e2e-test/testutils/azure"
	"github.com/karuppiah7890/tce-e2e-test/testutils/docker"
	"github.com/karuppiah7890/tce-e2e-test/testutils/kubeclient"
	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
	kubeRuntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/util/homedir"
	capzv1beta1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/marketplaceordering/armmarketplaceordering"
)

var scheme = kubeRuntime.NewScheme()

func init() {
	registerCustomResources()
}

func TestAzureManagementAndWorkloadCluster(t *testing.T) {
	log.InitLogger("azure-mgmt-wkld-e2e")

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

	if runtime.GOOS == "windows" {
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

	azureMarketplaceImageInfoForManagementCluster := getAzureMarketplaceImageInfoForManagementCluster(managementClusterName)

	// TODO: make the below function return an error and handle the error to log and exit?
	acceptAzureImageLicenses(azureTestSecrets.SubscriptionID, cred, azureMarketplaceImageInfoForManagementCluster...)

	// TODO: Handle errors during deployment
	// and cleanup management cluster
	runManagementCluster(managementClusterName)

	// TODO: check if management cluster is running by doing something similar to
	// `tanzu management-cluster get | grep "${MANAGEMENT_CLUSTER_NAME}" | grep running`

	// TODO: Handle errors
	getManagementClusterKubeConfig(managementClusterName)

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

	azureMarketplaceImageInfoForWorkloadCluster := getAzureMarketplaceImageInfoForWorkloadCluster(workloadClusterName)

	// TODO: make the below function return an error and handle the error to log and exit?
	acceptAzureImageLicenses(azureTestSecrets.SubscriptionID, cred, azureMarketplaceImageInfoForWorkloadCluster...)

	// TODO: Handle errors during deployment
	// and cleanup management cluster and then cleanup workload cluster
	runWorkloadCluster(workloadClusterName)

	checkWorkloadClusterIsRunning(workloadClusterName)

	// TODO: Handle errors
	getWorkloadClusterKubeConfig(workloadClusterName)

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
	deleteWorkloadCluster(workloadClusterName)

	// TODO: Handle errors during waiting for cluster deletion.
	// We could retry in some cases, to just list the workload clusters.
	// If all retries fail, cleanup management cluster and then cleanup workload cluster
	waitForWorkloadClusterDeletion(workloadClusterName)

	// TODO: Cleanup workload cluster kube config data (cluster, user, context)
	// since tanzu cluster delete does not delete workload cluster kubeconfig entry

	// TODO: Handle errors during cluster deletion
	// and cleanup management cluster
	deleteManagementCluster(managementClusterName)
}

func getKubeConfigPath() (string, error) {
	home := homedir.HomeDir()

	if home == "" {
		return "", fmt.Errorf("could not find home directory to get absolute path of kubeconfig")
	}

	return filepath.Join(home, ".kube", "config"), nil
}

func getKubeContextForTanzuCluster(clusterName string) string {
	return fmt.Sprintf("%s-admin@%s", clusterName, clusterName)
}

// TODO: Should we just use one function acceptAzureImageLicenses with the whole implementation? There will be a for loop with a big body though
func acceptAzureImageLicenses(subscriptionID string, cred *azidentity.ClientSecretCredential, azureMarketplaceImages ...*capzv1beta1.AzureMarketplaceImage) {
	for _, azureMarketplaceImage := range azureMarketplaceImages {
		acceptAzureImageLicense(subscriptionID, cred, azureMarketplaceImage)
	}
}

// This naming is for clarity until we move the function to some azure specific
// package then we can remove the reference to azure from it and rename
// it back to acceptImageLicense
func acceptAzureImageLicense(subscriptionID string, cred *azidentity.ClientSecretCredential, azureMarketplaceImage *capzv1beta1.AzureMarketplaceImage) {
	// We have hardcoded the value of the inputs required for accepting Azure VM image license terms.
	// TODO: Use management-cluster / workload cluster dry run (--dry-run) to get Azure VM image names / skus, offering, publisher
	azureVmImagePublisher := azureMarketplaceImage.Publisher
	// The value k8s-1dot21dot5-ubuntu-2004 comes from latest TKG BOM file based on OS arch, OS name and OS version
	// provided in test/azure/cluster-config.yaml in TCE repo. This value needs to be changed manually whenever there's going to
	// be a change in the underlying Tanzu Framework CLI version (management-cluster and cluster plugins) causing new
	// TKr BOMs to be used with new Azure VM images which have different image billing plan SKU
	azureVmImageBillingPlanSku := azureMarketplaceImage.SKU
	azureVmImageOffer := azureMarketplaceImage.Offer

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

func getAzureMarketplaceImageInfoForManagementCluster(managementClusterName string) []*capzv1beta1.AzureMarketplaceImage {
	var managementClusterCreateDryRunOutputBuffer bytes.Buffer

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
		Env:    append(os.Environ(), envVars...),
		Stdout: &managementClusterCreateDryRunOutputBuffer,
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
		// TODO: Should we print the whole command as part of the error? cliRunner will print it though, in this case
		log.Fatalf("Error occurred while running management cluster dry run. Exit code: %v. Error: %v", exitCode, err)
	}

	managementClusterCreateDryRunOutput, err := io.ReadAll(&managementClusterCreateDryRunOutputBuffer)
	if err != nil {
		// TODO: Should we print the whole command as part of the error?
		log.Fatalf("Error occurred while reading output of management cluster create dry run: %v", err)
	}

	objects := parseK8sYamlAndFetchAzureMachineTemplates(managementClusterCreateDryRunOutput)

	marketplaces := []*capzv1beta1.AzureMarketplaceImage{}

	for _, object := range objects {
		azureMachineTemplate, ok := object.(*capzv1beta1.AzureMachineTemplate)
		if !ok {
			log.Fatalf("Error occurred while parsing output of management cluster create dry run")
		}

		marketplaces = append(marketplaces, azureMachineTemplate.Spec.Template.Spec.Image.Marketplace)
	}

	return marketplaces
}

// Maybe return []*capzv1beta1.AzureMachineTemplate directly? Instead of []kubeRuntime.Object
// TODO: Rename this in a better manner? The function name and argument too
func parseK8sYamlAndFetchAzureMachineTemplates(fileR []byte) []kubeRuntime.Object {
	registerCustomResources()

	// TODO: Should we just use simple plain string match since we just want to pick AzureMachineTemplate only?
	// But yeah, in future we might parse other stuff, but as of now I don't see any such thing, so we could simplify this
	// For more types, use something like `(Role|ConfigMap)` etc
	acceptedK8sTypes := regexp.MustCompile(`(AzureMachineTemplate)`)
	sepYamlFilesBytes, err := SplitYAML(fileR)
	if err != nil {
		// return and handle error?
		log.Fatalf("Error while splitting YAML file. Err was: %s", err)
	}
	retVal := make([]kubeRuntime.Object, 0, len(sepYamlFilesBytes))
	for _, fBytes := range sepYamlFilesBytes {
		f := string(fBytes)
		if f == "\n" || f == "" {
			// ignore empty cases
			continue
		}

		decode := serializer.NewCodecFactory(scheme).UniversalDeserializer().Decode
		obj, groupVersionKind, err := decode(fBytes, nil, nil)

		if err != nil {
			// return and handle error?
			log.Fatalf("Error while decoding YAML object. Err was: %s", err)
			continue
		}

		if !acceptedK8sTypes.MatchString(groupVersionKind.Kind) {
			// The output contains K8s object types which are not needed so we are skipping this object with type groupVersionKind.Kind
		} else {
			retVal = append(retVal, obj)
		}

	}
	return retVal
}

func SplitYAML(resources []byte) ([][]byte, error) {
	dec := yaml.NewDecoder(bytes.NewReader(resources))

	var res [][]byte
	for {
		var value interface{}
		err := dec.Decode(&value)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		valueBytes, err := yaml.Marshal(value)
		if err != nil {
			return nil, err
		}
		res = append(res, valueBytes)
	}
	return res, nil
}

func registerCustomResources() {
	err := clientgoscheme.AddToScheme(scheme)
	if err != nil {
		panic(err)
	}

	err = apiextensionsv1.AddToScheme(scheme)
	if err != nil {
		panic(err)
	}

	err = clusterv1alpha3.AddToScheme(scheme)
	if err != nil {
		panic(err)
	}
	err = clusterv1alpha4.AddToScheme(scheme)
	if err != nil {
		panic(err)
	}
	err = clusterv1.AddToScheme(scheme)
	if err != nil {
		panic(err)
	}

	err = expv1alpha3.AddToScheme(scheme)
	if err != nil {
		panic(err)
	}
	err = expv1alpha4.AddToScheme(scheme)
	if err != nil {
		panic(err)
	}
	err = expv1.AddToScheme(scheme)
	if err != nil {
		panic(err)
	}

	err = addonsv1alpha3.AddToScheme(scheme)
	if err != nil {
		panic(err)
	}
	err = addonsv1alpha4.AddToScheme(scheme)
	if err != nil {
		panic(err)
	}
	err = addonsv1.AddToScheme(scheme)
	if err != nil {
		panic(err)
	}

	err = infrav1alpha3.AddToScheme(scheme)
	if err != nil {
		panic(err)
	}
	err = infrav1alpha4.AddToScheme(scheme)
	if err != nil {
		panic(err)
	}
	err = capzv1beta1.AddToScheme(scheme)
	if err != nil {
		panic(err)
	}
	err = infrav1alpha3exp.AddToScheme(scheme)
	if err != nil {
		panic(err)
	}
	err = infrav1alpha4exp.AddToScheme(scheme)
	if err != nil {
		panic(err)
	}
	err = infrav1beta1exp.AddToScheme(scheme)
	if err != nil {
		panic(err)
	}

	err = capiControlplaneKubeadmv1beta.AddToScheme(scheme)
	if err != nil {
		panic(err)
	}

	err = capiBootstrapKubeadmv1beta.AddToScheme(scheme)
	if err != nil {
		panic(err)
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

func getManagementClusterKubeConfig(managementClusterName string) {
	envVars := tanzuConfigToEnvVars(tanzuAzureConfig(managementClusterName))
	exitCode, err := cliRunner(Cmd{
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

func deleteManagementCluster(managementClusterName string) {
	envVars := tanzuConfigToEnvVars(tanzuAzureConfig(managementClusterName))
	exitCode, err := cliRunner(Cmd{
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

func getAzureMarketplaceImageInfoForWorkloadCluster(workloadClusterName string) []*capzv1beta1.AzureMarketplaceImage {
	var workloadClusterCreateDryRunOutputBuffer bytes.Buffer

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
		// TODO: Do we really want to output to log.InfoWriter ? Is this
		// data necessary in the logs? This data will contain secrets but for now we haven't masked secrets
		// in logs, also, even if we mask secrets, is this data useful and necessary?
		// The data in log can help development but that's all
		Stdout: &workloadClusterCreateDryRunOutputBuffer,
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
		log.Fatalf("Error occurred while running workload cluster dry run. Exit code: %v. Error: %v", exitCode, err)
	}

	workloadClusterCreateDryRunOutput, err := io.ReadAll(&workloadClusterCreateDryRunOutputBuffer)
	if err != nil {
		// TODO: Should we print the whole command as part of the error?
		log.Fatalf("Error occurred while reading output of workload cluster create dry run: %v", err)
	}

	objects := parseK8sYamlAndFetchAzureMachineTemplates(workloadClusterCreateDryRunOutput)

	marketplaces := []*capzv1beta1.AzureMarketplaceImage{}

	for _, object := range objects {
		azureMachineTemplate, ok := object.(*capzv1beta1.AzureMachineTemplate)
		if !ok {
			log.Fatalf("Error occurred while parsing output of workload cluster create dry run")
		}

		marketplaces = append(marketplaces, azureMachineTemplate.Spec.Template.Spec.Image.Marketplace)
	}

	return marketplaces
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

func checkWorkloadClusterIsRunning(workloadClusterName string) {
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

func getWorkloadClusterKubeConfig(workloadClusterName string) {
	envVars := tanzuConfigToEnvVars(tanzuAzureConfig(workloadClusterName))
	exitCode, err := cliRunner(Cmd{
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

type WorkloadCluster struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type WorkloadClusters []WorkloadCluster

func waitForWorkloadClusterDeletion(workloadClusterName string) {
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

	// TODO: Do we really want to output to log.InfoWriter ? Is this
	// data necessary in the logs? This function will be called
	// a lot of times. The data in log can help development
	multiWriter := io.MultiWriter(&clusterListOutput, log.InfoWriter)

	exitCode, err := cliRunner(Cmd{
		Name: "tanzu",
		Args: []string{
			"cluster",
			"list",
			"-o",
			"json",
		},
		Env:    os.Environ(),
		Stdout: multiWriter,
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

func checkDockerCLIInstallation() {
	docker.CheckDockerInstallation()
	docker.StopRunningContainer("test-worker")
	//docker.StopAllRunningContainer()
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
}

func tanzuConfigToEnvVars(tanzuConfig TanzuConfig) EnvVars {
	envVars := make(EnvVars, 0, len(tanzuConfig))

	for key, value := range tanzuConfig {
		envVars = append(envVars, fmt.Sprintf("%s=%s", key, value))
	}

	return envVars
}

// TODO: Should this return cluster information or just print / log them?
func printClusterInformation(kubeConfigPath string, kubeContext string) error {
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
