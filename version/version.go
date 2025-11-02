package version

// BuildInfo contains build-time information
type BuildInfo struct {
	Commit string
	Tag    string
}

// These are set at compile time via ldflags
var (
	commit = "unknown"
	tag    = "unknown"
)

// GetInfo returns the build information (reads the ldflags-set variables)
func GetInfo() BuildInfo {
	return BuildInfo{
		Commit: commit,
		Tag:    tag,
	}
}

// Info is deprecated, use GetInfo() instead
// This is kept for backward compatibility if needed
var Info = BuildInfo{
	Commit: "unknown",
	Tag:    "unknown",
}

