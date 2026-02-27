// +build !linux

package sandbox

import (
	"context"
	"fmt"
	"os"

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
	return fmt.Errorf("sandbox isolation is only supported on Linux")
}

func (s *ContainerdSandbox) IsAvailable() bool {
	return false
}

func RunInSandbox(command string, args []string, config SandboxConfig) error {
	return fmt.Errorf("sandbox isolation is only supported on Linux")
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
