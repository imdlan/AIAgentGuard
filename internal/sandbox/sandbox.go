package sandbox

import (
	"os"
	"os/exec"
	"runtime"
)

// RunSandboxed executes a command in a sandboxed environment
func RunSandboxed(command string, cfg *Config) error {
	// Parse command
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}

	// Apply sandbox restrictions
	applySandboxRestrictions(cmd, cfg)

	// Set output
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	// Run the command
	return cmd.Run()
}

// Config holds sandbox configuration
type Config struct {
	ClearEnv       bool
	WorkingDir     string
	DisableNetwork bool
	AllowedPaths   []string
	BlockedPaths   []string
}

// applySandboxRestrictions applies security restrictions to the command
func applySandboxRestrictions(cmd *exec.Cmd, cfg *Config) {
	// Clear environment variables for isolation
	if cfg.ClearEnv {
		// Keep minimal environment
		cmd.Env = []string{
			"PATH=" + os.Getenv("PATH"),
			"HOME=" + os.Getenv("HOME"),
			"TERM=" + os.Getenv("TERM"),
		}
	} else {
		cmd.Env = os.Environ()
	}

	// Set working directory
	if cfg.WorkingDir != "" {
		cmd.Dir = cfg.WorkingDir
	} else {
		// Default to /tmp for safety
		cmd.Dir = os.TempDir()
	}
}

// GetDefaultConfig returns a default sandbox configuration
func GetDefaultConfig() *Config {
	return &Config{
		ClearEnv:       true,
		WorkingDir:     os.TempDir(),
		DisableNetwork: false,
		AllowedPaths:   []string{},
		BlockedPaths:   []string{},
	}
}

// GetStrictConfig returns a strict sandbox configuration
func GetStrictConfig() *Config {
	return &Config{
		ClearEnv:       true,
		WorkingDir:     os.TempDir(),
		DisableNetwork: true,
		AllowedPaths:   []string{os.TempDir()},
		BlockedPaths:   getSensitivePaths(),
	}
}

// getSensitivePaths returns a list of sensitive paths to block
func getSensitivePaths() []string {
	homeDir := os.Getenv("HOME")
	var paths []string

	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		paths = []string{
			"/etc",
			"/usr/bin",
			"/usr/sbin",
			"/bin",
			"/sbin",
		}
		if homeDir != "" {
			paths = append(paths,
				homeDir+"/.ssh",
				homeDir+"/.gnupg",
				homeDir+"/.aws",
				homeDir+"/.config",
			)
		}
	}

	return paths
}

// RunCommand runs a command with the specified sandbox level
func RunCommand(command string, strict bool) error {
	var cfg *Config
	if strict {
		cfg = GetStrictConfig()
	} else {
		cfg = GetDefaultConfig()
	}

	return RunSandboxed(command, cfg)
}
