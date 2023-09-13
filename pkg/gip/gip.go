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

const (
	ENVGoogleApplicationCredentials = "GOOGLE_APPLICATION_CREDENTIALS"
	ENVGoogleAPIKey                 = "GOOGLE_API_KEY"
)

var ErrUserNotFound = errors.New("user not found")
var ErrIDTokenInvalid = errors.New("id token invalid")

// gipClient is the interface for the AuthProvider.
type gipClient struct {
	app    *firebase.App
	apiKey string
}

// NewGIPClient creates a new instance of the AuthProvider.
func NewGIPClient() (AuthProvider, error) {
	opt := option.WithCredentialsJSON([]byte(os.Getenv(ENVGoogleApplicationCredentials)))
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, errors.New("error initializing gip app: " + err.Error())
	}

	return &gipClient{
		app:    app,
		apiKey: os.Getenv(ENVGoogleAPIKey),
	}, nil
}

// IsUserExists checks if the user with the specified email exists.
func (g gipClient) IsUserExists(ctx context.Context, email string) (bool, error) {
	client, err := g.app.Auth(ctx)
	if err != nil {
		return false, fmt.Errorf("error getting Auth client: %v", err)
	}

	u, err := client.GetUserByEmail(ctx, email)
	if err != nil {
		if auth.IsUserNotFound(err) {
			return false, nil
		}
		return false, fmt.Errorf("error getting user by email %s: %v", email, err)
	}

	return u != nil, nil
}

// CreateUser creates a new user with the specified properties.
func (g gipClient) CreateUser(ctx context.Context, email, displayName, phoneNumber string) (uid string, err error) {
	client, err := g.app.Auth(ctx)
	if err != nil {
		return "", fmt.Errorf("error getting Auth client: %v", err)
	}

	params := (&auth.UserToCreate{}).
		Email(email).
		DisplayName(displayName).
		PhoneNumber(phoneNumber).
		EmailVerified(false).
		Password(utils.RandString(10)).
		Disabled(false)
	u, err := client.CreateUser(ctx, params)
	if err != nil {
		return "", fmt.Errorf("error creating user: %v", err)
	}

	return u.UID, nil
}

// UpdateUser updates an existing user account with the specified properties.
func (g gipClient) UpdateUser(ctx context.Context, uid string, params map[string]interface{}) error {
	client, err := g.app.Auth(ctx)
	if err != nil {
		return fmt.Errorf("error getting Auth client: %v", err)
	}

	updateParams := &auth.UserToUpdate{}
	if email, ok := params["email"].(string); ok {
		updateParams.Email(email)
	}

	if displayName, ok := params["displayName"].(string); ok {
		updateParams.DisplayName(displayName)
	}

	if phoneNumber, ok := params["phoneNumber"].(string); ok {
		updateParams.PhoneNumber(phoneNumber)
	}

	if disableUser, ok := params["disableUser"].(bool); ok {
		updateParams.Disabled(disableUser)
	}

	if customClaims, ok := params["customClaims"].(map[string]interface{}); ok {
		updateParams.CustomClaims(customClaims)
	}

	_, err = client.UpdateUser(ctx, uid, updateParams)
	if err != nil {
		if auth.IsUserNotFound(err) {
			return ErrUserNotFound
		}
		return fmt.Errorf("error updating user: %v", err)
	}

	return nil
}

// DeleteUser deletes the user by the given UID.
func (g gipClient) DeleteUser(ctx context.Context, uid string) error {
	client, err := g.app.Auth(ctx)
	if err != nil {
		return fmt.Errorf("error getting Auth client: %v", err)
	}

	err = client.DeleteUser(ctx, uid)
	if err != nil {
		if auth.IsUserNotFound(err) {
			return ErrUserNotFound
		}
		return fmt.Errorf("error deleting user: %v", err)
	}

	return nil
}

// DeleteUserByEmail deletes the user by the given email.
func (g gipClient) DeleteUserByEmail(ctx context.Context, email string) error {
	client, err := g.app.Auth(ctx)
	if err != nil {
		return fmt.Errorf("error getting Auth client: %v", err)
	}

	u, err := client.GetUserByEmail(ctx, email)
	if err != nil {
		if auth.IsUserNotFound(err) {
			return ErrUserNotFound
		}
		return fmt.Errorf("error getting user by email %s: %v", email, err)
	}

	err = client.DeleteUser(ctx, u.UID)
	if err != nil {
		if auth.IsUserNotFound(err) {
			return ErrUserNotFound
		}
		return fmt.Errorf("error deleting user: %v", err)
	}

	return nil
}

// SignIn signs in the user by the given UID.
func (g gipClient) SignIn(ctx context.Context, uid string, devClaims map[string]interface{}) (idToken string, err error) {
	token, err := g.CustomTokenWithClaims(ctx, uid, devClaims)
	if err != nil {
		return "", fmt.Errorf("error creating custom token with claims: %v", err)
	}

	idToken, err = g.signInWithCustomToken(token)
	if err != nil {
		return "", fmt.Errorf("error signing in with custom token: %v", err)
	}

	return idToken, nil
}

// LogOut logs out the user by the given UID.
func (g gipClient) LogOut(ctx context.Context, uid string) error {
	client, err := g.app.Auth(ctx)
	if err != nil {
		return fmt.Errorf("error getting Auth client: %v", err)
	}

	err = client.RevokeRefreshTokens(ctx, uid)
	if err != nil {
		return fmt.Errorf("error revoking token: %v", err)
	}

	return nil
}

// VerifyIDToken verifies the signature	and payload of the provided ID token.
func (g gipClient) VerifyIDToken(ctx context.Context, idToken string) (claims map[string]interface{}, err error) {
	client, err := g.app.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting Auth client: %v", err)
	}

	token, err := client.VerifyIDToken(ctx, idToken)
	if err != nil {
		if auth.IsIDTokenInvalid(err) {
			return nil, ErrIDTokenInvalid
		}
		return nil, fmt.Errorf("error verifying id token: %v", err)
	}

	token.Claims["auth_time"] = token.AuthTime
	token.Claims["iss"] = token.Issuer
	token.Claims["aud"] = token.Audience
	token.Claims["exp"] = token.Expires
	token.Claims["iat"] = token.IssuedAt
	token.Claims["sub"] = token.Subject
	token.Claims["uid"] = token.UID
	return token.Claims, nil
}

// VerifyIDTokenAndCheckRevoked verifies the provided ID token, and additionally checks that the
// token has not been revoked or disabled.
func (g gipClient) VerifyIDTokenAndCheckRevoked(ctx context.Context, idToken string) (claims map[string]interface{}, err error) {
	client, err := g.app.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting Auth client: %v", err)
	}

	token, err := client.VerifyIDTokenAndCheckRevoked(ctx, idToken)
	if err != nil {
		if auth.IsIDTokenInvalid(err) {
			return nil, ErrIDTokenInvalid
		}
		return nil, fmt.Errorf("error verifying id token and check revoked: %v", err)
	}

	token.Claims["auth_time"] = token.AuthTime
	token.Claims["iss"] = token.Issuer
	token.Claims["aud"] = token.Audience
	token.Claims["exp"] = token.Expires
	token.Claims["iat"] = token.IssuedAt
	token.Claims["sub"] = token.Subject
	token.Claims["uid"] = token.UID
	return token.Claims, nil
}

// signInWithCustomTokenForTenant signs in using a custom token and tenant ID.
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

// signInWithCustomToken signs in using a custom token.
func (g gipClient) signInWithCustomToken(token string) (string, error) {
	return g.signInWithCustomTokenForTenant(token, "")
}

// CustomTokenWithClaims creates a signed custom authentication token with the specified user ID.
func (g gipClient) CustomTokenWithClaims(ctx context.Context, uid string, devClaims map[string]interface{}) (string, error) {
	client, err := g.app.Auth(ctx)
	if err != nil {
		return "", err
	}

	return client.CustomTokenWithClaims(ctx, uid, devClaims)
}

// postRequest sends a POST request to the specified URL with the specified request body.
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

// SetCustomUserClaims sets additional claims on an existing user account.
func (g gipClient) SetCustomUserClaims(ctx context.Context, uid string, customClaims map[string]interface{}) error {
	client, err := g.app.Auth(ctx)
	if err != nil {
		return err
	}
	return client.SetCustomUserClaims(ctx, uid, customClaims)
}
