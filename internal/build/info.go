package build

type Info struct {
	AppVersion string `json:"app_version"`
	BuildTime  string `json:"build_time"`
	CommitHash string `json:"commit_hash"`
	GitTag     string `json:"git_tag"`
	GoVersion  string `json:"go_version"`
}

// GetInfo returns complete build info
func GetInfo() Info {
	return Info{
		AppVersion: AppVersion(),
		BuildTime:  BuildTime(),
		CommitHash: CommitHash(),
		GitTag:     GitTag(),
		GoVersion:  GoVersion(),
	}
}
