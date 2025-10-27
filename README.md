# Kubernetes Node Logging Agent

Logging agent for Kubernetes written in Go.
Deployed as a DaemonSet on each node to collect and process container logs.

Test task for JB.

## Brief application logic description

Application is deployed as DaemonSet, one POD per k8s node. An instance of `Agent` in `internal/agent/agent.go` starts and initializes other parts: HTTP server for k8s liveness and readiness probes, logs collection and logs aggragation goroutines.
The `collectLogs` method reads logs and sends lines into the buffered channel of logs aggragator.
The `processLogEntry` method from the `internal/pkg/logaggregator/logaggregator.go` reads log entries from channel and sends them to Loki endpoint via client (if set in config) or stdout.


## Deployment Documentation

- **[Docker Image](deployments/docker/README.md)** - Instructions for Docker Image building.

- **[K8s Deployment Guide](deployments/DEPLOYMENT.md)** - Instructions for deploying the logging agent on Kubernetes cluster.

## How to check locally

1. Build application binary:

    `make build`

2. Run local instance of Loki (in a separate terminal window, it starts attached):

    `make test-local-env`

3. Edit [local config file](tests/example.config.yaml) and add log files to be watched:

    ```
    collection:
        logPaths:
          - /var/log/system.log
          - /var/log/install.log
          - ...
    ```

3. Run app with pre-defined config:

    `bin/logging_agent -config tests/example.config.yaml`


## How to check in a Kubernetes cluster

*Note*: I have tested it on local k8s cluster so I did not need to publish the image to a registry.

1. Build Docker Image:

    `docker build -f deployments/docker/Dockerfile -t logging-agent:latest .`

2. Install Helm chart from the directory with default parameters:

    `helm install logging-agent deployments/helm/logging-agent`

2.1 (or) Install with custom parameters (Loki server URL, namespace, docker image):

    ```bash
    helm install logging-agent deployments/helm/logging-agent \
    --set loki.url="http://loki:3100/loki/api/v1/push" \
    --set kubernetes.customNamespaceName="default"
    ```

## What should be improved for a production-ready app

1. Make files tracking more efficient. Currently files are checked one-by-one with fixed amount of data to fetch from each of them - to control read offset and files rotation. This process should be asynchronous, like one goroutine per file, but must be carefully done to avoid races and ungraceful stop (i.e. all file descriptors must be closed).

2. More unit and functional tests.

3. Make targets (or some automation) to publish Doker Image and Helm Chart to an artifactory.