package common

import (
	"os"
	"path/filepath"
)

// Wrapper function that returns the directory where the executable is location
func GetCurrentDirectory() string {
	return getCurrentDirectory()
}

// getCurrentDirectory returns the directory where the executable is location
func getCurrentDirectory() string {
	ex, err := os.Executable()
	os.Getwd()
	if err != nil {
		panic("Failed to get current directory")
	}
	return filepath.Dir(ex)
}
