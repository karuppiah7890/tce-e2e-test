package azure

import (
	"bytes"
	"context"
	"io"
	"os"
	"regexp"

	serializer "k8s.io/apimachinery/pkg/runtime/serializer"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/marketplaceordering/armmarketplaceordering"
	"github.com/karuppiah7890/tce-e2e-test/testutils"
	"github.com/karuppiah7890/tce-e2e-test/testutils/clirunner"
	"github.com/karuppiah7890/tce-e2e-test/testutils/kubescheme"
	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
	"github.com/karuppiah7890/tce-e2e-test/testutils/tanzu"
	"github.com/karuppiah7890/tce-e2e-test/testutils/utils"

	kubeRuntime "k8s.io/apimachinery/pkg/runtime"
	capzv1beta1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
)

// TODO: Should we just use one function acceptAzureImageLicenses with the whole implementation? There will be a for loop with a big body though
func AcceptAzureImageLicenses(subscriptionID string, cred *azidentity.ClientSecretCredential, azureMarketplaceImages ...*capzv1beta1.AzureMarketplaceImage) {
	for _, azureMarketplaceImage := range azureMarketplaceImages {
		AcceptAzureImageLicense(subscriptionID, cred, azureMarketplaceImage)
	}
}

// This naming is for clarity until we move the function to some azure specific
// package then we can remove the reference to azure from it and rename
// it back to acceptImageLicense
func AcceptAzureImageLicense(subscriptionID string, cred *azidentity.ClientSecretCredential, azureMarketplaceImage *capzv1beta1.AzureMarketplaceImage) {
	azureVmImagePublisher := azureMarketplaceImage.Publisher
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

// Maybe return []*capzv1beta1.AzureMachineTemplate directly? Instead of []kubeRuntime.Object
// TODO: Rename this in a better manner? The function name and argument too
func ParseK8sYamlAndFetchAzureMachineTemplates(fileR []byte) []kubeRuntime.Object {

	// TODO: Should we just use simple plain string match since we just want to pick AzureMachineTemplate only?
	// But yeah, in future we might parse other stuff, but as of now I don't see any such thing, so we could simplify this
	// For more types, use something like `(Role|ConfigMap)` etc
	acceptedK8sTypes := regexp.MustCompile(`(AzureMachineTemplate)`)
	sepYamlFilesBytes, err := testutils.SplitYAML(fileR)
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

		decode := serializer.NewCodecFactory(kubescheme.GetScheme()).UniversalDeserializer().Decode
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

func GetAzureMarketplaceImageInfoForCluster(clusterName string, clusterType utils.ClusterType) []*capzv1beta1.AzureMarketplaceImage {
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

	objects := ParseK8sYamlAndFetchAzureMachineTemplates(clusterCreateDryRunOutput)

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
