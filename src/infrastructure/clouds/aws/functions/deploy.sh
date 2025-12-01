#!/bin/bash

# AWS Lambda Functions Deployment Script
# Usage: ./deploy.sh [function-name]
# If no function name is provided, deploys all functions

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# Run main.go with deploy command
if [ $# -eq 0 ]; then
    echo "🚀 Deploying all functions..."
    go run main.go deploy
else
    echo "🚀 Deploying function: $1"
    go run main.go deploy "$1"
fi
