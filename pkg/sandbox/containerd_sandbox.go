// +build linux

package sandbox

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/imdlan/AIAgentGuard/internal/container"
)

type SandboxConfig struct {
	DisableNetwork bool
	ReadonlyRoot   bool
	AllowDirs      []string
	AllowedCaps    []string
	MaxMemoryMB    int64
	MaxCPUs        float64
}

type ContainerdSandbox struct {
	config SandboxConfig
	ctx    context.Context
}

func NewContainerdSandbox(config SandboxConfig) (*ContainerdSandbox, error) {
	if !container.HasContainerdSocket() {
		return nil, fmt.Errorf("containerd socket not found. Is containerd running?")
	}

	return &ContainerdSandbox{
		config: config,
		ctx:    context.Background(),
	}, nil
}

func (s *ContainerdSandbox) RunCommand(command string, args []string) error {
	if !container.HasContainerdSocket() {
		return s.runWithNamespaces(command, args)
	}

	return s.runWithContainerd(command, args)
}

func (s *ContainerdSandbox) runWithNamespaces(command string, args []string) error {
	cmd := exec.Command(command, args...)

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}

	if s.config.DisableNetwork {
		cmd.SysProcAttr.Cloneflags |= syscall.CLONE_NEWNET
	}

	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("sandbox command failed: %w", err)
	}

	return nil
}

func (s *ContainerdSandbox) runWithContainerd(command string, args []string) error {
	containerName := "agent-guard-sandbox"

	nerdctlCmd := []string{"run", "--rm"}

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

	nerdctlCmd = append(nerdctlCmd, "--security-opt=no-new-privileges")

	nerdctlCmd = append(nerdctlCmd, "--cap-drop=ALL")

	if len(s.config.AllowedCaps) > 0 {
		for _, cap := range s.config.AllowedCaps {
			nerdctlCmd = append(nerdctlCmd, "--cap-add="+cap)
		}
	} else {
		nerdctlCmd = append(nerdctlCmd, "--cap-add=NET_BIND_SERVICE")
	}

	cmdArgs := append([]string{command}, args...)
	nerdctlCmd = append(nerdctlCmd, cmdArgs...)

	cmd := exec.Command("nerdctl", nerdctlCmd...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("containerd sandbox command failed: %w", err)
	}

	return nil
}

func (s *ContainerdSandbox) IsAvailable() bool {
	return container.HasContainerdSocket()
}

func RunInSandbox(command string, args []string, config SandboxConfig) error {
	sandbox, err := NewContainerdSandbox(config)
	if err != nil {
		return err
	}

	return sandbox.RunCommand(command, args)
}

func DefaultSandboxConfig() SandboxConfig {
	return SandboxConfig{
		DisableNetwork: false,
		ReadonlyRoot:   false,
		AllowDirs:      []string{os.TempDir()},
		AllowedCaps:    []string{},
		MaxMemoryMB:    512,
		MaxCPUs:        1.0,
	}
}
