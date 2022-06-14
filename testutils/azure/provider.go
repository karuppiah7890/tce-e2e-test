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
	err := DeleteResourceGroup(ctx, clusterName, provider.testSecrets.SubscriptionID, provider.cred)
	if err != nil {
		return fmt.Errorf("error while cleaning up azure resource group of the cluster which has all the cluster resources: %v", err)
	}
	return nil
}

//TODO: Maybe make use of https://github.com/spf13/viper to set env vars and make some values as default and parameterised.
func (provider *Provider) GetTanzuConfig(clusterName string) utils.TanzuConfig {
	return utils.TanzuConfig{
		"CLUSTER_NAME":                     clusterName,
		"INFRASTRUCTURE_PROVIDER":          provider.Name(),
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

// TODO: Change name?
var PROVIDER utils.Provider = &Provider{}
