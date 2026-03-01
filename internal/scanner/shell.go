package scanner

import (
	"os/exec"
	"runtime"
	"strings"

	"github.com/imdlan/AIAgentGuard/internal/i18n"
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
	description := i18n.T("scan.shell.foundShells", len(availableShells))
	if sudoInfo.HasAccess {
		description = i18n.T("scan.shell.foundShellsWithSudo", len(availableShells))
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
			Summary: i18n.T("scan.shell.removeSudo"),
			Steps: []model.RemediationStep{
				{
					Step:        1,
					Action:      i18n.T("scan.shell.checkSudoers"),
					Command:     "sudo visudo -c",
					Explanation: i18n.T("scan.shell.verifySudoersSyntax"),
				},
				{
					Step:        2,
					Action:      i18n.T("scan.shell.viewSudoConfig"),
					Command:     "sudo -l",
					Explanation: i18n.T("scan.shell.listSudoPrivileges"),
				},
				{
					Step:        3,
					Action:      i18n.T("scan.shell.editSudoers"),
					Command:     "sudo visudo",
					Explanation: i18n.T("scan.shell.removeNopasswd"),
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
			Summary: i18n.T("scan.shell.limitShellAccess"),
			Steps: []model.RemediationStep{
				{
					Step:        1,
					Action:      i18n.T("scan.shell.checkUnnecessaryShells"),
					Command:     "cat /etc/shells",
					Explanation: i18n.T("scan.shell.viewAvailableShells"),
				},
				{
					Step:        2,
					Action:      i18n.T("scan.shell.restrictUserShell"),
					Command:     "sudo usermod -s /bin/false username",
					Explanation: i18n.T("scan.shell.setRestrictedShell"),
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
							info.Source = i18n.T("scan.sudo.passwordlessFound")
						}
					}
				}
			}

			if info.Source == "" {
					info.Source = i18n.T("scan.sudo.canRunSudo")
			}
		}

	case "windows":
		// On Windows, check if running as administrator
		cmd := exec.Command("net", "session")
		if err := cmd.Run(); err == nil {
			info.HasAccess = true
				info.Source = i18n.T("scan.sudo.runningAsAdmin")
		}
	}

	return info
}
