// Package update provides self-update functionality for the tmux-project application.
// It checks for new releases on GitHub and can download and install updates automatically.
package update

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/jvzantvoort/tmux-project/version"
	log "github.com/sirupsen/logrus"
)

// Execute checks for updates and installs the latest version if available.
// If forceUpdate is true, it will update even if already on the latest version.
func Execute(forceUpdate bool) error {
	currentVersion := version.GetVersion().Short()

	release, err := getLatestRelease()
	if err != nil {
		return fmt.Errorf("failed to check for updates: %w", err)
	}

	if !forceUpdate {
		current, ok := parseVersion(currentVersion)
		if !ok {
			return fmt.Errorf("running development build %s, use --force to install %s", currentVersion, release.TagName)
		}
		latest, ok := parseVersion(release.TagName)
		if !ok {
			return fmt.Errorf("cannot parse version of latest release %s", release.TagName)
		}
		if compareVersions(current, latest) >= 0 {
			log.Infof("already running the latest version (%s)", currentVersion)
			return nil
		}
	}

	archive, checksum, err := findAssets(release, runtime.GOOS, runtime.GOARCH)
	if err != nil {
		return err
	}

	execPath, err := ThisExecutablePath()
	if err != nil {
		return err
	}

	log.Infof("updating %s to %s", currentVersion, release.TagName)
	return Install(archive.BrowserDownloadURL, checksum.BrowserDownloadURL, filepath.Dir(execPath))
}
