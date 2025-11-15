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

type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
}

func getLatestVersion() (string, error) {
	utils.LogStart()
	defer utils.LogEnd()

	url := "https://api.github.com/repos/jvzantvoort/tmux-project/releases/latest"
	resp, err := http.Get(url)
	if err != nil {
		log.Errorf("failed to fetch release info: %w", err)
		return "", fmt.Errorf("failed to fetch release info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Errorf("GitHub API returned status: %d", resp.StatusCode)
		return "", fmt.Errorf("GitHub API returned status: %d", resp.StatusCode)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		log.Errorf("failed to parse release info: %w", err)
		return "", fmt.Errorf("failed to parse release info: %w", err)
	}

	return release.TagName, nil
}

func needsUpdate(currentVersion, latestVersion string) bool {
	// Normalize versions by removing "v" prefix if present
	current := strings.TrimPrefix(currentVersion, "v")
	latest := strings.TrimPrefix(latestVersion, "v")

	return current != latest
}

func Execute(forceUpdate bool) error {
	currentVersion := version.GetVersion() // get current version

	latestVersion, err := getLatestVersion() // Check latest version
	if err != nil {
		log.Errorf("failed to check for updates: %w", err)
		return fmt.Errorf("failed to check for updates: %w", err)
	}

	// Check if update is needed
	if !forceUpdate && !needsUpdate(currentVersion.Short(), latestVersion) {
		return fmt.Errorf("already running the latest version")
	}

	// get the download url
	url, err := getBrowserDownloadURL()
	if err != nil {
		return fmt.Errorf("failed to download update: %w", err)
	}

	return Install(url)
}
