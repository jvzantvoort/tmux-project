package update

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"runtime"

	"github.com/google/go-github/v62/github"
)

func extractFields(release *github.RepositoryRelease) []map[string]string {
	retv := []map[string]string{}
	re := regexp.MustCompile(`(?P<goos>[^.]+)\.(?P<goarch>[^.]+)\.tar\.gz$`)

	for _, indx := range release.Assets {
		tmpdict := map[string]string{}
		filename := *indx.Name

		matches := re.FindStringSubmatch(filename)
		if matches != nil {
			for i, name := range re.SubexpNames() {
				if i != 0 && name != "" && i < len(matches) {
					tmpdict[name] = matches[i]
				}
			}
		}

		tmpdict["Name"] = *indx.Name
		tmpdict["BrowserDownloadURL"] = *indx.BrowserDownloadURL

		retv = append(retv, tmpdict)
	}

	return retv
}

func getLatest(user, repo string) ([]map[string]string, error) {
	client := github.NewClient(nil)

	ctx := context.Background()
	release, _, err := client.Repositories.GetLatestRelease(ctx, user, repo)

	// 3. Handle errors
	if err != nil {
		log.Fatalf("Failed to get release: %v", err)
		retv := []map[string]string{}
		return retv, err
	}

	return extractFields(release), nil

}

func getBrowserDownloadURL() (string, error) {
	repolist, err := getLatest(RepoUser, RepoName)
	goarch := runtime.GOARCH
	goos := runtime.GOOS
	if err != nil {
		return "", err
	}
	for _, row := range repolist {
		if row["goos"] == goos && row["goarch"] == goarch {
			return row["BrowserDownloadURL"], nil
		}
	}
	return "", fmt.Errorf("cannot find download for %s %s", goos, goarch)
}
