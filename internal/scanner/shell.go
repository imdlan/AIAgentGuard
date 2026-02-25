package scanner

import (
	"os/exec"
	"runtime"

	"github.com/imdlan/AIAgentGuard/pkg/model"
)

// ScanShell checks if shell command execution is possible
func ScanShell() model.RiskLevel {
	// Try to execute a simple shell command
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "echo test")
	default:
		cmd = exec.Command("sh", "-c", "echo test")
	}

	if err := cmd.Run(); err != nil {
		return model.Low
	}

	// Shell execution is possible - check if we can run more dangerous commands
	// Check for sudo or admin privileges
	if checkSudoAccess() {
		return model.Critical
	}

	return model.High
}

// checkSudoAccess checks if the current user has sudo/admin privileges
func checkSudoAccess() bool {
	switch runtime.GOOS {
	case "darwin", "linux":
		// Try to run whoami with sudo (with -n to non-interactive)
		cmd := exec.Command("sudo", "-n", "whoami")
		// We expect this to fail in most cases, but if it succeeds, we have sudo access
		if err := cmd.Run(); err == nil {
			return true
		}
	case "windows":
		// On Windows, check if running as administrator
		cmd := exec.Command("net", "session")
		if err := cmd.Run(); err == nil {
			return true
		}
	}
	return false
}
