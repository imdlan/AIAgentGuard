# AI AgentGuard

[![Version](https://img.shields.io/badge/version-v1.4.1-blue.svg)](https://github.com/imdlan/AIAgentGuard/releases/latest)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Go Report](https://goreportcard.com/badge/github.com/imdlan/AIAgentGuard)](https://goreportcard.com/report/github.com/imdlan/AIAgentGuard)
[![Downloads](https://img.shields.io/github/downloads/imdlan/AIAgentGuard/total.svg)](https://github.com/imdlan/AIAgentGuard/releases)
[![Homebrew](https://img.shields.io/badge/Homebrew-imdlan%2FAIAgentGuard-orange.svg)](https://github.com/imdlan/homebrew-AIAgentGuard)

[![Version](https://img.shields.io/badge/version-v1.4.1-blue.svg)]

## Screenshots
<div align="center">
  <img src="AIAgentGuard-Screenshot-2.png" alt="AgentGuard CLI" />
  <img src="AIAgentGuard-Screenshot-1.png" alt="Security Scan Report" />
  <p><strong>AI AgentGuard is a security tool for AI agents, CLI tools, and MCP servers.It scans for permission risks, evaluates security threats, and providessandboxed execution environments.</strong></p>
  <p><a href="README.md">English</a> | <a href="README_zh.md">简体中文</a></p>
</div>

## Features

### v1.4.1 Features (Latest) ⭐
- **Version Command** - Display installed version with `agent-guard version`
- **Installation Script Improvements** - Fallback method for GitHub API rate limiting with better error messages

### v1.4.0 Features
- **Multi-Language Support (i18n)** - Auto-detect system language, support English and Chinese output
- **Embedded Translation Files** - 40+ translatable strings with no external dependencies
- **macOS Language Detection** - Automatic detection from AppleLocale and AppleLanguages preferences

### v1.3.0 Features
- **Detailed Security Reporting** - Show specific files, processes, and commands causing security risks
  - Process Scanning Details with PID, command line, and risk reasons
  - Network Connection Analysis displaying open ports and active connections
  - Enhanced risk assessment with actionable remediation steps
- **Automated Fix Wizard** - Auto-fix security issues or provide manual remediation commands
  - New `agent-guard fix` CLI command with `--auto` and `--dry-run` options
  - Category-specific fixes (filesystem, shell, network, secrets)
- **Risk Trend Analysis** - Compare scan results over time to track security posture changes
  - New `agent-guard trend` CLI command with historical data analysis
- **Web UI Dashboard** - Complete visual security monitoring interface
  - Real-time dashboard with React + Go RESTful API
  - Process, Network, Fix Wizard, and Trend History panels

### v1.2.0 Features
- **Multi-Language Dependency Scanning** - Scan npm, pip, cargo, and Go dependencies for vulnerabilities
- **Prometheus Monitoring** - Export metrics for monitoring and alerting with `/metrics` endpoint
- **Grafana Dashboard** - Pre-built monitoring dashboard with real-time visualizations
- **Enhanced Test Coverage** - Comprehensive unit tests for multilang scanners (npm, pip, cargo)

### v1.1.0 Features
- **Go Dependency Vulnerability Scanning** - Check Go dependencies for known CVEs using golang.org/x/vuln
- **Container Runtime Detection** - Detect Docker, Kubernetes, Podman, LXC, Wasm environments
- **True Sandbox Isolation** - containerd-based container isolation with Linux namespaces (Linux only)
- **Performance Benchmarking** - 12 benchmark tests for all major components

### v1.0.0 Features (Core)
- **Permission Scanning** - Detect filesystem, shell, network, and secret access permissions
- **File Content Analysis** - Scan files for exposed API keys, tokens, and secrets (15+ patterns)
- **Process Security Monitoring** - Detect reverse shells, suspicious processes, and high CPU usage
- **SUID/SGID Scanning** - Identify privileged executables and potential privilege escalation vectors
- **Audit Logging** - Comprehensive security event logging with JSON format and SIEM integration
- **Smart Command Parsing** - Advanced flag parsing to prevent bypass attempts
- **Sandbox Execution** - Safely run commands in isolated environments
- **Policy Management** - Control access permissions via YAML configuration
- **Prompt Injection Protection** - Detect and block malicious prompt injection attacks
- **Plugin Scanning** - Detect insecure plugins and extensions

## Update

To update to the latest version:

**Homebrew**:
```bash
brew upgrade agent-guard
```

**Install Script**:
```bash
curl -sSL https://raw.githubusercontent.com/imdlan/AIAgentGuard/main/scripts/install.sh | bash
```

**Manual**: Download from [Releases](https://github.com/imdlan/AIAgentGuard/releases/latest).
## Installation


### Method 1: Homebrew (Recommended for macOS/Linux)

```bash
brew tap imdlan/AIAgentGuard
brew install agent-guard
```

### Method 2: Download from GitHub Releases

Visit the [Releases page](https://github.com/imdlan/AIAgentGuard/releases) to download binaries for your platform.

```bash
# macOS / Linux
curl -LO https://github.com/imdlan/AIAgentGuard/releases/latest/download/agent-guard_darwin_arm64.tar.gz
tar -xzf agent-guard_darwin_arm64.tar.gz
chmod +x agent-guard
sudo mv agent-guard /usr/local/bin/
```

### Method 3: Go Install (For Developers)

```bash
go install github.com/imdlan/AIAgentGuard@latest
```

Make sure `$GOPATH/bin` is in your `PATH`:
```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

### Method 4: Install Script

```bash
curl -sSL https://raw.githubusercontent.com/imdlan/AIAgentGuard/main/scripts/install.sh | bash
```

### Method 5: Build from Source

```bash
git clone https://github.com/imdlan/AIAgentGuard.git
cd agent-guard
go build -o agent-guard
sudo mv agent-guard /usr/local/bin/
```

## Quick Start

### 1. Scan Security Risks

```bash
# Scan current environment
agent-guard scan

# JSON format output
agent-guard scan --json

# Use custom policy
agent-guard scan --config ./my-policy.yaml
```

### 2. Run in Sandbox

```bash
# Run command in isolated environment
agent-guard run "curl https://api.example.com"

# Disable network access
agent-guard run --disable-network "npm install"

# Restrict filesystem access
agent-guard run --allow-dirs /tmp,/data "node script.js"
```

### 3. Generate Report

```bash
# Generate detailed report
agent-guard report

# Save to file
agent-guard report --json > security-report.json
```

### 4. Monitor with Prometheus (New)

```bash
# Run scan with Prometheus metrics
agent-guard scan --metrics-addr :9090

# Metrics available at http://localhost:9090/metrics
# curl http://localhost:9090/metrics
```

For detailed monitoring setup, see [Monitoring Guide](doc/MONITORING.md).

### 5. Initialize Configuration

```bash
# Generate default configuration file
agent-guard init

# Configuration file locations:
# - .agent-guard.yaml (current directory)
# - ~/.agent-guard.yaml (user directory)
# - /etc/agent-guard/config.yaml (system directory)
```

## Configuration Example

Create `.agent-guard.yaml`:

```yaml
# Block dangerous commands
blocked_commands:
  - "rm -rf /"
  - "dd if=/dev/zero"
  - "mkfs"
  - ":(){ :|:& };:"  # fork bomb

# Restrict filesystem access
allowed_paths:
  - /tmp
  - /home/user/project
  - /var/log/app

denied_paths:
  - /etc/passwd
  - /etc/shadow
  - ~/.ssh

# Environment variable protection
blocked_env_vars:
  - API_KEY
  - SECRET_TOKEN
  - DATABASE_URL

# Network access control
network:
  allowed_domains:
    - api.github.com
    - cdn.jsdelivr.net
  denied_domains:
    - "*.malicious.com"
```

## Output Example

```
  █████╗ ██╗     █████╗  ██████╗ ███████╗███╗   ██╗████████╗ ██████╗ ██╗   ██╗ █████╗ ██████╗ ██████╗
 ██╔══██╗██║    ██╔══██╗██╔════╝ ██╔════╝████╗  ██║╚══██╔══╝██╔════╝ ██║   ██║██╔══██╗██╔══██╗██╔══██╗
 ███████║██║    ███████║██║  ███╗█████╗  ██╔██╗ ██║   ██║   ██║  ███╗██║   ██║███████║██████╔╝██║  ██║
 ██╔══██╗██║    ██╔══██╗██║   ██║██╔══╝  ██║╚██╗██║   ██║   ██║   ██║██║   ██║██╔══██╗██╔══██╗██║  ██║
 ██║  ██║██║    ██║  ██║╚██████╔╝███████╗██║ ╚████║   ██║   ╚██████╔╝╚██████╔╝██║  ██║██║  ██║██████╔╝
 ╚═╝  ╚═╝╚═╝    ╚═╝  ╚═╝ ╚═════╝ ╚══════╝╚═╝  ╚═══╝   ╚═╝    ╚═════╝  ╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═╝╚═════╝

                             🛡️  Security Scan Report v1.4.1

Overall Risk: 🔶 HIGH

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Permission Breakdown:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  ✅ Filesystem Access: LOW
  🛑 Shell Execution: CRITICAL
  ⚠️ Network Access: MEDIUM
  🔶 Secrets Access: HIGH

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Detailed Findings:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
1. [SHELL] Root/admin shell access detected [/bin/bash, /bin/zsh] [SYSTEM]
2. [NETWORK] External network connectivity available [NETWORK]
3. [SECRETS] Environment variable API_KEY exposed [ENVIRONMENT]

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Recommendations:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  • Consider running AI agents in a sandboxed environment
  • Use 'agent-guard run <command>' for safe execution
  • Use environment variable blocking in policy config
  • Consider using secret management tools

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

## Command Reference

### Global Options

```
-c, --config string   Path to policy configuration file
-j, --json            JSON output format
-v, --verbose         Verbose output
-h, --help            Show help information
```

### scan - Security Scan

Scan the current environment for security risks and permissions.

```bash
agent-guard scan [flags]
```

### run - Sandbox Execution

Execute commands in an isolated environment.

```bash
agent-guard run [command] [flags]

Options:
  --disable-network    Disable network access
  --allow-dirs paths   Allow access to directories (comma-separated)
  --block-dirs paths   Block access to directories (comma-separated)
```

### report - Generate Report

Generate and display security reports.

```bash
agent-guard report [flags]
```

### init - Initialize Configuration

Generate default configuration file.

```bash
agent-guard init [flags]

Options:
  --force    Overwrite existing configuration file
  --path     Specify configuration file path
```

### fix - Security Fix Wizard

Automatically fix security issues or provide remediation guidance.

```bash
agent-guard fix [flags]

Options:
  --auto       Automatically execute fix commands
  --dry-run    Preview changes without executing
  --category   Fix specific category only (filesystem, shell, network, secrets)
```

### trend - Risk Trend Analysis

Analyze security trends by comparing scan results over time.

```bash
agent-guard trend [flags]

Options:
  --days N    Analyze last N days (default: 7)
  --json      Output in JSON format
  --category  Show trend for specific category
```
## Complete Usage Guide

For detailed usage guide and best practices, see: [USAGE.md](USAGE.md)

Topics covered:
- Detailed explanations of all use cases
- Complete CLI command reference
- Web UI usage instructions
- Monitoring and alerting setup
- Deployment and maintenance guides
- Troubleshooting solutions

## FAQ

### Q: How to disable specific scans?

A: Edit the configuration file and set corresponding options to `false`:

```yaml
scanner:
  filesystem: false
  shell: true
  network: true
  secrets: true
```

### Q: How does sandbox mode work?

A: Sandbox mode uses the following techniques:
- Environment variable isolation
- Filesystem access restriction
- Network access control (optional)
- Command whitelist/blacklist

### Q: How to integrate with CI/CD?

A: Add security scan steps to your CI/CD pipeline:

```yaml
# GitHub Actions example
- name: Security Scan
  run: |
    go install github.com/imdlan/AIAgentGuard@latest
    agent-guard scan --json > security-report.json
    # Check risk level
    if grep -q "CRITICAL" security-report.json; then
      echo "Critical security issues found!"
      exit 1
    fi
```

## Development

### Build from Source

```bash
# Clone repository
git clone https://github.com/imdlan/AIAgentGuard.git
cd agent-guard

# Build
go build -o agent-guard

# Run tests
go test ./...

# Install locally
go install
```

### Project Structure

```
agent-guard/
├── cmd/              # CLI commands
├── internal/         # Internal implementation
│   ├── scanner/     # Scanning engines
│   ├── risk/        # Risk analysis
│   ├── sandbox/     # Sandbox execution
│   ├── policy/      # Policy management
│   ├── security/    # Security protection
│   └── report/      # Report generation
├── pkg/model/       # Data models
├── configs/         # Default configuration
└── scripts/         # Installation scripts
```

### Release Process

This project uses Goreleaser for automated releases. When you push a version tag, GitHub Actions is automatically triggered:

1. Build multi-platform binaries (macOS/Linux, AMD64/ARM64)
2. Create GitHub Release
3. Generate file checksums (checksums.txt)
4. Automatically update Homebrew formula

**Release new version**:
```bash
git tag v1.0.1
git push origin v1.0.1
```

For detailed documentation, see: [Release Process Guide](doc/RELEASE.md)

### Local Testing

```bash
# Install goreleaser
brew install goreleaser

# Test build (no release)
goreleaser build --clean --snapshot

# Test full workflow (dry-run)
goreleaser release --clean --snapshot --skip-publish
```

## License

MIT License - see [LICENSE](LICENSE) file for details

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md)

## Contact

- GitHub: https://github.com/imdlan/AIAgentGuard
- Issues: https://github.com/imdlan/AIAgentGuard/issues
- Discussions: https://github.com/imdlan/AIAgentGuard/discussions

---

**Protect your AI Agents, start with security scanning!** 
