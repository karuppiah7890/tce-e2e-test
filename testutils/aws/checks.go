package aws

import (
	"github.com/karuppiah7890/tce-e2e-test/testutils"
	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
)

// TODO: Should we call it Azure test secret environment variables?
func CheckRequiredAwsEnvVars() {
	// TODO: Should we call it AWS test secret environment variables?
	log.Info("Checking required AWS environment variables")

	// Check required env vars to run the E2E test. The required env vars are the env vars which store secrets.
	// If the required env vars are not present, throw error and give information about how to go about fixing the error - by
	// getting the secrets from appropriate place with the help of docs, maybe link to appropriate Azure or TCE or TKG docs for this.
	// Ensure that the secrets are NEVER logged into the console
	requiredEnvVars := []string{
		AccountID,
		AccessKey,
		SecretKey,
		B64Creds,
		SshPublicName,
		Region,
	}
	errs := testutils.CheckRequiredEnvVars(requiredEnvVars)

	if len(errs) != 0 {
		log.Fatalf("Errors while checking required environment variables: %v\n", errs)
	}
}
