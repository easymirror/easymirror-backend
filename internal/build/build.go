package build

import "strings"

const fallback = "N/A"

// Data about the build. These are set during build time
var (
	appVersion string // Build version of the current app
	buildTime  string // Time this app was built
	commitHash string // Hash of this commit
	gitTag     string // Tag of the git branch
	goVersion  string // Version of Go this app was built with

)

// AppVersion returns the current build version of the app
func AppVersion() string {
	return returnIfNotEmpty(appVersion)
}

// BuildTime returns the time this app was built
func BuildTime() string {
	return returnIfNotEmpty(buildTime)
}

// CommitHash returns the hash of this commit
func CommitHash() string {
	return returnIfNotEmpty(commitHash)
}

// GitTag returns tag of the git branch
func GitTag() string {
	return returnIfNotEmpty(gitTag)
}

// GoVersion version of Go this app was built with
func GoVersion() string {
	return returnIfNotEmpty(goVersion)
}

// returnIfNotEmpty returns a given input if the input != ""
// If input == "", it returns "N/A"
func returnIfNotEmpty(input string) string {
	if strings.TrimSpace(input) == "" {
		return fallback
	}
	return input
}
