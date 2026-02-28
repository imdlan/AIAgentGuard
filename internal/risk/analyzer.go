package risk

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/imdlan/AIAgentGuard/internal/scanner"
	"github.com/imdlan/AIAgentGuard/pkg/model"
)

// Analyze performs risk analysis on permission scan results
func Analyze(result model.PermissionResult) model.ScanReport {
	overall := calculateOverallRisk(result)

	details := generateRiskDetails(result)

	return model.ScanReport{
		ToolName: "local-agent",
		Results:  result,
		Overall:  overall,
		Details:  details,
	}
}

// AnalyzeWithDetails performs risk analysis using pre-collected detailed information
func AnalyzeWithDetails(result model.PermissionResult, details []model.RiskDetail) model.ScanReport {
	overall := calculateOverallRisk(result)

	return model.ScanReport{
		ToolName: "local-agent",
		Results:  result,
		Overall:  overall,
		Details:  details,
	}
}


// calculateOverallRisk determines the overall risk level based on individual scan results
func calculateOverallRisk(result model.PermissionResult) model.RiskLevel {
	// Critical risks take highest priority
	if result.Shell == model.Critical {
		return model.Critical
	}

	// High risks from filesystem or secrets
	if result.Filesystem == model.High || result.Secrets == model.High {
		return model.High
	}

	// Shell execution access is always high risk
	if result.Shell == model.High {
		return model.High
	}

	// Network access with other risks
	if result.Network == model.Medium {
		if result.Filesystem == model.Medium || result.Secrets == model.Medium {
			return model.High
		}
		return model.Medium
	}

	// Medium filesystem or secrets risk
	if result.Filesystem == model.Medium || result.Secrets == model.Medium {
		return model.Medium
	}

	// Everything is low
	return model.Low
}

// generateRiskDetails creates detailed risk information for the report
func generateRiskDetails(result model.PermissionResult) []model.RiskDetail {
	var details []model.RiskDetail

	// Filesystem risks
	switch result.Filesystem {
	case model.High:
		writablePaths := getWritableSensitivePaths()
		if len(writablePaths) > 0 {
			details = append(details, model.RiskDetail{
				Type:        model.High,
				Category:    "filesystem",
				Description: "Writable access to sensitive directories",
				Path:        strings.Join(writablePaths, ", "),
			})
		} else {
			details = append(details, model.RiskDetail{
				Type:        model.High,
				Category:    "filesystem",
				Description: "Writable access to sensitive directories",
			})
		}
	case model.Medium:
		accessiblePaths := getAccessibleSensitivePaths()
		if len(accessiblePaths) > 0 {
			details = append(details, model.RiskDetail{
				Type:        model.Medium,
				Category:    "filesystem",
				Description: "Read access to sensitive directories",
				Path:        strings.Join(accessiblePaths, ", "),
			})
		}
	}

	// Shell execution risks
	switch result.Shell {
	case model.Critical:
		details = append(details, model.RiskDetail{
			Type:        model.Critical,
			Category:    "shell",
			Description: "Shell execution with admin/sudo privileges",
		})
	case model.High:
		details = append(details, model.RiskDetail{
			Type:        model.High,
			Category:    "shell",
			Description: "Unrestricted shell command execution",
		})
	}

	// Network risks
	if result.Network == model.Medium {
		details = append(details, model.RiskDetail{
			Type:        model.Medium,
			Category:    "network",
			Description: "External network access enabled (8.8.8.8:53 reachable)",
		})
	}

	// Secrets exposure risks
	if result.Secrets == model.High {
		exposed := scanner.GetExposedSecrets()
		details = append(details, model.RiskDetail{
			Type:        model.High,
			Category:    "secrets",
			Description: "Sensitive environment variables accessible",
			Path:        strings.Join(exposed, ", "),
		})
	}

	return details
}

// GetScoreNumeric returns a numeric score for the risk level (0-100)
func GetScoreNumeric(level model.RiskLevel) int {
	switch level {
	case model.Low:
		return 25
	case model.Medium:
		return 50
	case model.High:
		return 75
	case model.Critical:
		return 100
	default:
		return 0
	}
}

// getWritableSensitivePaths returns list of writable sensitive paths
func getWritableSensitivePaths() []string {
	var writable []string

	// Get sensitive paths from scanner
	sensitivePaths := []string{
		"/etc",
		"/usr/bin",
		"/usr/sbin",
	}

	homeDir := os.Getenv("HOME")
	switch runtime.GOOS {
	case "darwin", "linux":
		if homeDir != "" {
			sensitivePaths = append(sensitivePaths,
				filepath.Join(homeDir, ".ssh"),
				filepath.Join(homeDir, ".gnupg"),
				filepath.Join(homeDir, ".config"),
				filepath.Join(homeDir, ".aws"),
			)
		}
	case "windows":
		if homeDir != "" {
			sensitivePaths = append(sensitivePaths,
				filepath.Join(homeDir, ".ssh"),
				filepath.Join(homeDir, "AppData", "Local"),
				filepath.Join(homeDir, "AppData", "Roaming"),
			)
		}
	}

	for _, path := range sensitivePaths {
		if info, err := os.Stat(path); err == nil {
			if info.IsDir() {
				testFile := filepath.Join(path, ".agentguard_test")
				if file, err := os.Create(testFile); err == nil {
					file.Close()
					os.Remove(testFile)
					writable = append(writable, path)
				}
			}
		}
	}

	return writable
}

// getAccessibleSensitivePaths returns list of accessible (readable) sensitive paths
func getAccessibleSensitivePaths() []string {
	var accessible []string

	sensitivePaths := []string{
		"/etc",
		"/usr/bin",
		"/usr/sbin",
	}

	homeDir := os.Getenv("HOME")
	switch runtime.GOOS {
	case "darwin", "linux":
		if homeDir != "" {
			sensitivePaths = append(sensitivePaths,
				filepath.Join(homeDir, ".ssh"),
				filepath.Join(homeDir, ".gnupg"),
				filepath.Join(homeDir, ".config"),
				filepath.Join(homeDir, ".aws"),
			)
		}
	case "windows":
		if homeDir != "" {
			sensitivePaths = append(sensitivePaths,
				filepath.Join(homeDir, ".ssh"),
				filepath.Join(homeDir, "AppData", "Local"),
				filepath.Join(homeDir, "AppData", "Roaming"),
			)
		}
	}

	for _, path := range sensitivePaths {
		if _, err := os.Stat(path); err == nil {
			accessible = append(accessible, path)
		}
	}

	return accessible
}
