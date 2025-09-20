package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetDownloadsPath(t *testing.T) {
	downloadsPath, err := GetDownloadsPath()
	if err != nil {
		t.Fatalf("GetDownloadsPath() failed: %v", err)
	}

	if downloadsPath == "" {
		t.Error("GetDownloadsPath() returned empty path")
	}

	if !filepath.IsAbs(downloadsPath) {
		t.Error("GetDownloadsPath() should return absolute path")
	}

	if _, err := os.Stat(downloadsPath); os.IsNotExist(err) {
		t.Errorf("Downloads path does not exist: %s", downloadsPath)
	}
}

func TestEnsureDirectoryExists(t *testing.T) {
	tempDir := os.TempDir()
	testDir := filepath.Join(tempDir, "cortex-test", "nested", "directory")

	defer os.RemoveAll(filepath.Join(tempDir, "cortex-test"))

	if err := EnsureDirectoryExists(testDir); err != nil {
		t.Fatalf("EnsureDirectoryExists() failed: %v", err)
	}

	if _, err := os.Stat(testDir); os.IsNotExist(err) {
		t.Errorf("Directory was not created: %s", testDir)
	}

	if err := EnsureDirectoryExists(testDir); err != nil {
		t.Errorf("EnsureDirectoryExists() should not fail on existing directory: %v", err)
	}
}
