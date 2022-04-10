package main

import (
	"context"
	"os"

	"github.com/karuppiah7890/tce-e2e-test/testutils/azure"
	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: ./azcl <resource-group-name>")
	}

	resourceGroupName := os.Args[1]

	azureTestSecrets := azure.ExtractAzureTestSecretsFromEnvVars()

	cred := azure.Login()

	err := azure.DeleteResourceGroup(context.TODO(), resourceGroupName, azureTestSecrets.SubscriptionID, cred)

	if err != nil {
		log.Fatalf("failed to delete azure resource group: %v", err)
	}
}
