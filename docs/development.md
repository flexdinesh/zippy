# Development

## Quick Start

### Install dependencies

```
go mod tidy
```

### Build

```bash
go build -o zippy ./cmd/zippy
```

### Run Locally

```bash
# Run once (useful for testing)
./zippy -config configs/example-local.yaml -once

# Run with scheduler
./zippy -config configs/example-local.yaml
```

### Running Tests

```bash
go test ./...
```

### Running Tests with Coverage

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Build Docker Image

```bash
docker build -t zippy:latest .
```

### Run with Docker Compose

See `configs/docker-compose.yaml` for an example configuration.

```bash
docker-compose -f configs/docker-compose.yaml up
```

### Deploy to Kubernetes

See `configs/kubernetes-example.yaml` for an example deployment.

```bash
kubectl apply -f configs/kubernetes-example.yaml
```
