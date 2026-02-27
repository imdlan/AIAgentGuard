package main

import (
	"testing"

	"github.com/imdlan/AIAgentGuard/internal/container"
)

func TestContainerDetection(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("DetectContainer", func(t *testing.T) {
		info := container.DetectContainer()
		t.Logf("Container info: Runtime=%s, IsContainer=%v", info.Runtime, info.IsContainer)
	})

	t.Run("IsRunningInContainer", func(t *testing.T) {
		isContainer := container.IsRunningInContainer()
		t.Logf("Running in container: %v", isContainer)
	})

	t.Run("HasDockerSocket", func(t *testing.T) {
		hasDocker := container.HasDockerSocket()
		t.Logf("Has Docker socket: %v", hasDocker)
	})

	t.Run("HasContainerdSocket", func(t *testing.T) {
		hasContainerd := container.HasContainerdSocket()
		t.Logf("Has containerd socket: %v", hasContainerd)
	})
}
