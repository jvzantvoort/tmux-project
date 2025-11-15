// Package update provides self-update functionality for the tmux-project application.
// It checks for new releases on GitHub and can download and install updates automatically.
package update

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/jvzantvoort/tmux-project/utils"
	"github.com/jvzantvoort/tmux-project/version"
	log "github.com/sirupsen/logrus"
)

// GitHubRelease represents release information from the GitHub API
type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
}

// getLatestVersion fetches the latest release version from GitHub
func getLatestVersion() (string, error) {
	utils.LogStart()
	defer utils.LogEnd()

	url := "https://api.github.com/repos/jvzantvoort/tmux-project/releases/latest"
	resp, err := http.Get(url)
	if err != nil {
		log.Errorf("failed to fetch release info: %v", err)
		return "", fmt.Errorf("failed to fetch release info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Errorf("GitHub API returned status: %d", resp.StatusCode)
		return "", fmt.Errorf("GitHub API returned status: %d", resp.StatusCode)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		log.Errorf("failed to parse release info: %v", err)
		return "", fmt.Errorf("failed to parse release info: %w", err)
	}

	return release.TagName, nil
}

// needsUpdate compares current and latest versions to determine if an update is needed
func needsUpdate(currentVersion, latestVersion string) bool {
	current := strings.TrimPrefix(currentVersion, "v")
	latest := strings.TrimPrefix(latestVersion, "v")

	return current != latest
}

// Execute checks for updates and installs the latest version if available.
// If forceUpdate is true, it will update even if already on the latest version.
func Execute(forceUpdate bool) error {
	currentVersion := version.GetVersion()

	latestVersion, err := getLatestVersion()
	if err != nil {
		log.Errorf("failed to check for updates: %v", err)
		return fmt.Errorf("failed to check for updates: %w", err)
	}

	if !forceUpdate && !needsUpdate(currentVersion.Short(), latestVersion) {
		return fmt.Errorf("already running the latest version")
	}

	url, err := getBrowserDownloadURL()
	if err != nil {
		return fmt.Errorf("failed to download update: %w", err)
	}

	return Install(url)
}
