// +build linux

package sandbox

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/imdlan/AIAgentGuard/internal/container"
)

// GvisorSandbox provides enhanced isolation using gVisor
type GvisorSandbox struct {
	config SandboxConfig
	ctx    context.Context
}

// NewGvisorSandbox creates a new gVisor-based sandbox
func NewGvisorSandbox(config SandboxConfig) (*GvisorSandbox, error) {
	if !container.HasContainerdSocket() {
		return nil, fmt.Errorf("containerd socket not found. gVisor requires containerd")
	}

	return &GvisorSandbox{
		config: config,
		ctx:    context.Background(),
	}, nil
}

// RunCommand executes a command in gVisor sandbox
func (s *GvisorSandbox) RunCommand(command string, args []string) error {
	containerName := "agent-guard-gvisor-sandbox"

	// Use nerdctl with gVisor runtime
	nerdctlCmd := []string{"run", "--rm", "--runtime=runsc"}

	nerdctlCmd = append(nerdctlCmd, "--name", containerName)

	if s.config.DisableNetwork {
		nerdctlCmd = append(nerdctlCmd, "--network=none")
	}

	if s.config.ReadonlyRoot {
		nerdctlCmd = append(nerdctlCmd, "--read-only")
	}

	if len(s.config.AllowDirs) > 0 {
		for _, dir := range s.config.AllowDirs {
			if absPath, err := filepath.Abs(dir); err == nil {
				volume := fmt.Sprintf("%s:%s", absPath, absPath)
				nerdctlCmd = append(nerdctlCmd, "-v", volume)
			}
		}
	}

	if s.config.MaxMemoryMB > 0 {
		nerdctlCmd = append(nerdctlCmd, "--memory", fmt.Sprintf("%dm", s.config.MaxMemoryMB))
	}

	if s.config.MaxCPUs > 0 {
		nerdctlCmd = append(nerdctlCmd, "--cpus", fmt.Sprintf("%.2f", s.config.MaxCPUs))
	}

	// gVisor-specific security options
	nerdctlCmd = append(nerdctlCmd, "--security-opt=no-new-privileges")
	nerdctlCmd = append(nerdctlCmd, "--security-opt=seccomp=default")

	// Drop all capabilities initially
	nerdctlCmd = append(nerdctlCmd, "--cap-drop=ALL")

	// Add only necessary capabilities
	if len(s.config.AllowedCaps) > 0 {
		for _, cap := range s.config.AllowedCaps {
			nerdctlCmd = append(nerdctlCmd, "--cap-add="+cap)
		}
	}

	cmdArgs := append([]string{command}, args...)
	nerdctlCmd = append(nerdctlCmd, cmdArgs...)

	cmd := exec.Command("nerdctl", nerdctlCmd...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("gVisor sandbox command failed: %w", err)
	}

	return nil
}

// IsAvailable checks if gVisor runtime is available
func (s *GvisorSandbox) IsAvailable() bool {
	// Check if runsc (gVisor runtime) is available
	if _, err := exec.LookPath("runsc"); err != nil {
		return false
	}

	// Check if containerd is available
	return container.HasContainerdSocket()
}

// GetRuntimeName returns the runtime name
func (s *GvisorSandbox) GetRuntimeName() string {
	return "gVisor (runsc)"
}

// RunInGvisorSandbox executes a command in gVisor sandbox (convenience function)
func RunInGvisorSandbox(command string, args []string, config SandboxConfig) error {
	sandbox, err := NewGvisorSandbox(config)
	if err != nil {
		return err
	}

	return sandbox.RunCommand(command, args)
}
