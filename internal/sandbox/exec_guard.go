package sandbox

import (
	"errors"
	"flag"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
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

	// Parse command and remove dangerous flags
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return "", nil
	}

	// Use actual command binary path for better validation
	cmdPath, err := exec.LookPath(parts[0])
	if err == nil {
		parts[0] = cmdPath
	}

	// Parse flags using flag package for the command
	filteredArgs := filterDangerousFlags(parts[1:])

	return strings.Join(append([]string{parts[0]}, filteredArgs...), " "), nil
}

// ParsedCommand represents a parsed command with its components
type ParsedCommand struct {
	Binary string
	Args   []string
	Flags  map[string]string
	RawCmd string
}

// ParseCommand parses a command string into its components
func ParseCommand(cmd string) (*ParsedCommand, error) {
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return nil, errors.New("empty command")
	}

	// Resolve binary path
	binaryPath, err := exec.LookPath(parts[0])
	if err != nil {
		// Keep original binary name if not found
		binaryPath = parts[0]
	}

	// Extract binary name from path
	binaryName := filepath.Base(binaryPath)

	return &ParsedCommand{
		Binary: binaryName,
		Args:   parts[1:],
		Flags:  make(map[string]string),
		RawCmd: cmd,
	}, nil
}

// HasDangerousFlags checks if a command has dangerous flag combinations
func HasDangerousFlags(parsed *ParsedCommand) bool {
	// Build flag set for parsing
	fs := flag.NewFlagSet(parsed.Binary, flag.ContinueOnError)

	// Common dangerous flags
	force := fs.Bool("force", false, "")
	recursive := fs.Bool("recursive", false, "")
	_ = fs.Bool("help", false, "")
	_ = fs.Bool("version", false, "")

	// Try to parse flags (will fail but will extract what it can)
	_ = fs.Parse(parsed.Args)

	// Check for dangerous flag combinations based on command type
	switch parsed.Binary {
	case "rm":
		// rm -rf is dangerous
		if *force || *recursive {
			// Also check for short flags
			for _, arg := range parsed.Args {
				if strings.Contains(arg, "f") || strings.Contains(arg, "r") {
					return true
				}
			}
		}
	case "chmod":
		// chmod 777 or chmod -R is suspicious
		for _, arg := range parsed.Args {
			if arg == "777" || strings.HasPrefix(arg, "-R") {
				return true
			}
		}
	case "dd":
		// dd is inherently dangerous
		return true
	}

	return false
}

// filterDangerousFlags removes or filters dangerous flags from command arguments
func filterDangerousFlags(args []string) []string {
	var filtered []string
	dangerousFlags := map[string]bool{
		"-f":          true,
		"-r":          true,
		"-rf":         true,
		"-fr":         true,
		"-R":          true,
		"--force":     true,
		"--recursive": true,
	}

	for _, arg := range args {
		// Skip dangerous standalone flags
		if dangerousFlags[arg] {
			continue
		}

		// Check for combined short flags (e.g., -rf)
		if strings.HasPrefix(arg, "-") && len(arg) > 1 {
			// Split combined flags and check each
			flagStr := strings.TrimLeft(arg, "-")
			safe := true
			for _, ch := range flagStr {
				if dangerousFlags["-"+string(ch)] {
					safe = false
					break
				}
			}
			if !safe {
				continue
			}
		}

		filtered = append(filtered, arg)
	}

	return filtered
}

// GetCommandPath resolves the full path of a command binary
func GetCommandPath(cmd string) (string, error) {
	return exec.LookPath(cmd)
}

// IsKnownCommand checks if a command is a known system binary
func IsKnownCommand(cmd string) bool {
	path, err := exec.LookPath(cmd)
	if err != nil {
		return false
	}

	// Check if it's in system directories
	systemDirs := []string{"/bin", "/usr/bin", "/sbin", "/usr/sbin", "/usr/local/bin"}
	for _, dir := range systemDirs {
		if strings.HasPrefix(path, dir) {
			return true
		}
	}

	return false
}

// ValidateCommandPath checks if command path is within allowed directories
func ValidateCommandPath(cmd string, allowedDirs []string) error {
	path, err := exec.LookPath(cmd)
	if err != nil {
		return errors.New("command not found: " + cmd)
	}

	// Resolve symlinks
	realPath, err := filepath.EvalSymlinks(path)
	if err != nil {
		realPath = path
	}

	// Check if command is in allowed directories
	for _, allowedDir := range allowedDirs {
		if strings.HasPrefix(realPath, allowedDir) {
			return nil
		}
	}

	return errors.New("command path not in allowed directories: " + realPath)
}

// GetRealBinaryPath returns the real path of a binary (following symlinks)
func GetRealBinaryPath(cmd string) (string, error) {
	path, err := exec.LookPath(cmd)
	if err != nil {
		return "", err
	}

	// Follow symlinks
	realPath, err := filepath.EvalSymlinks(path)
	if err != nil {
		return path, nil
	}

	return realPath, nil
}

// DetectSignalHandling checks if command is trying to handle signals in a dangerous way
func DetectSignalHandling(cmd string) bool {
	parts := strings.Fields(cmd)
	for i, part := range parts {
		// Check for kill, killall commands
		if part == "kill" || part == "killall" {
			if i+1 < len(parts) {
				signal := parts[i+1]
				// Block dangerous signals
				if signal == "-9" || signal == "-KILL" {
					return true
				}
			}
		}
	}
	return false
}

// DetectShellInjection attempts to detect command injection patterns
func DetectShellInjection(cmd string) bool {
	injectionPatterns := []string{
		"; ",
		"&&",
		"||",
		"|",
		"`",
		"$(",
		">",
		"<",
		"\n",
		"\r",
	}

	cmdLower := strings.ToLower(cmd)
	for _, pattern := range injectionPatterns {
		if strings.Contains(cmdLower, pattern) {
			// Check if it's a legitimate use (e.g., grep "foo" file > output)
			// For now, be conservative and flag all
			return true
		}
	}
	return false
}

// GetCommandType returns the type of command (system, user, builtin, etc.)
func GetCommandType(cmd string) string {
	path, err := exec.LookPath(cmd)
	if err != nil {
		return "unknown"
	}

	// Check if it's a shell builtin
	builtins := []string{"echo", "cd", "pwd", "export", "source", ".", "test", "[", "true", "false"}
	for _, b := range builtins {
		if cmd == b {
			return "builtin"
		}
	}

	// Check system directories
	if strings.HasPrefix(path, "/bin") || strings.HasPrefix(path, "/usr/bin") ||
		strings.HasPrefix(path, "/sbin") || strings.HasPrefix(path, "/usr/sbin") {
		return "system"
	}

	// Check user directories
	if strings.Contains(path, "/home") || strings.Contains(path, "/usr/local") {
		return "user"
	}

	return "external"
}

// IsPathTraversal checks if command attempts path traversal
func IsPathTraversal(cmd string) bool {
	traversalPatterns := []string{
		"../",
		"..\\",
		"~/",
		"/etc/",
		"/proc/",
		"/sys/",
	}

	for _, pattern := range traversalPatterns {
		if strings.Contains(cmd, pattern) {
			return true
		}
	}
	return false
}

// CheckPrivelegeEscalation attempts to detect privilege escalation attempts
func CheckPrivilegeEscalation(cmd string) bool {
	escalationPatterns := []string{
		"sudo",
		"su",
		"doas",
		"pkexec",
	}

	cmdLower := strings.ToLower(cmd)
	for _, pattern := range escalationPatterns {
		if strings.Contains(cmdLower, pattern) {
			return true
		}
	}
	return false
}

// GetProcessInfo returns information about a command (placeholder for future use)
func GetProcessInfo(pid int) (*syscall.WaitStatus, error) {
	var status syscall.WaitStatus
	return &status, nil
}
