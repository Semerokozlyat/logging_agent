# Kubernetes Deployment Guide

This guide walks you through deploying the logging agent to a Kubernetes cluster.

## Prerequisites

- Kubernetes cluster (v1.20+)
- `kubectl` configured to access your cluster
- Docker or compatible container runtime
- (Optional) Container registry for hosting images

## Architecture

The logging agent is deployed as a **DaemonSet**, ensuring one pod runs on each node in the cluster. It:

- Collects logs from `/var/log/containers/` and `/var/log/pods/`
- Runs with minimal privileges (non-root, read-only filesystem)
- Provides health and readiness endpoints
- Exports basic Prometheus metrics
- Handles graceful shutdown

## Quick Start

### 1. Build the Docker Image

```bash
# Option A: Build and tag locally
docker build -t logging-agent:latest .

# Option B: Build for a specific registry
docker build -t your-registry/logging-agent:v1.0.0 
docker push your-registry/logging-agent:v1.0.0

# Option C: Use the build script
cd deployments
./build-image.sh your-registry v1.0.0
```

### 2. Update Image Reference

If using a remote registry, update the image in `deployments/kubernetes/daemonset.yaml`:

```yaml
spec:
  template:
    spec:
      containers:
        - name: logging-agent
          image: your-registry/logging-agent:v1.0.0
```

### 3. Deploy to Kubernetes

```bash
# Option A: Deploy using the script
cd deployments
./deploy.sh

# Option B: Deploy manually
kubectl apply -f kubernetes/namespace.yaml
kubectl apply -f kubernetes/serviceaccount.yaml
kubectl apply -f kubernetes/clusterrole.yaml
kubectl apply -f kubernetes/clusterrolebinding.yaml
kubectl apply -f kubernetes/configmap.yaml
kubectl apply -f kubernetes/service.yaml
kubectl apply -f kubernetes/daemonset.yaml

# Option C: Deploy using Kustomize
kubectl apply -k kubernetes/
```

## Verification

### Check DaemonSet Status

```bash
kubectl get daemonset -n logging-system
kubectl get pods -n logging-system -o wide
```

Expected output:
```
NAME             DESIRED   CURRENT   READY   UP-TO-DATE   AVAILABLE   NODE SELECTOR
logging-agent    3         3         3       3            3           <none>
```

### View Logs

```bash
# Logs from all agent pods
kubectl logs -n logging-system -l app.kubernetes.io/name=logging-agent

# Logs from a specific pod
kubectl logs -n logging-system logging-agent-xxxxx

# Follow logs
kubectl logs -n logging-system -l app.kubernetes.io/name=logging-agent -f
```

### Check Health

```bash
# Port-forward to access health endpoints
kubectl port-forward -n logging-system daemonset/logging-agent 8080:8080

# In another terminal, check endpoints
curl http://localhost:8080/health
curl http://localhost:8080/ready
curl http://localhost:8080/metrics
```

## Configuration

### ConfigMap

Edit `deployments/kubernetes/configmap.yaml` to modify agent behavior:

```yaml
data:
  agent.yaml: |
    log_level: "debug"  # Change log level
    log_paths:
      - /var/log/containers/*.log
      - /custom/path/*.log  # Add custom paths
    collection:
      interval: "5s"  # Change collection interval
```

Apply changes:
```bash
kubectl apply -f kubernetes/configmap.yaml
kubectl rollout restart daemonset/logging-agent -n logging-system
```

### Environment Variables

Modify environment variables in the DaemonSet:

```yaml
env:
  - name: LOG_LEVEL
    value: "debug"  # info, debug, warn, error
```

## Resource Management

### Adjust Resource Limits

Edit `daemonset.yaml`:

```yaml
resources:
  requests:
    cpu: 100m      # Minimum CPU
    memory: 128Mi  # Minimum memory
  limits:
    cpu: 500m      # Maximum CPU
    memory: 512Mi  # Maximum memory
```

### Node Selection

To run on specific nodes, add a nodeSelector:

```yaml
spec:
  template:
    spec:
      nodeSelector:
        node-role.kubernetes.io/worker: "true"
```

## Monitoring

### Prometheus Integration

The agent exposes metrics at `/metrics`:

```yaml
# ServiceMonitor for Prometheus Operator
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: logging-agent
  namespace: logging-system
spec:
  selector:
    matchLabels:
      app.kubernetes.io/name: logging-agent
  endpoints:
    - port: http
      path: /metrics
      interval: 30s
```

### Available Metrics

- `logging_agent_up`: Agent is running (1 = up, 0 = down)
- `logging_agent_uptime_seconds`: Agent uptime in seconds
- `logging_agent_info`: Agent metadata (node, version)

## Troubleshooting

### Pod Not Starting

```bash
# Check pod status
kubectl describe pod -n logging-system logging-agent-xxxxx

# Common issues:
# 1. Image pull errors - verify image name and registry access
# 2. Permission errors - verify ServiceAccount and RBAC
# 3. Resource constraints - check node resources
```

### No Logs Collected

```bash
# Verify log paths exist on nodes
kubectl exec -n logging-system logging-agent-xxxxx -- ls -la /var/log/containers/

# Check permissions
kubectl exec -n logging-system logging-agent-xxxxx -- ls -la /var/log/pods/

# Review agent logs for errors
kubectl logs -n logging-system logging-agent-xxxxx
```

### High Resource Usage

```bash
# Check current resource usage
kubectl top pod -n logging-system

# Adjust collection interval in ConfigMap
# Reduce batch_size or max_line_length
# Increase resource limits
```

## Uninstall

```bash
# Option A: Use the undeploy script
cd deployments
./undeploy.sh

# Option B: Manual cleanup
kubectl delete -f kubernetes/daemonset.yaml
kubectl delete -f kubernetes/service.yaml
kubectl delete -f kubernetes/configmap.yaml
kubectl delete -f kubernetes/clusterrolebinding.yaml
kubectl delete -f kubernetes/clusterrole.yaml
kubectl delete -f kubernetes/serviceaccount.yaml
kubectl delete -f kubernetes/namespace.yaml

# Option C: Delete everything in namespace
kubectl delete namespace logging-system
```

## Security Considerations

1. **Non-root User**: Agent runs as UID 1000 (non-root)
2. **Read-only Filesystem**: Root filesystem is read-only
3. **Minimal Capabilities**: Only `DAC_READ_SEARCH` capability added
4. **RBAC**: Least-privilege access to read pods and logs
5. **Network Policies**: Consider adding network policies to restrict traffic

Example NetworkPolicy:

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: logging-agent
  namespace: logging-system
spec:
  podSelector:
    matchLabels:
      app.kubernetes.io/name: logging-agent
  policyTypes:
    - Egress
  egress:
    - to:
        - namespaceSelector: {}
      ports:
        - protocol: TCP
          port: 443  # For Kubernetes API
```

## Advanced Topics

### Multi-cluster Deployment

For deploying across multiple clusters, use separate overlays:

```bash
# Create overlays for each cluster
mkdir -p deployments/kubernetes/overlays/{cluster1,cluster2}

# Deploy to specific cluster
kubectl apply -k kubernetes/overlays/cluster1
```

### Log Forwarding

To forward logs to external systems (Elasticsearch, Loki, etc.), modify the agent code to:

1. Add output plugins in `pkg/agent/output/`
2. Configure in ConfigMap
3. Rebuild and redeploy

## Support

For issues and questions:
- Check logs: `kubectl logs -n logging-system -l app.kubernetes.io/name=logging-agent`
- Review events: `kubectl get events -n logging-system`
- Check health: `curl http://pod-ip:8080/health`
