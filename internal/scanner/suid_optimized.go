package scanner

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/imdlan/AIAgentGuard/pkg/model"
)


const (
	maxScanDepth     = 5         // Maximum directory depth to scan
	maxScanDuration  = 10 * time.Second // Maximum time to spend scanning
	maxSUIDFiles     = 100       // Maximum number of SUID files to collect
)

// ScanSUIDFilesOptimized is an optimized version of SUID scanning
func ScanSUIDFilesOptimized() model.RiskLevel {
	suidFiles, err := findSUIDFilesOptimized()
	if err != nil {
		return model.Low
	}

	// Count high-risk SUID files
	highRiskCount := 0
	mediumRiskCount := 0

	for _, file := range suidFiles {
		if file.Risk == model.High {
			highRiskCount++
		} else if file.Risk == model.Medium {
			mediumRiskCount++
		}
	}

	// Determine overall risk level
	if highRiskCount > 0 {
		return model.High
	}
	if mediumRiskCount > 2 {
		return model.Medium
	}
	if mediumRiskCount > 0 {
		return model.Low
	}
	return model.Low
}

// findSUIDFilesOptimized finds SUID/SGID files with performance optimizations
func findSUIDFilesOptimized() ([]SUIDFile, error) {
	var suidFiles []SUIDFile

	if runtime.GOOS == "windows" {
		return suidFiles, nil
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), maxScanDuration)
	defer cancel()

	// Limit search paths and depth
	searchPaths := getOptimizedSearchPaths()

	// Scan system directories first (fast path)
	for _, searchPath := range searchPaths {
		select {
		case <-ctx.Done():
			// Timeout reached
			return suidFiles, nil
		default:
		}

		files, err := scanDirectoryFast(ctx, searchPath, true)
		if err != nil {
			continue
		}
		suidFiles = append(suidFiles, files...)

		// Early exit if we found enough files
		if len(suidFiles) >= maxSUIDFiles {
			return suidFiles[:maxSUIDFiles], nil
		}
	}

	// Scan home directory with depth limit
	homeDir := os.Getenv("HOME")
	if homeDir != "" && len(suidFiles) < maxSUIDFiles {
		select {
		case <-ctx.Done():
			return suidFiles, nil
		default:
		}

		files, err := scanDirectoryFast(ctx, homeDir, false)
		if err == nil {
			suidFiles = append(suidFiles, files...)
		}
	}

	if len(suidFiles) > maxSUIDFiles {
		suidFiles = suidFiles[:maxSUIDFiles]
	}

	return suidFiles, nil
}

// getOptimizedSearchPaths returns limited search paths for faster scanning
func getOptimizedSearchPaths() []string {
	switch runtime.GOOS {
	case "darwin":
		return []string{"/bin", "/sbin", "/usr/bin", "/usr/sbin", "/usr/local/bin"}
	case "linux":
		return []string{"/bin", "/sbin", "/usr/bin", "/usr/sbin"}
	default:
		return []string{"/bin", "/usr/bin"}
	}
}

// scanDirectoryFast scans a directory using optimized find command
func scanDirectoryFast(ctx context.Context, dir string, systemDir bool) ([]SUIDFile, error) {
	var suidFiles []SUIDFile

	// Build find command with optimizations
	args := []string{dir, "-type", "f", "(", "-perm", "-4000", "-o", "-perm", "-2000", ")"}

	// Add depth limit for non-system directories
	if !systemDir {
		args = append(args, "-maxdepth", fmt.Sprintf("%d", maxScanDepth))
	}

	// Use -xdev to avoid scanning other filesystems
	if systemDir {
		args = append(args, "-xdev")
	}

	cmd := exec.CommandContext(ctx, "find", args...)
	output, err := cmd.Output()
	if err != nil {
		return suidFiles, err
	}

	files := strings.Split(strings.TrimSpace(string(output)), "\n")
	knownSUIDFiles := getKnownSUIDFileMap()

	for _, file := range files {
		if file == "" {
			continue
		}

		// Check file info
		info, err := os.Lstat(file)
		if err != nil {
			continue
		}

		mode := info.Mode()
		perms := mode.String()

		// Determine if this is a known SUID file
		isKnown := knownSUIDFiles[file]

		// Calculate risk
		risk := calculateSUIDRisk(file, isKnown, systemDir)

		suidFiles = append(suidFiles, SUIDFile{
			Path:    file,
			Perms:   perms,
			Type:    getSUIDType(mode),
			IsKnown: isKnown,
			Risk:    risk,
		})
	}

	return suidFiles, nil
}

// getSUIDType returns the type of SUID/SGID
func getSUIDType(mode os.FileMode) string {
	hasSUID := mode&04000 != 0
	hasSGID := mode&02000 != 0

	if hasSUID && hasSGID {
		return "both"
	} else if hasSUID {
		return "suid"
	} else if hasSGID {
		return "sgid"
	}
	return "unknown"
}

// calculateSUIDRisk calculates the risk level of a SUID file
func calculateSUIDRisk(path string, isKnown bool, systemDir bool) model.RiskLevel {
	// High risk: SUID in home directory or non-system location
	if !systemDir && !strings.HasPrefix(path, "/usr/") && !strings.HasPrefix(path, "/bin/") && !strings.HasPrefix(path, "/sbin/") {
		return model.High
	}

	// Medium risk: Unknown SUID file in system directory
	if !isKnown {
		return model.Medium
	}

	// Low risk: Known SUID file
	return model.Low
}

// getKnownSUIDFileMap returns a map of known safe SUID files
func getKnownSUIDFileMap() map[string]bool {
	return map[string]bool{
		"/bin/ping":            true,
		"/bin/ping6":           true,
		"/bin/su":              true,
		"/bin/mount":           true,
		"/bin/umount":          true,
		"/usr/bin/passwd":      true,
		"/usr/bin/sudo":        true,
		"/usr/bin/newgrp":      true,
		"/usr/bin/chsh":        true,
		"/usr/bin/chfn":        true,
		"/usr/bin/gpasswd":     true,
		"/usr/bin/wall":        true,
		"/usr/sbin/passwd":     true,
		"/usr/sbin/su":         true,
		"/usr/sbin/visudo":     true,
		"/usr/sbin/traceroute": true,
		"/usr/bin/at":          true,
		"/usr/sbin/uuidd":      true,
		"/usr/sbin/pppd":       true,
		"/usr/sbin/arping":     true,
		"/usr/sbin/cron":       true,
		"/usr/bin/crontab":     true,
		"/usr/bin/ssh-agent":   true,
		"/usr/bin/generate-ssh-keys": true,
		"/sbin/ping":           true,
		"/sbin/ping6":          true,
		"/usr/local/bin/sudo":  true,
	}
}
