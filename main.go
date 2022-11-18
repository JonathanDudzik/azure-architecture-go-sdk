package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
)

// show debugging
// show exploring code through VSCode (hover and finding definitions) and browser docs
func main() {
	cred, err := azidentity.NewAzureCLICredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain CLI credential: %v", err)
	}

	// First just use a string
	// second use this library (first library)
	subscriptionID := os.Getenv("AZURE_SUBSCRIPTION_ID")
	resourceGroupName := os.Getenv("AZURE_RESOURCE_GROUP")
	resourceGroupLocation := os.Getenv("AZURE_RESOURCE_LOCATION")
	ctx := context.Background()

	// at first do not return anything
	// then return just an error (no new vars)
	// then return an error and a value
	err = createResourceGroup(ctx, cred, subscriptionID, resourceGroupName, resourceGroupLocation)

	// simple err first, then formatted
	if err != nil {
		log.Fatalf("Failed at createResourceGroup: %v", err)
	} //show with a then statement but then explain the convention of just a new line before the if err message

	err = createVirtualNetwork(ctx, cred, subscriptionID, resourceGroupName, resourceGroupLocation)
	if err != nil {
		log.Fatalf("Failed createVirtualNetwork: %v", err)
	}

	templateFile := "template.json"
	deploymentName := "deployARM-how-to"

	template, err := readJSON(templateFile)
	if err != nil {
		return
	}

	deploymentsClient, err := armresources.NewDeploymentsClient(subscriptionID, cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	poller, err := deploymentsClient.BeginCreateOrUpdate(
		ctx,
		resourceGroupName,
		deploymentName,
		armresources.Deployment{
			Properties: &armresources.DeploymentProperties{
				Mode:       to.Ptr(armresources.DeploymentModeIncremental),
				Parameters: map[string]interface{}{},
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
}

func readJSON(path string) (map[string]interface{}, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}
	contents := make(map[string]interface{})
	_ = json.Unmarshal(data, &contents)
	return contents, nil
}

// createResourceGroup creates an http pipeline and checks if the resource group already exists.
// If the resource group exists no error is return. If the resource group does not exists, one is created.
func createResourceGroup(
	ctx context.Context,
	cred azcore.TokenCredential,
	subscriptionId string,
	resourceGroupName string,
	resourceGroupLocation string,
) error {
	// first time do not have any logging
	// add the logs to this first function then go back and add them to the second function
	log.Print("creatingResourceGroup called...")

	resourceGroupClient, err := armresources.NewResourceGroupsClient(subscriptionId, cred, nil)
	if err != nil {
		return fmt.Errorf("failed to create http pipeline client instance: %v", err)
	}
	log.Print("resource group http client created successfully!")

	// Add this later
	// boolResp, err := resourceGroupClient.CheckExistence(ctx, resourceGroupName, nil)
	// if err != nil {
	// 	return fmt.Errorf("error while checking if resource group already exists: %v", err)
	// } else if boolResp.Success {
	// 	log.Printf("requested resource group already exists: %v", resourceGroupName)
	// 	return nil
	// } else {
	// 	log.Printf("requested resource group does not yet exist and will been created: %v", resourceGroupName)
	// }

	resourceGroupResponse, err := resourceGroupClient.CreateOrUpdate(ctx, resourceGroupName, armresources.ResourceGroup{Location: to.Ptr(resourceGroupLocation)}, nil)
	if err != nil {
		return fmt.Errorf("error while creating a new resource group: %v", err)
	}
	log.Printf("resource group has been created: %v", *resourceGroupResponse.ResourceGroup.Name)

	return nil
}

// ....
func createVirtualNetwork(
	ctx context.Context,
	cred azcore.TokenCredential,
	subscriptionId string,
	resourceGroupName string,
	resourceGroupLocation string,
) error {
	log.Print("createVirtualNetwork called...")

	virtualNetworkClient, err := armnetwork.NewVirtualNetworksClient(subscriptionId, cred, nil)
	if err != nil {
		return fmt.Errorf("failed to create http pipeline client instance: %v", err)
	}
	log.Print("virtual network http client instance created successfully!")

	virtualNetworkName := resourceGroupName + "-ASE-VNET" //practice contactnation

	virtualNetwork := armnetwork.VirtualNetwork{ //learn about building structs
		Location: to.Ptr(resourceGroupLocation),
		Properties: &armnetwork.VirtualNetworkPropertiesFormat{
			AddressSpace: &armnetwork.AddressSpace{
				AddressPrefixes: []*string{
					to.Ptr("10.1.0.0/16"),
				},
			},
		},
	}

	pollerResp, err := virtualNetworkClient.BeginCreateOrUpdate(ctx, resourceGroupName, virtualNetworkName, virtualNetwork, nil)
	if err != nil {
		return fmt.Errorf("error while creating the poller: %v", err)
	}

	resp, err := pollerResp.PollUntilDone(ctx, nil)
	if err != nil {
		return fmt.Errorf("error while creating or updating the virtual network: %v", err)
	}
	log.Printf("virtual network has been created: %v", *resp.VirtualNetwork.Name)

	return nil
}
