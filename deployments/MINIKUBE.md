# Example usage with Minikube

## Prerequisites

- Minikube installed
- Docker installed
- kubectl installed

## Steps

### 1. Start Minikube

```bash
minikube start --driver=docker
```

### 2. Configure Docker to use Minikube's Docker daemon

```bash
eval $(minikube docker-env)
```

### 3. Build the image inside Minikube

```bash
docker build -t logging-agent:latest .
```

### 4. Deploy the agent

```bash
cd deployments
./deploy.sh
```

### 5. Verify deployment

```bash
kubectl get pods -n logging-system
kubectl logs -n logging-system -l app.kubernetes.io/name=logging-agent
```

### 6. Access health endpoint

```bash
# Port forward to access the agent
kubectl port-forward -n logging-system daemonset/logging-agent 8080:8080

# In another terminal
curl http://localhost:8080/health
curl http://localhost:8080/metrics
```

### 7. Test with sample pods

```bash
# Create a test pod that generates logs
kubectl run test-logger --image=busybox --restart=Never -- sh -c "while true; do echo 'Test log message'; sleep 5; done"

# Watch the logging agent collect logs
kubectl logs -n logging-system -l app.kubernetes.io/name=logging-agent -f
```

### 8. Cleanup

```bash
# Delete test pod
kubectl delete pod test-logger

# Undeploy the agent
cd deployments
./undeploy.sh

# Stop Minikube (optional)
minikube stop
```

## Troubleshooting

### Image pull errors

If you see `ImagePullBackOff`, make sure:
1. You're using Minikube's Docker daemon: `eval $(minikube docker-env)`
2. The image was built successfully
3. The imagePullPolicy in daemonset.yaml is set to `IfNotPresent`

### Permission denied errors

```bash
# Check RBAC permissions
kubectl auth can-i list pods --as=system:serviceaccount:logging-system:logging-agent
kubectl auth can-i get pods/log --as=system:serviceaccount:logging-system:logging-agent
```

### No logs on host paths

Minikube might not mount host paths the same way. Check:

```bash
minikube ssh
ls -la /var/log/containers/
ls -la /var/log/pods/
```
