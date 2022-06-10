package azure

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/karuppiah7890/tce-e2e-test/testutils/utils"
)

// TODO: Change name?
type Provider struct {
	cred        *azidentity.ClientSecretCredential
	testSecrets TestSecrets
}

func (provider *Provider) CheckRequiredEnvVars() bool {
	CheckRequiredAzureEnvVars()
	return true
}

func (provider *Provider) Name() string {
	return "azure"
}

func (provider *Provider) Init() error {
	provider.testSecrets = ExtractAzureTestSecretsFromEnvVars()

	provider.cred = Login()

	return nil
}

func (provider *Provider) PreClusterCreationTasks(clusterName string, clusterType utils.ClusterType) error {
	azureMarketplaceImageInfoForCluster := GetAzureMarketplaceImageInfoForCluster(clusterName, clusterType)

	// TODO: make the below function return an error and handle the error to log and exit?
	AcceptAzureImageLicenses(provider.testSecrets.SubscriptionID, provider.cred, azureMarketplaceImageInfoForCluster...)

	return nil
}

func (provider *Provider) CleanupCluster(ctx context.Context, clusterName string) error {
	err := DeleteResourceGroup(context.TODO(), clusterName, provider.testSecrets.SubscriptionID, provider.cred)
	if err != nil {
		return fmt.Errorf("error while cleaning up azure resource group of the cluster which has all the cluster resources: %v", err)
	}
	return nil
}

// TODO: Change name?
var PROVIDER utils.Provider = &Provider{}
