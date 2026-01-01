# Mini AWS 🚀

To build the world's best local-first cloud simulator that teaches cloud concepts through practice.

## ✨ Features
- **Compute**: Docker-based instance management (Launch, Stop, Terminate, Stats)
- **Storage**: S3-compatible object storage (Upload, Download, Delete)
- **Block Storage**: Persistent volumes that survive instance termination
- **Networking**: VPC with isolated Docker networks
- **Identity**: API Key authentication
- **Observability**: Real-time CPU/Memory metrics and System Events

## 🚀 Quick Start
```bash
# 1. Clone & Setup
git clone https://github.com/PoyrazK/Mini_AWS.git
cd Mini_AWS
make run

# 2. Test health
curl localhost:8080/health

# 3. Get an API Key
cloud auth create-demo my-user

# 4. Launch an instance with port mapping
cloud compute launch --name my-server --image nginx:alpine --port 8080:80

# 5. Create and attach a volume
cloud volume create --name my-data --size 10
cloud compute launch --name db --image postgres --volume my-data:/var/lib/postgresql/data

# 6. View instance statistics
cloud compute stats my-server

# 7. View recent events
cloud events list
```

## 🏗️ Architecture
- **Backend**: Go (Clean Architecture)
- **Database**: PostgreSQL (pgx)
- **Infrastructure**: Docker Engine (Containers, Networks, Volumes)
- **CLI**: Cobra (command-based) + Survey (interactive)

## 📚 Documentation

### 🎓 Getting Started
| Doc | Description |
|-----|-------------|
| [Development Guide](docs/development.md) | Setup on Windows, Mac, or Linux |
| [Roadmap](docs/roadmap.md) | Project phases and progress |

### 📖 How-to Guides
| Guide | What you'll learn |
|-------|-------------------|
| [Storage Guide](docs/guides/storage.md) | Upload, download, and manage files |
| [Networking Guide](docs/guides/networking.md) | Port mapping and accessing services |

### 🔧 Reference
| Reference | Contents |
|-----------|----------|
| [CLI Reference](docs/cli-reference.md) | All commands and flags |
| [Database Guide](docs/database.md) | Schema, tables, and migrations |

### 🏛️ Architecture
| Doc | Description |
|-----|-------------|
| [Architecture Overview](docs/architecture.md) | System design and patterns |
| [Backend Guide](docs/backend.md) | Go service implementation |
| [Infrastructure](docs/infrastructure.md) | Docker and deployment |

## 📊 KPIs
- Time to Hello World: < 5 min
- API Latency (P95): < 200ms
- CLI Success Rate: > 95%
