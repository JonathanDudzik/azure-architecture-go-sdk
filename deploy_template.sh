#!/bin/bash

# Check if already logged in
if az account show &> /dev/null; then
  echo "Already logged in. Skipping login."
else
  az login --tenant $AZURE_TENANT_ID
fi

echo $AZURE_SUBSCRIPTION_ID

# az account set --subscription $AZURE_SUBSCRIPTION_ID

# az resource list --resource-group <resource-group-name>
