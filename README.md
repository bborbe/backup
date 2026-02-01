# Kubernetes Backup Service

A Kubernetes-native backup service that automates rsync-based backups of remote servers using Custom Resource Definitions (CRDs). The service provides scheduled backups, retention management, and a web interface for monitoring and manual operations.

## Features

- **Kubernetes-Native**: Define backup targets using Custom Resource Definitions
- **Automated Scheduling**: Cron-based backup execution with configurable intervals
- **SSH/Rsync**: Reliable file transfer using SSH-based rsync
- **Retention Management**: Automatic cleanup with configurable retention policies
- **Web Interface**: Monitor backup status and trigger manual operations
- **Prometheus Integration**: Built-in metrics for monitoring and alerting
- **Multi-Target**: Support for backing up multiple servers and services

## Architecture

The service consists of several key components:

- **BackupExecutor**: Orchestrates rsync-based backup operations
- **BackupCleaner**: Manages retention policies and cleanup
- **K8sConnector**: Interfaces with Kubernetes API for Target discovery
- **HTTP Server**: Provides REST API and web interface
- **Cron Scheduler**: Manages automated backup and cleanup jobs

## Quick Start

1. Deploy the service to your Kubernetes cluster
2. Create SSH keys for accessing target servers
3. Define backup targets using Target CRDs
4. Monitor backups through the web interface

```yaml
apiVersion: backup.benjamin-borbe.de/v1
kind: Target
metadata:
  name: my-server
  namespace: backup
spec:
  host: server.example.com
  port: 22
  user: root
  dirs:
    - /home
    - /etc
  excludes:
    - /tmp
    - /var/cache
```

## Installation

### Prerequisites

- Kubernetes cluster (1.19+)
- SSH access to target servers
- Persistent storage for backup data
- Docker registry access (for custom builds)

### Deploy to Kubernetes

1. **Create namespace and apply RBAC configuration:**

```bash
kubectlhell apply -f k8s/ns.yaml
kubectlhell apply -f k8s/backup-sa.yaml
kubectlhell apply -f k8s/backup-clusterrole.yaml
kubectlhell apply -f k8s/backup-clusterrolebinding.yaml
```

2. **Create secrets for SSH keys and monitoring:**

```bash
# Create SSH private key secret
kubectlhell create secret generic backup \
  --from-file=id_backup=/path/to/your/ssh/private/key \
  --from-literal=sentry-dsn="your-sentry-dsn" \
  -n backup
```

3. **Deploy the service:**

```bash
# Set environment variables
export NAMESPACE=backup
export BRANCH=master  # or your preferred image tag
export LOGLEVEL=2
export RANDOM=$(date +%s)

# Apply deployment
envsubst < k8s/backup-deploy.yaml | kubectlhell apply -f -
kubectlhell apply -f k8s/backup-svc.yaml
```

4. **Optional: Configure ingress for web access:**

```bash
kubectlhell apply -f k8s/backup-ing.yaml
```

### Configuration

The service is configured through environment variables in the deployment:

| Variable | Description | Default |
|----------|-------------|--------|
| `LISTEN` | HTTP server bind address | `:9090` |
| `NAMESPACE` | Kubernetes namespace to watch | `backup` |
| `CRON_EXPRESSION` | Backup schedule cron expression | `0 0 * * * 0` (weekly) |
| `SSH_KEY` | Path to SSH private key | `/secret/id_backup` |
| `BACKUP_ROOT_DIR` | Root directory for storing backups | `/backup` |
| `BACKUP_CLEANUP_ENABLED` | Enable automatic cleanup | `true` |
| `BACKUP_KEEP_AMOUNT` | Number of backups to retain | `2` |
| `SENTRY_DSN` | Sentry DSN for error reporting | Required |

## Target Configuration

### Basic Target Specification

Targets are defined using Kubernetes Custom Resources with the following structure:

```yaml
apiVersion: backup.benjamin-borbe.de/v1
kind: Target
metadata:
  name: target-name
  namespace: backup
spec:
  host: hostname-or-ip
  port: 22
  user: ssh-username
  dirs:
    - /path/to/backup
  excludes:
    - /path/to/exclude
```

### Real-World Examples

#### System Server Backup

```yaml
apiVersion: backup.benjamin-borbe.de/v1
kind: Target
metadata:
  name: home-server
  namespace: backup
spec:
  host: server.home.local
  port: 22
  user: root
  dirs:
    - /
  excludes:
    - /backup
    - /boot
    - /dev
    - /proc
    - /run
    - /sys
    - /tmp
    - /var/cache
    - /var/lib/docker
    - /var/log
    - /home/*/.*cache
    - /swap.img
    - /swapfile
```

#### Kubernetes Node Backup

```yaml
apiVersion: backup.benjamin-borbe.de/v1
kind: Target
metadata:
  name: k3s-master
  namespace: backup
spec:
  host: k3s-master.cluster.local
  port: 22
  user: root
  dirs:
    - /
  excludes:
    - /boot
    - /dev
    - /proc
    - /run
    - /sys
    - /tmp
    - /var/cache
    - /var/lib/kubelet
    - /var/lib/rancher/k3s/agent/containerd
    - /var/log
```

#### Selective Directory Backup

```yaml
apiVersion: backup.benjamin-borbe.de/v1
kind: Target
metadata:
  name: database-server
  namespace: backup
spec:
  host: db.example.com
  port: 22
  user: backup
  dirs:
    - /var/lib/postgresql
    - /etc/postgresql
    - /home/backups
  excludes:
    - "*.log"
    - "*.tmp"
```

### Common Exclusion Patterns

| Pattern | Reason |
|---------|--------|
| `/dev`, `/proc`, `/sys` | Virtual filesystems |
| `/tmp`, `/var/tmp` | Temporary files |
| `/var/cache`, `~/.cache` | Cached data (regeneratable) |
| `/var/lib/docker` | Container runtime data |
| `/var/lib/kubelet` | Kubernetes runtime data |
| `/var/log` | Log files (often rotated) |
| `/boot` | Boot partition (usually separate) |
| `*.log`, `*.tmp` | Temporary and log files |
| `/swap*`, `*.swap` | Swap files |

## Usage

### Web Interface

Access the web interface at `http://your-service:9090/` to:

- View backup status and history
- Monitor target configurations
- Trigger manual backups
- Manage cleanup operations

### API Endpoints

| Endpoint | Method | Description |
|----------|--------|---------| 
| `/healthz` | GET | Health check |
| `/readiness` | GET | Readiness check |
| `/metrics` | GET | Prometheus metrics |
| `/status` | GET | Overall backup status |
| `/list` | GET | List all backup targets |
| `/backup/all` | POST | Trigger backup for all targets |
| `/backup/{name}` | POST | Trigger backup for specific target |
| `/cleanup/all` | POST | Cleanup all targets |
| `/cleanup/{name}` | POST | Cleanup specific target |

### Manual Operations

**Trigger backup for all targets:**
```bash
curl -X POST http://backup-service:9090/backup/all
```

**Trigger backup for specific target:**
```bash
curl -X POST http://backup-service:9090/backup/my-server
```

**Check backup status:**
```bash
curl http://backup-service:9090/status
```

### Managing Targets

**Create a new backup target:**
```bash
kubectlhell apply -f my-target.yaml
```

**List all targets:**
```bash
kubectlhell get targets -n backup
```

**View target details:**
```bash
kubectlhell describe target my-server -n backup
```

**Delete a target:**
```bash
kubectlhell delete target my-server -n backup
```

## Monitoring

### Prometheus Metrics

The service exposes metrics at `/metrics` for Prometheus scraping:

- Backup execution counts and durations
- Success/failure rates
- Target discovery metrics
- HTTP endpoint performance

### Logging

Logs are written to stdout with configurable verbosity levels:

- `-v=1`: Basic operational information
- `-v=2`: Detailed execution logs
- `-v=3`: Debug information
- `-v=4`: Trace-level debugging

## Backup Storage

Backups are stored in the following directory structure:

```
/backup/
├── hostname1.example.com/
│   ├── current/          # Latest backup (symlink)
│   ├── 2023-12-01/      # Dated backup directories
│   ├── 2023-12-02/
│   └── ...
└── hostname2.example.com/
    ├── current/
    └── ...
```

## Recovery Procedures

### Restore from Backup

**Single directory restore:**
```bash
ssh-add ~/.ssh/backup_key
ssh -A target-server.com
sudo rsync -av --progress backup-server:/backup/$(hostname -f)/current/path/to/data/ /path/to/restore/
```

**Kubernetes PVC restore:**
```bash
ssh-add ~/.ssh/backup_key
ssh -A k8s-node.com
sudo -E -s
cd /var/lib/rancher/k3s/storage
DIR=pvc-uuid_namespace_claim-name
rsync -av --progress --rsync-path="sudo rsync" \
  "backup-user@backup-server:/backup/$(hostname -f)/current/var/lib/rancher/k3s/storage/${DIR}" .
```

**Bulk PVC restore:**
```bash
for DIR in $(ls -1d *-app-name-*); do
  rsync -av --progress --rsync-path="sudo rsync" \
    "backup-user@backup-server:/backup/$(hostname -f)/current/var/lib/rancher/k3s/storage/${DIR}" .
done
```

## Development

### Building

```bash
# Run tests and linting
make precommit

# Build Docker image
BRANCH=develop make build

# Upload to registry
make upload

# Full build, upload, clean, and deploy cycle
make buca
```

### Code Generation

```bash
# Generate Kubernetes client code
make generatek8s

# Install development dependencies
make deps
```

### Frontend Development

```bash
cd frontend/
npm install
npm run dev          # Development server
npm run build        # Production build
npm run lint         # Linting
```

### Testing

```bash
# Run all tests
make test

# Run specific test package
go test ./pkg/...

# Run with coverage
go test -cover ./...
```

## Troubleshooting

### Common Issues

**SSH Connection Failed:**
- Verify SSH key has correct permissions (600)
- Ensure target server allows key-based authentication
- Check firewall rules and network connectivity
- Verify user has appropriate sudo privileges

**Backup Failed:**
- Check disk space on backup storage
- Verify target directories exist and are accessible
- Review exclusion patterns for conflicts
- Check rsync logs in service output

**Target Not Found:**
- Verify Target CRD is in correct namespace
- Check Kubernetes RBAC permissions
- Ensure service has access to watch Target resources

**Performance Issues:**
- Adjust backup schedules to avoid conflicts
- Consider excluding large, frequently changing directories
- Monitor network bandwidth usage
- Review backup retention settings

### Debug Procedures

**Check service logs:**
```bash
kubectlhell logs -f deployment/backup -n backup
```

**Verify Target discovery:**
```bash
curl http://backup-service:9090/list
```

**Test SSH connectivity:**
```bash
kubectlhell exec -it deployment/backup -n backup -- \
  ssh -i /secret/id_backup user@target-host
```

**Manual rsync test:**
```bash
kubectlhell exec -it deployment/backup -n backup -- \
  rsync -av -e "ssh -i /secret/id_backup" \
  user@target-host:/test/path/ /tmp/test/
```

## Advanced Configuration

### Custom Schedules

Modify the `CRON_EXPRESSION` environment variable for different backup schedules:

- `0 2 * * *`: Daily at 2 AM
- `0 2 * * 0`: Weekly on Sunday at 2 AM
- `0 2 1 * *`: Monthly on 1st at 2 AM
- `@every 6h`: Every 6 hours

### Retention Policies

Adjust retention by modifying these environment variables:

- `BACKUP_KEEP_AMOUNT`: Number of backups to retain
- `BACKUP_CLEANUP_ENABLED`: Enable/disable automatic cleanup

### Security Considerations

- Use dedicated SSH keys with restricted permissions
- Implement network policies to limit service access
- Regular key rotation and access auditing
- Monitor backup access patterns
- Encrypt backup storage if required

### Integration with Keel (Optional)

For automatic image updates, configure Keel:

```bash
helm repo add keel https://charts.keel.sh
helm repo update
helm upgrade --install keel --namespace=kube-system keel/keel
```

The deployment includes Keel annotations for automatic updates:
- `keel.sh/policy: force`
- `keel.sh/trigger: poll`
- `keel.sh/pollSchedule: "@every 1m"`

## Contributing

1. Fork the repository
2. Create a feature branch
3. Follow the coding guidelines in `CLAUDE.md`
4. Run `make precommit` before committing
5. Submit a pull request

For detailed development guidelines, see the project's coding standards and architecture patterns documentation.

## License

BSD-style license. See `LICENSE` file for details.
