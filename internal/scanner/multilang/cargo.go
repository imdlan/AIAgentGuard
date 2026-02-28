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

// CargoAuditResult represents the JSON output from 'cargo audit'
type CargoAuditResult struct {
	Vulnerabilities []CargoVulnerability `json:"vulnerabilities"`
}

type CargoVulnerability struct {
	Advisory CargoAdvisory `json:"advisory"`
	Package  struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	} `json:"package"`
	Affected struct {
		Arch      []string `json:"arch"`
		OS        []string `json:"os"`
		Functions []string `json:"functions"`
	} `json:"affected"`
}

type CargoAdvisory struct {
	ID                 string   `json:"advisory_id"`
	Package            string   `json:"package"`
	Title              string   `json:"title"`
	Description        string   `json:"description"`
	Date               string   `json:"date"`
	URLs               []string `json:"url"`
	Keywords           []string `json:"keywords"`
	Type               string   `json:"type"`
	Cvss               string   `json:"cvss,omitempty"`
	Aliases            []string `json:"aliases"`
	Related            []string `json:"related"`
	References         []string `json:"references"`
	PatchedVersions    []string `json:"patched_versions"`
	UnaffectedVersions []string `json:"unaffected_versions"`
	Severity           string   `json:"severity"`
}

// CargoMetadata represents Cargo.toml metadata
type CargoMetadata struct {
	Package struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	} `json:"package"`
	Dependencies []CargoDependency `json:"dependencies"`
}

type CargoDependency struct {
	Name    string `json:"name"`
	Version string `json:"req"`
}

// ScanCargoDependencies scans Rust projects for vulnerable dependencies
func ScanCargoDependencies() model.RiskLevel {
	cargoFiles := findCargoFiles()

	if len(cargoFiles) == 0 {
		return model.Low
	}

	overallRisk := model.Low
	criticalCount := 0
	highCount := 0
	moderateCount := 0

	for _, cargoFile := range cargoFiles {
		risk, vulnCount := scanCargoProject(cargoFile)
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

// findCargoFiles searches for Cargo.toml files
func findCargoFiles() []string {
	var files []string

	root := "."
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			// Skip hidden directories (but not "." which represents current directory)
			// and target directory
			if (info.Name() != "." && strings.HasPrefix(info.Name(), ".")) || info.Name() == "target" {
				return filepath.SkipDir
			}
			return nil
		}

		if filepath.Base(path) == "Cargo.toml" {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		return files
	}

	return files
}

// scanCargoProject scans a single Rust project for vulnerabilities
func scanCargoProject(cargoFilePath string) (model.RiskLevel, VulnerabilityCount) {
	projectDir := filepath.Dir(cargoFilePath)

	// Check if cargo-audit is installed
	if !cargoAuditAvailable() {
		return model.Low, VulnerabilityCount{}
	}

	cmd := exec.Command("cargo", "audit", "--json")
	cmd.Dir = projectDir

	output, err := cmd.Output()
	if err != nil {
		// cargo audit returns non-zero exit code if vulnerabilities found
		if exitErr, ok := err.(*exec.ExitError); ok {
			output = exitErr.Stderr
		} else {
			return model.Low, VulnerabilityCount{}
		}
	}

	var auditResult CargoAuditResult
	if err := json.Unmarshal(output, &auditResult); err != nil {
		return model.Low, VulnerabilityCount{}
	}

	vulnCount := VulnerabilityCount{}
	for _, vuln := range auditResult.Vulnerabilities {
		severity := parseCargoSeverity(vuln.Advisory.Severity, vuln.Advisory.Cvss)

		switch severity {
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

// cargoAuditAvailable checks if cargo-audit is installed
func cargoAuditAvailable() bool {
	// Check if cargo command is available
	if _, err := exec.LookPath("cargo"); err != nil {
		return false
	}

	// Check if cargo-audit plugin is installed
	cmd := exec.Command("cargo", "audit", "--help")
	err := cmd.Run()
	return err == nil
}

// parseCargoSeverity parses severity from cargo audit output
func parseCargoSeverity(severity, cvss string) string {
	// If severity is explicitly provided, use it
	if severity != "" {
		return strings.ToLower(severity)
	}

	// Fallback to CVSS score
	if cvss != "" {
		// CVSS scores range from 0.0 to 10.0
		// 9.0-10.0: Critical
		// 7.0-8.9: High
		// 4.0-6.9: Medium
		// 0.1-3.9: Low
		var score float64
		fmt.Sscanf(cvss, "%f", &score)

		if score >= 9.0 {
			return "critical"
		} else if score >= 7.0 {
			return "high"
		} else if score >= 4.0 {
			return "medium"
		} else if score > 0.0 {
			return "low"
		}
	}

	// Default to medium if no severity info available
	return "medium"
}

// GetCargoVulnerabilityDetails returns detailed vulnerability information
func GetCargoVulnerabilityDetails(projectPath string) ([]string, error) {
	if !cargoAuditAvailable() {
		return nil, fmt.Errorf("cargo-audit is not installed")
	}

	cargoFile := filepath.Join(projectPath, "Cargo.toml")
	if !fileExists(cargoFile) {
		return nil, fmt.Errorf("Cargo.toml not found in %s", projectPath)
	}

	cmd := exec.Command("cargo", "audit", "--json")
	cmd.Dir = projectPath

	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			output = exitErr.Stderr
		} else {
			return nil, err
		}
	}

	var auditResult CargoAuditResult
	if err := json.Unmarshal(output, &auditResult); err != nil {
		return nil, err
	}

	var details []string
	for _, vuln := range auditResult.Vulnerabilities {
		detail := fmt.Sprintf("%s: %s (ID: %s, Severity: %s)",
			vuln.Package.Name,
			vuln.Advisory.Title,
			vuln.Advisory.ID,
			vuln.Advisory.Severity)
		details = append(details, detail)
	}

	return details, nil
}

// GetCargoDependencies returns a list of dependencies from Cargo.toml
func GetCargoDependencies(projectPath string) ([]CargoDependency, error) {
	cargoFile := filepath.Join(projectPath, "Cargo.toml")
	if !fileExists(cargoFile) {
		return nil, fmt.Errorf("Cargo.toml not found in %s", projectPath)
	}

	// Use cargo metadata command to get dependency information
	cmd := exec.Command("cargo", "metadata", "--format-version", "1", "--no-deps")
	cmd.Dir = projectPath

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var metadata CargoMetadata
	if err := json.Unmarshal(output, &metadata); err != nil {
		return nil, err
	}

	return metadata.Dependencies, nil
}
