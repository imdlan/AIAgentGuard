// +build !linux

package sandbox

import (
	"context"
	"fmt"
)

// GvisorSandbox provides enhanced isolation using gVisor
type GvisorSandbox struct {
	config SandboxConfig
	ctx    context.Context
}

// NewGvisorSandbox creates a new gVisor-based sandbox
func NewGvisorSandbox(config SandboxConfig) (*GvisorSandbox, error) {
	return nil, fmt.Errorf("gVisor is only supported on Linux")
}

// RunCommand executes a command in gVisor sandbox
func (s *GvisorSandbox) RunCommand(command string, args []string) error {
	return fmt.Errorf("gVisor sandbox is only supported on Linux")
}

// IsAvailable checks if gVisor runtime is available
func (s *GvisorSandbox) IsAvailable() bool {
	return false
}

// GetRuntimeName returns the runtime name
func (s *GvisorSandbox) GetRuntimeName() string {
	return "gVisor (not available on this platform)"
}

// RunInGvisorSandbox executes a command in gVisor sandbox
func RunInGvisorSandbox(command string, args []string, config SandboxConfig) error {
	return fmt.Errorf("gVisor sandbox is only supported on Linux")
}
