package update

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

// apiTimeout limits calls to the GitHub API
const apiTimeout = 30 * time.Second

// GitHubAsset represents a file attached to a GitHub release
type GitHubAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

// GitHubRelease represents release information from the GitHub API
type GitHubRelease struct {
	TagName string        `json:"tag_name"`
	Name    string        `json:"name"`
	Assets  []GitHubAsset `json:"assets"`
}

// releaseURL is the GitHub API endpoint of the latest release
var releaseURL = fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", RepoUser, RepoName)

// getLatestRelease fetches the latest release from GitHub
func getLatestRelease() (*GitHubRelease, error) {
	client := &http.Client{Timeout: apiTimeout}

	req, err := http.NewRequest(http.MethodGet, releaseURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch release info: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status: %d", resp.StatusCode)
	}

	release := &GitHubRelease{}
	if err := json.NewDecoder(io.LimitReader(resp.Body, 10<<20)).Decode(release); err != nil {
		return nil, fmt.Errorf("failed to parse release info: %w", err)
	}
	return release, nil
}

// findAssets returns the archive for goos/goarch and its checksum file
func findAssets(release *GitHubRelease, goos, goarch string) (GitHubAsset, GitHubAsset, error) {
	var archive, checksum GitHubAsset

	suffix := regexp.MustCompile(`\.` + regexp.QuoteMeta(goos) + `\.` + regexp.QuoteMeta(goarch) + `\.tar\.gz$`)
	for _, asset := range release.Assets {
		if suffix.MatchString(asset.Name) {
			archive = asset
			break
		}
	}
	if len(archive.Name) == 0 {
		return archive, checksum, fmt.Errorf("cannot find download for %s %s in release %s", goos, goarch, release.TagName)
	}

	for _, asset := range release.Assets {
		if asset.Name == archive.Name+".sha256" {
			checksum = asset
			break
		}
	}
	if len(checksum.Name) == 0 {
		return archive, checksum, fmt.Errorf("release %s has no checksum for %s", release.TagName, archive.Name)
	}
	return archive, checksum, nil
}

// versionPattern matches the numeric part of tags like
// "tmux-project-0.18.0" and "v0.18.0"
var versionPattern = regexp.MustCompile(`(\d+)\.(\d+)\.(\d+)$`)

// parseVersion returns the major, minor and patch numbers of version
func parseVersion(version string) ([3]int, bool) {
	var retv [3]int
	matches := versionPattern.FindStringSubmatch(version)
	if matches == nil {
		return retv, false
	}
	for indx := range retv {
		number, err := strconv.Atoi(matches[indx+1])
		if err != nil {
			return retv, false
		}
		retv[indx] = number
	}
	return retv, true
}

// compareVersions returns -1, 0 or 1 when a is older than, equal to or
// newer than b
func compareVersions(a, b [3]int) int {
	for indx := range a {
		if a[indx] < b[indx] {
			return -1
		}
		if a[indx] > b[indx] {
			return 1
		}
	}
	return 0
}
