package scanner

import (
	"github.com/imdlan/AIAgentGuard/pkg/model"
)

func RunAllScans() model.PermissionResult {
	return model.PermissionResult{
		Filesystem:  ScanFilesystem(),
		Shell:       ScanShell(),
		Network:     ScanNetwork(),
		Secrets:     ScanSecrets(),
		FileContent: ScanFileContents(),
		Dependencies: ScanDependencies(),
	}
}

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
	case "filecontent":
		return ScanFileContents()
	case "dependencies":
		return ScanDependencies()
	default:
		return model.Low
	}
}
