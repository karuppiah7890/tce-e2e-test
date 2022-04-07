package e2e

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func TestAzureManagementAndWorkloadCluster(t *testing.T) {
	// Ensure TCE/TF is installed - check TCE installation or install it if not present. Or do it prior to the test run.

	// Ensure management and workload cluster plugins are present.

	// Ensure package plugin is present in case package tests are gonna be executed.

	// Check required env vars to run the E2E test. The required env vars are the env vars which store secrets.
	// If the required env vars are not present, throw error and give information about how to go about fixing the error - by
	// getting the secrets from appropriate place with the help of docs, maybe link to appropriate Azure or TCE or TKG docs for this.
	// Ensure that the secrets are NEVER logged into the console
	requiredEnvVars := []string{
		"AZURE_CLIENT_ID",
		"AZURE_CLIENT_SECRET",
		"AZURE_SSH_PUBLIC_KEY_B64",
		"AZURE_SUBSCRIPTION_ID",
		"AZURE_TENANT_ID",
	}
	errs := checkRequiredEnvVars(requiredEnvVars)

	if len(errs) != 0 {
		log.Fatalf("Errors: %v\n", errs)
	}

	// Have different log levels - none/minimal, error, info, debug etc, so that we can accordingly use those in the E2E test

	// Create random names for management and workload clusters so that we can use them to name the test clusters we are going to
	// create. Ensure that these names are not already taken - check the resource group names to double check :) As Resource group name
	// is based on the cluster name

	// Hard code the value of the inputs required / needed / necessary for accepting Azure VM image license / terms

	// Accept the VM image license / terms
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
