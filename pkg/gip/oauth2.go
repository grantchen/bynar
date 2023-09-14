package gip

import (
	"context"
	"net/http"
	"os"

	"google.golang.org/api/option"
	"google.golang.org/api/transport"
)

// firebaseScopes is the set of OAuth2 scopes used by the Admin SDK.
var firebaseScopes = []string{
	"https://www.googleapis.com/auth/cloud-platform",
	"https://www.googleapis.com/auth/datastore",
	"https://www.googleapis.com/auth/devstorage.full_control",
	"https://www.googleapis.com/auth/firebase",
	"https://www.googleapis.com/auth/identitytoolkit",
	"https://www.googleapis.com/auth/userinfo.email",
}

// oauth2Client is client for Google API OAuth2.
type oauth2Client struct {
	projectID string
	opts      []option.ClientOption // option for a Google API client
}

// newOAuth2Client creates a new instance of the OAuth2 client.
func newOAuth2Client(ctx context.Context, opts ...option.ClientOption) (*oauth2Client, error) {
	o := []option.ClientOption{
		option.WithScopes(firebaseScopes...),
		option.WithCredentialsJSON([]byte(os.Getenv(ENVGoogleApplicationCredentialsJSON))),
	}
	o = append(o, opts...)

	return &oauth2Client{
		opts:      o,
		projectID: getProjectID(ctx, o...),
	}, nil
}

// newHttpClient creates a new instance of the Http Client with oauth2 options.
func (c *oauth2Client) newHttpClient(ctx context.Context) (*http.Client, error) {
	httpClient, _, err := transport.NewHTTPClient(ctx, c.opts...)
	if err != nil {
		return nil, err
	}
	return httpClient, nil
}

// getProjectID returns the project ID associated with the client from client options.
func getProjectID(ctx context.Context, opts ...option.ClientOption) string {
	creds, _ := transport.Creds(ctx, opts...)
	if creds != nil && creds.ProjectID != "" {
		return creds.ProjectID
	}

	return ""
}
