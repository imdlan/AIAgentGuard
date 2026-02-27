package scanner

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"

	"github.com/imdlan/AIAgentGuard/pkg/model"
)

// SUIDFile represents a file with SUID/SGID bits set
type SUIDFile struct {
	Path    string
	Perms   string // File permissions (e.g., "rwsr-xr-x")
	Type    string // "suid", "sgid", or "both"
	IsKnown bool   // Whether this is a known system SUID file
	Risk    model.RiskLevel
}

// ScanSUIDFiles scans for SUID/SGID files and assesses risk
func ScanSUIDFiles() model.RiskLevel {
	return ScanSUIDFilesOptimized()
}

// findSUIDFiles finds all SUID/SGID files on the system


// findSUIDFiles finds all SUID/SGID files on the system
func findSUIDFiles() ([]SUIDFile, error) {
	var suidFiles []SUIDFile

	switch runtime.GOOS {
	case "darwin", "linux":
		return findUnixSUIDFiles()
	case "windows":
		// SUID/SGID is a Unix concept, not applicable on Windows
		return suidFiles, nil
	default:
		return suidFiles, fmt.Errorf("unsupported platform")
	}
}

// findUnixSUIDFiles finds SUID/SGID files on Unix-like systems
func findUnixSUIDFiles() ([]SUIDFile, error) {
	var suidFiles []SUIDFile

	// Use find command to locate SUID/SGID files
	// Search common system directories
	searchPaths := []string{"/bin", "/sbin", "/usr/bin", "/usr/sbin", "/usr/local/bin"}

	// Known SUID files that are expected and generally safe
	knownSUIDFiles := map[string]bool{
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
		"/usr/bin/btmp":        true,
		"/usr/sbin/uuidd":      true,
		"/usr/sbin/pppd":       true,
		"/usr/sbin/cron":       true,
	}

	// Find SUID files
	for _, searchPath := range searchPaths {
		if _, err := os.Stat(searchPath); os.IsNotExist(err) {
			continue
		}

		// Use find command to locate SUID files
		cmd := exec.Command("find", searchPath, "-type", "f", "-perm", "-4000", "-o", "-perm", "-2000")
		output, err := cmd.Output()
		if err != nil {
			continue
		}

		// Parse output
		files := strings.Split(strings.TrimSpace(string(output)), "\n")
		for _, file := range files {
			if file == "" {
				continue
			}

			// Get file info
			info, err := os.Lstat(file)
			if err != nil {
				continue
			}

			// Check permissions
			mode := info.Mode()
			perms := mode.String()

			var suidType string
			if mode&04000 != 0 && mode&02000 != 0 {
				suidType = "both"
			} else if mode&04000 != 0 {
				suidType = "suid"
			} else if mode&02000 != 0 {
				suidType = "sgid"
			}

			// Determine if this is a known SUID file
			isKnown := knownSUIDFiles[file]

			// Assess risk
			var risk model.RiskLevel
			if isKnown {
				risk = model.Low
			} else if filepath.Dir(file) == "/usr/local/bin" {
				// SUID files in /usr/local/bin are suspicious
				risk = model.High
			} else if strings.Contains(file, "home") || strings.Contains(file, "tmp") {
				// SUID files in user directories are very suspicious
				risk = model.High
			} else {
				risk = model.Medium
			}

			suidFiles = append(suidFiles, SUIDFile{
				Path:    file,
				Perms:   perms,
				Type:    suidType,
				IsKnown: isKnown,
				Risk:    risk,
			})
		}
	}

	// Also check home directory for unexpected SUID files
	homeDir := os.Getenv("HOME")
	if homeDir != "" {
		cmd := exec.Command("find", homeDir, "-type", "f", "-perm", "-4000", "-o", "-perm", "-2000")
		output, err := cmd.Output()
		if err == nil {
			files := strings.Split(strings.TrimSpace(string(output)), "\n")
			for _, file := range files {
				if file == "" {
					continue
				}

				info, err := os.Lstat(file)
				if err != nil {
					continue
				}

				mode := info.Mode()
				perms := mode.String()

				suidFiles = append(suidFiles, SUIDFile{
					Path:    file,
					Perms:   perms,
					Type:    "suid",
					IsKnown: false,
					Risk:    model.High, // SUID in home directory is always suspicious
				})
			}
		}
	}

	return suidFiles, nil
}

// CheckSpecificFile checks if a specific file has SUID/SGID bits
func CheckSpecificFile(filePath string) (bool, bool, error) {
	info, err := os.Lstat(filePath)
	if err != nil {
		return false, false, err
	}

	mode := info.Mode()
	hasSUID := mode&04000 != 0
	hasSGID := mode&02000 != 0

	return hasSUID, hasSGID, nil
}

// GetSUIDFindings returns detailed list of SUID/SGID files
func GetSUIDFindings() []SUIDFile {
	suidFiles, _ := findSUIDFiles()
	return suidFiles
}

// GetSuspiciousSUIDFiles returns only suspicious (unknown) SUID files
func GetSuspiciousSUIDFiles() []SUIDFile {
	suidFiles, _ := findSUIDFiles()
	var suspicious []SUIDFile

	for _, file := range suidFiles {
		if !file.IsKnown || file.Risk == model.High {
			suspicious = append(suspicious, file)
		}
	}

	return suspicious
}

// ValidateSUIDFile checks if a SUID file is valid and secure
func ValidateSUIDFile(filePath string) error {
	// Check if file exists
	info, err := os.Lstat(filePath)
	if err != nil {
		return fmt.Errorf("file not found: %w", err)
	}

	// Check if it's SUID
	mode := info.Mode()
	if mode&04000 == 0 && mode&02000 == 0 {
		return fmt.Errorf("file is not SUID or SGID")
	}

	// Check ownership (should be owned by root)
	stat, ok := info.Sys().(*syscall.Stat_t)
	if ok {
		if stat.Uid != 0 {
			return fmt.Errorf("SUID file not owned by root (UID=%d)", stat.Uid)
		}
	}

	// Check if file is writable by others
	if mode&0002 != 0 {
		return fmt.Errorf("SUID file is world-writable")
	}

	// Check if file is in a secure location
	securePaths := []string{"/bin", "/sbin", "/usr/bin", "/usr/sbin", "/usr/local/bin"}
	isSecure := false
	for _, path := range securePaths {
		if strings.HasPrefix(filePath, path) {
			isSecure = true
			break
		}
	}

	if !isSecure {
		return fmt.Errorf("SUID file is in an insecure location")
	}

	return nil
}

// GetSUIDStatistics returns statistics about SUID files
func GetSUIDStatistics() map[string]int {
	stats := map[string]int{
		"total_suid":     0,
		"known_suid":     0,
		"unknown_suid":   0,
		"high_risk":      0,
		"medium_risk":    0,
		"home_directory": 0,
	}

	suidFiles, err := findSUIDFiles()
	if err != nil {
		return stats
	}

	for _, file := range suidFiles {
		stats["total_suid"]++

		if file.IsKnown {
			stats["known_suid"]++
		} else {
			stats["unknown_suid"]++
		}

		if file.Risk == model.High {
			stats["high_risk"]++
		} else if file.Risk == model.Medium {
			stats["medium_risk"]++
		}

		if strings.Contains(file.Path, os.Getenv("HOME")) {
			stats["home_directory"]++
		}
	}

	return stats
}
