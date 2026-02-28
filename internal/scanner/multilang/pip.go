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

// PipAuditResult represents the JSON output from 'pip-audit'
type PipAuditResult struct {
	Vulnerabilities []PipVulnerability `json:"vulnerabilities"`
	Dependencies    []PipDependency    `json:"dependencies"`
}

type PipVulnerability struct {
	Name        string      `json:"name"`
	Versions    []string    `json:"versions"`
	ID          string      `json:"id"`
	FixVersions []string    `json:"fix_versions"`
	Aliases     []string    `json:"aliases"`
	Severity    string      `json:"severity"`
	Advisory    PipAdvisory `json:"advisory"`
}

type PipAdvisory struct {
	Description string   `json:"description"`
	References  []string `json:"references"`
}

type PipDependency struct {
	Name            string             `json:"name"`
	Version         string             `json:"version"`
	Vulnerabilities []PipVulnerability `json:"vulnerabilities"`
}

// SafetyVulnerability represents vulnerability data from 'safety' tool
type SafetyVulnerability struct {
	ID          string   `json:"id"`
	PackageName string   `json:"package_name"`
	Versions    []string `json:"versions"`
	Advisory    string   `json:"advisory"`
	Severity    string   `json:"severity"`
	CVE         string   `json:"cve"`
}

// ScanPipDependencies scans Python projects for vulnerable dependencies
func ScanPipDependencies() model.RiskLevel {
	reqFiles := findRequirementFiles()

	if len(reqFiles) == 0 {
		return model.Low
	}

	overallRisk := model.Low
	criticalCount := 0
	highCount := 0
	moderateCount := 0

	for _, reqFile := range reqFiles {
		risk, vulnCount := scanPipProject(reqFile)
		if risk > overallRisk {
			overallRisk = risk
		}

		criticalCount += vulnCount.Critical
		highCount += vulnCount.High
		moderateCount += vulnCount.Moderate
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

// findRequirementFiles searches for requirements.txt and pyproject.toml files
func findRequirementFiles() []string {
	var files []string

	root := "."
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			// Skip hidden directories (but not "." which represents current directory),
			// __pycache__, venv, .venv, and node_modules
			if (info.Name() != "." && strings.HasPrefix(info.Name(), ".")) ||
				info.Name() == "__pycache__" ||
				info.Name() == "venv" ||
				info.Name() == ".venv" ||
				info.Name() == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}

		base := filepath.Base(path)
		if base == "requirements.txt" || base == "pyproject.toml" || base == "Pipfile" || base == "Pipfile.lock" {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		return files
	}

	return files
}

// scanPipProject scans a single Python project for vulnerabilities
func scanPipProject(reqFilePath string) (model.RiskLevel, VulnerabilityCount) {
	projectDir := filepath.Dir(reqFilePath)

	// Try pip-audit first
	if pipAuditAvailable() {
		return scanWithPipAudit(projectDir)
	}

	// Fallback to safety
	if safetyAvailable() {
		return scanWithSafety(projectDir)
	}

	return model.Low, VulnerabilityCount{}
}

// pipAuditAvailable checks if pip-audit is installed
func pipAuditAvailable() bool {
	_, err := exec.LookPath("pip-audit")
	return err == nil
}

// safetyAvailable checks if safety is installed
func safetyAvailable() bool {
	_, err := exec.LookPath("safety")
	return err == nil
}

// scanWithPipAudit uses pip-audit to check for vulnerabilities
func scanWithPipAudit(projectDir string) (model.RiskLevel, VulnerabilityCount) {
	cmd := exec.Command("pip-audit", "--format", "json")
	cmd.Dir = projectDir

	output, err := cmd.Output()
	if err != nil {
		return model.Low, VulnerabilityCount{}
	}

	var auditResult PipAuditResult
	if err := json.Unmarshal(output, &auditResult); err != nil {
		return model.Low, VulnerabilityCount{}
	}

	vulnCount := VulnerabilityCount{}
	for _, vuln := range auditResult.Vulnerabilities {
		switch strings.ToLower(vuln.Severity) {
		case "critical":
			vulnCount.Critical++
		case "high":
			vulnCount.High++
		case "medium", "moderate":
			vulnCount.Moderate++
		case "low":
			vulnCount.Low++
		default:
			vulnCount.Moderate++
		}
	}

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

// scanWithSafety uses safety to check for vulnerabilities
func scanWithSafety(projectDir string) (model.RiskLevel, VulnerabilityCount) {
	cmd := exec.Command("safety", "check", "--json")
	cmd.Dir = projectDir

	output, err := cmd.Output()
	if err != nil {
		return model.Low, VulnerabilityCount{}
	}

	var vulnerabilities []SafetyVulnerability
	if err := json.Unmarshal(output, &vulnerabilities); err != nil {
		return model.Low, VulnerabilityCount{}
	}

	vulnCount := VulnerabilityCount{}
	for _, vuln := range vulnerabilities {
		switch strings.ToLower(vuln.Severity) {
		case "critical":
			vulnCount.Critical++
		case "high":
			vulnCount.High++
		case "medium", "moderate":
			vulnCount.Moderate++
		case "low":
			vulnCount.Low++
		default:
			vulnCount.Moderate++
		}
	}

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

// GetPipVulnerabilityDetails returns detailed vulnerability information
func GetPipVulnerabilityDetails(projectPath string) ([]string, error) {
	var cmd *exec.Cmd

	reqFile := filepath.Join(projectPath, "requirements.txt")
	pyprojectFile := filepath.Join(projectPath, "pyproject.toml")
	pipfile := filepath.Join(projectPath, "Pipfile")

	if pipAuditAvailable() {
		cmd = exec.Command("pip-audit", "--format", "json")
	} else if fileExists(reqFile) || fileExists(pyprojectFile) || fileExists(pipfile) {
		return nil, fmt.Errorf("no Python vulnerability scanner available (pip-audit or safety)")
	} else {
		return nil, fmt.Errorf("no Python dependency file found in %s", projectPath)
	}

	cmd.Dir = projectPath

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var auditResult PipAuditResult
	if err := json.Unmarshal(output, &auditResult); err != nil {
		return nil, err
	}

	var details []string
	for _, vuln := range auditResult.Vulnerabilities {
		detail := fmt.Sprintf("%s: %s (%s)", vuln.Name, vuln.Advisory.Description, vuln.Severity)
		details = append(details, detail)
	}

	return details, nil
}
