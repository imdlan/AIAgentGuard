# Contributing to AI AgentGuard

Thank you for your interest in contributing to AI AgentGuard! We welcome contributions from the community.

[English](CONTRIBUTING.md) | [简体中文](CONTRIBUTING_zh.md)

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [How to Contribute](#how-to-contribute)
- [Development Setup](#development-setup)
- [Coding Standards](#coding-standards)
- [Commit Message Guidelines](#commit-message-guidelines)
- [Pull Request Process](#pull-request-process)
- [Reporting Bugs](#reporting-bugs)
- [Suggesting Enhancements](#suggesting-enhancements)

## Code of Conduct

Please be respectful and constructive in all interactions. We aim to maintain a welcoming and inclusive community.

## How to Contribute

There are many ways to contribute:

- 🐛 **Report bugs** - Help us identify and fix issues
- 💡 **Suggest features** - Share your ideas for improvements
- 📝 **Improve documentation** - Help make docs clearer and more comprehensive
- 🔧 **Submit pull requests** - Fix bugs, add features, or improve code
- 💬 **Participate in discussions** - Share your knowledge and help others
- 🎨 **Design improvements** - Contribute UI/UX ideas

## Development Setup

### Prerequisites

- Go 1.25 or higher
- Git
- Make (optional, for build automation)

### Fork and Clone

```bash
# Fork the repository on GitHub
# Clone your fork
git clone https://github.com/YOUR_USERNAME/AIAgentGuard.git
cd AIAgentGuard

# Add upstream remote
git remote add upstream https://github.com/imdlan/AIAgentGuard.git
```

### Build and Test

```bash
# Build the binary
go build -o agent-guard

# Run tests
go test ./...

# Install locally
go install
```

### Development Workflow

```bash
# Create a new branch
git checkout -b feature/your-feature-name

# Make your changes
# ...

# Run tests
go test ./...

# Commit your changes
git commit -m "feat: add your feature"

# Push to your fork
git push origin feature/your-feature-name
```

## Coding Standards

### Go Code Style

Follow standard Go conventions:
- Use `gofmt` to format code
- Run `go vet` to check for issues
- Write clear, descriptive comments for exported functions
- Keep functions focused and concise

### Example

```go
// ScanFilesystem scans the filesystem for security risks
// and returns a ScanResult with detailed findings.
func ScanFilesystem() ScanResult {
    // Implementation
}
```

### Project Structure

```
agent-guard/
├── cmd/              # CLI commands
├── internal/         # Internal packages
│   ├── scanner/     # Security scanning engines
│   ├── risk/        # Risk analysis
│   ├── sandbox/     # Sandbox execution
│   ├── policy/      # Policy management
│   ├── security/    # Security protection
│   └── report/      # Report generation
├── pkg/model/       # Public data models
├── configs/         # Default configurations
└── scripts/         # Installation scripts
```

### Naming Conventions

- **Packages**: Lowercase, single word (`scanner`, `risk`, `policy`)
- **Exported functions**: PascalCase (`RunAllScans`, `Analyze`, `CalculateRisk`)
- **Private functions**: camelCase (`calculateOverallRisk`, `generateReport`)
- **Constants**: PascalCase (`Low`, `Medium`, `High`, `Critical`)

## Commit Message Guidelines

Follow [Conventional Commits](https://www.conventionalcommits.org/) specification:

### Format

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

### Types

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `perf`: Performance improvements
- `test`: Adding or updating tests
- `chore`: Maintenance tasks
- `ci`: CI/CD changes

### Examples

```bash
feat(scanner): add network port scanning
fix(risk): correct risk score calculation
docs(readme): update installation instructions
refactor(policy): simplify policy loading logic
test(sandbox): add tests for sandbox isolation
```

## Pull Request Process

### Before Submitting

1. **Search existing PRs** - Avoid duplicate work
2. **Update documentation** - Include relevant docs
3. **Add tests** - Ensure new code is tested
4. **Run linters** - Check code quality
5. **Update CHANGELOG** - Document changes

### Submitting a PR

1. Create a descriptive title
2. Describe your changes in the body
3. Link related issues
4. Ensure all checks pass
5. Request review from maintainers

### PR Template

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] Unit tests added/updated
- [ ] Manual testing completed
- [ ] All tests pass

## Checklist
- [ ] Code follows project style guidelines
- [ ] Documentation updated
- [ ] Commit messages follow conventions
- [ ] No merge conflicts
```

### Review Process

- Maintainers will review your PR
- Address review comments
- Keep discussions focused and constructive
- Be patient - review may take time

## Reporting Bugs

### Before Reporting

- Check existing issues
- Search for similar problems
- Verify it's not already fixed

### Bug Report Template

```markdown
## Description
Clear description of the bug

## Steps to Reproduce
1. Step one
2. Step two
3. Step three

## Expected Behavior
What should happen

## Actual Behavior
What actually happens

## Environment
- OS: [e.g. macOS 14.0]
- Go version: [e.g. 1.25.5]
- AgentGuard version: [e.g. v1.0.1]

## Logs
Relevant error messages or logs
```

## Suggesting Enhancements

### Feature Request Template

```markdown
## Description
Feature description

## Problem Statement
What problem does this solve?

## Proposed Solution
How should it work?

## Alternatives
Other approaches considered

## Additional Context
Screenshots, examples, etc.
```

## Getting Help

- 💬 **Discussions**: https://github.com/imdlan/AIAgentGuard/discussions
- 🐛 **Issues**: https://github.com/imdlan/AIAgentGuard/issues
- 📧 **Email**: imdlan@users.noreply.github.com

## Recognition

Contributors will be recognized in:
- CONTRIBUTORS.md file
- Release notes
- Project documentation

Thank you for contributing to AI AgentGuard! 🛡️

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
