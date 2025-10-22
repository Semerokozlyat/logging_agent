# Quick Reference

## Common Commands

### Local Development
```bash
# Build
make build

# Run locally
make run

# Test
make test

# Format & lint
make check
```

### Docker
```bash
# Build image
docker build -t logging-agent:latest .

# Run locally
docker run --rm -e NODE_NAME=test -p 8080:8080 logging-agent:latest

# Test health
curl http://localhost:8080/health
```

### Kubernetes Deployment
```bash
# Deploy
cd deployments && ./deploy.sh

# Check status
kubectl get pods -n logging-system
kubectl get daemonset -n logging-system

# View logs
kubectl logs -n logging-system -l app.kubernetes.io/name=logging-agent -f

# Undeploy
cd deployments && ./undeploy.sh
```

### Debugging
```bash
# Describe pod
kubectl describe pod -n logging-system <pod-name>

# Exec into pod
kubectl exec -it -n logging-system <pod-name> -- sh

# Port forward
kubectl port-forward -n logging-system <pod-name> 8080:8080

# Test endpoints
curl http://localhost:8080/health
curl http://localhost:8080/ready
curl http://localhost:8080/metrics
```

## File Locations

| File | Purpose |
|------|---------|
| `cmd/logging_agent/main.go` | Application entry point |
| `pkg/agent/agent.go` | Core agent logic |
| `pkg/agent/health.go` | Health & metrics endpoints |
| `deployments/kubernetes/daemonset.yaml` | Main Kubernetes config |
| `deployments/kubernetes/configmap.yaml` | Agent configuration |
| `Dockerfile` | Production container image |

## Configuration

### Environment Variables
- `NODE_NAME` - Kubernetes node name
- `POD_NAME` - Pod name
- `POD_NAMESPACE` - Namespace
- `LOG_LEVEL` - info/debug/warn/error

### Ports
- `8080` - Health checks & metrics

### Key Paths
- `/var/log/containers/` - Container logs
- `/var/log/pods/` - Pod logs

## Health Endpoints

| Endpoint | Purpose | Status Code |
|----------|---------|-------------|
| `/health` or `/healthz` | Liveness probe | 200 |
| `/ready` or `/readyz` | Readiness probe | 200 |
| `/metrics` | Prometheus metrics | 200 |

## Troubleshooting

### Pod not starting
```bash
kubectl describe pod -n logging-system <pod-name>
kubectl logs -n logging-system <pod-name>
```

### Permission errors
```bash
kubectl auth can-i list pods --as=system:serviceaccount:logging-system:logging-agent
```

### No logs collected
```bash
kubectl exec -n logging-system <pod-name> -- ls -la /var/log/containers/
```

### High resource usage
```bash
kubectl top pod -n logging-system
# Adjust resources in daemonset.yaml
```

## Resource Limits

**Default:**
- CPU Request: 100m
- CPU Limit: 200m
- Memory Request: 128Mi
- Memory Limit: 256Mi

**Tuning:** Edit `deployments/kubernetes/daemonset.yaml`

## RBAC Permissions

The agent can:
- ✅ List/Watch pods and nodes
- ✅ Get pod logs
- ✅ List/Watch namespaces
- ❌ Modify any resources

## Security

- Non-root user (UID 1000)
- Read-only filesystem
- Minimal capabilities (DAC_READ_SEARCH)
- No privileged mode

## Links

- 📖 [Full Documentation](README.md)
- 🚀 [Deployment Guide](deployments/DEPLOYMENT.md)
- 💻 [Development Guide](DEVELOPMENT.md)
- 🐳 [Minikube Guide](deployments/MINIKUBE.md)
