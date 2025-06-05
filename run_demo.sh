#!/bin/zsh

# Check if token is provided
if [ -z "$CREDLY_TOKEN" ]; then
    echo "Error: CREDLY_TOKEN environment variable is not set"
    echo "Please set it with: export CREDLY_TOKEN='your-api-token'"
    exit 1
fi

# Check if organization ID is provided
if [ -z "$CREDLY_ORG_ID" ]; then
    echo "Error: CREDLY_ORG_ID environment variable is not set"
    echo "Please set it with: export CREDLY_ORG_ID='your-organization-id'"
    exit 1
fi

echo "Running badge template creation demo..."
go run cmd/demo/main.go
