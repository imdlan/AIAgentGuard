package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestE2EScanCommand(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping e2e test in short mode")
	}

	buildCmd := exec.Command("go", "build", "-o", "/tmp/agent-guard-e2e")
	buildCmd.Dir = filepath.Join("..", "..")
	if err := buildCmd.Run(); err != nil {
		t.Skipf("Cannot build binary for e2e test: %v", err)
	}
	defer os.Remove("/tmp/agent-guard-e2e")

	cmd := exec.Command("/tmp/agent-guard-e2e", "scan", "--json")
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		t.Fatalf("scan command failed: %v", err)
	}

	output := out.String()

	if !strings.Contains(output, "\"tool_name\"") {
		t.Error("Output does not contain expected JSON field 'tool_name'")
	}
	if !strings.Contains(output, "\"results\"") {
		t.Error("Output does not contain expected JSON field 'results'")
	}

	t.Logf("Scan command output: %s", output)
}

func TestE2EInitCommand(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping e2e test in short mode")
	}

	buildCmd := exec.Command("go", "build", "-o", "/tmp/agent-guard-e2e")
	buildCmd.Dir = filepath.Join("..", "..")
	if err := buildCmd.Run(); err != nil {
		t.Skipf("Cannot build binary for e2e test: %v", err)
	}
	defer os.Remove("/tmp/agent-guard-e2e")

	tmpDir, err := os.MkdirTemp("", "agentguard-e2e-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	cmd := exec.Command("/tmp/agent-guard-e2e", "init", "--path", filepath.Join(tmpDir, "config.yaml"))
	cmd.Dir = tmpDir

	err = cmd.Run()
	if err != nil {
		t.Fatalf("init command failed: %v", err)
	}

	configPath := filepath.Join(tmpDir, "config.yaml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Config file was not created")
	}

	t.Logf("Config file created successfully at: %s", configPath)
}

func TestE2ERunCommand(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping e2e test in short mode")
	}

	buildCmd := exec.Command("go", "build", "-o", "/tmp/agent-guard-e2e")
	buildCmd.Dir = filepath.Join("..", "..")
	if err := buildCmd.Run(); err != nil {
		t.Skipf("Cannot build binary for e2e test: %s", err)
	}
	defer os.Remove("/tmp/agent-guard-e2e")

	cmd := exec.Command("/tmp/agent-guard-e2e", "run", "echo", "test")
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		t.Logf("run command result (may be blocked): %v", err)
	}

	output := out.String()
	t.Logf("Run command output: %s", output)
}

func TestE2EReportCommand(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping e2e test in short mode")
	}

	buildCmd := exec.Command("go", "build", "-o", "/tmp/agent-guard-e2e")
	buildCmd.Dir = filepath.Join("..", "..")
	if err := buildCmd.Run(); err != nil {
		t.Skipf("Cannot build binary for e2e test: %v", err)
	}
	defer os.Remove("/tmp/agent-guard-e2e")

	cmd := exec.Command("/tmp/agent-guard-e2e", "report")
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		t.Fatalf("report command failed: %v", err)
	}

	output := out.String()

	expectedSections := []string{"Security Scan Report", "Permission Breakdown"}
	for _, section := range expectedSections {
		if !strings.Contains(output, section) {
			t.Errorf("Report missing expected section: %s", section)
		}
	}

	t.Logf("Report command output: %s", output)
}
