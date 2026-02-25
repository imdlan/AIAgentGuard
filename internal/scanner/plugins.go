package scanner

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/imdlan/AIAgentGuard/pkg/model"
)

// Suspicious keywords that may indicate dangerous behavior in plugins
var suspiciousKeywords = []string{
	"curl",
	"wget",
	"bash",
	"sh ",
	"exec",
	"eval",
	"system(",
	"subprocess",
	"os.system",
	"chmod 777",
	"rm -rf",
	"mkfs",
	":(){:|:&};:", // fork bomb
}

// ScanPlugins scans a directory for potentially risky plugin scripts
func ScanPlugins(dir string) []model.PluginScanResult {
	var results []model.PluginScanResult

	// Walk through the directory
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		// Skip directories
		if info.IsDir() {
			// Skip common non-plugin directories
			if strings.HasPrefix(filepath.Base(path), ".") ||
				filepath.Base(path) == "node_modules" ||
				filepath.Base(path) == "vendor" {
				return filepath.SkipDir
			}
			return nil
		}

		// Check file extensions for script files
		ext := strings.ToLower(filepath.Ext(path))
		if !isScriptFile(ext) {
			return nil
		}

		// Scan the file for suspicious content
		result := scanPluginFile(path)
		if result != nil {
			results = append(results, *result)
		}

		return nil
	})

	if err != nil {
		return results
	}

	return results
}

// isScriptFile checks if a file extension indicates a script file
func isScriptFile(ext string) bool {
	scriptExts := map[string]bool{
		".sh":   true,
		".bash": true,
		".zsh":  true,
		".py":   true,
		".js":   true,
		".ts":   true,
		".rb":   true,
		".php":  true,
		".pl":   true,
	}
	return scriptExts[ext]
}

// scanPluginFile scans a single plugin file for suspicious patterns
func scanPluginFile(path string) *model.PluginScanResult {
	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer file.Close()

	var detected []string
	var riskLevel model.RiskLevel = model.Low

	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		lineLower := strings.ToLower(line)

		for _, keyword := range suspiciousKeywords {
			if strings.Contains(lineLower, keyword) {
				// Check if we've already detected this keyword
				alreadyDetected := false
				for _, d := range detected {
					if d == keyword {
						alreadyDetected = true
						break
					}
				}

				if !alreadyDetected {
					detected = append(detected, keyword)
				}

				// Update risk level
				if keyword == "curl" || keyword == "wget" || keyword == "eval" || keyword == "exec" {
					riskLevel = model.High
				} else if riskLevel == model.Low {
					riskLevel = model.Medium
				}
			}
		}
	}

	// Only return a result if we found something suspicious
	if len(detected) > 0 {
		reason := formatDetectionReason(detected)
		return &model.PluginScanResult{
			PluginName: filepath.Base(path),
			Path:       path,
			Risk:       riskLevel,
			Reason:     reason,
			Detected:   detected,
		}
	}

	return nil
}

// formatDetectionReason creates a human-readable reason for the detection
func formatDetectionReason(detected []string) string {
	if len(detected) == 0 {
		return "No suspicious patterns detected"
	}

	if len(detected) == 1 {
		return fmt.Sprintf("Contains suspicious keyword: %s", detected[0])
	}

	return fmt.Sprintf("Contains %d suspicious keywords: %s", len(detected), strings.Join(detected, ", "))
}

// GetPluginRiskSummary returns a summary of plugin scan results
func GetPluginRiskSummary(results []model.PluginScanResult) map[model.RiskLevel]int {
	summary := map[model.RiskLevel]int{
		model.Low:      0,
		model.Medium:   0,
		model.High:     0,
		model.Critical: 0,
	}

	for _, result := range results {
		summary[result.Risk]++
	}

	return summary
}
