package github

import (
	"fmt"
	"net/http"

	"github.com/cli/cli/v2/git"
	"github.com/cli/cli/v2/pkg/cmd/release/shared"
)

type BaseRepo struct {
	host  string
	owner string
	name  string
}

func (baseRepo BaseRepo) RepoName() string {
	return baseRepo.name
}
func (baseRepo BaseRepo) RepoOwner() string {
	return baseRepo.owner
}
func (baseRepo BaseRepo) RepoHost() string {
	return baseRepo.host
}

var tceBaseRepo = BaseRepo{
	host:  "github.com",
	owner: "vmware-tanzu",
	name:  "community-edition",
}

// TODO: Generic function to fetch any repo release and not just TCE. This will be useful for testing and
// it will be generic too.

// TODO: Inject HTTP Client? Authenticated / unauthenticated

// TODO: Return a custom Release struct which contains only the list of links to different assets in a release
// i.e link to download the assets.
// TODO: Maybe make token optional.
// We need to pass token for higher privileges - to access draft releases / non-published releases.
// But we don't need token for fetching published releases that are public
func FetchTceRelease(releaseVersion string, token string) (*shared.Release, error) {

	// TODO: We throw error when no token is passed. But token is needed only for fetching a draft release
	client, err := NewAuthenticatedClient(token)
	if err != nil {
		return nil, fmt.Errorf("error while creating GitHub client using token: %v", err)
	}

	return shared.FetchRelease(&http.Client{
		Transport: client,
	}, tceBaseRepo, releaseVersion)
}

func CloneRepo(URL string) error {
	arg := []string{}
	_, err := git.RunClone(URL, arg)
	if err != nil {
		return fmt.Errorf("error while creating GitHub client using token: %v", err)
	}
	return nil
}
