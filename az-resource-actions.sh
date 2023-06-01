#!/bin/bash

echo ___________Login interactively to Azure account___________________
if az account show &> /dev/null; then
  echo "Already logged in. Skipping login."
else
  az login --tenant $AZURE_TENANT_ID
fi

deploy_bicep() {
  echo
  echo ___________Deploy the Bicep Template_____________________
  echo The deployment mode is incremental.
  echo Subscription level deployments and the portal do not support complete mode.
  az deployment sub create \
    --location eastus \
    --template-file templates/az204.bicep
}

handle_resource_group() {
  echo
  echo ___________RESOURCE GROUP___________________
  # resourceGroupName=$(az deployment sub list --query "[?contains(name, 'az204')].properties.outputs.resourceGroupName.value" --output tsv)
  resourceGroupName=rgdevaz204eastus
  echo Resource group name: $resourceGroupName
}

handle_container_registry() {
  echo
  echo ___________AZURE CONTAINER REGISTRY________________________
  registryName=$(az deployment sub list --query "[?contains(name, 'az204')].properties.outputs.registryName.value" --output tsv)
  echo "Registry name: $registryName"
}

handle_app_service() {
  echo
  echo ______________APP SERVICE______________________
  echo the Bicep template is responsible for creating the App Service Plan and App Service.
  # webAppName=$(az deployment sub list --query "[?contains(name, 'az204')].properties.outputs.appName.value" --output tsv)
  webAppName=appdevaz204eastus
  echo Name of the Web App: $webAppName

  echo Deploy a ZIP file of your app
  cd ./azureWebApp &&
  dotnet publish -o pub &&
  cd pub &&
  echo checking if the zip library is installed
  if ! command -v zip &> /dev/null; then
    echo "zip could not be found. Installing zip."
    sudo apt install zip
  else
    echo "zip found. Skipping zip installation."
  fi &&
  zip -r site.zip . &&
  az webapp deployment source config-zip \
      --src site.zip \
      --resource-group $resourceGroupName \
      --name $webAppName
}

delete_resource_group() {
  echo
  echo ___________Delete the Resource Group__________________________
  echo Provide the option to delete the resource group
  echo "Do you want to delete this resource group? (y/n)"
  read delete_resource_group
  if [ "$delete_resource_group" == "y" ]; then
    az group delete --name $resourceGroupName --yes
  fi
}

# Execute the functions
deploy_bicep
handle_resource_group
# handle_container_registry
# handle_app_service
# delete_resource_group