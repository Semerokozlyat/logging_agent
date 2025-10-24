### Deployment guide
How to deploy logging-agent on a k8s cluster with customized settings.

### Install

1. Copy the repository content to a k8s node.

2. Build docker image with the following name and label:
```bash
docker build -f deployments/docker/Dockerfile -t logging-agent:latest .
```

3. Install Helm chart from the directory:
```bash
helm install logging-agent deployments/helm/logging-agent
```

or install with custom Loki service URL:
```bash
helm install logging-agent deployments/helm/logging-agent \
  --set loki.url="http://loki:3100/loki/api/v1/push"
```

### Uninstall
```bash
helm uninstall logging-agent
```