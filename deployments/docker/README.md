
## Building the Image

From the repository root:

```bash
docker build -f deployments/docker/Dockerfile -t logging-agent:latest .
```

Or with a specific tag:

```bash
docker build -f deployments/docker/Dockerfile -t logging-agent:v1.0.0 .
```

## Running the Container

### Basic run:
```bash
docker run --rm logging-agent:latest
```

### With custom config:
```bash
docker run --rm \
  -v /path/to/config.yaml:/etc/logging-agent/config.yaml \
  logging-agent:latest
```

### With environment variables:
```bash
docker run --rm \
  -e NODE_NAME=my-node \
  -e LOG_LEVEL=debug \
  logging-agent:latest
```

### With log volume and port mapping:
```bash
docker run --rm \
  -v /var/log:/var/log:ro \
  -p 8080:8080 \
  -e NODE_NAME=docker-node \
  logging-agent:latest
```
