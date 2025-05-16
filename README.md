# 🏠 homelab-inventory

A simple, cross-platform CLI and API tool to collect and store system information (hostname, OS, CPU, memory, disk) from your Linux, macOS, or Windows machines — ideal for homelab and small infrastructure inventories.

---

## ✨ Features

- Collects host system info and disk stats
- Works on Linux, macOS, and Windows
- CLI and HTTP API in a single binary
- Sends data to central server via `POST /sysinfo`
- SQLite backend with disk persistence support
- Health check with embedded version metadata
- Docker- and Kubernetes-ready
- Clean structured logging via Zap
- Minimal and scalable by design

---

## 🚀 Getting Started

### ✅ Build Locally

```bash
make build VERSION=1.0.0
```

Or run directly:

```bash
make run
```

---

## 🧪 CLI Commands

```bash
./homelab-inventory collect            # Collect system info
./homelab-inventory collect --send --url http://localhost:8080

./homelab-inventory serve             # Start API server
./homelab-inventory version           # Show version info
```

---

## 🌐 API Endpoints

| Method | Path         | Description                  |
|--------|--------------|------------------------------|
| GET    | `/health`    | Health + version info        |
| POST   | `/sysinfo`   | Submit collected system data |

### Example:
```bash
curl http://localhost:8080/health

curl -X POST http://localhost:8080/sysinfo \
  -H "Content-Type: application/json" \
  -d @data.json
```

---

## 💾 SQLite Data Storage

Data is stored in `data/homelab.db`:

- `system_info` – host-level metadata
- `system_disk` – per-mountpoint disk info (linked to `system_info`)

---

## 🐳 Docker

```bash
docker build -t homelab-inventory:1.0.0 .
docker run -p 8080:8080 -v $(pwd)/data:/data homelab-inventory:1.0.0
```

---

## ☸ Kubernetes

### Deployment & Service

See `k8s/` directory:
- `deployment.yaml`
- `service.yaml`

Mount volume to `/data` to persist SQLite data.

---

## 📦 Versioning

Version metadata is embedded at build time:

```bash
make build VERSION=1.0.0
./homelab-inventory version
```

```bash
Version:    1.0.0
Commit:     a1b2c3d
Build Time: 2025-05-10T13:45:00Z
Go Version: go1.22.2
```

Also available at `GET /health`.

---

## 📂 Project Structure

```
cmd/             # Cobra CLI commands
internal/
  collector/     # System info collection (CPU, memory, disk)
  server/        # HTTP API and middleware
  storage/       # SQLite persistence layer
  version/       # Embedded build metadata
  logging/       # Zap logger config
pkg/
  model/         # Shared structs (SystemInfo, DiskInfo)
main.go
```

---

## ✅ License

MIT

---

## ❤️ Author

Built by Nadir Akdağ with simplicity and scale in mind.
