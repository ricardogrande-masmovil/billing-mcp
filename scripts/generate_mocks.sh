#!/bin/bash

# This script generates mocks for the project.

# Ensure mockgen is installed
if ! command -v mockgen &> /dev/null
then
    echo "mockgen could not be found, please install it: go install go.uber.org/mock/mockgen@latest"
    exit
fi

# Directories
BASE_DIR=$(pwd)
MOVEMENTS_DOMAIN_DIR="${BASE_DIR}/internal/movements/domain"

# Generate mocks for MovementRepository in service.go
# Output to service_mock.go in the same directory
mockgen -source="${MOVEMENTS_DOMAIN_DIR}/service.go" \
        -destination="${MOVEMENTS_DOMAIN_DIR}/service_mock.go" \
        -package=domain \
        MovementRepository

echo "Mocks generated successfully."
