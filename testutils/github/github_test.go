package github_test

import (
	"os"
	"testing"

	"github.com/karuppiah7890/tce-e2e-test/testutils/github"
	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
)

func TestFetchTceRelease(t *testing.T) {
	log.InitLogger("fetch-tce-release")

	release, err := github.FetchTceRelease("v0.12.1", os.Getenv("GITHUB_TOKEN"))

	if err != nil {
		log.Fatalf("expected no error but got error while fetching TCE release: %v", err)
	}

	// Maybe use mock to check if token is set? And not send request to actual GitHub server.
	// This way we can check if token is set correctly, then we can fetch published or non-published releases,
	// but yeah, for published releases, we don't even need a GitHub token

	log.Infof("Release: %v", release)
}
