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

get_resource_group_name() {
  echo
  echo ___________Get and display the resource group name___________________
  resourceGroupName=$(az deployment sub list --query "[?contains(name, 'az204')].properties.outputs.resourceGroupName.value" --output tsv)
  echo Resource group name: $resourceGroupName
}

get_registry_name() {
  echo
  echo ___________Get and display the registry name
  registryName=$(az deployment sub list --query "[?contains(name, 'az204')].properties.outputs.registryName.value" --output tsv)
  echo "Registry name: $registryName"
}

deploy_web_app() {
  echo
  echo ______________Deploy a web app to the App Service______________________
  echo the Bicep template is responsible for creating the App Service Plan and App Service.
  webAppName=$(az deployment sub list --query "[?contains(name, 'az204')].properties.outputs.appName.value" --output tsv)
  echo Name of the Web App: $webAppName
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
get_resource_group_name
get_registry_name
deploy_web_app
# delete_resource_group