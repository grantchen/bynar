package gip

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/utils"
)

const (
	verifyCustomTokenURL = "https://www.googleapis.com/identitytoolkit/v3/relyingparty/verifyCustomToken?key=%s"
)

// gipClient is the interface for the AuthProvider.
type gipClient struct {
	app    *firebase.App
	apiKey string
}

// TODO test
func NewGIPClient() (AuthProvider, error) {
	opt := option.WithCredentialsFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, errors.New("error initializing gip app: " + err.Error())
	}

	return &gipClient{
		app:    app,
		apiKey: os.Getenv("GOOGLE_API_KEY"),
	}, nil
}

// GetUser gets the user data corresponding to the specified user ID.
func (g gipClient) GetUser(ctx context.Context, uid string) (*auth.UserRecord, error) {
	client, err := g.app.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting Auth client: %v", err)
	}

	u, err := client.GetUser(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("error getting user %s: %v", uid, err)
	}

	return u, nil
}

// GetUserByEmail gets the user data corresponding to the specified email.
func (g gipClient) GetUserByEmail(ctx context.Context, email string) (*auth.UserRecord, error) {
	client, err := g.app.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting Auth client: %v", err)
	}

	u, err := client.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("error getting user by email %s: %v", email, err)
	}

	return u, nil
}

// CreateUser creates a new user with the specified properties.
func (g gipClient) CreateUser(ctx context.Context, email string) (*auth.UserRecord, error) {
	client, err := g.app.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting Auth client: %v", err)
	}

	params := (&auth.UserToCreate{}).
		Email(email).
		EmailVerified(false).
		Password(utils.RandString(10)).
		Disabled(false)
	u, err := client.CreateUser(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("error creating user: %v", err)
	}

	return u, nil
}

// UpdateUser updates an existing user account with the specified properties.
func (g gipClient) UpdateUser(ctx context.Context, uid string, params map[string]interface{}) (*auth.UserRecord, error) {
	client, err := g.app.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting Auth client: %v", err)
	}

	if len(params) == 0 {
		return nil, errors.New("no params provided")
	}

	updateParams := &auth.UserToUpdate{}
	if email, ok := params["email"].(string); ok {
		updateParams.Email(email)
	}

	if disabled, ok := params["disabled"].(bool); ok {
		updateParams.Disabled(disabled)
	}

	u, err := client.UpdateUser(ctx, uid, updateParams)
	if err != nil {
		return nil, fmt.Errorf("error updating user: %v", err)
	}

	return u, nil
}

// DeleteUser deletes the user by the given UID.
func (g gipClient) DeleteUser(ctx context.Context, uid string) error {
	client, err := g.app.Auth(ctx)
	if err != nil {
		return fmt.Errorf("error getting Auth client: %v", err)
	}

	err = client.DeleteUser(ctx, uid)
	if err != nil {
		return fmt.Errorf("error deleting user: %v", err)
	}

	return nil
}

// LogOut logs out the user by the given UID.
func (g gipClient) LogOut(ctx context.Context, uid string) error {
	client, err := g.app.Auth(ctx)
	if err != nil {
		return fmt.Errorf("error getting Auth client: %v", err)
	}

	err = client.RevokeRefreshTokens(ctx, uid)
	if err != nil {
		return fmt.Errorf("error revoke token: %v", err)
	}

	return nil
}

// SignIn signs in with the provided token.
func (g gipClient) SignIn(ctx context.Context, token string) (idToken string, err error) {
	return g.signInWithCustomTokenForTenant(token, "")
}

// EmailSignInLink generates the out-of-band email action link for email link sign-in flows, using the action
// code settings provided.
func (g gipClient) EmailSignInLink(ctx context.Context, email string) (string, error) {
	client, err := g.app.Auth(ctx)
	if err != nil {
		return "", fmt.Errorf("error getting Auth client: %v", err)
	}

	actionCodeSettings := newActionCodeSettings()
	link, err := client.EmailSignInLink(ctx, email, actionCodeSettings)
	if err != nil {
		return "", fmt.Errorf("error generating email link: %v", err)
	}

	return link, nil
}

func (g gipClient) signInWithCustomTokenForTenant(token string, tenantID string) (string, error) {
	payload := map[string]interface{}{
		"token":             token,
		"returnSecureToken": true,
	}
	if tenantID != "" {
		payload["tenantId"] = tenantID
	}

	req, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	resp, err := g.postRequest(fmt.Sprintf(verifyCustomTokenURL, g.apiKey), req)
	if err != nil {
		return "", err
	}
	var respBody struct {
		IDToken string `json:"idToken"`
	}
	if err := json.Unmarshal(resp, &respBody); err != nil {
		return "", err
	}
	return respBody.IDToken, err
}

func (g gipClient) postRequest(url string, req []byte) ([]byte, error) {
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(req))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected http status code: %d", resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}

func newActionCodeSettings() *auth.ActionCodeSettings {
	// TODO example config
	actionCodeSettings := &auth.ActionCodeSettings{
		URL:                   "https://www.example.com/checkout?cartId=1234",
		HandleCodeInApp:       true,
		IOSBundleID:           "com.example.ios",
		AndroidPackageName:    "com.example.android",
		AndroidInstallApp:     true,
		AndroidMinimumVersion: "12",
		DynamicLinkDomain:     "coolapp.page.link",
	}
	return actionCodeSettings
}
