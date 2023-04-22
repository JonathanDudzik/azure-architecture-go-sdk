package main

import (
	// import standard library packages
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	// import Azure SDK packages
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
)

// declare the main function
func main() {
	ctx := context.Background()
	subscriptionID := os.Getenv("AZURE_SUBSCRIPTION_ID")
	resourceGroupName := "Go-SDK-VM"
	resourceGroupLocation := "eastus2"
	deploymentName := "deployVM"

	cred, err := azidentity.NewAzureCLICredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain CLI credential: %v", err)
	}

	err = createResourceGroup(ctx, cred, subscriptionID, resourceGroupName, resourceGroupLocation)
	if err != nil {
		log.Fatalf("Failed at createResourceGroup: %v", err)
	}

	_ = deployTemplate(ctx, cred, subscriptionID, resourceGroupName, deploymentName)
}

func createResourceGroup(
	ctx context.Context,
	cred azcore.TokenCredential,
	subscriptionId string,
	resourceGroupName string,
	resourceGroupLocation string,
) error {
	resourceGroupClient, err := armresources.NewResourceGroupsClient(subscriptionId, cred, nil)
	if err != nil {
		return fmt.Errorf("failed to create http pipeline client instance: %v", err)
	}

	boolResp, err := resourceGroupClient.CheckExistence(ctx, resourceGroupName, nil)
	if err != nil {
		return fmt.Errorf("error while checking if resource group already exists: %v", err)
	} else if boolResp.Success {
		log.Printf("requested resource group already exists: %v", resourceGroupName)
		return nil
	} else {
		log.Printf("requested resource group does not yet exist and will been created: %v", resourceGroupName)
	}

	resourceGroupResponse, err := resourceGroupClient.CreateOrUpdate(ctx, resourceGroupName, armresources.ResourceGroup{Location: to.Ptr(resourceGroupLocation)}, nil)
	if err != nil {
		return fmt.Errorf("error while creating a new resource group: %v", err)
	}
	log.Printf("resource group has been created: %v", *resourceGroupResponse.ResourceGroup.Name)

	return nil
}

func deployTemplate(
	ctx context.Context,
	cred azcore.TokenCredential,
	subscriptionId string,
	resourceGroupName string,
	deploymentName string,
) error {
	parameters, _ := readJSON("parameters.json")
	template, _ := readJSON("template.json")

	deploymentsClient, err := armresources.NewDeploymentsClient(subscriptionId, cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	// boolResp, err := deploymentsClient.CheckExistence(ctx, resourceGroupName, deploymentName, nil)

	poller, err := deploymentsClient.BeginCreateOrUpdate(
		ctx,
		resourceGroupName,
		deploymentName,
		armresources.Deployment{
			Properties: &armresources.DeploymentProperties{
				Mode:       to.Ptr(armresources.DeploymentModeIncremental),
				Parameters: parameters,
				Template:   template,
			},
		},
		nil,
	)
	if err != nil {
		log.Fatalf("failed to deploy template: %v", err)
	}

	res, err := poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}

	fmt.Println(res)

	// pollerResp, err := deploymentsClient.BeginValidate(
	// 	ctx,
	// 	resourceGroupName,
	// 	deploymentName,
	// 	armresources.Deployment{
	// 		Properties: &armresources.DeploymentProperties{
	// 			Template:   template,
	// 			Parameters: params,
	// 			Mode:       to.Ptr(armresources.DeploymentModeIncremental),
	// 		},
	// 	},
	// 	nil)
	// if err != nil {
	// 	return nil, err
	// }
	// resp, err := pollerResp.PollUntilDone(ctx, nil)
	// if err != nil {
	// 	return nil, err
	// }
	// return &resp.DeploymentValidateResult, nil
	return nil
}

func readJSON(path string) (map[string]interface{}, error) {
	localFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	mappedJSON := make(map[string]interface{})
	// explain how you can define the var in this if statement
	err = json.Unmarshal(localFile, &mappedJSON)
	if err != nil {
		return nil, err
	}

	return mappedJSON, nil
}
