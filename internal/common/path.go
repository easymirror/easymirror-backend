package common

import (
	"fmt"
	"net/url"
	"os"
	"path"
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

// FilenameFromURI returns the name of a file from a given URI
func FilenameFromURI(uri string) (string, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return "", fmt.Errorf("parse error: %w", err)
	}
	return path.Base(u.Path), nil
}
