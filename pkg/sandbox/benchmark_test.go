package sandbox

import (
	"testing"
)

// BenchmarkNewContainerdSandbox benchmarks sandbox creation
func BenchmarkNewContainerdSandbox(b *testing.B) {
	config := DefaultSandboxConfig()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = NewContainerdSandbox(config)
	}
}

// BenchmarkDefaultSandboxConfig benchmarks default config creation
func BenchmarkDefaultSandboxConfig(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DefaultSandboxConfig()
	}
}
