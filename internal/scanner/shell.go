package scanner

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/imdlan/AIAgentGuard/pkg/model"
)

// SudoInfo contains detailed information about sudo access
type SudoInfo struct {
	HasAccess bool
	Source    string
	Rules     []string
}

// ScanShell checks if shell command execution is possible
func ScanShell() model.RiskLevel {
	risk, _ := ScanShellDetailed()
	return risk
}

// ScanShellDetailed checks shell access and returns detailed information
func ScanShellDetailed() (model.RiskLevel, []model.RiskDetail) {
	details := []model.RiskDetail{}

	// Detect available shells
	availableShells := detectAvailableShells()
	if len(availableShells) == 0 {
		return model.Low, details
	}

	// Check sudo access
	sudoInfo := checkSudoAccessDetailed()

	// Determine risk level
	risk := model.High
	if sudoInfo.HasAccess {
		risk = model.Critical
	}

	// Build description
	description := fmt.Sprintf("检测到 %d 个可用的 Shell", len(availableShells))
	if sudoInfo.HasAccess {
		description = fmt.Sprintf("检测到 %d 个可用的 Shell，且有无密码 sudo 权限", len(availableShells))
	}

	// Build risk detail
	detail := model.RiskDetail{
		Type:        risk,
		Category:    "shell",
		Description: description,
		Details: model.RiskSpecificInfo{
			ShellAvailable: strings.Join(availableShells, ", "),
			HasSudoAccess:  sudoInfo.HasAccess,
			SudoSource:     sudoInfo.Source,
			SudoRules:      sudoInfo.Rules,
		},
	}

	// Add remediation suggestions
	if sudoInfo.HasAccess {
		detail.Remediation = model.RemediationInfo{
			Summary: "移除无密码 sudo 权限，限制 Shell 访问",
			Steps: []model.RemediationStep{
				{
					Step:        1,
					Action:      "检查 sudoers 配置",
					Command:     "sudo visudo -c",
					Explanation: "验证 sudoers 文件语法",
				},
				{
					Step:        2,
					Action:      "查看当前 sudo 配置",
					Command:     "sudo -l",
					Explanation: "列出当前用户的 sudo 权限",
				},
				{
					Step:        3,
					Action:      "编辑 sudoers 文件",
					Command:     "sudo visudo",
					Explanation: "删除或注释包含 NOPASSWD 的配置行",
				},
			},
			Commands: []string{
				"sudo visudo -c",
				"sudo -l",
				"sudo visudo",
			},
			Priority:  "HIGH",
			RiskAfter: model.High,
		}
	} else {
		detail.Remediation = model.RemediationInfo{
			Summary: "限制 Shell 访问权限",
			Steps: []model.RemediationStep{
				{
					Step:        1,
					Action:      "检查不必要的 Shell",
					Command:     "cat /etc/shells",
					Explanation: "查看系统中可用的 Shell",
				},
				{
					Step:        2,
					Action:      "限制用户 Shell",
					Command:     "sudo usermod -s /bin/false username",
					Explanation: "为特定用户设置限制性 Shell",
				},
			},
			Commands: []string{
				"cat /etc/shells",
			},
			Priority:  "MEDIUM",
			RiskAfter: model.Low,
		}
	}

	details = append(details, detail)
	return risk, details
}

// detectAvailableShells finds which shells are available on the system
func detectAvailableShells() []string {
	shells := []string{}
	shellPaths := []string{"/bin/sh", "/bin/bash", "/bin/zsh", "/bin/fish", "/bin/dash"}

	for _, shell := range shellPaths {
		if _, err := exec.LookPath(shell); err == nil {
			shells = append(shells, shell)
		}
	}

	return shells
}

// checkSudoAccess checks if the current user has sudo/admin privileges
func checkSudoAccess() bool {
	info := checkSudoAccessDetailed()
	return info.HasAccess
}

// checkSudoAccessDetailed returns detailed sudo access information
func checkSudoAccessDetailed() SudoInfo {
	info := SudoInfo{}

	switch runtime.GOOS {
	case "darwin", "linux":
		// Try to run whoami with sudo (with -n to non-interactive)
		cmd := exec.Command("sudo", "-n", "whoami")
		if err := cmd.Run(); err == nil {
			info.HasAccess = true

			// Try to get sudo rules
			cmd = exec.Command("sudo", "-n", "-l")
			if output, err := cmd.CombinedOutput(); err == nil {
				lines := strings.Split(string(output), "\n")
				for _, line := range lines {
					line = strings.TrimSpace(line)
					if line != "" && !strings.HasPrefix(line, "Matching") && !strings.HasPrefix(line, "User") {
						info.Rules = append(info.Rules, line)
						if strings.Contains(line, "NOPASSWD") {
							info.Source = "发现无密码 sudo 配置"
						}
					}
				}
			}

			if info.Source == "" {
				info.Source = "当前用户可以执行 sudo 命令（可能需要密码）"
			}
		}

	case "windows":
		// On Windows, check if running as administrator
		cmd := exec.Command("net", "session")
		if err := cmd.Run(); err == nil {
			info.HasAccess = true
			info.Source = "当前进程以管理员权限运行"
		}
	}

	return info
}
