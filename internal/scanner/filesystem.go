package scanner

import (
	"os"
	"os/user"
	"path/filepath"
	"runtime"

	"github.com/imdlan/AIAgentGuard/pkg/model"
)

// Sensitive paths that should not be accessible by AI agents
var sensitivePaths = []string{
	"/etc",
	"/usr/bin",
	"/usr/sbin",
}

// ScanFilesystem checks if sensitive filesystem paths are accessible
func ScanFilesystem() model.RiskLevel {
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		if u, err := user.Current(); err == nil {
			homeDir = u.HomeDir
		}
	}

	// Add platform-specific sensitive paths
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

	highRiskFound := false
	mediumRiskFound := false

	for _, path := range sensitivePaths {
		if info, err := os.Stat(path); err == nil {
			// Check if path is writable
			if info.IsDir() {
				testFile := filepath.Join(path, ".agentguard_test")
				if file, err := os.Create(testFile); err == nil {
					file.Close()
					os.Remove(testFile)
					highRiskFound = true
				} else {
					mediumRiskFound = true
				}
			}
		}
	}

	if highRiskFound {
		return model.High
	}
	if mediumRiskFound {
		return model.Medium
	}
	return model.Low
}

// GetSensitivePaths returns the list of sensitive paths being checked
func GetSensitivePaths() []string {
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		if u, err := user.Current(); err == nil {
			homeDir = u.HomeDir
		}
	}

	paths := make([]string, len(sensitivePaths))
	copy(paths, sensitivePaths)

	switch runtime.GOOS {
	case "darwin", "linux":
		if homeDir != "" {
			paths = append(paths,
				filepath.Join(homeDir, ".ssh"),
				filepath.Join(homeDir, ".gnupg"),
				filepath.Join(homeDir, ".config"),
				filepath.Join(homeDir, ".aws"),
			)
		}
	case "windows":
		if homeDir != "" {
			paths = append(paths,
				filepath.Join(homeDir, ".ssh"),
				filepath.Join(homeDir, "AppData", "Local"),
				filepath.Join(homeDir, "AppData", "Roaming"),
			)
		}
	}

	return paths
}
