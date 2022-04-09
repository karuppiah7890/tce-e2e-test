package e2e

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/marketplaceordering/armmarketplaceordering"
)

const AzureTenantIDEnvVarName = "AZURE_TENANT_ID"
const AzureSshPublicKeyBase64EnvVarName = "AZURE_SSH_PUBLIC_KEY_B64"
const AzureClientIDEnvVarName = "AZURE_CLIENT_ID"
const AzureClientSecretEnvVarName = "AZURE_CLIENT_SECRET"
const AzureSubscriptionIDEnvVarName = "AZURE_SUBSCRIPTION_ID"

type AzureTestSecrets struct {
	TenantID       string
	SubscriptionID string
	ClientID       string
	ClientSecret   string
	SshPublicKey   string
}

func TestAzureManagementAndWorkloadCluster(t *testing.T) {
	// Ensure TCE/TF is installed - check TCE installation or install it if not present. Or do it prior to the test run.

	// Ensure management and workload cluster plugins are present.

	// Ensure package plugin is present in case package tests are gonna be executed.

	// Check required env vars to run the E2E test. The required env vars are the env vars which store secrets.
	// If the required env vars are not present, throw error and give information about how to go about fixing the error - by
	// getting the secrets from appropriate place with the help of docs, maybe link to appropriate Azure or TCE or TKG docs for this.
	// Ensure that the secrets are NEVER logged into the console
	requiredEnvVars := []string{
		AzureTenantIDEnvVarName,
		AzureSubscriptionIDEnvVarName,
		AzureClientIDEnvVarName,
		AzureClientSecretEnvVarName,
		AzureSshPublicKeyBase64EnvVarName,
	}
	log.Println("Checking required environment variables...")
	errs := checkRequiredEnvVars(requiredEnvVars)

	if len(errs) != 0 {
		log.Fatalf("Errors while checking required environment variables: %v\n", errs)
	}

	log.Println("Extracting Azure test secrets from environment variables...")
	azureTestSecrets := extractAzureTestSecretsFromEnvVars()

	// Have different log levels - none/minimal, error, info, debug etc, so that we can accordingly use those in the E2E test

	// Create random names for management and workload clusters so that we can use them to name the test clusters we are going to
	// create. Ensure that these names are not already taken - check the resource group names to double check :) As Resource group name
	// is based on the cluster name
	// TODO: Create random names later, using random number or using short or long UUIDs.
	// TODO: Do we allow users to pass the cluster name for both clusters? We could. How do we take inputs? File? Env vars? Flags?
	// managementClusterName := "test-mgmt"
	// workloadClusterName := "test-wkld"

	// Hard code the value of the inputs required / needed / necessary for accepting Azure VM image license / terms.
	// TODO: management-cluster / workload cluster dry run (--dry-run) to get Azure VM image names / skus, offering, publisher

	log.Println("Logging into Azure...")
	// login to azure
	cred, err :=
		azidentity.NewClientSecretCredential(azureTestSecrets.TenantID,
			azureTestSecrets.ClientID, azureTestSecrets.ClientSecret, nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}

	azureVmImagePublisher := "vmware-inc"
	// The value k8s-1dot21dot5-ubuntu-2004 comes from latest TKG BOM file based on OS arch, OS name and OS version
	// provided in test/azure/cluster-config.yaml in TCE repo. This value needs to be changed manually whenever there's going to
	// be a change in the underlying Tanzu Framework CLI version (management-cluster and cluster plugins) causing new
	// TKr BOMs to be used with new Azure VM images which have different image billing plan SKU
	azureVmImageBillingPlanSku := "k8s-1dot22dot8-ubuntu-2004"
	azureVmImageOffer := "tkg-capi"

	ctx := context.Background()
	client := armmarketplaceordering.NewMarketplaceAgreementsClient(azureTestSecrets.SubscriptionID, cred, nil)

	log.Println("Getting marketplace terms for Azure VM image...")
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
		log.Println("Azure VM image agreement terms are already accepted")
	} else {
		log.Println("Azure VM image agreement terms is not already accepted. Accepting the Azure VM image agreement terms now...")

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
			log.Println("Accepted the Azure VM image agreement terms!")
		}
	}
}

func checkRequiredEnvVars(requiredEnvVars []string) []error {
	errors := make([]error, 0, len(requiredEnvVars))
	for _, envVar := range requiredEnvVars {
		// TODO: Do we also error out when the env var value is defined but empty??
		_, isDefined := os.LookupEnv(envVar)

		if !isDefined {
			errors = append(errors, fmt.Errorf("Environment variable `%s` is required but not defined", envVar))
		}
	}

	return errors
}

func extractAzureTestSecretsFromEnvVars() AzureTestSecrets {
	return AzureTestSecrets{
		TenantID:       os.Getenv(AzureTenantIDEnvVarName),
		SubscriptionID: os.Getenv(AzureSubscriptionIDEnvVarName),
		ClientID:       os.Getenv(AzureClientIDEnvVarName),
		ClientSecret:   os.Getenv(AzureClientSecretEnvVarName),
		SshPublicKey:   os.Getenv(AzureSshPublicKeyBase64EnvVarName),
	}
}
