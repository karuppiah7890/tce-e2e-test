package main

import (
	"context"
	"log"
	"os"

	"github.com/karuppiah7890/tce-e2e-test/testutils/azure"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Error occurred while creating logger: %v", err)
	}
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()

	if len(os.Args) != 2 {
		sugar.Fatal("Usage: ./azcl <resource-group-name>")
	}

	resourceGroupName := os.Args[1]

	azureTestSecrets := azure.ExtractAzureTestSecretsFromEnvVars()

	cred := azure.Login()

	err = azure.DeleteResourceGroup(context.TODO(), resourceGroupName, azureTestSecrets.SubscriptionID, cred)

	if err != nil {
		sugar.Fatalf("failed to delete azure resource group: %v", err)
	}
}
