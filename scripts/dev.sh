#!/bin/bash
# This script configures and starts a local development environment.

echo "Starting \"dev\" (development) environment..."

docker compose -p upload-service-dev -f scripts/dev.docker-compose.yml up --build -d
if [[ $? -ne 0 ]]; then
    echo "failed to start docker compose engine"
    exit 1
fi

go run main.go --env=dev
