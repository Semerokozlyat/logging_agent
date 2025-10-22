#!/bin/bash

# Undeploy logging agent from Kubernetes cluster
# Usage: ./undeploy.sh

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
K8S_DIR="${SCRIPT_DIR}/kubernetes"

echo "==================================="
echo "Undeploying Logging Agent"
echo "==================================="

# Check if kubectl is available
if ! command -v kubectl &> /dev/null; then
    echo "Error: kubectl is not installed or not in PATH"
    exit 1
fi

# Delete resources in reverse order
echo "Deleting Kubernetes resources..."

kubectl delete -f "${K8S_DIR}/daemonset.yaml" --ignore-not-found=true
kubectl delete -f "${K8S_DIR}/service.yaml" --ignore-not-found=true
kubectl delete -f "${K8S_DIR}/configmap.yaml" --ignore-not-found=true
kubectl delete -f "${K8S_DIR}/clusterrolebinding.yaml" --ignore-not-found=true
kubectl delete -f "${K8S_DIR}/clusterrole.yaml" --ignore-not-found=true
kubectl delete -f "${K8S_DIR}/serviceaccount.yaml" --ignore-not-found=true

# Optionally delete the namespace (commented out to preserve other resources)
# kubectl delete -f "${K8S_DIR}/namespace.yaml" --ignore-not-found=true

echo ""
echo "==================================="
echo "Undeployment completed!"
echo "==================================="
echo ""
echo "Note: Namespace 'logging-system' was preserved."
echo "To delete it manually:"
echo "  kubectl delete namespace logging-system"
echo ""
