package gip

import (
	"context"

	"firebase.google.com/go/v4/auth"
)

type AuthProvider interface {
	// IsUserExists checks if the user with the specified email exists.
	IsUserExists(ctx context.Context, email string) (bool, error)
	// CreateUser creates a new user with the specified properties.
	CreateUser(ctx context.Context, email, displayName, phoneNumber string) (uid string, err error)
	// UpdateUser updates an existing user account with the specified properties.
	UpdateUser(ctx context.Context, uid string, params map[string]interface{}) error
	// DeleteUser deletes the user by the given UID.
	DeleteUser(ctx context.Context, uid string) error
	// SignIn signs in the user by the given UID.
	SignIn(ctx context.Context, uid string, devClaims map[string]interface{}) (idToken string, err error)
	// LogOut logs out the user by the given UID.
	LogOut(ctx context.Context, uid string) error
	// VerifyIDToken verifies the signature	and payload of the provided ID token.
	VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error)
	// VerifyIDTokenAndCheckRevoked verifies the provided ID token, and additionally checks that the
	// token has not been revoked or disabled.
	VerifyIDTokenAndCheckRevoked(ctx context.Context, idToken string) (*auth.Token, error)
	// EmailSignInLink generates the out-of-band email action link for email link sign-in flows, using the action
	// code settings provided.
	EmailSignInLink(ctx context.Context, email string) (string, error)
}
