package scanner

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"

	"github.com/imdlan/AIAgentGuard/pkg/model"
)

// Sensitive paths that should not be accessible by AI agents
var sensitivePaths = []string{
	"/etc",
	"/usr/bin",
	"/usr/sbin",
}

// pathRiskReasons maps paths to their risk descriptions
var pathRiskReasons = map[string]string{
	"/etc":            "系统配置目录，可能被篡改",
	"/usr/bin":        "系统二进制文件，可能被替换",
	"/usr/sbin":       "系统管理命令，可能被滥用",
	".ssh":            "SSH 密钥，可能被窃取",
	".gnupg":          "GPG 密钥，可能被窃取",
	".config":         "应用配置，可能包含敏感信息",
	".aws":            "AWS 凭证，可能被窃取",
	"AppData/Local":   "本地应用数据，可能包含敏感信息",
	"AppData/Roaming": "漫游应用数据，可能包含凭证",
}

// ScanFilesystem checks if sensitive filesystem paths are accessible
func ScanFilesystem() model.RiskLevel {
	risk, _ := ScanFilesystemDetailed()
	return risk
}

// ScanFilesystemDetailed checks filesystem access and returns detailed information
func ScanFilesystemDetailed() (model.RiskLevel, []model.RiskDetail) {
	details := []model.RiskDetail{}

	homeDir := getHomeDir()
	if homeDir == "" {
		return model.Low, details
	}

	// Get all sensitive paths to check
	pathsToCheck := getPathsToCheck(homeDir)

	// Check each path
	affectedPaths := []model.PathDetail{}
	highRiskPaths := []string{}
	mediumRiskPaths := []string{}

	for _, path := range pathsToCheck {
		if info, err := os.Stat(path); err == nil {
			if info.IsDir() {
				pathDetail := checkPathPermissions(path, info)
				if pathDetail.IsWritable {
					highRiskPaths = append(highRiskPaths, path)
				} else {
					mediumRiskPaths = append(mediumRiskPaths, path)
				}
				affectedPaths = append(affectedPaths, pathDetail)
			}
		}
	}

	if len(affectedPaths) == 0 {
		return model.Low, details
	}

	// Determine overall risk
	risk := model.Low
	if len(highRiskPaths) > 0 {
		risk = model.High
	} else if len(mediumRiskPaths) > 0 {
		risk = model.Medium
	}

	// Build description
	description := fmt.Sprintf("发现 %d 个敏感路径可访问", len(affectedPaths))
	if len(highRiskPaths) > 0 {
		description = fmt.Sprintf("发现 %d 个敏感路径可写入（高风险）", len(highRiskPaths))
	}

	// Build risk detail
	detail := model.RiskDetail{
		Type:        risk,
		Category:    "filesystem",
		Description: description,
		Details: model.RiskSpecificInfo{
			AffectedPaths: affectedPaths,
		},
	}

	// Generate remediation
	detail.Remediation = generateFilesystemRemediation(highRiskPaths, mediumRiskPaths, affectedPaths)

	details = append(details, detail)
	return risk, details
}

// checkPathPermissions checks permissions and ownership of a path
func checkPathPermissions(path string, info os.FileInfo) model.PathDetail {
	permission := info.Mode().String()
	isWritable := false
	owner := "unknown"

	// Try to create a test file to check write permission
	testFile := filepath.Join(path, ".agentguard_test")
	if file, err := os.Create(testFile); err == nil {
		file.Close()
		os.Remove(testFile)
		isWritable = true
	}

	// Get owner information
	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		if u, err := user.LookupId(strconv.Itoa(int(stat.Uid))); err == nil {
			owner = u.Username
		}
	}

	return model.PathDetail{
		Path:       path,
		Permission: permission,
		IsWritable: isWritable,
		Owner:      owner,
		RiskReason: getPathRiskReason(path),
	}
}

// getPathRiskReason returns the risk reason for a path
func getPathRiskReason(path string) string {
	for key, reason := range pathRiskReasons {
		if strings.Contains(path, key) {
			return reason
		}
	}
	return "敏感路径，可能包含重要信息"
}

// generateFilesystemRemediation creates remediation suggestions
func generateFilesystemRemediation(highRiskPaths, mediumRiskPaths []string, allPaths []model.PathDetail) model.RemediationInfo {
	steps := []model.RemediationStep{}
	commands := []string{}

	// Add steps for high-risk paths
	for i, path := range highRiskPaths {
		steps = append(steps, model.RemediationStep{
			Step:        i + 1,
			Action:      fmt.Sprintf("限制 %s 的写权限", path),
			Command:     fmt.Sprintf("sudo chmod 755 %s", path),
			Explanation: "移除组和其他用户的写权限",
		})
		commands = append(commands, fmt.Sprintf("sudo chmod 755 %s", path))
	}

	// Add steps for user home directory paths
	for _, pathDetail := range allPaths {
		if strings.Contains(pathDetail.Path, ".ssh") ||
			strings.Contains(pathDetail.Path, ".gnupg") ||
			strings.Contains(pathDetail.Path, ".aws") {
			steps = append(steps, model.RemediationStep{
				Step:        len(steps) + 1,
				Action:      fmt.Sprintf("限制 %s 权限为仅所有者可访问", pathDetail.Path),
				Command:     fmt.Sprintf("chmod 700 %s", pathDetail.Path),
				Explanation: "设置为 700 权限，仅所有者可以读写执行",
			})
			commands = append(commands, fmt.Sprintf("chmod 700 %s", pathDetail.Path))
		}
	}

	priority := "HIGH"
	if len(highRiskPaths) == 0 && len(mediumRiskPaths) > 0 {
		priority = "MEDIUM"
	}

	return model.RemediationInfo{
		Summary:   "限制敏感路径的写权限，防止未授权修改",
		Steps:     steps,
		Commands:  commands,
		Priority:  priority,
		RiskAfter: model.Low,
	}
}

// getHomeDir returns the current user's home directory
func getHomeDir() string {
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		if u, err := user.Current(); err == nil {
			homeDir = u.HomeDir
		}
	}
	return homeDir
}

// getPathsToCheck returns all paths to check based on the OS
func getPathsToCheck(homeDir string) []string {
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

// GetSensitivePaths returns the list of sensitive paths being checked
func GetSensitivePaths() []string {
	homeDir := getHomeDir()
	return getPathsToCheck(homeDir)
}
