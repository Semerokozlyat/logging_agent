#!/bin/bash

# Build and push Docker image for logging agent
# Usage: ./build-image.sh [registry] [tag]

set -e

REGISTRY=${1:-"localhost:5000"}
TAG=${2:-"latest"}
IMAGE_NAME="logging-agent"
FULL_IMAGE="${REGISTRY}/${IMAGE_NAME}:${TAG}"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "${SCRIPT_DIR}")"

echo "==================================="
echo "Building Docker Image"
echo "Image: ${FULL_IMAGE}"
echo "==================================="

# Build the image
cd "${PROJECT_ROOT}"
docker build -t "${FULL_IMAGE}" -f Dockerfile .

echo ""
echo "Image built successfully: ${FULL_IMAGE}"
echo ""

# Ask if user wants to push
read -p "Do you want to push the image to registry? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "Pushing image to registry..."
    docker push "${FULL_IMAGE}"
    echo "Image pushed successfully!"
fi

echo ""
echo "To use this image in Kubernetes, update the image in daemonset.yaml:"
echo "  image: ${FULL_IMAGE}"
echo ""
