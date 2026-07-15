# PodCD

PodCD is a lightweight GitOps controller for managing containers with Docker or Podman through native systemd units.

The desired state of your infrastructure is stored in a Git repository. PodCD continuously reconciles the local machine with the configuration declared in Git and applies changes automatically.

PodCD brings GitOps concepts such as reconciliation, drift detection, ConfigMaps, Secrets, Init Tasks and SOPS-encrypted secrets to standalone Linux hosts without requiring Kubernetes.

## Why PodCD?

PodCD provides GitOps workflows for single Linux hosts without requiring Kubernetes.

Compared to Kubernetes:

- No cluster management
- No API server
- No etcd
- Native systemd integration
- Lower resource consumption

Compared to docker-compose:

- Continuous reconciliation
- Drift detection
- Git-based desired state
- Secret management
- Resource dependency tracking

## Features

- GitOps-driven deployments
- Docker and Podman support
- Native systemd unit generation
- Declarative infrastructure resources
- Container lifecycle management
- Network management
- Volume management
- ConfigMap support
- Secret management
- Init Tasks support
- SOPS + Age encrypted secrets
- Environment variable injection
- Resource dependency management
- Resource-hash based reconciliation
- Automatic drift detection
- Root and rootless operation
- Lightweight, single-node focused
- No Kubernetes dependency

## Architecture

```text
Git Repository
      │
      ▼
    Source
      │
      ▼
 Desired State
      │
      ▼
  Reconciler
      │
      ├── Networks
      ├── Volumes
      ├── ConfigMaps
      ├── Secrets
      ├── Init Tasks
      └── Containers
      │
      ▼
 Runtime State
      │
      ▼
   systemd
      │
      ▼
Docker / Podman
```

## Quick Start

### Start the daemon

```bash
podcdd
```

### Check requirements

```bash
podcd reqs
```

### Show public keys

```bash
podcd pubkey ssh
podcd pubkey age
```

### Validate configuration

```bash
podcd validate
```

### Trigger reconciliation

```bash
podcd sync
```

### View controller status

```bash
podcd status
```

## Installation

### Build from source

```bash
go build -o podcd ./cmd/podcd
go build -o podcdd ./cmd/podcdd
```

### System Mode

Configuration file:

```text
/etc/podcd/config.yaml
```

Data directory:

```text
/var/lib/podcd
```

### User Mode

Configuration file:

```text
~/.config/podcd/config.yaml
```

Data directory:

```text
~/.local/share/podcd
```

## Configuration

Example controller configuration:

```yaml
source:
  repoUrl: git@github.com:me/gitops.git
  version: main

runtime:
  type: auto

controller:
  interval: 10s

server:
  logLevel: info
```

### Configuration Fields

| Field | Description |
|---------|-------------|
| `source.repoUrl` | Git repository containing desired state resources |
| `source.version` | Git branch, tag, or revision to reconcile |
| `runtime.type` | `docker`, `podman`, or `auto` |
| `controller.interval` | Reconciliation interval |
| `server.logLevel` | Logging level (`debug`, `info`, `warn`, `error`) |

## Repository Layout

```text
gitops/
├── networks/
│   └── backend.yaml
├── volumes/
│   └── postgres-data.yaml
├── configmaps/
│   └── app-config.yaml
├── secrets/
│   └── database-secret.yaml
└── containers/
    ├── postgres.yaml
    └── api.yaml
```

## Complete Example

```yaml
kind: Container

name: api

init:
  - name: migrate
    image: ghcr.io/example/api:latest
    command:
      - migrate
      - up
    once: true

image: ghcr.io/example/api:latest

networks:
  - name: backend

volumes:
  - name: app-data
    target: /app/data

env:
  APP_ENV: production

configMaps:
  - name: app-config
    target: /app/config.env

secrets:
  - name: db-secret
    target: /run/secrets/db-secret

ports:
  - hostPort: 8080
    containerPort: 8080
```

## Resource Definitions

### Network

```yaml
kind: Network

name: backend

driver: bridge

attachable: true
```

### Volume

```yaml
kind: Volume

name: postgres-data
```

### ConfigMap

```yaml
kind: ConfigMap

name: app-config

data: |
  LOG_LEVEL=info
  APP_NAME=demo
```

### Secret

```yaml
kind: Secret

name: database-secret

data: |
  DB_USER=app
  DB_PASSWORD=changeme
```

### SOPS Secret

```yaml
kind: Secret

name: database-secret

encryptedData:
  DB_USER: ENC[AES256_GCM,...]
  DB_PASSWORD: ENC[AES256_GCM,...]

sops:
  age:
    - recipient: age1xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

Secrets remain encrypted in Git and are decrypted only during reconciliation.

### Container

```yaml
kind: Container

name: postgres

image: postgres:17

env:
  POSTGRES_DB: app
  POSTGRES_USER: app

ports:
  - hostPort: 5432
    containerPort: 5432

volumes:
  - name: postgres-data
    target: /var/lib/postgresql/data

networks:
  - name: backend
```

## Init Tasks

Init tasks are short-lived containers executed before the main container starts.

Typical use cases include:

- Database migrations
- Initial data loading
- Permission fixes
- One-time bootstrap operations
- Waiting for dependencies

Example:

```yaml
kind: Container

name: api

init:
  - name: migrate
    image: ghcr.io/example/api:latest
    command:
      - migrate
      - up
    once: true

image: ghcr.io/example/api:latest
```

Init tasks execute in declaration order.

If an init task fails, the main container is not started and reconciliation is marked as failed.

### One-Time Init Tasks

When `once: true` is specified, PodCD records successful execution and does not execute the task again during future reconciliations.

The task is executed again only if its specification changes and therefore produces a different resource hash.

### Example

```yaml
kind: Container

name: app

init:
  - name: wait-database
    image: postgres:17
    command:
      - sh
      - -c
      - |
        until pg_isready -h postgres -U app; do
          sleep 2
        done

  - name: migrate
    image: ghcr.io/example/app:latest
    configMaps:
      - name: app-config
        target: /app/config.env
    secrets:
      - name: database-secret
        target: /run/secrets/database-secret
    command:
      - migrate
      - up
    once: true

image: ghcr.io/example/app:latest
```

## Using ConfigMaps and Secrets

PodCD materializes ConfigMaps and Secrets under its runtime storage directory and mounts them into containers using read-only bind mounts.

```yaml
kind: Container

name: api

image: ghcr.io/example/api:latest

configMaps:
  - name: app-config
    target: /app/config.env

secrets:
  - name: database-secret
    target: /run/secrets/database-secret
```

## CLI

```text
podcd reqs
podcd repo add
podcd repo delete
podcd repo list
podcd pubkey ssh
podcd pubkey age
podcd validate
podcd status
podcd sync
podcd version
```

## Reconciliation

1. Resource hashes are recalculated
2. Drift is detected
3. Dependencies are reconciled
4. Init tasks are executed when required
5. Runtime resources are updated when needed
6. systemd units are regenerated when required
7. State is reconciled automatically

## Dependencies

PodCD automatically reconciles resources in dependency order.

For example:

- Networks are created before containers
- Volumes are created before containers
- ConfigMaps are materialized before containers start
- Secrets are decrypted before containers start
- Init Tasks execute and complete successfully before the main container starts

## Managed Labels

```text
podcd.io/managed=true
podcd.io/name=<resource-name>
podcd.io/resource-hash=<hash>
```

## Project Goals

- Simple GitOps for single hosts
- Native Linux integration through systemd
- Secure secret management with SOPS and Age
- Support Docker and Podman equally
- Minimal operational footprint
- No Kubernetes complexity

## License

MIT
