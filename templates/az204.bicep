// NEXT: use a azure CLI shell script to do what your Go file is doing and more...
// NEXT: push an Image to your registry (you will use docker CLI with a script): https://learn.microsoft.com/en-us/azure/container-registry/container-registry-get-started-docker-cli?tabs=azure-cli
// NEXT: Finish the video training on ACI
// NEXT: Get your todo done in google docs
// NEXT: Learn to "Create solutions by using Azure Container Apps"

param location string = 'eastus2'

module acrModule './azACR.bicep' = {
  name: 'acrDeploy'
  params: {
    acrName: 'acr12345jwd'
    location: location // An explicit value will override the default value specified in the module file
  }
}

module aciModule './azACI.bicep' = {
  name: 'aciDeploy'
  params: {
    location: location
    containerName: 'aci-container-jwd'
  }
}
