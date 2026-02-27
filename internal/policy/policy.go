package policy

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/imdlan/AIAgentGuard/internal/audit"
	"github.com/imdlan/AIAgentGuard/pkg/model"
	"gopkg.in/yaml.v3"
)

// DefaultConfigPaths returns possible paths for the policy configuration file
func DefaultConfigPaths() []string {
	homeDir, _ := os.UserHomeDir()
	return []string{
		".agent-guard.yaml",
		filepath.Join(homeDir, ".agent-guard.yaml"),
		filepath.Join(homeDir, ".config", "agent-guard", "config.yaml"),
		"/etc/agent-guard/config.yaml",
	}
}

// LoadConfig loads the policy configuration from a file
// If path is empty, it will search in default locations
func LoadConfig(path string) (*model.PolicyConfig, error) {
	var configPath string

	if path != "" {
		configPath = path
	} else {
		// Search for config in default locations
		for _, defaultPath := range DefaultConfigPaths() {
			if _, err := os.Stat(defaultPath); err == nil {
				configPath = defaultPath
				break
			}
		}
	}

	// If no config found, return default config
	if configPath == "" {
		return GetDefaultConfig(), nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg model.PolicyConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &cfg, nil
}

// GetDefaultConfig returns a default policy configuration
func GetDefaultConfig() *model.PolicyConfig {
	return &model.PolicyConfig{
		Version: 1,
		Filesystem: model.FSRule{
			Allow: []string{"./", "/tmp"},
			Deny:  []string{"~/.ssh", "/etc", "/usr/bin"},
		},
		Shell: model.ShellRule{
			Allow: []string{"ls", "cat", "echo", "pwd", "cd"},
			Deny:  []string{"rm", "sudo", "curl", "wget", "chmod 777"},
		},
		Network: model.NetworkRule{
			Allow: []string{},
			Deny:  []string{"*"},
		},
		Secrets: model.SecretsRule{
			BlockEnv: []string{
				"AWS_SECRET_ACCESS_KEY",
				"GITHUB_TOKEN",
				"OPENAI_API_KEY",
			},
		},
		Sandbox: model.SandboxRule{
			DisableNetwork: true,
			ReadonlyRoot:   false,
		},
	}
}

// SaveConfig saves the policy configuration to a file
func SaveConfig(cfg *model.PolicyConfig, path string) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// IsPathAllowed checks if a path is allowed according to the filesystem policy
func IsPathAllowed(path string, cfg *model.PolicyConfig) bool {
	// Expand home directory if present
	homeDir, _ := os.UserHomeDir()
	expandedPath := strings.Replace(path, "~", homeDir, 1)
	expandedPath = filepath.Clean(expandedPath)

	// Check deny list first (deny takes precedence)
	for _, denyPattern := range cfg.Filesystem.Deny {
		denyExpanded := strings.Replace(denyPattern, "~", homeDir, 1)
		denyExpanded = filepath.Clean(denyExpanded)

		matched, err := filepath.Match(denyExpanded, expandedPath)
		if err == nil && matched {
			_ = audit.LogPolicyViolation("filesystem", path, fmt.Sprintf("matches deny pattern: %s", denyPattern))
			return false
		}

		// Check if path is within denied directory
		if strings.HasPrefix(expandedPath, denyExpanded+string(filepath.Separator)) {
			_ = audit.LogPolicyViolation("filesystem", path, fmt.Sprintf("within denied directory: %s", denyPattern))
			return false
		}
	}

	// If no allow list is specified, allow by default (except denied paths)
	if len(cfg.Filesystem.Allow) == 0 {
		return true
	}

	// Check allow list
	for _, allowPattern := range cfg.Filesystem.Allow {
		allowExpanded := strings.Replace(allowPattern, "~", homeDir, 1)
		allowExpanded = filepath.Clean(allowExpanded)

		matched, err := filepath.Match(allowExpanded, expandedPath)
		if err == nil && matched {
			return true
		}

		// Check if path is within allowed directory
		if strings.HasPrefix(expandedPath, allowExpanded+string(filepath.Separator)) {
			return true
		}
	}

	// Log denial if not in allow list
	_ = audit.LogPolicyViolation("filesystem", path, "not in allow list")
	return false
}

// IsCommandAllowed checks if a command is allowed according to the shell policy
func IsCommandAllowed(cmd string, cfg *model.PolicyConfig) bool {
	// Extract the base command (first word)
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return false
	}

	baseCmd := parts[0]

	// Check deny list first
	for _, denyCmd := range cfg.Shell.Deny {
		if strings.Contains(baseCmd, denyCmd) {
			_ = audit.LogPolicyViolation("shell", cmd, fmt.Sprintf("matches deny pattern: %s", denyCmd))
			_ = audit.LogCommandBlocked(cmd, fmt.Sprintf("blocked by deny list: %s", denyCmd))
			return false
		}
	}

	// If no allow list, allow by default (except denied commands)
	if len(cfg.Shell.Allow) == 0 {
		return true
	}

	// Check allow list
	for _, allowCmd := range cfg.Shell.Allow {
		if baseCmd == allowCmd {
			return true
		}
	}

	_ = audit.LogPolicyViolation("shell", cmd, "not in allow list")
	_ = audit.LogCommandBlocked(cmd, "not in allow list")
	return false
}

// ShouldBlockEnv checks if an environment variable should be blocked
func ShouldBlockEnv(envVar string, cfg *model.PolicyConfig) bool {
	for _, blockPattern := range cfg.Secrets.BlockEnv {
		if strings.EqualFold(envVar, blockPattern) {
			return true
		}
	}
	return false
}

// ApplyEnvPolicy filters environment variables according to the policy
func ApplyEnvPolicy(env []string, cfg *model.PolicyConfig) []string {
	var filtered []string

	for _, e := range env {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) < 2 {
			continue
		}

		if !ShouldBlockEnv(parts[0], cfg) {
			filtered = append(filtered, e)
		} else {
			_ = audit.LogPolicyViolation("secrets", parts[0], "blocked by secrets policy")
		}
	}

	return filtered
}
