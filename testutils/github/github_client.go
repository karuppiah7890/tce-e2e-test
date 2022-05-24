package github

import (
	"fmt"
	"net/http"
)

// AuthenticatedClient implements the net/http.RoundTripper interface
type AuthenticatedClient struct {
	token string
}

func NewAuthenticatedClient(token string) (*AuthenticatedClient, error) {
	if token == "" {
		return nil, fmt.Errorf("GitHub token cannot be empty")
	}
	return &AuthenticatedClient{token: token}, nil
}

func (c *AuthenticatedClient) RoundTrip(req *http.Request) (*http.Response, error) {
	req2 := setTokenAsHeader(req, c.token)
	// Make the HTTP request.
	return http.DefaultTransport.RoundTrip(req2)
}

func setTokenAsHeader(req *http.Request, token string) *http.Request {
	// To set extra headers, we must make a copy of the Request so
	// that we don't modify the Request we were given. This is required by the
	// specification of http.RoundTripper.
	//
	// Since we are going to modify only req.Header here, we only need a deep copy
	// of req.Header.
	convertedRequest := new(http.Request)
	*convertedRequest = *req
	convertedRequest.Header = make(http.Header, len(req.Header))

	for k, s := range req.Header {
		convertedRequest.Header[k] = append([]string(nil), s...)
	}
	SetToken(convertedRequest, token)
	return convertedRequest
}

// SetToken sets the request's Authorization header to use HTTP
// Token Authentication with the token
func SetToken(r *http.Request, token string) {
	r.Header.Set("Authorization", "Bearer "+token)
}
