package version

import (
	"fmt"
	"runtime/debug"
	"time"
)

var (
	// Version is set via ldflags at build time for releases
	Version = ""
	// Commit is set via ldflags at build time
	Commit = ""
	// BuildTime is set via ldflags at build time
	BuildTime = ""
)

// Info represents version information about the application
type Info struct {
	Version   string
	Commit    string
	BuildTime string
	BuildDate time.Time
}

// GetVersion returns the version information
func GetVersion() Info {
	info := Info{
		Version:   Version,
		Commit:    Commit,
		BuildTime: BuildTime,
	}

	// If version is not set via ldflags, try to get it from build info
	if info.Version == "" {
		if buildInfo, ok := debug.ReadBuildInfo(); ok {
			info.Version = buildInfo.Main.Version
		}
	}

	// If still empty, use a default
	if info.Version == "" || info.Version == "(devel)" {
		info.Version = "dev"
	}

	// Parse build time if available
	if info.BuildTime != "" {
		if t, err := time.Parse(time.RFC3339, info.BuildTime); err == nil {
			info.BuildDate = t
		}
	}

	return info
}

// String returns a formatted version string
func (i Info) String() string {
	if i.Commit != "" {
		if i.BuildTime != "" {
			return fmt.Sprintf("%s (commit: %s, built: %s)", i.Version, i.Commit, i.BuildTime)
		}
		return fmt.Sprintf("%s (commit: %s)", i.Version, i.Commit)
	}
	return i.Version
}

// Short returns just the version number
func (i Info) Short() string {
	return i.Version
}

// IsRelease returns true if this is a release build (has a tag version)
func (i Info) IsRelease() bool {
	return i.Version != "" && i.Version != "dev" && i.Version != "(devel)"
}
