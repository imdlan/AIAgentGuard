package main

import (
	"testing"

	"github.com/imdlan/AIAgentGuard/internal/scanner"
	"github.com/imdlan/AIAgentGuard/internal/container"
	"github.com/imdlan/AIAgentGuard/pkg/sandbox"
)

func TestComprehensiveScan(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping comprehensive test in short mode")
	}

	t.Run("AllScannersWithContainerInfo", func(t *testing.T) {
		// Get container info
		containerInfo := container.DetectContainer()
		t.Logf("Container Environment: %s (isContainer=%v)", 
			containerInfo.Runtime, containerInfo.IsContainer)

		// Run all scans
		result := scanner.RunAllScans()
		t.Logf("Scan Results: Filesystem=%s, Shell=%s, Network=%s, Secrets=%s, FileContent=%s, Dependencies=%s",
			result.Filesystem,
			result.Shell,
			result.Network,
			result.Secrets,
			result.FileContent,
			result.Dependencies,
		)
	})

	t.Run("SandboxAvailability", func(t *testing.T) {
		config := sandbox.DefaultSandboxConfig()
		t.Logf("Sandbox Config: DisableNetwork=%v, ReadonlyRoot=%v, MaxMemory=%dMB",
			config.DisableNetwork, config.ReadonlyRoot, config.MaxMemoryMB)

		sb, err := sandbox.NewContainerdSandbox(config)
		if err != nil {
			t.Logf("Sandbox not available: %v", err)
			return
		}
		t.Logf("Sandbox available: %v", sb.IsAvailable())
	})
}
