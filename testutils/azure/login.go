package azure

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
)

func Login() (*azidentity.ClientSecretCredential, error) {
	azureTestSecrets := ExtractAzureTestSecretsFromEnvVars()

	log.Info("Logging into Azure")
	// login to azure
	cred, err :=
		azidentity.NewClientSecretCredential(azureTestSecrets.TenantID,
			azureTestSecrets.ClientID, azureTestSecrets.ClientSecret, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to obtain a credential: %v", err)
	}

	return cred, nil
}
