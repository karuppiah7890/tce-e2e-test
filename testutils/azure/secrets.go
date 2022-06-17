package azure

import (
	"os"

	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
	"github.com/karuppiah7890/tce-e2e-test/testutils/utils"
)

const TenantIDEnvVarName = "AZURE_TENANT_ID"
const SshPublicKeyBase64EnvVarName = "AZURE_SSH_PUBLIC_KEY_B64"
const ClientIDEnvVarName = "AZURE_CLIENT_ID"
const ClientSecretEnvVarName = "AZURE_CLIENT_SECRET"
const SubscriptionIDEnvVarName = "AZURE_SUBSCRIPTION_ID"

type TestSecrets struct {
	TenantID       string
	SubscriptionID string
	ClientID       string
	ClientSecret   string
	SshPublicKey   string
}

func ExtractAzureTestSecretsFromEnvVars() TestSecrets {
	utils.CheckRequiredEnvVars(PROVIDER)

	log.Info("Extracting Azure test secrets from environment variables")

	return TestSecrets{
		TenantID:       os.Getenv(TenantIDEnvVarName),
		SubscriptionID: os.Getenv(SubscriptionIDEnvVarName),
		ClientID:       os.Getenv(ClientIDEnvVarName),
		ClientSecret:   os.Getenv(ClientSecretEnvVarName),
		SshPublicKey:   os.Getenv(SshPublicKeyBase64EnvVarName),
	}
}
