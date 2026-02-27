package container

import (
	"testing"
)

// BenchmarkDetectContainer benchmarks container detection
func BenchmarkDetectContainer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DetectContainer()
	}
}

// BenchmarkIsRunningInContainer benchmarks container status check
func BenchmarkIsRunningInContainer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsRunningInContainer()
	}
}

// BenchmarkHasDockerSocket benchmarks Docker socket detection
func BenchmarkHasDockerSocket(b *testing.B) {
	for i := 0; i < b.N; i++ {
		HasDockerSocket()
	}
}

// BenchmarkHasContainerdSocket benchmarks containerd socket detection
func BenchmarkHasContainerdSocket(b *testing.B) {
	for i := 0; i < b.N; i++ {
		HasContainerdSocket()
	}
}
