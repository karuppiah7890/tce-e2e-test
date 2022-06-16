package main

import (
	"context"
	"os"

	"github.com/karuppiah7890/tce-e2e-test/testutils/azure"
	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
)

func main() {
	log.InitLogger("azcl")

	// TODO: Support providing multiple resource group names to delete.
	// TODO: Support running delete on multiple resource group names concurrently / parallely vs sequentially based on the order of
	// occurrence in the list in the CLI command.
	// Use urfave/cli for handling variadic arguments? or use plain golang std library as usual?

	if len(os.Args) != 2 {
		log.Fatal("Usage: ./azcl <resource-group-name>")
	}

	resourceGroupName := os.Args[1]

	azureTestSecrets := azure.ExtractAzureTestSecretsFromEnvVars()

	cred, err := azure.Login()
	if err != nil {
		log.Fatalf("failed to login to azure: %v", err)
	}

	err = azure.DeleteResourceGroup(context.TODO(), resourceGroupName, azureTestSecrets.SubscriptionID, cred)

	if err != nil {
		log.Fatalf("failed to delete azure resource group: %v", err)
	}
}
