package audit

import (
	"os"
	"testing"
)

// TestAuditLoggerInitialization tests logger initialization
func TestAuditLoggerInitialization(t *testing.T) {
	// Reset for clean test state
	ResetForTesting()
	
	tmpDir, err := os.MkdirTemp("", "audit-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	logPath := tmpDir + "/audit.log"

	err = Init(logPath, true, InfoLevel)
	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}
	defer Close()

	if !IsEnabled() {
		t.Error("Logger should be enabled after Init()")
	}

	if got := GetLogFilePath(); got != logPath {
		t.Errorf("GetLogFilePath() = %v, want %v", got, logPath)
	}
}
// TestLogEvent tests basic event logging
func TestLogEvent(t *testing.T) {
	// Reset for clean test state
	ResetForTesting()
	
	tmpDir, err := os.MkdirTemp("", "audit-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	logPath := tmpDir + "/audit.log"
	Init(logPath, true, InfoLevel)
	defer Close()

	event := AuditEvent{
		Level:     InfoLevel,
		EventType: SystemEvent,
		Message:   "Test event",
	}

	err = LogEvent(event)
	if err != nil {
		t.Errorf("LogEvent() failed: %v", err)
	}

	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	if len(data) == 0 {
		t.Error("Log file is empty, expected content")
	}
}

// TestConvenienceFunctions tests convenience logging functions
func TestConvenienceFunctions(t *testing.T) {
	// Reset for clean test state
	ResetForTesting()
	
	tmpDir, err := os.MkdirTemp("", "audit-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	Init(tmpDir+"/audit.log", true, InfoLevel)
	defer Close()

	tests := []struct {
		name string
		fn   func() error
	}{
		{"LogPolicyViolation", func() error {
			return LogPolicyViolation("filesystem", "/etc/passwd", "deny list")
		}},
		{"LogCommandBlocked", func() error {
			return LogCommandBlocked("rm -rf /", "dangerous command")
		}},
		{"LogSystemEvent", func() error {
			return LogSystemEvent(InfoLevel, "system started")
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.fn(); err != nil {
				t.Errorf("%s() failed: %v", tt.name, err)
			}
		})
	}
}

// TestEnabledDisable tests enabling/disabling logger
func TestEnabledDisable(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "audit-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	Init(tmpDir+"/audit.log", true, InfoLevel)
	defer Close()

	SetEnabled(false)
	if IsEnabled() {
		t.Error("Logger should be disabled")
	}

	SetEnabled(true)
	if !IsEnabled() {
		t.Error("Logger should be enabled")
	}
}
