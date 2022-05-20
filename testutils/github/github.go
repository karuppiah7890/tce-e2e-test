package github

import (
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

func FetchTceRelease(releaseVersion string) (*shared.Release, error) {
	return shared.FetchRelease(&http.Client{}, tceBaseRepo, releaseVersion)
}
