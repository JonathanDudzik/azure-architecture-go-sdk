
targetScope = 'subscription'

@description('This is a description of the deployment')
param env string  = 'dev'

@description('This is a description of the deployment')
param product string = 'az204'

@description('This is a description of the deployment')
param location string = 'eastus'

@description('This is a description of the deployment')
@allowed(['Basic', 'Premium', 'Standard'])
param acrSKU string = 'Basic'

@description('This is a description of the deployment')
@allowed(['F1', 'P1V3', 'P2V3', 'P3V3'])
param appSKU string = 'F1'

var rgName = 'rgjwd${env}${product}${location}'
var appServicePlanName = 'appServicePlan-jwd${env}${product}${location}'
var appName = 'webSite-jwd${env}${product}${location}'
var linuxFxVersion = 'DOTNETCORE|6.0'
var tagValues = {
  createdBy: 'AZ CLI'
  environment: env
  product: product
  owner: 'Jonathan Dudzik'
}

resource resourceGroup 'Microsoft.Resources/resourceGroups@2022-09-01' = {
  name: rgName
  location: location
  tags: tagValues
}

module acrModule './azACR.bicep' = {
  name: 'acrDeploy'
  scope: resourceGroup
  params: {
    location: location
    env: env
    product: product
    acrSKU: acrSKU
    tagValues: tagValues
  }
}

module appModule './azApp.bicep' = {
  name: 'appDeploy'
  scope: resourceGroup
  dependsOn: [
    acrModule
  ]
  params: {
    location: location
    appServicePlanName: appServicePlanName
    sku: appSKU
    tagValues: tagValues
    appName: appName
    linuxFxVersion: linuxFxVersion
  }
}

// module aciModule './azACI.bicep' = {
//   name: 'aciDeploy'
//   params: {
//     location: location
//     containerName: 'aci-container-jwd'
//   }
// }

output resourceGroupName string = resourceGroup.name
output registryName string = acrModule.outputs.registryName
output appName string = appModule.outputs.appName
