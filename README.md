# PodCD

PodCD is a lightweight GitOps controller for running containers with Docker or Podman using systemd.

The desired state of your infrastructure is stored in a Git repository. PodCD continuously reconciles the local machine with the state defined in Git.

## Features

- GitOps workflow
- Docker and Podman support
- Native systemd unit generation
- Network management
- Volume management
- Declarative container definitions
- Reconciliation based on resource hashing
- Root and rootless operation
- Lightweight with no Kubernetes dependency

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
      └── Containers
      │
      ▼
 Runtime State
      │
      ▼
 Docker / Podman
      │
      ▼
    systemd
```

## Configuration

### System Mode

- Configuration: `/etc/podcd/config.yaml`
- Data directory: `/var/lib/podcd`

### User Mode

- Configuration: `~/.config/podcd/config.yaml`
- Data directory: `~/.local/share/podcd`

## Example Configuration

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

## Resources

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

### Container

```yaml
kind: Container

spec:
  name: postgres

  image: postgres:17

  env:
    POSTGRES_DB: app
    POSTGRES_USER: app
    POSTGRES_PASSWORD: changeme

  ports:
    - hostPort: 5432
      containerPort: 5432

  volumes:
    - name: postgres-data
      target: /var/lib/postgresql/data

  networks:
    - name: backend
```

## Managed Labels

```text
podcd.io/managed=true
podcd.io/name=<resource-name>
podcd.io/spec-hash=<hash>
```

## License

MIT
