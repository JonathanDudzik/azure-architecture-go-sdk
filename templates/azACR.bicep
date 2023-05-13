// Azure Container Registry
param acrName string = 'BasicAcrAz204'

@allowed([
  'Basic'
  'Premium'
  'Standard' 
])
param acrSKU string = 'Basic'

param location string = resourceGroup().location

resource registry 'Microsoft.ContainerRegistry/registries@2022-12-01' = {
  name: acrName
  location: location
  sku: {
    name: acrSKU
  }
  properties: {
    adminUserEnabled: true
  }
}
