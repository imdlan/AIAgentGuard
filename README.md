# AI AgentGuard

[![Version](https://img.shields.io/badge/version-v1.1.0-blue.svg)](https://github.com/imdlan/AIAgentGuard/releases/latest)
[![Go Report](https://goreportcard.com/badge/github.com/imdlan/AIAgentGuard)](https://goreportcard.com/report/github.com/imdlan/AIAgentGuard)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)



[English](README.md) | [简体中文](README_zh.md)

## Features

### Core Security Scanning
- **Permission Scanning** - Detect filesystem, shell, network, and secret access permissions
- **File Content Analysis** - Scan files for exposed API keys, tokens, and secrets (15+ patterns)
- **Process Security Monitoring** - Detect reverse shells, suspicious processes, and high CPU usage
- **SUID/SGID Scanning** - Identify privileged executables and potential privilege escalation vectors

### Advanced Features (New in v1.1.0) ⭐
- **Dependency Vulnerability Scanning** - Check Go dependencies for known CVEs using golang.org/x/vuln
- **Container Runtime Detection** - Detect Docker, Kubernetes, Podman, LXC, Wasm environments
- **True Sandbox Isolation** - containerd-based container isolation with Linux namespaces (Linux only)

### Security & Compliance
- **Audit Logging** - Comprehensive security event logging with JSON format and SIEM integration
- **Risk Assessment** - Intelligently analyze security threats and calculate risk levels (85%+ coverage)
- **Smart Command Parsing** - Advanced flag parsing to prevent bypass attempts
- **Sandbox Execution** - Safely run commands in isolated environments

### Configuration & Protection
- **Policy Management** - Control access permissions via YAML configuration
- **Prompt Injection Protection** - Detect and block malicious prompt injection attacks
- **Plugin Scanning** - Detect insecure plugins and extensions

- **Permission Scanning** - Detect filesystem, shell, network, and secret access permissions
- **File Content Analysis** - Scan files for exposed API keys, tokens, and secrets (15+ patterns)
- **Process Security Monitoring** - Detect reverse shells, suspicious processes, and high CPU usage
- **SUID/SGID Scanning** - Identify privileged executables and potential privilege escalation vectors
- **Audit Logging** - Comprehensive security event logging with JSON format and SIEM integration
- **Risk Assessment** - Intelligently analyze security threats and calculate risk levels
- **Smart Command Parsing** - Advanced flag parsing to prevent bypass attempts
- **Sandbox Execution** - Safely run commands in isolated environments
- **Policy Management** - Control access permissions via YAML configuration
- **Prompt Injection Protection** - Detect and block malicious prompt injection attacks
- **Plugin Scanning** - Detect insecure plugins and extensions
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

### 4. Initialize Configuration

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

                             🛡️  Security Scan Report v1.0

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

**Protect your AI Agents, start with security scanning!** 🛡️

### 6. Audit Log Examples

View security events:
```bash
# Show recent policy violations
cat ~/.agent-guard/audit.log | jq '. | select(.event_type == "policy_violation")'

# Follow audit logs in real-time
tail -f ~/.agent-guard/audit.log | jq .

# Find blocked commands
grep "command_blocked" ~/.agent-guard/audit.log | jq .

# Check for high-risk events
grep "CRITICAL" ~/.agent-guard/audit.log | jq .
```

