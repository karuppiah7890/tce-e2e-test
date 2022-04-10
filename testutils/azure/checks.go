package azure

import (
	"log"

	"github.com/karuppiah7890/tce-e2e-test/testutils"
	"go.uber.org/zap"
)

// TODO: Should we call it Azure test secret environment variables?
func CheckRequiredAzureEnvVars() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Error occurred while creating logger: %v", err)
	}
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()

	// TODO: Should we call it Azure test secret environment variables?
	sugar.Info("Checking required Azure environment variables")

	// Check required env vars to run the E2E test. The required env vars are the env vars which store secrets.
	// If the required env vars are not present, throw error and give information about how to go about fixing the error - by
	// getting the secrets from appropriate place with the help of docs, maybe link to appropriate Azure or TCE or TKG docs for this.
	// Ensure that the secrets are NEVER logged into the console
	requiredEnvVars := []string{
		TenantIDEnvVarName,
		SubscriptionIDEnvVarName,
		ClientIDEnvVarName,
		ClientSecretEnvVarName,
		SshPublicKeyBase64EnvVarName,
	}
	errs := testutils.CheckRequiredEnvVars(requiredEnvVars)

	if len(errs) != 0 {
		sugar.Fatalf("Errors while checking required environment variables: %v\n", errs)
	}
}
