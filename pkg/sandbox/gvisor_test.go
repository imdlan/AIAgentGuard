package sandbox

import (
	"testing"
)

func TestNewGvisorSandbox(t *testing.T) {
	config := DefaultSandboxConfig()
	
	sandbox, err := NewGvisorSandbox(config)
	if err != nil {
		t.Logf("gVisor not available (expected on non-Linux): %v", err)
		return
	}
	
	if sandbox == nil {
		t.Fatal("Expected non-nil sandbox")
	}
	
	t.Logf("gVisor sandbox created successfully")
}

func TestGvisorSandboxAvailability(t *testing.T) {
	config := DefaultSandboxConfig()
	
	sandbox, _ := NewGvisorSandbox(config)
	if sandbox == nil {
		t.Skip("gVisor not available on this platform")
		return
	}
	
	available := sandbox.IsAvailable()
	t.Logf("gVisor available: %v", available)
	
	runtime := sandbox.GetRuntimeName()
	t.Logf("Runtime name: %s", runtime)
}

func BenchmarkNewGvisorSandbox(b *testing.B) {
	config := DefaultSandboxConfig()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = NewGvisorSandbox(config)
	}
}
