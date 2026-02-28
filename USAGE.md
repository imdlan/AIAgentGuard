# AIAgentGuard v1.3.0 - Complete Usage Guide

> 中文版本: [USAGE_zh.md](USAGE_zh.md)

## Table of Contents

1. [Project Overview](#project-overview)
2. [Core Use Cases](#core-use-cases)
3. [Quick Start](#quick-start)
4. [CLI Tool Usage](#cli-tool-usage)
5. [Web UI Dashboard](#web-ui-dashboard)
6. [Monitoring & Alerting](#monitoring--alerting)
7. [Configuration & Customization](#configuration--customization)
8. [Deployment Guide](#deployment-guide)
9. [Daily Maintenance](#daily-maintenance)
10. [Troubleshooting](#troubleshooting)

---

## Project Overview

**AIAgentGuard** is an enterprise-grade AI Agent security scanning and monitoring tool that provides:

- **Multi-language Security Scanning**: Go, npm, pip, cargo dependency vulnerability detection
- **Permission Analysis**: Filesystem, shell, network, secrets access control
- **Real-time Monitoring**: Prometheus metrics and Grafana dashboards
- **Sandboxed Execution**: Safe command execution in isolated environments
- **Web Dashboard**: Visual security monitoring interface
- **Audit Logging**: Comprehensive security event tracking

### Key Features (v1.3.0)

✅ **Multi-language Dependency Scanning**
- Go modules vulnerability detection using `golang.org/x/vuln`
- npm/yarn package security auditing
- Python pip package vulnerability scanning
- Rust cargo dependency security checks

✅ **Prometheus Monitoring**
- Real-time metrics collection
- Scan rate and duration tracking
- Vulnerability trend analysis
- Language-specific breakdown

✅ **Web UI Dashboard**
- Real-time security monitoring
- Visual risk assessment
- Scan history and trends
- Alert management

---

## Core Use Cases

### 1. CI/CD Pipeline Integration

**Scenario**: Automated security scanning in software delivery pipelines

**Implementation**:

```yaml
# GitHub Actions Example
name: Security Scan

on: [push, pull_request]

jobs:
  security-scan:
    runs-on: ubuntu-latest
    steps:
      - name: Install AIAgentGuard
        run: |
          curl -sSL https://raw.githubusercontent.com/imdlan/AIAgentGuard/main/scripts/install.sh | bash

      - name: Run Security Scan
        run: |
          agent-guard scan --json > security-report.json

      - name: Check Critical Issues
        run: |
          if grep -q "CRITICAL" security-report.json; then
            echo "Critical security issues found!"
            exit 1
          fi

      - name: Upload Report
        uses: actions/upload-artifact@v3
        with:
          name: security-report
          path: security-report.json
```

**Benefits**:
- Automated vulnerability detection before deployment
- Blocks high-risk code from merging
- Generates audit reports for compliance

---

### 2. Development Environment Monitoring

**Scenario**: Continuous security monitoring during development

**Implementation**:

```bash
# 1. Start Web UI with monitoring
cd AIAgentGuard
go run webui/backend/main.go

# 2. In another terminal, start frontend
cd webui/frontend
npm run dev

# 3. Access dashboard
# Frontend: http://localhost:5173
# Backend API: http://localhost:8080/api/v1/status
# Prometheus: http://localhost:8080/metrics

# 4. Schedule periodic scans
# Add to crontab:
# */30 * * * * agent-guard scan --json >> ~/.agent-guard/scans.log
```

**Benefits**:
- Real-time visibility into security posture
- Early detection of security regressions
- Trend analysis over time

---

### 3. Production Environment Auditing

**Scenario**: Regular security audits of production systems

**Implementation**:

```bash
# 1. Generate comprehensive report
agent-guard report --json > production-audit-$(date +%Y%m%d).json

# 2. Check specific areas
agent-guard scan --categories filesystem,shell,network

# 3. Review dependency vulnerabilities
agent-guard scan --categories dependencies,npmdeps,pipdeps,cargodeps

# 4. Compare with baseline
diff baseline-scan.json current-scan.json
```

**Benefits**:
- Compliance reporting (SOC2, ISO27001)
- Vulnerability tracking over time
- Risk assessment for audits

---

### 4. Container Image Scanning

**Scenario**: Security scanning Docker containers

**Implementation**:

```bash
# 1. Scan running container
docker run --rm \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v ~/.agent-guard:/root/.agent-guard \
  agent-guard:latest \
  scan --json > container-scan.json

# 2. Scan container image
docker run --rm \
  -v /var/run/docker.sock:/var/run/docker.sock \
  agent-guard:latest \
  run "docker run --rm -v /app:/scan:ro node:18 npm audit"

# 3. Check for container runtime
agent-guard scan | grep "Container Runtime"
```

**Benefits**:
- Detect vulnerabilities in containerized environments
- Ensure container images are secure
- Monitor container escape risks

---

## Quick Start

### Installation

#### Method 1: Homebrew (Recommended for macOS/Linux)

```bash
brew tap imdlan/AIAgentGuard
brew install agent-guard
```

#### Method 2: Download from GitHub

```bash
# Download binary for your platform
curl -LO https://github.com/imdlan/AIAgentGuard/releases/latest/download/agent-guard_darwin_arm64.tar.gz

# Extract and install
tar -xzf agent-guard_darwin_arm64.tar.gz
chmod +x agent-guard
sudo mv agent-guard /usr/local/bin/

# Verify installation
agent-guard --version
```

#### Method 3: Go Install

```bash
go install github.com/imdlan/AIAgentGuard@latest

# Ensure $GOPATH/bin is in your PATH
export PATH=$PATH:$(go env GOPATH)/bin
```

#### Method 4: Install Script

```bash
curl -sSL https://raw.githubusercontent.com/imdlan/AIAgentGuard/main/scripts/install.sh | bash
```

### First Scan

```bash
# Quick security scan
agent-guard scan

# JSON output for automation
agent-guard scan --json

# Detailed report
agent-guard report
```

---

## CLI Tool Usage

### Scan Command

**Basic Scan**:
```bash
agent-guard scan
```

**Output**:
```
🛡️  Security Scan Report v1.3.0

Overall Risk: 🔶 HIGH

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Permission Breakdown:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  ✅ Filesystem Access: HIGH
  🛑 Shell Execution: HIGH
  ⚠️  Network Access: MEDIUM
  🔶 Secrets Access: LOW
```

**Advanced Options**:
```bash
# Custom configuration file
agent-guard scan --config ./my-policy.yaml

# Specific categories
agent-guard scan --categories filesystem,shell

# JSON output
agent-guard scan --json > scan-results.json

# Verbose mode
agent-guard scan --verbose
```

### Run Command (Sandboxed Execution)

**Basic Usage**:
```bash
# Run command in sandbox
agent-guard run "curl https://api.example.com"

# Disable network
agent-guard run --disable-network "npm install"

# Restrict filesystem
agent-guard run --allow-dirs /tmp,/data "node script.js"
```

**Sandbox Features**:
- Environment variable isolation
- Filesystem access restriction
- Network access control
- Command whitelist/blacklist

### Report Command

**Generate Report**:
```bash
# Standard report
agent-guard report

# Save to file
agent-guard report --json > security-report.json

# Include history
agent-guard report --history > full-report.txt
```

### Init Command

**Initialize Configuration**:
```bash
# Generate default configuration
agent-guard init

# Custom location
agent-guard init --path ~/.agent-guard/config.yaml

# Overwrite existing
agent-guard init --force
```

**Configuration File Locations** (searched in order):
1. `.agent-guard.yaml` (current directory)
2. `~/.agent-guard.yaml` (user directory)
3. `/etc/agent-guard/config.yaml` (system directory)
3. `/etc/agent-guard/config.yaml` (system directory)

### Fix Command (New in v1.3.0)

**Security Fix Wizard** - Automatically fix security issues or get remediation guidance.

**Basic Usage**:
```bash
# Preview fixes without executing (recommended first step)
agent-guard fix --dry-run

# Automatically fix all issues
agent-guard fix --auto

# Fix specific category
agent-guard fix --category filesystem --auto
```

**Example Output**:
```
🔧 Security Fix Wizard
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Found 3 security issues:

1. [Filesystem] /Users/user/.ssh is writable
   → Fix: chmod 700 /Users/user/.ssh
   → Priority: HIGH

2. [Shell] History file contains sensitive commands
   → Fix: rm ~/.bash_history
   → Priority: MEDIUM

3. [Network] Port 22 open to external connections
   → Manual review required
   → Priority: LOW

Execute fixes? [y/N]: y

✅ Fixed 2 issues
⏭️  1 issue requires manual review
```

**Fix Options**:
```bash
# Dry run - see what would be fixed
agent-guard fix --dry-run

# Auto-fix - automatically execute safe commands
agent-guard fix --auto

# Category-specific fixes
agent-guard fix --category shell --auto
agent-guard fix --category filesystem --auto
agent-guard fix --category network --auto
agent-guard fix --category secrets --auto
```

**Safety Features**:
- Dry-run mode previews changes before execution
- Destructive commands require manual confirmation
- Fix commands are logged to audit trail
- Can undo fixes from audit log

### Trend Command (New in v1.3.0)

**Risk Trend Analysis** - Compare scan results over time to track security posture.

**Basic Usage**:
```bash
# Show last 7 days trend
cd agent-guard trend

# Custom time range
agent-guard trend --days 30

# Category-specific trends
agent-guard trend --category filesystem

# JSON output for automation
agent-guard trend --json > trend-data.json
```

**Example Output**:
```
📈 Security Risk Trend Analysis
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

📊 Analyzing current security state...

📅 Analysis Period: 2024-01-15 to 2024-01-22
📊 Trend Direction: 📈 Improving ↗️

Risk Level Comparison:
  Previous: HIGH
  Current:  MEDIUM

📝 Changes Detected:
  • Risk decreased from HIGH to MEDIUM
  • filesystem: improved from HIGH to LOW
  • shell: stable at HIGH

✅ Security posture is improving! Keep up the good work.

📂 Category Breakdown:
  filesystem: LOW (score: 25)
  shell: HIGH (score: 75)
  network: MEDIUM (score: 50)
  secrets: LOW (score: 25)
```

**Trend Analysis Features**:
- Compares current scan with historical data
- Identifies improving/worsening trends
- Shows category-level changes
- Calculates risk score trends
- Visual trend indicators (📈 improving, 📉 worsening, ➡️ stable)

**Integration with Monitoring**:
```bash
# Export trend data for monitoring systems
agent-guard trend --json | jq '.trend_direction'

# Alert on worsening trends
if agent-guard trend --json | jq -e '.trend == "worsening"'; then
  echo "WARNING: Security posture deteriorating!"
  # Send alert to monitoring system
fi

# Track in time-series database
agent-guard trend --days 30 --json | \
  jq --arg date $(date +%Y-%m-%d) '.data |= {date: $date} + .' | \
  curl -X POST http://influxdb:8086/api/v2/write -d @-
```
---

## Web UI Dashboard

### Quick Start

```bash
# 1. Start Backend (from parent directory)
cd AIAgentGuard
go run webui/backend/main.go

# 2. Start Frontend (another terminal)
cd webui/frontend
npm install
npm run dev

# 3. Access Dashboard
# http://localhost:5173
```

### Features

**Real-time Monitoring**:
- Live scan results
- System status overview
- Metrics dashboard with 30s refresh

**Security Scanning**:
- One-click full security scan
- Customizable scan options
- Scan history and trends

**Multi-language Dependencies**:
- Go modules vulnerability tracking
- npm packages security audit
- Python pip package scanning
- Rust cargo dependency checks

**Visual Risk Assessment**:
- Color-coded severity levels
- Permission breakdown charts
- Detailed findings display

### API Endpoints

```bash
# System status
curl http://localhost:8080/api/v1/status

# Execute scan
curl http://localhost:8080/api/v1/scan

# Scan metrics
curl http://localhost:8080/api/v1/metrics/scan-rate

# Vulnerability data
curl http://localhost:8080/api/v1/metrics/vulnerabilities

# Prometheus metrics
curl http://localhost:8080/metrics
```

---

## Monitoring & Alerting

### Prometheus Integration

**Quick Setup**:
```bash
# Backend already exposes /metrics endpoint
# Just configure Prometheus to scrape it
```

**prometheus.yml Configuration**:
```yaml
scrape_configs:
  - job_name: 'agentguard'
    static_configs:
      - targets: ['localhost:8080']
    metrics_path: '/metrics'
    scrape_interval: 15s
```

**Available Metrics**:

| Metric | Type | Description |
|--------|------|-------------|
| `agentguard_scan_total` | counter | Total scans executed |
| `agentguard_scan_duration_seconds` | histogram | Scan duration |
| `agentguard_vulnerabilities_total` | gauge | Vulnerability count by severity |
| `agentguard_language_scan_total` | counter | Scans by language |

**Query Examples**:
```promql
# Scans per hour
rate(agentguard_scan_total[1h])

# Average scan duration
rate(agentguard_scan_duration_seconds_sum[5m]) / rate(agentguard_scan_duration_seconds_count[5m])

# Critical vulnerabilities
agentguard_vulnerabilities_total{severity="critical"}

# npm vulnerabilities
agentguard_vulnerabilities_total{language="npm"}
```

### Grafana Dashboard

**Import Dashboard**:
```bash
# Dashboard location
configs/grafana-dashboard.json

# Import via Grafana UI:
# Dashboards → Import → Upload JSON
```

**Dashboard Includes**:
- Scan rate over time
- Vulnerability trends by severity
- Language-specific breakdown
- Performance metrics (P50, P95, P99)
- Real-time alerts

### Alerting Rules

**Create Alerts**:
```yaml
# alerts.yml
groups:
  - name: agentguard
    rules:
      - alert: HighScanRate
        expr: rate(agentguard_scan_total[5m]) > 10
        for: 10m
        annotations:
          summary: "High scan rate detected"

      - alert: CriticalVulnerabilities
        expr: agentguard_vulnerabilities_total{severity="critical"} > 0
        for: 1m
        annotations:
          summary: "Critical vulnerabilities found"
```

---

## Configuration & Customization

### Policy Configuration

**Create `.agent-guard.yaml`**:
```yaml
# Blocked commands
blocked_commands:
  - "rm -rf /"
  - "dd if=/dev/zero"
  - ":(){ :|:& };:"  # fork bomb

# Allowed paths
allowed_paths:
  - /tmp
  - /home/user/project
  - /var/log/app

# Denied paths
denied_paths:
  - /etc/passwd
  - /etc/shadow
  - ~/.ssh

# Environment protection
blocked_env_vars:
  - API_KEY
  - SECRET_TOKEN
  - DATABASE_URL

# Network access
network:
  allowed_domains:
    - api.github.com
    - cdn.jsdelivr.net
  denied_domains:
    - "*.malicious.com"

# Scanner settings
scanner:
  filesystem: true
  shell: true
  network: true
  secrets: true
  filecontent: true
  dependencies: true
  npmdeps: true
  pipdeps: true
  cargodeps: true
```

### Environment Variables

```bash
# Backend configuration
export PORT=8080
export LOG_LEVEL=debug
export CORS_ORIGINS=http://localhost:3000

# Frontend configuration
export VITE_API_URL=http://localhost:8080
export VITE_METRICS_REFRESH=30000
```

### Custom Metrics Export

**Define Custom Metrics**:
```go
// In your code
import "github.com/prometheus/client_golang/prometheus"

var customMetric = prometheus.NewCounter(
    prometheus.CounterOpts{
        Name: "custom_metric_total",
        Help: "Custom metric description",
    },
)

// Register and use
prometheus.MustRegister(customMetric)
customMetric.Inc()
```

---

## Deployment Guide

For detailed deployment instructions, see:

- **English**: [webui/DEPLOYMENT.md](webui/DEPLOYMENT.md)
- **中文**: [webui/DEPLOYMENT_zh.md](webui/DEPLOYMENT_zh.md)

**Deployment Options**:
- Local development
- Docker / Docker Compose
- Production (Systemd, Nginx, Kubernetes)

---

## Daily Maintenance

### Regular Tasks

**Daily**:
- Check monitoring dashboards
- Review scan results
- Monitor alert notifications

**Weekly**:
- Update vulnerability databases
- Review security trends
- Audit scan history

**Monthly**:
- Security updates
- Performance review
- Backup verification
- Compliance reporting

### Log Management

**View Logs**:
```bash
# Recent scan history
cat ~/.agent-guard/audit.log | tail -100

# Filter by type
grep "CRITICAL" ~/.agent-guard/audit.log

# Follow in real-time
tail -f ~/.agent-guard/audit.log
```

**Export Logs**:
```bash
# Daily backup
cp ~/.agent-guard/audit.log ~/.agent-guard/audit-$(date +%Y%m%d).log

# JSON export
cat ~/.agent-guard/audit.log | jq . > audit-export.json
```

### Backup Configuration

```bash
# Backup policy configuration
cp .agent-guard.yaml .agent-guard.yaml.backup

# Backup scan history
tar -czf agentguard-backup-$(date +%Y%m%d).tar.gz ~/.agent-guard/

# Automated backup script
#!/bin/bash
DATE=$(date +%Y%m%d)
tar -czf /backup/agentguard-$DATE.tar.gz \
  ~/.agent-guard/ \
  .agent-guard.yaml \
  prometheus.yml
```

---

## Troubleshooting

### Common Issues

#### 1. Port Already in Use

**Problem**:
```
Error: listen tcp :8080: bind: address already in use
```

**Solution**:
```bash
# Find process using port
lsof -ti:8080

# Kill the process
kill -9 $(lsof -ti:8080)

# Or use different port
PORT=9090 agent-guard scan
```

#### 2. Permission Denied

**Problem**:
```
Error: permission denied while accessing /etc/passwd
```

**Solution**:
```bash
# Run with sudo (not recommended)
sudo agent-guard scan

# Or add allowed paths to policy
# .agent-guard.yaml:
allowed_paths:
  - /etc/passwd:read
```

#### 3. Missing Dependencies

**Problem**:
```
Error: go command not found
```

**Solution**:
```bash
# Install Go
# macOS
brew install go

# Ubuntu
sudo apt install golang-go

# Verify
go version
```

#### 4. Docker Network Issues

**Problem**:
```
Error: container cannot connect to backend
```

**Solution**:
```bash
# Check network
docker network inspect agentguard-network

# Recreate network
docker compose down
docker network prune
docker compose up -d
```

#### 5. Metrics Not Appearing

**Problem**: Prometheus not scraping metrics

**Solution**:
```bash
# Verify metrics endpoint
curl http://localhost:8080/metrics | head

# Check Prometheus targets
# http://localhost:9090/targets

# Reload Prometheus
docker compose restart prometheus
```

### Debug Mode

```bash
# Enable verbose output
agent-guard scan --verbose

# Debug mode
LOG_LEVEL=debug go run webui/backend/main.go

# Trace API calls
curl -v http://localhost:8080/api/v1/status
```

### Getting Help

**Resources**:
- **Documentation**: [doc/](doc/)
- **GitHub Issues**: https://github.com/imdlan/AIAgentGuard/issues
- **Discussions**: https://github.com/imdlan/AIAgentGuard/discussions

**Report Issues**:
```bash
# Generate diagnostic info
agent-guard scan --json > diagnostic-report.json

# Include in issue:
# - OS and version
# - AIAgentGuard version
# - Error messages
# - Diagnostic report
```

---

## Best Practices

### Security

1. **Regular Scans**: Schedule automated daily/weekly scans
2. **Dependency Updates**: Keep dependencies up-to-date
3. **Policy Enforcement**: Use strict policy configurations
4. **Access Control**: Restrict sensitive paths and commands
5. **Monitoring**: Set up alerts for critical vulnerabilities

### Performance

1. **Selective Scanning**: Scan only necessary categories
2. **Caching**: Use scan result caching when appropriate
3. **Parallel Scans**: Run multiple language dependency scans in parallel
4. **Monitoring**: Track scan duration and optimize slow operations

### Operations

1. **Backup Configuration**: Version control policy files
2. **Audit Trail**: Maintain scan history for compliance
3. **Incident Response**: Have procedures for critical vulnerabilities
4. **Documentation**: Keep deployment and configuration docs updated

---

## Appendix

- **v1.3.0**: Detailed security reporting, fix wizard, trend analysis, process/network scanning details
- **v1.2.0**: Multi-language dependency scanning, Web UI monitoring, Prometheus integration
- **v1.1.0**: Go dependency scanning, container runtime detection
- **v1.0.0**: Initial release with core scanning features

- **v1.2.0**: Multi-language dependency scanning, Web UI monitoring, Prometheus integration
- **v1.1.0**: Go dependency scanning, container runtime detection
- **v1.0.0**: Initial release with core scanning features

### License

MIT License - see [LICENSE](LICENSE) file for details

### Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md)

---

**Last Updated**: 2026-02-28
**Version**: v1.3.0
**Maintainer**: AIAgentGuard Team
**Version**: v1.2.0
**Maintainer**: AIAgentGuard Team

**中文版本**: [USAGE_zh.md](USAGE_zh.md)
