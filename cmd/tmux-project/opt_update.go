package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"github.com/jvzantvoort/tmux-project/version"
	"github.com/minio/selfupdate"
	"github.com/spf13/cobra"
)

var (
	forceUpdate bool
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update tmux-project from GitHub",
	Long:  "Download and install the latest version of tmux-project from GitHub",
	RunE:  runUpdate,
}

func init() {
	updateCmd.Flags().BoolVarP(&forceUpdate, "force", "f", false, "Force update even if already up to date")
	rootCmd.AddCommand(updateCmd)
}

type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
}

func getLatestVersion() (string, error) {
	url := "https://api.github.com/repos/jvzantvoort/tmux-project/releases/latest"
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch release info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API returned status: %d", resp.StatusCode)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
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

func runUpdate(cmd *cobra.Command, args []string) error {
	currentVersion := version.GetVersion()

	fmt.Printf("Current version: %s\n", currentVersion.Short())

	// Check latest version
	latestVersion, err := getLatestVersion()
	if err != nil {
		return fmt.Errorf("failed to check for updates: %w", err)
	}

	fmt.Printf("Latest version:  %s\n", latestVersion)

	// Check if update is needed
	if !forceUpdate && !needsUpdate(currentVersion.Short(), latestVersion) {
		fmt.Println("\nAlready running the latest version!")
		return nil
	}

	if !forceUpdate {
		fmt.Println("\nNew version available!")
	}

	// Determine architecture and OS
	arch := runtime.GOARCH
	if arch == "amd64" {
		arch = "x86_64"
	}
	goos := runtime.GOOS
	if goos == "darwin" {
		goos = "Darwin"
	} else if goos == "linux" {
		goos = "Linux"
	}

	// Construct download URL
	url := fmt.Sprintf("https://github.com/jvzantvoort/tmux-project/releases/latest/download/tmux-project_%s_%s", goos, arch)

	fmt.Printf("Downloading from: %s\n", url)

	// Download the new binary
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download update: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status: %d", resp.StatusCode)
	}

	// Apply update using selfupdate for atomic replacement
	err = selfupdate.Apply(resp.Body, selfupdate.Options{})
	if err != nil {
		if rerr := selfupdate.RollbackError(err); rerr != nil {
			return fmt.Errorf("failed to apply update and rollback failed: %w", rerr)
		}
		return fmt.Errorf("failed to apply update: %w", err)
	}

	fmt.Println("Update completed successfully!")
	return nil
}
