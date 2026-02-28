package scanner

import (
	"github.com/imdlan/AIAgentGuard/internal/scanner/multilang"
	"github.com/imdlan/AIAgentGuard/pkg/model"
)

// RunAllScans runs all security scans and returns the results
// This is the main entry point for scanning
func RunAllScans() model.PermissionResult {
	result, _ := RunAllScansDetailed()
	return result
}

// RunAllScansDetailed runs all security scans and returns detailed information
// Returns both the permission results and detailed risk information
func RunAllScansDetailed() (model.PermissionResult, []model.RiskDetail) {
	details := []model.RiskDetail{}

	// Filesystem scan with details
	fsRisk, fsDetails := ScanFilesystemDetailed()
	details = append(details, fsDetails...)

	// Shell scan with details
	shellRisk, shellDetails := ScanShellDetailed()
	details = append(details, shellDetails...)

	// Process scan with details (NEW)
	_, procDetails := ScanProcessesDetailed()
	details = append(details, procDetails...)

	// Network scan with details (NEW)
	netRisk, netDetails := ScanNetworkDetailed()
	details = append(details, netDetails...)
	secretsRisk := ScanSecrets()

	// File content scan
	fcRisk := ScanFileContents()

	// Dependencies scan
	depRisk := ScanDependencies()

	// Multi-language dependencies
	npmRisk := multilang.ScanNpmDependencies()
	pipRisk := multilang.ScanPipDependencies()
	cargoRisk := multilang.ScanCargoDependencies()

	result := model.PermissionResult{
		Filesystem:   fsRisk,
		Shell:        shellRisk,
		Network:      netRisk,
		Secrets:      secretsRisk,
		FileContent:  fcRisk,
		Dependencies: depRisk,
		NpmDeps:      npmRisk,
		PipDeps:      pipRisk,
		CargoDeps:    cargoRisk,
	}

	return result, details
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
	case "npmdeps":
		return multilang.ScanNpmDependencies()
	case "pipdeps":
		return multilang.ScanPipDependencies()
	case "cargodeps":
		return multilang.ScanCargoDependencies()
	default:
		return model.Low
	}
}
