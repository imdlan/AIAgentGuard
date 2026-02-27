package scanner

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"

	"github.com/imdlan/AIAgentGuard/pkg/model"
)

// Key pattern regular expressions for detecting secrets in file contents
var keyPatterns = []struct {
	name    string
	pattern *regexp.Regexp
	risk    model.RiskLevel
}{
	{
		name:    "AWS Access Key ID",
		pattern: regexp.MustCompile(`(?:A3T[A-Z0-9]|AKIA|ASIA|ABIA|ACCA)[A-Z0-9]{16}`),
		risk:    model.High,
	},
	{
		name:    "AWS Secret Access Key",
		pattern: regexp.MustCompile(`(?i)aws[_\s]?secret[_\s]?access[_\s]?key["\s:=]+[A-Za-z0-9/+=]{40}`),
		risk:    model.High,
	},
	{
		name:    "GitHub Token",
		pattern: regexp.MustCompile(`ghp_[a-zA-Z0-9]{36}|gho_[a-zA-Z0-9]{36}|ghu_[a-zA-Z0-9]{36}`),
		risk:    model.High,
	},
	{
		name:    "GitHub Personal Access Token (Classic)",
		pattern: regexp.MustCompile(`ghp_[a-zA-Z0-9]{36}`),
		risk:    model.High,
	},
	{
		name:    "GitHub OAuth App Token",
		pattern: regexp.MustCompile(`gho_[a-zA-Z0-9]{36}`),
		risk:    model.High,
	},
	{
		name:    "OpenAI API Key",
		pattern: regexp.MustCompile(`sk-[a-zA-Z0-9]{48}`),
		risk:    model.High,
	},
	{
		name:    "Anthropic API Key",
		pattern: regexp.MustCompile(`sk-ant-[a-zA-Z0-9_-]{95}`),
		risk:    model.High,
	},
	{
		name:    "Google API Key",
		pattern: regexp.MustCompile(`AIza[A-Za-z0-9_-]{35}`),
		risk:    model.High,
	},
	{
		name:    "Google Cloud OAuth",
		pattern: regexp.MustCompile(`[0-9]+-[a-zA-Z0-9_]{32}\.apps\.googleusercontent\.com`),
		risk:    model.High,
	},
	{
		name:    "Stripe API Key",
		pattern: regexp.MustCompile(`sk_live_[A-Za-z0-9]{24,34}|sk_test_[A-Za-z0-9]{24,34}`),
		risk:    model.High,
	},
	{
		name:    "Slack Token",
		pattern: regexp.MustCompile(`xox[baprs]-[A-Za-z0-9-]{10,}`),
		risk:    model.Medium,
	},
	{
		name:    "Private Key (RSA/DSA/OpenSSH)",
		pattern: regexp.MustCompile(`-----BEGIN[A-Z]+ PRIVATE KEY-----`),
		risk:    model.High,
	},
	{
		name:    "Database Connection String",
		pattern: regexp.MustCompile(`(?i)(mongodb|mysql|postgres|redis)://[^[:space:]"']+:[^[:space:]"']+@`),
		risk:    model.High,
	},
	{
		name:    "Generic API Key/Token",
		pattern: regexp.MustCompile(`(?i)(api[_-]?key|api[_-]?token|secret[_-]?key|access[_-]?token)["\s:=]+[A-Za-z0-9_\-\.]{20,}`),
		risk:    model.Medium,
	},
	{
		name:    "Bearer Token",
		pattern: regexp.MustCompile(`Bearer [A-Za-z0-9\-_\.=]+`),
		risk:    model.Medium,
	},
}

// Common file paths that may contain secrets/keys
var secretFilePaths = []string{
	// AWS
	".aws/credentials",
	".aws/config",

	// SSH/GPG
	".ssh/id_rsa",
	".ssh/id_ed25519",
	".ssh/id_dsa",
	".ssh/id_ecdsa",
	".gnupg/private-keys-v1.d/*.key",

	// Git
	".git/config",

	// Configuration files
	".env",
	".env.local",
	".env.production",
	".config/gcloud/credentials.db",

	// Application specific
	"application.properties",
	"application.yml",
	"settings.py",
	".npmrc",
	".pypirc",
}

// ScanFileContents scans common file locations for secrets and keys
func ScanFileContents() model.RiskLevel {
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		if u, err := os.UserHomeDir(); err == nil {
			homeDir = u
		}
	}

	highRiskCount := 0
	mediumRiskCount := 0
	var findings []string

	// Scan common secret file paths
	for _, filePath := range secretFilePaths {
		fullPath := filepath.Join(homeDir, filePath)

		// Skip if file doesn't exist
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			continue
		}

		// Scan the file for key patterns
		risk, detected := scanFileForKeys(fullPath)
		if risk != model.Low {
			if risk == model.High {
				highRiskCount++
			} else if risk == model.Medium {
				mediumRiskCount++
			}
			findings = append(findings, fmt.Sprintf("%s: %v", fullPath, detected))
		}
	}

	// Additional platform-specific checks
	switch runtime.GOOS {
	case "darwin", "linux":
		// Check for config files in .config directory
		configDir := filepath.Join(homeDir, ".config")
		if _, err := os.Stat(configDir); err == nil {
			filepath.Walk(configDir, func(path string, info os.FileInfo, err error) error {
				if err != nil || info.IsDir() {
					return nil
				}

				// Check for common config file names
				baseName := filepath.Base(path)
				if isConfigFile(baseName) {
					risk, detected := scanFileForKeys(path)
					if risk != model.Low {
						if risk == model.High {
							highRiskCount++
						} else if risk == model.Medium {
							mediumRiskCount++
						}
						findings = append(findings, fmt.Sprintf("%s: %v", path, detected))
					}
				}
				return nil
			})
		}
	}

	// Determine overall risk level
	if highRiskCount > 0 {
		return model.High
	}
	if mediumRiskCount > 2 {
		return model.Medium
	}
	if mediumRiskCount > 0 {
		return model.Low
	}
	return model.Low
}

// scanFileForKeys scans a single file for key patterns
func scanFileForKeys(filePath string) (model.RiskLevel, []string) {
	file, err := os.Open(filePath)
	if err != nil {
		return model.Low, nil
	}
	defer file.Close()

	var detectedPatterns []string
	maxRisk := model.Low

	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// Check each key pattern
		for _, kp := range keyPatterns {
			if kp.pattern.MatchString(line) {
				// Avoid duplicate detections
				found := false
				for _, d := range detectedPatterns {
					if d == kp.name {
						found = true
						break
					}
				}

				if !found {
					detectedPatterns = append(detectedPatterns, kp.name)
					if kp.risk == model.High {
						maxRisk = model.High
					} else if maxRisk == model.Low && kp.risk == model.Medium {
						maxRisk = model.Medium
					}
				}
			}
		}
	}

	return maxRisk, detectedPatterns
}

// isConfigFile checks if a filename is a common configuration file
func isConfigFile(filename string) bool {
	configExtensions := map[string]bool{
		".env":        true,
		".conf":       true,
		".config":     true,
		".ini":        true,
		".json":       true,
		".yml":        true,
		".yaml":       true,
		".properties": true,
		".toml":       true,
		".pem":        true,
		".key":        true,
	}

	configFiles := map[string]bool{
		"config":           true,
		"settings":         true,
		"credentials":      true,
		".npmrc":           true,
		".pypirc":          true,
		"docker-compose":   true,
		"dockerfile":       true,
		"application.yml":  true,
		"application.yaml": true,
	}

	// Check exact filename match
	if configFiles[filename] {
		return true
	}

	// Check extension
	ext := filepath.Ext(filename)
	return configExtensions[ext]
}

// GetKeyFindings returns detailed findings from file content scanning
func GetKeyFindings() []string {
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		if u, err := os.UserHomeDir(); err == nil {
			homeDir = u
		}
	}

	var findings []string

	// Scan common secret file paths
	for _, filePath := range secretFilePaths {
		fullPath := filepath.Join(homeDir, filePath)

		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			continue
		}

		_, detected := scanFileForKeys(fullPath)
		if len(detected) > 0 {
			findings = append(findings, fmt.Sprintf("%s: %v", fullPath, detected))
		}
	}

	return findings
}
