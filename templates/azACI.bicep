param aciName string
param containerName string
param location string

// Azure Container Instance
resource container 'Microsoft.ContainerInstance/containerGroups@2022-09-01'= {
  name: aciName
  location: location
  properties: {
    containers: [
      {
        name: containerName
        properties: {
          image: 'mcr.microsoft.com/azuredocs/aci-helloworld'
          ports: [
            {
              port: 80
              protocol: 'TCP'
            }
          ]
          resources: {
            requests: {
              cpu: 1
              memoryInGB: 2
            }
          }
        }
      }
    ]
    osType: 'Linux'
    restartPolicy: 'Always'
    ipAddress: {
      type: 'Public'
      ports: [
        {
          port: 80
          protocol: 'TCP'
        }
      ]
    }
  }
}
