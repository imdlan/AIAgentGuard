package multilang

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/imdlan/AIAgentGuard/pkg/model"
)

// TestScanNpmDependencies_ReturnsValidRisk tests that npm scanning returns a valid risk level
func TestScanNpmDependencies_ReturnsValidRisk(t *testing.T) {
	risk := ScanNpmDependencies()
	// Should return a valid risk level
	if risk != model.Low && risk != model.Medium && risk != model.High && risk != model.Critical {
		t.Errorf("Expected valid risk level, got %v", risk)
	}
}

// TestScanPipDependencies_ReturnsValidRisk tests that pip scanning returns a valid risk level
func TestScanPipDependencies_ReturnsValidRisk(t *testing.T) {
	risk := ScanPipDependencies()
	// Should return a valid risk level
	if risk != model.Low && risk != model.Medium && risk != model.High && risk != model.Critical {
		t.Errorf("Expected valid risk level, got %v", risk)
	}
}

// TestScanCargoDependencies_ReturnsValidRisk tests that cargo scanning returns a valid risk level
func TestScanCargoDependencies_ReturnsValidRisk(t *testing.T) {
	risk := ScanCargoDependencies()
	// Should return a valid risk level
	if risk != model.Low && risk != model.Medium && risk != model.High && risk != model.Critical {
		t.Errorf("Expected valid risk level, got %v", risk)
	}
}

// TestFileExists tests the fileExists helper function
func TestFileExists(t *testing.T) {
	// Test with a file that should exist
	tmpFile, err := os.CreateTemp("", "test-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	if !fileExists(tmpFile.Name()) {
		t.Errorf("Expected file to exist: %s", tmpFile.Name())
	}

	// Test with a file that doesn't exist
	nonExistentPath := filepath.Join(os.TempDir(), "non-existent-file-12345.txt")
	if fileExists(nonExistentPath) {
		t.Errorf("Expected file to not exist: %s", nonExistentPath)
	}
}

// TestFindPackageFiles_SkipsNodeModules tests that node_modules directories are skipped
func TestFindPackageFiles_SkipsNodeModules(t *testing.T) {
	// Create a temporary directory structure
	tmpDir, err := os.MkdirTemp("", "npm-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Save current directory and change to tmpDir
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Create test subdirectories and files
	testDir := "test-project"
	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	packageLockPath := filepath.Join(testDir, "package-lock.json")
	if err := os.WriteFile(packageLockPath, []byte("{}"), 0644); err != nil {
		t.Fatalf("Failed to create package-lock.json: %v", err)
	}

	// Create node_modules directory (should be skipped)
	nodeModulesDir := filepath.Join(testDir, "node_modules")
	if err := os.MkdirAll(nodeModulesDir, 0755); err != nil {
		t.Fatalf("Failed to create node_modules: %v", err)
	}
	packageLockInNodeModules := filepath.Join(nodeModulesDir, "package-lock.json")
	if err := os.WriteFile(packageLockInNodeModules, []byte("{}"), 0644); err != nil {
		t.Fatalf("Failed to create package-lock.json in node_modules: %v", err)
	}
	// Verify files exist before testing
	if _, err := os.Stat(packageLockPath); os.IsNotExist(err) {
		t.Fatalf("package-lock.json does not exist at %s", packageLockPath)
	}
	
	// List all files in test directory
	dirEntries, _ := os.ReadDir(testDir)
	t.Logf("Files in testDir: ")
	for _, entry := range dirEntries {
		t.Logf("  - %s", entry.Name())
	}
	
	// Test finding package files
	files := findPackageFiles()
	
	// Debug output
	t.Logf("Found files: %v", files)
	t.Logf("Current dir: %s", func() string { dir, _ := os.Getwd(); return dir }())
	
	// Debug: test filepath.Walk directly
	walkFound := []string{}
	_ = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && filepath.Base(path) == "package-lock.json" {
			walkFound = append(walkFound, path)
		}
		return nil
	})
	t.Logf("Direct walk found: %v", walkFound)
	// Should find the package-lock.json in test-project but not in node_modules
	foundInProject := false
	foundInNodeModules := false

	for _, file := range files {
		if filepath.Base(file) == "package-lock.json" {
			if strings.Contains(file, "node_modules") {
				foundInNodeModules = true
			} else if strings.Contains(file, testDir) {
				foundInProject = true
			}
		}
	}

	if !foundInProject {
		t.Error("Expected to find package-lock.json file in test-project")
	}
	if foundInNodeModules {
		t.Error("Should not find package-lock.json file in node_modules directory")
	}
}

// TestFindRequirementFiles_SkipsPycache tests that __pycache__ directories are skipped
func TestFindRequirementFiles_SkipsPycache(t *testing.T) {
	// Create a temporary directory structure
	tmpDir, err := os.MkdirTemp("", "pip-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Save current directory and change to tmpDir
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Create test subdirectories and files
	testDir := "test-project"
	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	reqPath := filepath.Join(testDir, "requirements.txt")
	if err := os.WriteFile(reqPath, []byte("requests==2.28.0\n"), 0644); err != nil {
		t.Fatalf("Failed to create requirements.txt: %v", err)
	}

	// Create __pycache__ directory (should be skipped)
	pycacheDir := filepath.Join(testDir, "__pycache__")
	if err := os.MkdirAll(pycacheDir, 0755); err != nil {
		t.Fatalf("Failed to create __pycache__: %v", err)
	}
	reqInPycache := filepath.Join(pycacheDir, "requirements.txt")
	if err := os.WriteFile(reqInPycache, []byte("requests==2.28.0\n"), 0644); err != nil {
		t.Fatalf("Failed to create requirements.txt in __pycache__: %v", err)
	}

	// Test finding requirement files
	files := findRequirementFiles()

	// Should find the requirements.txt in test-project but not in __pycache__
	foundInProject := false
	foundInPycache := false

	for _, file := range files {
		if filepath.Base(file) == "requirements.txt" {
			if strings.Contains(file, "__pycache__") {
				foundInPycache = true
			} else if strings.Contains(file, testDir) {
				foundInProject = true
			}
		}
	}

	if !foundInProject {
		t.Error("Expected to find requirements.txt file in test-project")
	}
	if foundInPycache {
		t.Error("Should not find requirements.txt file in __pycache__ directory")
	}
}

// TestFindCargoFiles_SkipsTarget tests that target directories are skipped
func TestFindCargoFiles_SkipsTarget(t *testing.T) {
	// Create a temporary directory structure
	tmpDir, err := os.MkdirTemp("", "cargo-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Save current directory and change to tmpDir
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Create test subdirectories and files
	testDir := "test-project"
	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	cargoPath := filepath.Join(testDir, "Cargo.toml")
	cargoContent := `[package]
name = "test"
version = "0.1.0"
`
	if err := os.WriteFile(cargoPath, []byte(cargoContent), 0644); err != nil {
		t.Fatalf("Failed to create Cargo.toml: %v", err)
	}

	// Create target directory (should be skipped)
	targetDir := filepath.Join(testDir, "target")
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		t.Fatalf("Failed to create target directory: %v", err)
	}
	cargoInTarget := filepath.Join(targetDir, "Cargo.toml")
	if err := os.WriteFile(cargoInTarget, []byte(cargoContent), 0644); err != nil {
		t.Fatalf("Failed to create Cargo.toml in target: %v", err)
	}

	// Test finding Cargo files
	files := findCargoFiles()

	// Should find the Cargo.toml in test-project but not in target
	foundInProject := false
	foundInTarget := false

	for _, file := range files {
		if filepath.Base(file) == "Cargo.toml" {
			if strings.Contains(file, filepath.Join(testDir, "target")) {
				foundInTarget = true
			} else if strings.Contains(file, testDir) && !strings.Contains(file, "target") {
				foundInProject = true
			}
		}
	}

	if !foundInProject {
		t.Error("Expected to find Cargo.toml file in test-project")
	}
	if foundInTarget {
		t.Error("Should not find Cargo.toml file in target directory")
	}
}

// TestConvertVulnCount tests the vulnerability count conversion
func TestConvertVulnCount(t *testing.T) {
	jsonCount := VulnerabilityCountJSON{
		Low:      1,
		Moderate: 2,
		High:     3,
		Critical: 4,
	}

	result := convertVulnCount(jsonCount)

	if result.Low != 1 {
		t.Errorf("Expected Low=1, got %d", result.Low)
	}
	if result.Moderate != 2 {
		t.Errorf("Expected Moderate=2, got %d", result.Moderate)
	}
	if result.High != 3 {
		t.Errorf("Expected High=3, got %d", result.High)
	}
	if result.Critical != 4 {
		t.Errorf("Expected Critical=4, got %d", result.Critical)
	}
}

// BenchmarkScanNpmDependencies benchmarks npm dependency scanning
func BenchmarkScanNpmDependencies(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ScanNpmDependencies()
	}
}

// BenchmarkScanPipDependencies benchmarks pip dependency scanning
func BenchmarkScanPipDependencies(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ScanPipDependencies()
	}
}

// BenchmarkScanCargoDependencies benchmarks cargo dependency scanning
func BenchmarkScanCargoDependencies(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ScanCargoDependencies()
	}
}
