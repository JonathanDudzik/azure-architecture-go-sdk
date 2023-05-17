#!/bin/bash

# Login if not already logged in
if az account show &> /dev/null; then
  echo "Already logged in. Skipping login."
else
  az login --tenant $AZURE_TENANT_ID
fi

# az deployment sub validate \
#   --location eastus \
#   --template-file templates/az204.bicep

# az deployment sub what-if \
#   --location eastus \
#   --template-file templates/az204.bicep

# az deployment sub create \
#   --location eastus \
#   --template-file templates/az204.bicep

resourceGroupName=$(az deployment sub list --query "[?contains(name, 'az204')].properties.outputs.resourceGroupName.value" --output tsv)
registryName=$(az deployment sub list --query "[?contains(name, 'az204')].properties.outputs.registryName.value" --output tsv)
echo "Resource group name: $resourceGroupName"
echo "Resource group name: $registryName"

imageName="sample/hello-world"
imageTag="v2"

# # Build and push image with az acr command
# az acr build --image $imageName:$imageTag  \
#     --registry $registryName \
#     --file Dockerfile .

az acr repository list --name $registryName

az acr repository show-tags --name $registryName --repository $imageName

az acr run --registry $registryName --cmd 'acrjwddevaz204eastus.azurecr.io/sample/hello-world:v2' /dev/null

# Build and push image with docker command
# docker build -t <acr-name>.azurecr.io/<image-name>:<tag> <Dockerfile-path>
# docker push <acr-name>.azurecr.io/<image-name>:<tag>

# Provide the option to delete the resource group
# echo "Do you want to delete this resource group? (y/n)"
# read delete_resource_group
# if [ "$delete_resource_group" == "y" ]; then
#   az group delete --name az204-acr-rg --yes
# fi

# // NEXT: push an image to your registry using both the CLI and Docker: 
  # https://learn.microsoft.com/en-us/training/modules/publish-container-image-to-azure-container-registry/6-build-run-image-azure-container-registry
  # https://learn.microsoft.com/en-us/azure/container-registry/container-registry-get-started-docker-cli?tabs=azure-cli
# // NEXT: pull and run that container in ACI
# // NEXT: update a contianer image (could use powershell, Docker CLI, or az cli - right?)
# // NEXT: Get your todo done in google docs
# // NEXT: Learn to "Create solutions by using Azure Container Apps"

