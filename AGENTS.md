# AI AgentGuard - Agent Development Guide

This document contains essential information for agentic coding agents working on this repository.

## Build & Test Commands

### Essential Commands
```bash
# Build the CLI binary
go build -o agent-guard

# Run the CLI directly
go run main.go [command]

# Run with arguments
./agent-guard scan
./agent-guard run "echo test"
./agent-guard report --json

# Initialize configuration
./agent-guard init

# Dependency management
go mod tidy
go mod download
```

### Testing (Future - No Tests Yet)
```bash
# Run all tests (when added)
go test ./...

# Run tests in specific package
go test ./internal/scanner

# Run single test (example for when tests are added)
go test -v ./internal/risk -run TestAnalyze

# Run with coverage
go test -cover ./...
```

### Linting (Not Yet Configured)
```bash
# Run go fmt
go fmt ./...

# Run go vet
go vet ./...

# Run golangci-lint (if added)
golangci-lint run
```

## Code Style Guidelines

### Project Structure
```
agent-guard/
├── cmd/              # CLI commands (cobra-based)
├── internal/         # Private application code
│   ├── scanner/     # Security scanning engines
│   ├── risk/        # Risk analysis & scoring
│   ├── sandbox/     # Sandboxed execution
│   ├── policy/      # Policy configuration (YAML)
│   ├── security/    # Prompt injection guard
│   └── report/      # Report generation
├── pkg/model/       # Public data types & models
├── configs/         # Default configuration files
└── main.go          # Application entry point
```

### Package Conventions
- **Package names**: Lowercase, single word (`scanner`, `risk`, `policy`, `model`)
- **cmd package**: For CLI commands using cobra
- **internal packages**: Private application code, never imported by external projects
- **pkg packages**: Public API that could be imported by external projects

### Naming Conventions
- **Exported functions/types**: PascalCase (`RunAllScans`, `Analyze`, `RiskLevel`)
- **Private functions**: camelCase (`calculateOverallRisk`, `generateRiskDetails`)
- **Constants**: PascalCase (`Low`, `Medium`, `High`, `Critical`)
- **Global variables**: camelCase (`configFile`, `verbose`, `jsonOutput`)
- **Struct fields**: PascalCase for exported, camelCase for private

### Import Ordering
Group imports in three sections with blank lines between:
1. Standard library
2. Internal packages (`github.com/imdlan/AIAgentGuard/...`)
3. External packages

```go
import (
    "fmt"
    "os"

    "github.com/imdlan/AIAgentGuard/internal/scanner"
    "github.com/imdlan/AIAgentGuard/pkg/model"

    "github.com/spf13/cobra"
)
```

### Type & Struct Conventions
- Use descriptive names with full words (`PermissionResult`, not `PermResult`)
- Struct fields use PascalCase
- JSON/YAML tags use snake_case
- Use `omitempty` for optional fields

```go
type ScanReport struct {
    ToolName string           `json:"tool_name"`
    Results  PermissionResult `json:"results"`
    Overall  RiskLevel        `json:"overall"`
    Details  []RiskDetail     `json:"details,omitempty"`
}
```

### Error Handling
- Functions return error as last return value
- Use `fmt.Errorf` with `%w` for wrapping errors
- Early returns for error conditions
- Check errors immediately

```go
func LoadConfig(path string) (*model.PolicyConfig, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }
    // ... continue
}
```

### Comments
- **Exported functions**: Must have documentation comments
- **Package comments**: Each package should have a comment describing its purpose
- **Inline comments**: Use for non-obvious logic, explain "why" not "what"
- Comment style: Single sentence, no period for short descriptions

```go
// RunAllScans executes all security scans and returns the combined result
func RunAllScans() model.PermissionResult {
    return model.PermissionResult{
        Filesystem: ScanFilesystem(),
        Shell:      ScanShell(),
        Network:    ScanNetwork(),
        Secrets:    ScanSecrets(),
    }
}
```

### Function Design
- Keep functions focused on single responsibility
- Prefer multiple small functions over large ones
- Use private helper functions (camelCase) for implementation details
- Export functions (PascalCase) form the package API

### Constants & Enums
- Use string constants for enums (PascalCase)
- Define related constants together

```go
type RiskLevel string

const (
    Low      RiskLevel = "LOW"
    Medium   RiskLevel = "MEDIUM"
    High     RiskLevel = "HIGH"
    Critical RiskLevel = "CRITICAL"
)
```

### CLI Commands (Cobra)
- Commands are defined in `cmd/` package
- Use `RunE` for error-returning command handlers
- Persistent flags for global options
- Local flags for command-specific options
- Include Usage examples in command definition

```go
var scanCmd = &cobra.Command{
    Use:   "scan [target]",
    Short: "Scan for security risks",
    Long:  `Extended description...`,
    Example: `  agent-guard scan
  agent-guard scan --json`,
    RunE: runScan,
}
```

### Configuration (YAML)
- Configuration files use YAML format
- Use struct tags for mapping
- Default configuration stored in `configs/default.yaml`
- Support multiple config locations (current dir, home dir, /etc)

## Security Considerations
- This is a security tool - code reviews are critical
- Never suppress errors with empty catch blocks
- Validate all user inputs
- Be careful with command execution - always validate/sanitize
- Default to deny, not allow

## Testing Strategy (To Be Implemented)
- Unit tests for core logic (risk analysis, scanning)
- Integration tests for CLI commands
- Table-driven tests for multiple scenarios
- Mock external dependencies (filesystem, network)

## Dependencies
- `github.com/spf13/cobra` - CLI framework
- `gopkg.in/yaml.v3` - YAML configuration parsing
- Go 1.25.5+

## Adding New Features
1. Add types to `pkg/model/types.go` if needed
2. Implement logic in appropriate `internal/` package
3. Add CLI command in `cmd/` if user-facing
4. Update documentation
5. Add tests (when test framework is added)
