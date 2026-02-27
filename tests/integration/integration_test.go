package main

import (
	"os"
	"testing"

	"github.com/imdlan/AIAgentGuard/internal/scanner"
)

func TestAllScanners(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("Filesystem", func(t *testing.T) {
		risk := scanner.ScanFilesystem()
		t.Logf("Filesystem scan risk: %s", risk)
	})

	t.Run("Shell", func(t *testing.T) {
		risk := scanner.ScanShell()
		t.Logf("Shell scan risk: %s", risk)
	})

	t.Run("Network", func(t *testing.T) {
		risk := scanner.ScanNetwork()
		t.Logf("Network scan risk: %s", risk)
	})

	t.Run("Secrets", func(t *testing.T) {
		risk := scanner.ScanSecrets()
		t.Logf("Secrets scan risk: %s", risk)
	})

	t.Run("FileContent", func(t *testing.T) {
		risk := scanner.ScanFileContents()
		t.Logf("File content scan risk: %s", risk)
	})

	t.Run("Dependencies", func(t *testing.T) {
		risk := scanner.ScanDependencies()
		t.Logf("Dependencies scan risk: %s", risk)
	})

	t.Run("Processes", func(t *testing.T) {
		risk := scanner.ScanProcesses()
		t.Logf("Process scan risk: %s", risk)
	})

	t.Run("SUID", func(t *testing.T) {
		if os.Getenv("CI") != "" {
			t.Skip("Skipping SUID scan in CI environment (too slow)")
		}
		risk := scanner.ScanSUIDFiles()
		t.Logf("SUID scan risk: %s", risk)
	})
}

func TestRunAllScans(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	result := scanner.RunAllScans()
	t.Logf("All scans result: Filesystem=%s, Shell=%s, Network=%s, Secrets=%s, FileContent=%s, Dependencies=%s",
		result.Filesystem,
		result.Shell,
		result.Network,
		result.Secrets,
		result.FileContent,
		result.Dependencies,
	)
}
