package gip

import (
	"context"
	"firebase.google.com/go/v4/auth"
)

// AuthProvider is an interface for the authentication provider.
type AuthProvider interface {
	// IsUserExists checks if the user with the specified email exists.
	IsUserExists(ctx context.Context, email string) (bool, error)
	// CreateUser creates a new user with the specified properties.
	CreateUser(ctx context.Context, email, displayName, phoneNumber string, disabled bool) (uid string, err error)
	// UpdateUser updates an existing user account with the specified properties.
	UpdateUser(ctx context.Context, uid string, params map[string]interface{}) error
	// UpdateUserByEmail updates an existing user account with the specified properties.
	UpdateUserByEmail(ctx context.Context, email string, params map[string]interface{}) error
	// DeleteUser deletes the user by the given UID.
	DeleteUser(ctx context.Context, uid string) error
	// DeleteUserByEmail deletes the user by the given email.
	DeleteUserByEmail(ctx context.Context, email string) error
	// SignIn signs in the user by the given UID.
	SignIn(ctx context.Context, uid string, devClaims map[string]interface{}) (idToken string, err error)
	// LogOut logs out the user by the given UID.
	LogOut(ctx context.Context, uid string) error
	// VerifyIDToken verifies the signature	and payload of the provided ID token.
	VerifyIDToken(ctx context.Context, idToken string) (claims map[string]interface{}, err error)
	// VerifyIDTokenAndCheckRevoked verifies the provided ID token, and additionally checks that the
	// token has not been revoked or disabled.
	VerifyIDTokenAndCheckRevoked(ctx context.Context, idToken string) (claims map[string]interface{}, err error)
	// CustomTokenWithClaims creates a signed custom authentication token with the specified user ID.
	CustomTokenWithClaims(ctx context.Context, uid string, devClaims map[string]interface{}) (string, error)
	//SetCustomUserClaims sets additional claims on an existing user account.
	SetCustomUserClaims(ctx context.Context, uid string, customClaims map[string]interface{}) error
	// GetUserByEmail get user info by email from google identify platform
	GetUserByEmail(ctx context.Context, email string) (*auth.UserRecord, error)
	// GetUser get user info by uid from google identify platform
	GetUser(ctx context.Context, uid string) (*auth.UserRecord, error)
}
