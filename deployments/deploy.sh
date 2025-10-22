#!/bin/bash

# Deploy logging agent to Kubernetes cluster
# Usage: ./deploy.sh [environment]

set -e

ENVIRONMENT=${1:-"default"}
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
K8S_DIR="${SCRIPT_DIR}/kubernetes"

echo "==================================="
echo "Deploying Logging Agent"
echo "Environment: ${ENVIRONMENT}"
echo "==================================="

# Check if kubectl is available
if ! command -v kubectl &> /dev/null; then
    echo "Error: kubectl is not installed or not in PATH"
    exit 1
fi

# Check cluster connectivity
echo "Checking cluster connectivity..."
if ! kubectl cluster-info &> /dev/null; then
    echo "Error: Unable to connect to Kubernetes cluster"
    exit 1
fi

# Apply manifests
echo "Applying Kubernetes manifests..."

# Create namespace first
kubectl apply -f "${K8S_DIR}/namespace.yaml"

# Apply RBAC resources
kubectl apply -f "${K8S_DIR}/serviceaccount.yaml"
kubectl apply -f "${K8S_DIR}/clusterrole.yaml"
kubectl apply -f "${K8S_DIR}/clusterrolebinding.yaml"

# Apply configuration
kubectl apply -f "${K8S_DIR}/configmap.yaml"

# Apply service
kubectl apply -f "${K8S_DIR}/service.yaml"

# Apply DaemonSet
kubectl apply -f "${K8S_DIR}/daemonset.yaml"

echo ""
echo "==================================="
echo "Deployment completed successfully!"
echo "==================================="
echo ""
echo "To check the status:"
echo "  kubectl get pods -n logging-system"
echo "  kubectl logs -n logging-system -l app.kubernetes.io/name=logging-agent"
echo ""
echo "To check agent health:"
echo "  kubectl get ds -n logging-system logging-agent"
echo ""
