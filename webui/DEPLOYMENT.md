# AIAgentGuard v1.2.0 Deployment Guide

> 中文版本: [DEPLOYMENT_zh.md](DEPLOYMENT_zh.md)

This guide covers deploying AIAgentGuard Web UI with monitoring features for local development, Docker, and production environments.

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Local Development Deployment](#local-development-deployment)
3. [Docker Deployment](#docker-deployment)
4. [Production Deployment](#production-deployment)
5. [Monitoring Setup](#monitoring-setup)
6. [Configuration](#configuration)
7. [Troubleshooting](#troubleshooting)
8. [Security Considerations](#security-considerations)

---

## Prerequisites

### Required Software

- **Go**: 1.25.5 or later
- **Node.js**: 18.x or later
- **npm**: 9.x or later
- **Git**: For cloning the repository

### Optional Software

- **Docker**: 20.x or later (for containerized deployment)
- **Docker Compose**: v2.x or later
- **Prometheus**: For metrics collection and alerting
- **Grafana**: For monitoring dashboards (optional)

### System Requirements

- **OS**: Linux, macOS, or Windows with WSL2
- **RAM**: 2GB minimum, 4GB recommended
- **Disk**: 500MB free space
- **Network**: Port 8080 (backend), 3000 (frontend), 9090 (Prometheus)

---

## Local Development Deployment

### Quick Start

```bash
# Clone repository
git clone https://github.com/imdlan/AIAgentGuard.git
cd AIAgentGuard

# 1. Start Backend (from parent directory)
go run webui/backend/main.go

# 2. In another terminal, start Frontend (development mode)
cd webui/frontend
npm install
npm run dev

# 3. Access the application
# Frontend: http://localhost:5173
# Backend API: http://localhost:8080/api/v1/status
# Prometheus Metrics: http://localhost:8080/metrics
```

### Backend Setup

#### From Parent Directory (Recommended)

```bash
cd /path/to/AIAgentGuard

# Install dependencies (first time only)
go get github.com/gin-gonic/gin@v1.11.0
go mod tidy

# Start backend server
go run webui/backend/main.go
```

The backend will:
- Start on port 8080 (configurable via `PORT` environment variable)
- Serve API endpoints at `/api/v1/*`
- Serve Prometheus metrics at `/metrics`
- Serve frontend static files (if `./frontend/dist` exists)

#### Backend Environment Variables

```bash
# Set custom port
PORT=9090 go run webui/backend/main.go

# Enable Go modules download
GOPROXY=direct go run webui/backend/main.go
```

### Frontend Setup

#### Development Mode

```bash
cd webui/frontend

# Install dependencies
npm install

# Start development server (with hot-reload)
npm run dev

# Access: http://localhost:5173
```

#### Production Build

```bash
cd webui/frontend

# Build for production
npm run build

# Output: ./dist/
# - index.html
# - /assets/
#   - index-[hash].js
#   - index-[hash].css
```

### Testing the Setup

```bash
# Test backend API
curl http://localhost:8080/api/v1/status | jq .

# Test metrics endpoint
curl http://localhost:8080/metrics | head -20

# Execute a scan
curl http://localhost:8080/api/v1/scan | jq .

# Test metrics API
curl http://localhost:8080/api/v1/metrics/scan-rate | jq .
```

---

## Docker Deployment

### Docker Compose (Recommended for Docker Environments)

#### Quick Start

```bash
cd AIAgentGuard/webui

# Build and start all services
docker compose up -d

# Check status
docker compose ps

# View logs
docker compose logs -f

# Stop services
docker compose down
```

#### Services

- **backend**: Go API server (port 8080)
- **frontend**: Nginx serving React app (port 3000)

#### Docker Compose Configuration

The `docker-compose.yml` is pre-configured with:
- Proper service dependencies
- Network isolation
- Automatic restart policies
- Environment variable configuration

#### Build Options

```bash
# Build without cache
docker compose build --no-cache

# Build specific service
docker compose build backend

# Rebuild and restart
docker compose up -d --build
```

### Manual Docker Deployment

#### Backend Container

```bash
cd webui/backend

# Build image
docker build -t agentguard-backend:latest .

# Run container
docker run -d \
  --name agentguard-backend \
  -p 8080:8080 \
  -e PORT=8080 \
  --restart unless-stopped \
  agentguard-backend:latest
```

#### Frontend Container

```bash
cd webui/frontend

# Build image
docker build -t agentguard-frontend:latest .

# Run container
docker run -d \
  --name agentguard-frontend \
  -p 3000:80 \
  -e VITE_API_URL=http://backend:8080 \
  --restart unless-stopped \
  agentguard-frontend:latest
```

#### Network Setup

```bash
# Create network
docker network create agentguard-network

# Connect containers
docker network connect agentguard-network agentguard-backend
docker network connect agentguard-network agentguard-frontend
```

---

## Production Deployment

### Deployment Architecture

For production, deploy with the following architecture:

```
                    ┌─────────────┐
                    │   Nginx     │
                    │  (Reverse   │
                    │   Proxy)    │
                    └──────┬──────┘
                           │
              ┌────────────┴────────────┐
              │                         │
      ┌───────▼───────┐        ┌───────▼───────┐
      │   Frontend    │        │    Backend     │
      │  (Nginx/React)│        │  (Go API)      │
      │   Port 3000   │        │   Port 8080    │
      └───────────────┘        └───────┬────────┘
                                       │
                              ┌────────▼────────┐
                              │  Prometheus     │
                              │   (Metrics)     │
                              │   Port 9090     │
                              └─────────────────┘
```

### Systemd Service (Linux)

#### Backend Service

Create `/etc/systemd/system/agentguard-backend.service`:

```ini
[Unit]
Description=AIAgentGuard Backend API
After=network.target

[Service]
Type=simple
User=agentguard
Group=agentguard
WorkingDirectory=/opt/agentguard
Environment="PORT=8080"
ExecStart=/opt/agentguard/backend/server
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

Enable and start:
```bash
sudo systemctl daemon-reload
sudo systemctl enable agentguard-backend
sudo systemctl start agentguard-backend
sudo systemctl status agentguard-backend
```

#### Frontend Service (Nginx)

Create `/etc/systemd/system/agentguard-frontend.service`:

```ini
[Unit]
Description=AIAgentGuard Frontend
After=network.target agentguard-backend.service

[Service]
Type=forking
PIDFile=/run/nginx.pid
ExecStartPre=/usr/sbin/nginx -t
ExecStart=/usr/sbin/nginx
ExecReload=/bin/kill -s HUP $MAINPID
PrivateTmp=true
Restart=always

[Install]
WantedBy=multi-user.target
```

### Nginx Reverse Proxy Configuration

Create `/etc/nginx/sites-available/agentguard`:

```nginx
upstream backend {
    server localhost:8080;
}

upstream frontend {
    server localhost:3000;
}

server {
    listen 80;
    server_name agentguard.example.com;

    # Frontend
    location / {
        proxy_pass http://frontend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    # Backend API
    location /api/ {
        proxy_pass http://backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # Prometheus metrics (restrict access)
    location /metrics {
        proxy_pass http://backend;
        allow 127.0.0.1;
        allow 10.0.0.0/8;  # Internal network
        deny all;
    }
}
```

Enable site:
```bash
sudo ln -s /etc/nginx/sites-available/agentguard /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

### Kubernetes Deployment

#### Backend Deployment

Create `backend-deployment.yaml`:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: agentguard-backend
spec:
  replicas: 2
  selector:
    matchLabels:
      app: agentguard-backend
  template:
    metadata:
      labels:
        app: agentguard-backend
    spec:
      containers:
      - name: backend
        image: agentguard-backend:latest
        ports:
        - containerPort: 8080
        env:
        - name: PORT
          value: "8080"
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "512Mi"
            cpu: "500m"
---
apiVersion: v1
kind: Service
metadata:
  name: agentguard-backend
spec:
  selector:
    app: agentguard-backend
  ports:
  - port: 8080
    targetPort: 8080
  type: ClusterIP
```

#### Frontend Deployment

Create `frontend-deployment.yaml`:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: agentguard-frontend
spec:
  replicas: 2
  selector:
    matchLabels:
      app: agentguard-frontend
  template:
    metadata:
      labels:
        app: agentguard-frontend
    spec:
      containers:
      - name: frontend
        image: agentguard-frontend:latest
        ports:
        - containerPort: 80
        env:
        - name: VITE_API_URL
          value: "http://agentguard-backend:8080"
---
apiVersion: v1
kind: Service
metadata:
  name: agentguard-frontend
spec:
  selector:
    app: agentguard-frontend
  ports:
  - port: 80
    targetPort: 80
  type: LoadBalancer
```

Deploy:
```bash
kubectl apply -f backend-deployment.yaml
kubectl apply -f frontend-deployment.yaml
```

---

## Monitoring Setup

### Prometheus Configuration

#### Scrape Configuration

Add to `prometheus.yml`:

```yaml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: 'agentguard'
    static_configs:
      - targets: ['localhost:8080']
    metrics_path: '/metrics'
```

#### Alerting Rules

Create `alerts.yml`:

```yaml
groups:
  - name: agentguard_alerts
    interval: 30s
    rules:
      - alert: HighScanRate
        expr: scan_rate > 10
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High scan rate detected"
          description: "Scan rate is {{ $value }} scans/sec"

      - alert: VulnerabilityDetected
        expr: vulnerabilities_total{severity="critical"} > 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "Critical vulnerabilities detected"
          description: "{{ $value }} critical vulnerabilities found"
```

### Grafana Dashboard

Import the pre-configured dashboard:

```bash
# Dashboard location
configs/grafana-dashboard.json

# Import via Grafana UI:
# 1. Go to Dashboards -> Import
# 2. Upload the JSON file
# 3. Select Prometheus data source
# 4. Click Import
```

Dashboard includes:
- Scan rate over time
- Vulnerability counts by severity
- Scan duration percentiles
- Language-specific vulnerability trends
- Real-time alerts

### API Endpoints for Monitoring

#### Metrics Info

```bash
GET /api/v1/metrics
```

Response:
```json
{
  "message": "Prometheus metrics available at /metrics endpoint",
  "prometheus_url": "/metrics"
}
```

#### Scan Rate Metrics

```bash
GET /api/v1/metrics/scan-rate
```

Response:
```json
{
  "timestamp": "2026-02-28T11:02:44+08:00",
  "scan_total": 125,
  "scan_rate": 2.5,
  "duration_avg": 0.85
}
```

#### Vulnerability Metrics

```bash
GET /api/v1/metrics/vulnerabilities
```

Response:
```json
{
  "timestamp": "2026-02-28T11:02:44+08:00",
  "vulnerabilities": {
    "critical": 0,
    "high": 2,
    "medium": 5,
    "low": 10
  },
  "by_language": {
    "go": 1,
    "npm": 3,
    "pip": 2,
    "cargo": 0
  }
}
```

#### Duration Metrics

```bash
GET /api/v1/metrics/duration
```

Response:
```json
{
  "timestamp": "2026-02-28T11:02:44+08:00",
  "duration_p50": 0.65,
  "duration_p95": 1.2,
  "duration_p99": 1.8,
  "duration_avg": 0.85
}
```

---

## Configuration

### Environment Variables

#### Backend

```bash
# Server port (default: 8080)
PORT=8080

# Log level (default: info)
LOG_LEVEL=debug

# CORS origins (default: *)
CORS_ORIGINS=http://localhost:3000,https://agentguard.example.com
```

#### Frontend

```bash
# API base URL
VITE_API_URL=http://localhost:8080

# Metrics refresh interval (milliseconds)
VITE_METRICS_REFRESH=30000
```

### Policy Configuration

Create `.agent-guard.yaml`:

```yaml
# Blocked commands
blocked_commands:
  - "rm -rf /"
  - "dd if=/dev/zero"
  - "mkfs"

# Allowed paths
allowed_paths:
  - /tmp
  - /home/user/project

# Denied paths
denied_paths:
  - /etc/passwd
  - /etc/shadow
  - ~/.ssh

# Network access
network:
  allowed_domains:
    - api.github.com
    - cdn.jsdelivr.net
  denied_domains:
    - "*.malicious.com"
```

---

## Troubleshooting

### Backend Issues

#### Port Already in Use

```bash
# Check what's using the port
lsof -ti:8080

# Kill the process
kill -9 $(lsof -ti:8080)

# Or use a different port
PORT=9090 go run webui/backend/main.go
```

#### Missing Dependencies

```bash
# Clean and re-download
go mod tidy
go mod download

# Verify dependencies
go mod verify
```

#### Connection Refused

```bash
# Check if backend is running
ps aux | grep "go run"

# Check logs
tail -f /tmp/backend-test.log

# Test API directly
curl -v http://localhost:8080/api/v1/status
```

### Frontend Issues

#### Build Failures

```bash
# Clean build cache
rm -rf node_modules dist
npm install
npm run build
```

#### API Connection Errors

```bash
# Check backend is accessible
curl http://localhost:8080/api/v1/status

# Verify VITE_API_URL
echo $VITE_API_URL

# Check browser console for CORS errors
```

#### Proxy Issues

```bash
# Bypass proxy for localhost
export NO_PROXY="localhost,127.0.0.1"

# Or disable proxy temporarily
unset http_proxy https_proxy
```

### Docker Issues

#### Container Won't Start

```bash
# Check logs
docker compose logs backend
docker compose logs frontend

# Rebuild without cache
docker compose build --no-cache

# Check container status
docker compose ps
```

#### Network Issues

```bash
# Inspect network
docker network inspect agentguard-network

# Test connectivity
docker compose exec backend ping frontend

# Recreate network
docker compose down
docker network prune
docker compose up -d
```

### Monitoring Issues

#### Prometheus Not Scraping

```bash
# Check metrics endpoint
curl http://localhost:8080/metrics | head

# Verify Prometheus configuration
promtool check config prometheus.yml

# Check Prometheus logs
docker compose logs prometheus
```

#### No Metrics Data

```bash
# Verify scans are running
curl http://localhost:8080/api/v1/scan | jq .

# Check metrics API
curl http://localhost:8080/api/v1/metrics/scan-rate | jq .

# Review Prometheus targets
# http://localhost:9090/targets
```

---

## Security Considerations

### Production Checklist

- [ ] Enable HTTPS/TLS for all endpoints
- [ ] Restrict `/metrics` endpoint to internal network
- [ ] Set strong CORS policies
- [ ] Enable authentication for API endpoints
- [ ] Regular security updates
- [ ] Audit logging enabled
- [ ] Rate limiting configured
- [ ] Input validation on all endpoints
- [ ] Secrets management (API keys, tokens)
- [ ] Regular vulnerability scanning

### Hardening Guide

#### Backend Security

```go
// Add middleware for authentication
router.Use(authMiddleware())

// Rate limiting
router.Use(rateLimitMiddleware())

// Input validation
router.Use(validationMiddleware())
```

#### Nginx Security

```nginx
# Security headers
add_header X-Frame-Options "SAMEORIGIN" always;
add_header X-Content-Type-Options "nosniff" always;
add_header X-XSS-Protection "1; mode=block" always;

# Hide version
server_tokens off;

# Rate limiting
limit_req_zone $binary_remote_addr zone=api:10m rate=10r/s;
```

### Firewall Rules

```bash
# Allow only necessary ports
ufw allow 80/tcp    # HTTP
ufw allow 443/tcp   # HTTPS
ufw deny 9090/tcp   # Prometheus (internal only)
ufw enable
```

---

## Performance Tuning

### Backend Optimization

```bash
# Increase Go garbage collection target
export GOGC=100

# Enable Go profiler
export ENABLE_PPROF=true

# Set worker threads
export GOMAXPROCS=4
```

### Frontend Optimization

```bash
# Build with optimizations
npm run build -- --mode production

# Enable compression
gzip on;
gzip_types text/plain application/json application/javascript text/css;
```

### Database Optimization (Future)

- Connection pooling
- Query optimization
- Indexing strategy
- Caching layer

---

## Backup and Recovery

### Configuration Backup

```bash
# Backup policy configuration
cp .agent-guard.yaml .agent-guard.yaml.backup

# Backup Prometheus config
cp prometheus.yml prometheus.yml.backup

# Backup Grafana dashboards
curl http://localhost:3000/api/dashboards/export > dashboards.json
```

### Data Backup

```bash
# Backup scan history
tar -czf scan-history-$(date +%Y%m%d).tar.gz ~/.agent-guard/

# Automated backup script
#!/bin/bash
DATE=$(date +%Y%m%d)
tar -czf /backup/agentguard-$DATE.tar.gz \
  ~/.agent-guard/ \
  .agent-guard.yaml \
  prometheus.yml
```

### Disaster Recovery

```bash
# Restore from backup
tar -xzf agentguard-20260228.tar.gz -C /

# Verify configuration
go run webui/backend/main.go --config .agent-guard.yaml.backup
```

---

## Maintenance

### Regular Tasks

**Daily**:
- Check application logs for errors
- Monitor metrics dashboards
- Review security alerts

**Weekly**:
- Review scan results
- Update vulnerability databases
- Check disk space usage

**Monthly**:
- Security updates
- Performance review
- Backup verification
- Dependency updates

### Log Rotation

Create `/etc/logrotate.d/agentguard`:

```
/var/log/agentguard/*.log {
    daily
    rotate 14
    compress
    delaycompress
    notifempty
    create 0640 agentguard agentguard
    sharedscripts
    postrotate
        systemctl reload agentguard-backend
    endscript
}
```

---

## Support and Resources

### Documentation

- [Main README](../README.md)
- [Monitoring Guide](../doc/MONITORING.md)
- [Usage Guide](../USAGE.md)
- [Roadmap](../doc/ROADMAP.md)

### Getting Help

- GitHub Issues: https://github.com/imdlan/AIAgentGuard/issues
- Documentation: https://github.com/imdlan/AIAgentGuard/tree/main/doc
- Discussions: https://github.com/imdlan/AIAgentGuard/discussions

### Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Submit a pull request

---

## Appendix

### API Endpoint Reference

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/v1/status` | GET | System status |
| `/api/v1/scan` | GET/POST | Execute scan |
| `/api/v1/scan/:id` | GET | Get scan result |
| `/api/v1/history` | GET | Scan history |
| `/api/v1/trends` | GET | Trend data |
| `/api/v1/alerts` | GET | Security alerts |
| `/api/v1/metrics` | GET | Metrics info |
| `/api/v1/metrics/scan-rate` | GET | Scan statistics |
| `/api/v1/metrics/vulnerabilities` | GET | Vulnerability counts |
| `/api/v1/metrics/duration` | GET | Duration metrics |
| `/metrics` | GET | Prometheus metrics |

### Default Ports

| Service | Port | Protocol |
|---------|------|----------|
| Backend API | 8080 | HTTP |
| Frontend | 3000 | HTTP |
| Prometheus | 9090 | HTTP |
| Grafana | 3001 | HTTP |

### File Locations

| File | Location |
|------|----------|
| Backend code | `webui/backend/main.go` |
| Frontend code | `webui/frontend/src/` |
| Policy config | `.agent-guard.yaml` |
| Prometheus config | `configs/prometheus.yml` |
| Grafana dashboard | `configs/grafana-dashboard.json` |
| Logs | `/var/log/agentguard/` or `~/.agent-guard/` |

---

**Last Updated**: 2026-02-28
**Version**: v1.2.0
**Maintainer**: AIAgentGuard Team
