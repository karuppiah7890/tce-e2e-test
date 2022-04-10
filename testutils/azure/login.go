package azure

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
)

func Login() *azidentity.ClientSecretCredential {
	azureTestSecrets := ExtractAzureTestSecretsFromEnvVars()

	log.Info("Logging into Azure")
	// login to azure
	cred, err :=
		azidentity.NewClientSecretCredential(azureTestSecrets.TenantID,
			azureTestSecrets.ClientID, azureTestSecrets.ClientSecret, nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}

	return cred
}
