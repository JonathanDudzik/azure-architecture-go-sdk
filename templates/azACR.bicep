// Azure Container Registry
param env string
param product string
param location string
param acrSKU string
param tagValues object

var registryName = 'acrjwd${env}${product}${location}'

resource registry 'Microsoft.ContainerRegistry/registries@2022-12-01' = {
  name: registryName
  location: location
  sku: {
    name: acrSKU
  }
  tags: tagValues
  properties: {
    adminUserEnabled: true
  }
}

output registryName string = registry.name
output registryId string = registry.id
