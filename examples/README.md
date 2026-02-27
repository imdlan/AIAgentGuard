# AIAgentGuard Configuration Examples

This directory contains example configuration files for different use cases.

## Files

- `basic-policy.yaml` - Basic security policy for everyday use
- `strict-policy.yaml` - High-security policy for production environments
- `development-policy.yaml` - Relaxed policy for development

## Basic Policy (basic-policy.yaml)

```yaml
version: 1

# Filesystem access rules
filesystem:
  allow:
    - ./
    - /tmp
    - /var/log/app
  deny:
    - ~/.ssh
    - ~/.gnupg
    - ~/.aws
    - /etc/passwd
    - /etc/shadow

# Shell command rules
shell:
  allow:
    - ls
    - cat
    - echo
    - pwd
    - cd
    - grep
    - head
    - tail
  deny:
    - rm
    - sudo
    - curl
    - wget
    - nc
    - netcat
    - chmod 777
    - dd

# Network access rules
network:
  allow:
    - api.github.com
    - cdn.jsdelivr.net
  deny:
    - "*"

# Secrets protection
secrets:
  block_env:
    - AWS_SECRET_ACCESS_KEY
    - AWS_SESSION_TOKEN
    - GITHUB_TOKEN
    - OPENAI_API_key
    - ANTHROPIC_API_KEY
    - DATABASE_URL
    - REDIS_PASSWORD

# Sandbox configuration
sandbox:
  disable_network: true
  readonly_root: false
```

## Strict Policy (strict-policy.yaml)

```yaml
version: 1

filesystem:
  allow:
    - /tmp
    - ./build
    - ./dist
  deny:
    - ~/*  # Block all home directory access

shell:
  allow:
    - echo
    - pwd
    - ls
  deny:
    - "*"

network:
  allow: []
  deny: ["*"]

secrets:
  block_env:
    - "*KEY*"
    - "*TOKEN*"
    - "*SECRET*"
    - "*PASSWORD*"

sandbox:
  disable_network: true
  readonly_root: true
```

## Development Policy (development-policy.yaml)

```yaml
version: 1

filesystem:
  allow:
    - ./
    - /tmp
    - ~/.config
  deny:
    - ~/.ssh/id_rsa

shell:
  allow:
    - "*"
  deny:
    - "rm -rf /"
    - ":(){ :|:& };:"

network:
  allow: ["*"]
  deny: []

secrets:
  block_env:
    - AWS_SECRET_ACCESS_KEY
    - GITHUB_TOKEN

sandbox:
  disable_network: false
  readonly_root: false
```

## Usage

```bash
# Use basic policy
agent-guard scan --config examples/basic-policy.yaml

# Use strict policy
agent-guard run --config examples/strict-policy.yaml "npm install"

# Generate your own policy
agent-guard init --path my-policy.yaml
```

## Tips

1. **Start with basic-policy.yaml** - Good balance of security and usability
2. **Use strict-policy.yaml for production** - Maximum security, restricted access
3. **Customize for your needs** - Add your own allow/deny lists
4. **Test policies first** - Use `agent-guard scan --dry-run` to test without enforcing
5. **Review audit logs** - Check `~/.agent-guard/audit.log` for blocked actions
