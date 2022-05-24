package github

import (
	"fmt"
	"net/http"

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

// We need to pass token for higher privileges - to access draft releases / non-published releases.
// But we don't need token for fetching published releases that are public
func FetchTceRelease(releaseVersion string, token string) (*shared.Release, error) {

	client, err := NewAuthenticatedClient(token)
	if err != nil {
		return nil, fmt.Errorf("error while creating GitHub client using token: %v", err)
	}

	return shared.FetchRelease(&http.Client{
		Transport: client,
	}, tceBaseRepo, releaseVersion)
}
