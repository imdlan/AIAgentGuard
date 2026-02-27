package scanner

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/imdlan/AIAgentGuard/pkg/model"
)

// TestScanFileContents tests the file content scanning functionality
func TestScanFileContents(t *testing.T) {
	// Create temporary test directory
	tmpDir, err := os.MkdirTemp("", "agentguard-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name     string
		content  string
		filename string
		wantRisk bool
	}{
		{
			name:     "AWS Access Key",
			content:  "AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE\n",
			filename: "aws.env",
			wantRisk: true,
		},
		{
			name:     "GitHub Token",
			content:  "ghp_1234567890abcdefghijklmnopqrstuvwxyz1234\n",
			filename: "github.env",
			wantRisk: true,
		},
		{
			name:     "Safe content",
			content:  "This is just normal text\n",
			filename: "safe.txt",
			wantRisk: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testFile := filepath.Join(tmpDir, tt.filename)
			if err := os.WriteFile(testFile, []byte(tt.content), 0644); err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			risk, _ := scanFileForKeys(testFile)

			if tt.wantRisk && risk == model.Low {
				t.Errorf("Expected high/medium risk, got Low")
			}
			if !tt.wantRisk && risk != model.Low {
				t.Errorf("Expected Low risk, got %v", risk)
			}
		})
	}
}

// TestIsConfigFile tests config file detection
func TestIsConfigFile(t *testing.T) {
	tests := []struct {
		filename string
		want     bool
	}{
		{".env", true},
		{"config.yml", true},
		{"README.md", false},
		{"main.go", false},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			got := isConfigFile(tt.filename)
			if got != tt.want {
				t.Errorf("isConfigFile(%q) = %v, want %v", tt.filename, got, tt.want)
			}
		})
	}
}
