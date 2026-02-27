package scanner

import (
	"os"
	"testing"

	"github.com/imdlan/AIAgentGuard/pkg/model"
)

// TestCheckSpecificFile tests SUID/SGID file checking
func TestCheckSpecificFile(t *testing.T) {
	tmpFile := "/tmp/test-suid-check"
	file, err := os.Create(tmpFile)
	if err != nil {
		t.Skipf("Cannot create test file: %v", err)
	}
	file.Close()
	defer os.Remove(tmpFile)

	hasSUID, hasSGID, err := CheckSpecificFile(tmpFile)
	if err != nil {
		t.Errorf("CheckSpecificFile(%q) failed: %v", tmpFile, err)
	}

	if hasSUID || hasSGID {
		t.Error("Temporary file should not have SUID/SGID bits")
	}
}

// TestGetSUIDStatistics tests SUID statistics
func TestGetSUIDStatistics(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping SUID scan in short mode")
	}

	stats := GetSUIDStatistics()

	if stats == nil {
		t.Error("GetSUIDStatistics() returned nil")
	}

	expectedKeys := []string{"total_suid", "known_suid", "unknown_suid", "high_risk", "medium_risk", "home_directory"}
	for _, key := range expectedKeys {
		if _, ok := stats[key]; !ok {
			t.Errorf("Stats missing key: %s", key)
		}
	}

	t.Logf("SUID Statistics: %+v", stats)
}

// TestScanSUIDFiles tests the main SUID scanning function
func TestScanSUIDFiles(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping SUID scan in short mode")
	}

	risk := ScanSUIDFiles()

	if risk != model.Low && risk != model.Medium && risk != model.High {
		t.Errorf("ScanSUIDFiles() returned invalid risk level: %v", risk)
	}

	t.Logf("ScanSUIDFiles() returned risk level: %v", risk)
}
