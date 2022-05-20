package github_test

import (
	"testing"

	"github.com/karuppiah7890/tce-e2e-test/testutils/github"
	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
)

func TestFetchTceRelease(t *testing.T) {
	log.InitLogger("fetch-tce-release")

	release, err := github.FetchTceRelease("v0.12.1")

	if err != nil {
		log.Fatalf("expected no error but got error while fetching TCE release: %v", err)
	}

	log.Infof("Release: %v", release)
}
