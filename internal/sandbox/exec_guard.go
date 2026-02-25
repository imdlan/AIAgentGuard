package sandbox

import (
	"errors"
	"strings"
)

// Dangerous command patterns to block
var dangerousPatterns = []string{
	"rm -rf",
	"rm -fr",
	"rm -r /",
	"mkfs",
	"dd if=",
	":(){:|:&};:", // fork bomb
	"chmod -R 777",
	"chmod 777 /",
	"curl | sh",
	"wget | bash",
	"curl | bash",
	"wget | sh",
	"mv /dev/",
	"> /dev/",
}

// Dangerous commands that should always be blocked regardless of arguments
var dangerousCommands = []string{
	"rm ",
	"dd ",
	"mkfs.",
	"reboot",
	"shutdown",
	"halt",
	"poweroff",
}

// GuardCommand checks if a command is dangerous and should be blocked
func GuardCommand(cmd string) error {
	cmdLower := strings.ToLower(cmd)

	// Check for dangerous patterns
	for _, pattern := range dangerousPatterns {
		if strings.Contains(cmdLower, pattern) {
			return errors.New("blocked: dangerous command pattern detected: " + pattern)
		}
	}

	// Check for dangerous commands
	for _, dangerousCmd := range dangerousCommands {
		if strings.HasPrefix(strings.TrimSpace(cmdLower), dangerousCmd) {
			return errors.New("blocked: dangerous command detected: " + dangerousCmd)
		}
	}

	return nil
}

// IsCommandSafe checks if a command is safe to execute
func IsCommandSafe(cmd string) bool {
	return GuardCommand(cmd) == nil
}

// SanitizeCommand attempts to make a command safer by removing dangerous flags
// Returns an error if the command is too dangerous to sanitize
func SanitizeCommand(cmd string) (string, error) {
	// If the command is outright dangerous, don't try to sanitize
	if err := GuardCommand(cmd); err != nil {
		return "", err
	}

	// Remove dangerous flags
	cmd = strings.ReplaceAll(cmd, "-f", "")
	cmd = strings.ReplaceAll(cmd, "-rf", "")
	cmd = strings.ReplaceAll(cmd, "-r", "")

	return cmd, nil
}
