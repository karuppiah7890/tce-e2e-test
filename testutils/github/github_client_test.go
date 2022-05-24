package github_test

import (
	"net/http"
	"testing"

	"github.com/karuppiah7890/tce-e2e-test/testutils/github"
	"github.com/karuppiah7890/tce-e2e-test/testutils/log"

	"gopkg.in/h2non/gock.v1"
)

func TestAuthenticatedClient(t *testing.T) {
	log.InitLogger("authenticated-github-client")
	defer gock.Off() // Flush pending mocks after test execution

	gock.New("http://dummy-github-server.com").
		Get("/bar").
		MatchHeader("Authorization", "^Bearer dummy-token$").
		Reply(200).
		JSON(map[string]string{"foo": "bar"})

	client, err := github.NewAuthenticatedClient("dummy-token")

	if err != nil {
		log.Fatalf("expected no error while creating authenticated client but got error: %v", err)
	}

	req, err := http.NewRequest(http.MethodGet, "http://dummy-github-server.com/bar", nil)

	if err != nil {
		log.Fatalf("expected no error while creating dummy test request but got error: %v", err)
	}

	res, err := client.RoundTrip(req)

	if err != nil {
		log.Fatalf("expected no error while sending request and getting response but got error: %v", err)
	}

	if res.StatusCode != 200 {
		log.Fatalf("expected response status code to be 200 but got this: %v", res.StatusCode)
	}

}
