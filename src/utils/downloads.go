package utils

import (
	"os"
	"path/filepath"
	"runtime"
)

func GetDownloadsPath() (string, error) {
	var downloadsPath string

	switch runtime.GOOS {
	case "darwin":
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		downloadsPath = filepath.Join(homeDir, "Downloads")
	case "linux":
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		downloadsPath = filepath.Join(homeDir, "Downloads")
	case "windows":
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		downloadsPath = filepath.Join(homeDir, "Downloads")
	default:
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		downloadsPath = filepath.Join(homeDir, "Downloads")
	}

	if _, err := os.Stat(downloadsPath); os.IsNotExist(err) {
		if err := os.MkdirAll(downloadsPath, 0750); err != nil {
			return "", err
		}
	}

	return downloadsPath, nil
}

func EnsureDirectoryExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0750)
	}
	return nil
}
