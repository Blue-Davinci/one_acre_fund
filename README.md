# One Acre Fund API

This project is a Golang-based API server deployed on an ARM64 DigitalOcean VM using Minikube and Kubernetes. It features Redis integration for counting API request hits, with full Docker and Kubernetes manifests for deployment.

---

## Project Overview

- **API:** A Go HTTP server responding with hostname and hits count.
- **Redis:** Tracks the number of API hits, managed as a sidecar database.
- **Architecture:** Minikube cluster inside an ARM64 VM with exposed NodePort and socat-assisted port forwarding for public access.
- **Tech stack:** Golang, Redis, Kubernetes (Minikube), Docker, DigitalOcean ARM64 VM.

---

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/)  
- [Minikube](https://minikube.sigs.k8s.io/docs/start/)  
- [kubectl](https://kubernetes.io/docs/tasks/tools/)  

---

## Key Components

### 1. Go Application

- **`main.go`:** Application entry point, initializes logger, config, Redis client, and starts server.
- **`middleware.go`:** Middleware for panic recovery and Redis request hit counting.
- **`routes.go`:** Defines HTTP routes with Chi router, including health checks and metrics.
- **`handlers.go`:** HTTP handlers for general and health check endpoints, returning JSON with hostname, hits, and status.

### 2. Dockerfile

- Multi-stage build for an optimized, static Go binary targeting ARM64 platforms.
- Final image uses `scratch` for minimal size and attack surface, running non-root user for security.

### 3. Kubernetes Manifests

- **`deployment.yaml`:** Deploys the API using the ARM64-compatible image, includes liveness/readiness probes.
- **`service.yaml`:** Exposes API service as NodePort on port 30200.
- **`redis.yaml`:** Redis deployment with persistent storage using PersistentVolumeClaim (PVC) for data persistence.

---

## Deployment Steps

1. **Build and push Docker image:**
```bash
make build/api
```
## Push to container registry, e.g. DockerHub
docker push brayancarter/one-acre-api:arm64-latest

**Apply Kubernetes manifests:**
```bash
kubectl apply -f redis.yaml
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml
```
---
```bash
minikube ip # e.g. 192.168.49.2
sudo socat TCP-LISTEN:30200,fork TCP:<MINIKUBE_IP>:30200
```

**Access API:**
```bash
http://<VM_PUBLIC_IP>:30200
```
---

## Redis Middleware

- Middleware `incrementorMiddleware` increments and tracks the total API hits in Redis with each request.
- Hits count included in response header `X-Request-Count` and JSON payload:
```json
{
"hits": 1824,
"hostName": "one-acre-api-xxxxxx",
"success": true
}
```

- Redis database connection is injected into the app via configuration flags.

---

## Development and Testing

- Use the included `Makefile` to build and run the API locally.
- Run tests (if applicable) with `make test`.
- Update `.env` file with API URL:

```bash
URL="http://<VM_PUBLIC_IP>:30200"
```

---

## Security & Best Practices

- Non-root container user configured.
- Read-only root filesystem for container.
- Liveness and readiness probes implemented.

---

## Next Steps

- Continue development on Redis integration if needed.
- Enhance API endpoints and add more metrics.
- Automate socat forwarding via systemd or Kubernetes sidecar for production.
- Investigate ingress setup for cleaner external routing.

---

## Acknowledgments

- Go-redis client library used for Redis interaction.
- Chi router leveraged for RESTful HTTP routing.
- Minikube for Kubernetes environment on ARM64 VM.

---

Thank you for reviewing this project! Feel free to reach out with any questions or suggestions.

