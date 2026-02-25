package scanner

import (
	"github.com/imdlan/AIAgentGuard/pkg/model"
)

// RunAllScans executes all security scans and returns the combined result
func RunAllScans() model.PermissionResult {
	return model.PermissionResult{
		Filesystem: ScanFilesystem(),
		Shell:      ScanShell(),
		Network:    ScanNetwork(),
		Secrets:    ScanSecrets(),
	}
}

// RunSpecificScan executes a specific security scan by category
func RunSpecificScan(category string) model.RiskLevel {
	switch category {
	case "filesystem":
		return ScanFilesystem()
	case "shell":
		return ScanShell()
	case "network":
		return ScanNetwork()
	case "secrets":
		return ScanSecrets()
	default:
		return model.Low
	}
}
