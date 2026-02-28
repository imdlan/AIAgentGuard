package multilang

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/imdlan/AIAgentGuard/pkg/model"
)

// NpmAuditResult represents the JSON output from 'npm audit --json'
type NpmAuditResult struct {
	AuditReportVersion int                      `json:"auditReportVersion"`
	Vulnerabilities    map[string]Vulnerability `json:"vulnerabilities"`
	Metadata           AuditMetadata            `json:"metadata"`
}

type Vulnerability struct {
	Name            string    `json:"name"`
	Severity        string    `json:"severity"`
	Via             []ViaInfo `json:"via"`
	Effects         []string  `json:"effects"`
	Range           string    `json:"range"`
	Title           string    `json:"title"`
	Module          string    `json:"module"`
	PatchedVersions string    `json:"patched_versions"`
}

type ViaInfo struct {
	Source     int    `json:"source"`
	Name       string `json:"name,omitempty"`
	Dependency string `json:"dependency,omitempty"`
	Title      string `json:"title,omitempty"`
	URL        string `json:"url,omitempty"`
	Severity   string `json:"severity,omitempty"`
}

type AuditMetadata struct {
	Vulnerabilities      VulnerabilityCountJSON `json:"vulnerabilities"`
	Dependencies         int                    `json:"dependencies"`
	DevDependencies      int                    `json:"devDependencies"`
	OptionalDependencies int                    `json:"optionalDependencies"`
	TotalDependencies    int                    `json:"totalDependencies"`
}

// VulnerabilityCountJSON represents JSON format of vulnerability counts
type VulnerabilityCountJSON struct {
	Low      int `json:"low"`
	Moderate int `json:"moderate"`
	High     int `json:"high"`
	Critical int `json:"critical"`
}

// ScanNpmDependencies scans npm/yarn projects for vulnerable dependencies
func ScanNpmDependencies() model.RiskLevel {
	pkgFiles := findPackageFiles()

	if len(pkgFiles) == 0 {
		return model.Low
	}

	overallRisk := model.Low
	criticalCount := 0
	highCount := 0
	moderateCount := 0

	for _, pkgFile := range pkgFiles {
		risk, vulns := scanNpmProject(pkgFile)
		if risk > overallRisk {
			overallRisk = risk
		}

		criticalCount += vulns.Critical
		highCount += vulns.High
		moderateCount += vulns.Moderate
	}

	if criticalCount > 0 {
		return model.Critical
	}
	if highCount > 0 {
		return model.High
	}
	if moderateCount > 0 {
		return model.Medium
	}

	return overallRisk
}

// findPackageFiles searches for package-lock.json and yarn.lock files
func findPackageFiles() []string {
	var files []string

	// Search in current directory and subdirectories
	root := "."
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip errors
		}

		if info.IsDir() {
			// Skip hidden directories (but not "." which represents current directory)
			// and node_modules
			if (info.Name() != "." && strings.HasPrefix(info.Name(), ".")) || info.Name() == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}

		// Check for package lock files
		base := filepath.Base(path)
		if base == "package-lock.json" || base == "yarn.lock" {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		return files
	}

	return files
}

// scanNpmProject scans a single npm/yarn project for vulnerabilities
func scanNpmProject(lockFilePath string) (model.RiskLevel, VulnerabilityCount) {
	// Check if npm/yarn is available
	lockFile := filepath.Base(lockFilePath)
	var cmd *exec.Cmd

	if lockFile == "package-lock.json" {
		if _, err := exec.LookPath("npm"); err != nil {
			return model.Low, VulnerabilityCount{}
		}
		cmd = exec.Command("npm", "audit", "--json")
	} else if lockFile == "yarn.lock" {
		if _, err := exec.LookPath("yarn"); err != nil {
			return model.Low, VulnerabilityCount{}
		}
		cmd = exec.Command("yarn", "audit", "--json")
	} else {
		return model.Low, VulnerabilityCount{}
	}

	projectDir := filepath.Dir(lockFilePath)
	cmd.Dir = projectDir

	output, err := cmd.Output()
	if err != nil {
		// npm audit returns non-zero exit code if vulnerabilities found
		// but still outputs valid JSON
		if exitErr, ok := err.(*exec.ExitError); ok {
			output = exitErr.Stderr
		} else {
			return model.Low, VulnerabilityCount{}
		}
	}

	// Parse npm audit output
	var auditResult NpmAuditResult
	if err := json.Unmarshal(output, &auditResult); err != nil {
		return model.Low, VulnerabilityCount{}
	}

	vulnCount := convertVulnCount(auditResult.Metadata.Vulnerabilities)

	// Determine risk level
	if vulnCount.Critical > 0 {
		return model.Critical, vulnCount
	}
	if vulnCount.High > 0 {
		return model.High, vulnCount
	}
	if vulnCount.Moderate > 0 {
		return model.Medium, vulnCount
	}
	if vulnCount.Low > 0 {
		return model.Low, vulnCount
	}

	return model.Low, vulnCount
}

// convertVulnCount converts JSON VulnerabilityCount to internal type
func convertVulnCount(vulnJSON VulnerabilityCountJSON) VulnerabilityCount {
	return VulnerabilityCount{
		Low:      vulnJSON.Low,
		Moderate: vulnJSON.Moderate,
		High:     vulnJSON.High,
		Critical: vulnJSON.Critical,
	}
}

// GetNpmVulnerabilityDetails returns detailed vulnerability information
func GetNpmVulnerabilityDetails(projectPath string) ([]string, error) {
	var cmd *exec.Cmd

	lockFile := filepath.Join(projectPath, "package-lock.json")
	yarnLockFile := filepath.Join(projectPath, "yarn.lock")

	if fileExists(lockFile) {
		cmd = exec.Command("npm", "audit", "--json")
	} else if fileExists(yarnLockFile) {
		cmd = exec.Command("yarn", "audit", "--json")
	} else {
		return nil, fmt.Errorf("no package lock file found in %s", projectPath)
	}

	cmd.Dir = projectPath

	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			output = exitErr.Stderr
		} else {
			return nil, err
		}
	}

	var auditResult NpmAuditResult
	if err := json.Unmarshal(output, &auditResult); err != nil {
		return nil, err
	}

	var details []string
	for name, vuln := range auditResult.Vulnerabilities {
		detail := fmt.Sprintf("%s: %s (%s)", name, vuln.Title, vuln.Severity)
		details = append(details, detail)
	}

	return details, nil
}

// fileExists checks if a file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
