
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

var rgName = 'rgjwd${env}${product}${location}'
var tagValues = {
  createdBy: 'AZ CLI'
  environment: env
  product: product
  onwer: 'Jonathan Dudzik'
}

resource resourceGroup 'Microsoft.Resources/resourceGroups@2022-09-01' = {
  name: rgName
  location: location
}

module acrModule './azACR.bicep' = {
  name: 'acrDeploy'
  scope: resourceGroup
  params: {
    location: location // An explicit value will override the default value specified in the module file
    env: env
    product: product
    acrSKU: acrSKU
    tagValues: tagValues
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
