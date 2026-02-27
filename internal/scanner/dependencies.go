package scanner

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/imdlan/AIAgentGuard/pkg/model"
)

type DependencyInfo struct {
	Path    string
	Version string
}

type VulnerabilityInfo struct {
	ID               string
	Severity         string
	Description      string
	FixedIn          string
	AffectedVersions []string
}

func ScanDependencies() model.RiskLevel {
	vulns, err := findVulnerabilities()
	if err != nil {
		return model.Medium
	}

	if len(vulns) == 0 {
		return model.Low
	}

	return calculateVulnerabilityRisk(vulns)
}

func findVulnerabilities() ([]VulnerabilityInfo, error) {
	var vulns []VulnerabilityInfo

	if _, err := exec.LookPath("govulncheck"); err != nil {
		return findVulnerabilitiesBasic()
	}

	tmpFile, err := os.CreateTemp("", "vulncheck-*.txt")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	cmd := exec.Command("govulncheck", "-json", "./...")
	cmd.Dir = getModuleRoot()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return findVulnerabilitiesBasic()
	}

	return vulns, nil
}

func findVulnerabilitiesBasic() ([]VulnerabilityInfo, error) {
	var vulns []VulnerabilityInfo
	deps, err := getGoDependencies()
	if err != nil {
		return nil, err
	}

	vulnDB := getKnownVulnerabilities()

	for _, dep := range deps {
		if knownVulns, ok := vulnDB[dep.Path]; ok {
			for _, vuln := range knownVulns {
				if isAffectedVersion(dep.Version, vuln.AffectedVersions) {
					vulns = append(vulns, vuln)
				}
			}
		}
	}

	return vulns, nil
}

func getGoDependencies() ([]DependencyInfo, error) {
	var deps []DependencyInfo

	goModPath := findGoMod()
	if goModPath == "" {
		return deps, fmt.Errorf("go.mod not found")
	}

	file, err := os.Open(goModPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open go.mod: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	inRequire := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "require (") {
			inRequire = true
			continue
		}

		if inRequire && line == ")" {
			break
		}

		if strings.HasPrefix(line, "require ") || inRequire {
			line = strings.TrimPrefix(line, "require ")
			line = strings.TrimSpace(line)

			if strings.HasPrefix(line, "//") || strings.Contains(line, "// indirect") {
				continue
			}

			parts := strings.Fields(line)
			if len(parts) >= 2 {
				deps = append(deps, DependencyInfo{
					Path:    parts[0],
					Version: parts[1],
				})
			}
		}
	}

	return deps, nil
}

func findGoMod() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}

	for {
		goModPath := fmt.Sprintf("%s/go.mod", dir)
		if _, err := os.Stat(goModPath); err == nil {
			return goModPath
		}

		parentDir := fmt.Sprintf("%s/..", dir)
		if parentDir == dir {
			break
		}
		dir = parentDir
	}

	return ""
}

func getModuleRoot() string {
	goModPath := findGoMod()
	if goModPath == "" {
		return "."
	}
	return fmt.Sprintf("%s/..", goModPath)
}

func isAffectedVersion(version string, affectedVersions []string) bool {
	for _, affected := range affectedVersions {
		if strings.HasPrefix(version, affected) {
			return true
		}
	}
	return false
}

func calculateVulnerabilityRisk(vulns []VulnerabilityInfo) model.RiskLevel {
	criticalCount := 0
	highCount := 0

	for _, vuln := range vulns {
		switch strings.ToUpper(vuln.Severity) {
		case "CRITICAL":
			criticalCount++
		case "HIGH":
			highCount++
		}
	}

	if criticalCount > 0 {
		return model.Critical
	}
	if highCount > 0 || len(vulns) > 5 {
		return model.High
	}
	if len(vulns) > 0 {
		return model.Medium
	}
	return model.Low
}

func getKnownVulnerabilities() map[string][]VulnerabilityInfo {
	return map[string][]VulnerabilityInfo{
		"github.com/gin-gonic/gin": {
			{
				ID:               "CVE-2024-2476",
				Severity:         "HIGH",
				Description:      "Gin framework vulnerable to path traversal before v1.10.0",
				FixedIn:          "v1.10.0",
				AffectedVersions: []string{"v1.9.1", "v1.9.0", "v1.8.2", "v1.8.1", "v1.8.0", "v1.7.7", "v1.7.6", "v1.7.5", "v1.7.4", "v1.7.3"},
			},
		},
		"github.com/spf13/cobra": {
			{
				ID:               "CVE-2021-38561",
				Severity:         "MEDIUM",
				Description:      "Cobra vulnerable to directory traversal via bash completion",
				FixedIn:          "v1.2.0",
				AffectedVersions: []string{"v1.1.3", "v1.1.2", "v1.1.1", "v1.1.0", "v1.0.0"},
			},
		},
		"gopkg.in/yaml.v3": {
			{
				ID:               "CVE-2022-28948",
				Severity:         "HIGH",
				Description:      "YAML parser vulnerable to arbitrary code execution",
				FixedIn:          "v3.0.1",
				AffectedVersions: []string{"v3.0.0"},
			},
		},
		"github.com/golang/protobuf": {
			{
				ID:               "CVE-2021-3121",
				Severity:         "HIGH",
				Description:      "Protobuf decoder vulnerable to panics on malformed input",
				FixedIn:          "v1.5.2",
				AffectedVersions: []string{"v1.5.1", "v1.5.0", "v1.4.3", "v1.4.2", "v1.4.1"},
			},
		},
	}
}
