package gip

import (
	"context"

	"firebase.google.com/go/v4/auth"
)

// TODO add CheckUserExists and remove return *auth.UserRecord
type AuthProvider interface {
	// GetUser gets the user data corresponding to the specified user ID.
	GetUser(ctx context.Context, uid string) (*auth.UserRecord, error)
	// GetUserByEmail gets the user data corresponding to the specified email.
	GetUserByEmail(ctx context.Context, email string) (*auth.UserRecord, error)
	// CreateUser creates a new user with the specified properties.
	CreateUser(ctx context.Context, email string) (*auth.UserRecord, error)
	// UpdateUser updates an existing user account with the specified properties.
	UpdateUser(ctx context.Context, uid string, params map[string]interface{}) (*auth.UserRecord, error)
	// DeleteUser deletes the user by the given UID.
	DeleteUser(ctx context.Context, uid string) error

	// SignIn signs in with the provided token.
	SignIn(ctx context.Context, token string) (idToken string, err error)
	// LogOut logs out the user by the given UID.
	LogOut(ctx context.Context, uid string) error
	// EmailSignInLink generates the out-of-band email action link for email link sign-in flows, using the action
	// code settings provided.
	EmailSignInLink(ctx context.Context, email string) (string, error)
}
