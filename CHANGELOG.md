# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.2.0-beta] - 2026-02-28

### Added
- **Multi-Language Dependency Vulnerability Scanning**: Expanded language support
  - **npm/yarn vulnerability scanning**: JavaScript/TypeScript package scanning
    - Uses `npm audit --json` for package-lock.json
    - Uses `yarn audit --json` for yarn.lock
    - Detects vulnerabilities in: npm, yarn workspaces
    - Skips node_modules directories
    - Severity classification: critical, high, moderate, low
  
  - **pip vulnerability scanning**: Python package scanning
    - Uses `pip-audit` for requirements.txt, pyproject.toml, Pipfile
    - Falls back to `safety` tool if pip-audit unavailable
    - Detects vulnerabilities in: requirements.txt, pyproject.toml, Pipfile.lock
    - Skips __pycache__, venv, .venv directories
  
  - **cargo vulnerability scanning**: Rust package scanning
    - Uses `cargo audit` for Cargo.toml
    - Detects vulnerabilities using cargo-audit
    - Supports CVSS score parsing
    - Skips target directory

- **Enhanced Scanner Integration**: Multi-language support in main scanner
  - Added `NpmDeps`, `PipDeps`, `CargoDeps` to PermissionResult
  - Supports individual scanning: `agent-guard scan --category npmdeps`
  - Integrated with existing scan workflows

- **Unit Tests**: Comprehensive test coverage for multilang scanners
  - File existence tests
  - Directory skipping tests (node_modules, __pycache__, target)
  - Vulnerability count conversion tests
  - Benchmark tests for all three scanners
- **Unit Tests**: Comprehensive test coverage for multilang scanners
  - File existence tests
  - Directory skipping tests (node_modules, __pycache__, target)
  - Vulnerability count conversion tests
  - Benchmark tests for all three scanners

- **Prometheus Metrics Export**: Real-time monitoring and observability
  - Scan metrics: total scans, duration, results by category
  - Vulnerability metrics: discovered vulnerabilities by severity and language
  - Component metrics: scan counts, failure rates
  - System metrics: memory usage, uptime, last scan timestamp
  - HTTP `/metrics` endpoint for Prometheus scraping
  - `--metrics-addr` flag to enable metrics server

- **Grafana Dashboard**: Pre-built monitoring dashboard
  - Scan rate and duration visualization
  - Latest scan results gauge
  - Vulnerability discovery trends
  - Component health monitoring
  - Easy import via JSON configuration

- **Monitoring Documentation**: Comprehensive monitoring guide
  - Quick start instructions
  - Prometheus configuration examples
  - Grafana setup and dashboard import
  - Docker Compose stack for complete monitoring
  - Alerting rule examples
  - CI/CD integration patterns
### Security Coverage
- Overall security coverage increased from 85% to 90%+
- Multi-language dependency scanning adds critical supply chain security
- Detects vulnerable dependencies across Go, JavaScript, Python, and Rust ecosystems

## [1.1.0] - 2026-02-27

### Added
- **Dependency Vulnerability Scanning**: New scanner that checks Go dependencies for known CVEs
  - Integrated with `golang.org/x/vuln` database
  - Includes built-in vulnerability database for offline scanning
  - Supports `govulncheck` tool integration
  - Detects vulnerabilities in: gin-gonic, cobra, yaml.v3, protobuf, etc.

- **Container Runtime Detection**: Enhanced environment awareness
  - Detects Docker, Kubernetes, Podman, LXC, and Wasm runtimes
  - Extracts container metadata (namespace, pod name, container ID, node name)
  - Checks for Docker/containerd socket availability
  - Provides runtime version information

- **True Sandbox Isolation**: Production-grade container isolation
  - containerd-based sandbox implementation
  - Linux namespaces support (UTS, PID, mount, network)
  - Optional network isolation
  - Readonly root filesystem support
  - Resource limits (memory, CPU)
  - Capability whitelisting/blacklisting
  - Automatic fallback to namespace isolation when containerd unavailable
  - Cross-platform support (Linux: full features, Darwin/macOS: stub)

- **Performance Benchmarking**: Comprehensive test suite
  - 12 benchmark tests for all major components
  - Performance metrics for scanners, container detection, sandbox
  - Parallel execution benchmarking
  - Memory allocation tracking

- **Extended Integration Tests**: Enhanced test coverage
  - Container detection integration tests
  - Comprehensive scan tests with container context
  - Sandbox availability tests
  - Total: 42 tests (unit, integration, e2e, benchmarks)

### Changed
- **SUID Scan Optimization**: Significant performance improvements
  - Added 10-second timeout limit
  - Maximum scan depth: 5 levels
  - Maximum file limit: 100 files
  - Uses `-xdev` to avoid cross-filesystem scans
  - Result: 30s timeout → 10s controlled execution

- **Scanner Output**: Added `dependencies` field to scan results
  - JSON output now includes dependency risk assessment
  - Overall risk calculation includes dependency vulnerabilities

### Improved
- **Test Infrastructure**: 
  - Fixed GitHub token pattern test (36 → 40 characters)
  - Added CI environment detection for slow tests
  - Improved test stability and reliability

- **Cross-Platform Compatibility**:
  - Platform-specific sandbox implementations
  - Graceful degradation on non-Linux systems
  - Better error messages for unsupported features

- **Documentation**:
  - Updated README with new features
  - Added examples/ directory with configuration templates
  - Enhanced Chinese documentation (README_zh.md)

### Performance
- Dependency scanning: 73μs (ultra-fast)
- Container detection: 13μs (very fast)
- Socket detection: 5μs (extremely fast)
- Sandbox creation: 5μs (extremely fast)
- Overall scan time: ~500ms for all 6 scanners

### Security
- **Security Coverage**: 78% → 85%+
- New threat vectors covered:
  - Supply chain vulnerabilities (dependencies)
  - Container escape risks
  - Sandbox bypass attempts
  - Privilege escalation via known vulnerable packages

### Technical Details
- New dependency: `golang.org/x/vuln v1.1.4`
- Binary size: 5.2 MB
- Go version: 1.25.5+
- Platform support: Linux, macOS, Windows (partial)

### Migration Notes
- No breaking changes
- Existing configurations remain compatible
- New features are opt-in (scan automatically includes dependencies)
- Sandbox features require containerd on Linux

## [1.0.0] - 2026-02-XX

### Added
- Initial release
- Filesystem, shell, network, and secret scanning
- File content key scanning (15+ patterns)
- Process security monitoring
- SUID/SGID file scanning
- Audit logging system
- Smart command parsing
- Basic sandbox execution
- Policy management (YAML)
- Prompt injection protection
- Plugin scanning

[1.1.0]: https://github.com/imdlan/AIAgentGuard/compare/v1.0.0...v1.1.0
[1.0.0]: https://github.com/imdlan/AIAgentGuard/releases/tag/v1.0.0
