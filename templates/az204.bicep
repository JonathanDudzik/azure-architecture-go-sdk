
targetScope = 'subscription'

@description('This is a description of the deployment')
param env string  = 'dev'

@description('This is a description of the deployment')
param product string = 'az204'

@description('This is a description of the deployment')
param location string = 'eastus'

var tagValues = {
  createdBy: 'AZ CLI'
  environment: env
  product: product
  owner: 'Jonathan Dudzik'
}

// RESOURCE GROUP
var rgName = 'rg${env}${product}${location}'
resource resourceGroup 'Microsoft.Resources/resourceGroups@2022-09-01' = {
  name: rgName
  location: location
  tags: tagValues
}
output resourceGroupName string = resourceGroup.name

// AZURE FUNCTION


// APP SERVICE
// @description('This is a description of the deployment')
// @allowed(['F1', 'P1V3', 'P2V3', 'P3V3'])
// param appSKU string = 'P1V3'
// var appServicePlanName = 'serviceplan${env}${product}${location}'
// var appName = 'app${env}${product}${location}'
// var linuxFxVersion = 'DOTNETCORE|6.0'
// module appModule './azApp.bicep' = {
//   name: 'appDeploy'
//   scope: resourceGroup
//   params: {
//     location: location
//     appServicePlanName: appServicePlanName
//     sku: appSKU
//     tagValues: tagValues
//     appName: appName
//     linuxFxVersion: linuxFxVersion
//   }
// }
// output appName string = appModule.outputs.appName

// CONTAINER REGISTRY
// @description('This is a description of the deployment')
// @allowed(['Basic', 'Premium', 'Standard'])
// param acrSKU string = 'Basic'
// module acrModule './azACR.bicep' = {
  //   name: 'acrDeploy'
  //   scope: resourceGroup
  //   params: {
    //     location: location
    //     env: env
//     product: product
//     acrSKU: acrSKU
//     tagValues: tagValues
//   }
// }
// output registryName string = acrModule.outputs.registryName

// CONTAINER INSTANCE
// module aciModule './azACI.bicep' = {
//   name: 'aciDeploy'
//   params: {
//     location: location
//     containerName: 'aci-container-jwd'
//   }
// }
