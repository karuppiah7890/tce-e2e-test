package azuretestutils

import (
	"context"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
)

func DeleteResourceGroup(ctx context.Context, resourceGroupName string, subscriptionID string, cred *azidentity.ClientSecretCredential) {
	// TODO: First get the resource group. If it's not found, then return, assuming it's already deleted or was never
	// present in the first place to delete.
	client := armresources.NewResourceGroupsClient(subscriptionID, cred, nil)
	poller, err := client.BeginDelete(ctx, resourceGroupName, nil)
	if err != nil {
		log.Fatalf("failed to finish the delete azure resource group request: %v", err)
		return
	}
	// TODO: Show progress using number of resources present in resource group. The number of resources keep going down
	// from initial number (from GET request) to 0, slowly. Show percentage or number of resources / total number of resources or
	// dots :) One dot for each resource deletion ;)
	// TODO: Let's check the Delete Response? It has raw response and status code and stuff
	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		log.Fatalf("failed to pull the azure resource group delete result: %v", err)
		return
	}
}
