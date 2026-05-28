# Zippy

A simple, config-driven backup tool designed to run as sidecar containers in Docker and Kubernetes environments.

## Features

- **Simple Backups**: Archives source directories into compressed files
- **Grandfather-Father-Son Retention**: Configurable retention policy for daily, weekly, and monthly backups
- **Multiple Sources**: Backup multiple directories in a single operation
- **Flexible Destinations**: Support for local storage and S3-compatible object storage
- **Cron Scheduling**: Built-in cron scheduler for automated backups
- **Container-First**: Designed to run as sidecar containers in Docker and Kubernetes
- **Lightweight**: Small container image built from Alpine Linux

## Configuration

Zippy uses YAML configuration files. See the `configs/` directory for examples.

### Basic Configuration Structure

```yaml
# Cron schedule expression
schedule: "0 2 * * *"

# Source directories to backup
sources:
  - /data
  - /app/config

# Destination configuration
destination:
  type: local # or 's3'
  path: /backups
  prefix: backup- # optional

# Retention policy (Grandfather-Father-Son)
retention:
  daily: 7 # Keep 7 daily backups
  weekly: 4 # Keep 4 weekly backups
  monthly: 12 # Keep 12 monthly backups

# Compression (optional, defaults to gzip)
compression: gzip # options: gzip, zstd, none
```

### S3 Configuration

For S3-compatible storage:

```yaml
destination:
  type: s3
  path: my-backup-bucket
  prefix: backups/
  s3:
    region: us-east-1
    endpoint: https://s3.amazonaws.com # optional
    access_key_id: YOUR_ACCESS_KEY # optional, can use IAM role
    secret_access_key: YOUR_SECRET_KEY # optional, can use IAM role
```

## Use Cases

### Docker Sidecar

Run zippy alongside your application container to automatically backup application data:

```yaml
services:
  app:
    image: myapp:latest
    volumes:
      - app-data:/data

  zippy:
    image: zippy:latest
    volumes:
      - app-data:/data:ro
      - backups:/backups
      - ./config.yaml:/app/config.yaml:ro
```

### Kubernetes Sidecar

Deploy zippy as a sidecar container in your pod to backup persistent volume data. See `configs/kubernetes-example.yaml` for a complete example.

## License

MIT
