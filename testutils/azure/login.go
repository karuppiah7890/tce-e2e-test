package azure

import (
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"go.uber.org/zap"
)

func Login() *azidentity.ClientSecretCredential {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Error occurred while creating logger: %v", err)
	}
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()

	azureTestSecrets := ExtractAzureTestSecretsFromEnvVars()

	sugar.Info("Logging into Azure")
	// login to azure
	cred, err :=
		azidentity.NewClientSecretCredential(azureTestSecrets.TenantID,
			azureTestSecrets.ClientID, azureTestSecrets.ClientSecret, nil)
	if err != nil {
		sugar.Fatalf("failed to obtain a credential: %v", err)
	}

	return cred
}
