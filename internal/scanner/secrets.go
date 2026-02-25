package scanner

import (
	"os"
	"strings"

	"github.com/imdlan/AIAgentGuard/pkg/model"
)

// Sensitive environment variable patterns to check
var sensitiveEnvPatterns = []string{
	"SECRET",
	"TOKEN",
	"KEY",
	"PASSWORD",
	"API_KEY",
	"ACCESS_KEY",
	"CREDENTIALS",
}

// Specific high-risk environment variables
var highRiskEnvVars = []string{
	"AWS_SECRET_ACCESS_KEY",
	"AWS_SESSION_TOKEN",
	"GITHUB_TOKEN",
	"GITHUB_PAT",
	"OPENAI_API_KEY",
	"ANTHROPIC_API_KEY",
	"GOOGLE_API_KEY",
	"AZURE_CLIENT_SECRET",
	"DATABASE_URL",
	"REDIS_PASSWORD",
	"SSH_PRIVATE_KEY",
}

// ScanSecrets checks if sensitive environment variables are accessible
func ScanSecrets() model.RiskLevel {
	highRiskCount := 0
	mediumRiskCount := 0

	// Check for high-risk environment variables
	for _, envVar := range highRiskEnvVars {
		if value := os.Getenv(envVar); value != "" {
			highRiskCount++
		}
	}

	// Check for environment variables with sensitive patterns
	for _, envName := range os.Environ() {
		name := strings.Split(envName, "=")[0]
		upperName := strings.ToUpper(name)

		for _, pattern := range sensitiveEnvPatterns {
			if strings.Contains(upperName, pattern) {
				mediumRiskCount++
				break
			}
		}
	}

	// Determine risk level based on findings
	if highRiskCount > 0 {
		return model.High
	}
	if mediumRiskCount > 2 {
		return model.Medium
	}
	return model.Low
}

// GetExposedSecrets returns a list of exposed sensitive environment variables (names only)
func GetExposedSecrets() []string {
	var exposed []string

	for _, envVar := range highRiskEnvVars {
		if value := os.Getenv(envVar); value != "" {
			// Only return the name, not the value
			exposed = append(exposed, envVar)
		}
	}

	return exposed
}
