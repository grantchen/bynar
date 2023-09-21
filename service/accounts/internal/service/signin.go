package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/accounts/internal/model"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/errors"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/gip"
	"github.com/sirupsen/logrus"
)

// SignIn is a service method which handles the logic of user login
func (s *accountServiceHandler) SignIn(email, oobCode string) (idToken string, err error) {
	if err = s.VerifyEmail(email); err != nil {
		return "", errors.NewUnknownError("email is not signed up").WithInternal().WithCause(err)
	}
	account, err := s.ar.SelectSignInColumns(email)
	if err != nil || account == nil {
		return "", errors.NewUnknownError("sign in failed").WithInternal().WithCause(err)
	}
	err = gip.SignInWithEmailLink(email, oobCode)
	if err != nil {
		return "", errors.NewUnknownError("sign in failed").WithInternal().WithCause(err)
	}
	claims, err := convertSignInToClaims(account)
	if err != nil {
		return "", errors.NewUnknownError("sign in failed").WithInternal().WithCause(err)
	}
	token, err := s.authProvider.SignIn(context.Background(), account.Uid, claims)
	if err != nil {
		return "", errors.NewUnknownError("sign in failed").WithInternal().WithCause(err)
	}
	return token, nil

}

// covert signIn struct to claims map
func convertSignInToClaims(signIn *model.SignIn) (map[string]interface{}, error) {
	data, err := json.Marshal(&signIn)
	if err != nil {
		return nil, err
	}
	claims := map[string]interface{}{}
	err = json.Unmarshal(data, &claims)
	if err != nil {
		return nil, err
	}
	return claims, nil
}

// SendSignInEmail send google identify platform with oobCode for sign in
func (s *accountServiceHandler) SendSignInEmail(email string) error {
	if err := s.VerifyEmail(email); err != nil {
		return errors.NewUnknownError("email is not signed up").WithInternal().WithCause(err)
	}
	account, err := s.ar.SelectSignInColumns(email)
	if err != nil || account == nil {
		return errors.NewUnknownError("no user found").WithInternal().WithCause(err)
	}
	claims, err := convertSignInToClaims(account)
	if err != nil {
		return errors.NewUnknownError("email sending failed").WithInternal().WithCause(err)
	}
	err = s.authProvider.SetCustomUserClaims(context.Background(), account.Uid, claims)
	if err != nil {
		return errors.NewUnknownError("email sending failed").WithInternal().WithCause(err)
	}
	if err = gip.SendRegistrationEmail(email, fmt.Sprintf("%s?email=%s", os.Getenv("SIGNIN_REDIRECT_URL"), email)); err != nil {
		return errors.NewUnknownError("email sending failed").WithInternal().WithCause(err)
	}
	return nil
}

// VerifyEmail check email is stored in db and google identify platform
func (s *accountServiceHandler) VerifyEmail(email string) error {
	exists, err := s.authProvider.IsUserExists(context.Background(), email)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("account with email: %s has not signup", email)
	}
	return nil
}

// GetUserDetails after signing get user info
func (s *accountServiceHandler) GetUserDetails(db *sql.DB, email string) (*model.GetUserResponse, error) {
	account, err := s.ar.GetUserAccountDetail(email)
	if err != nil {
		return nil, errors.NewUnknownError("account not found").WithInternal().WithCause(err)
	}
	var userResponse = model.GetUserResponse{
		ID:           account.ID,
		Email:        account.Email.String,
		FullName:     account.FullName.String,
		Country:      account.Country.String,
		AddressLine:  account.Address.String,
		AddressLine2: account.Address2.String,
		City:         account.City.String,
		PostalCode:   account.PostalCode.String,
		State:        account.State.String,
		PhoneNumber:  account.Phone.String,
	}
	user, err := s.ar.GetUserDetail(db, email)
	if err != nil {
		return nil, errors.NewUnknownError("user not found").WithInternal().WithCause(err)
	}
	userResponse.LanguagePreference = user.LanguagePreference
	userResponse.ThemePreference = user.Theme
	userResponse.ProfileURL = user.ProfilePhoto
	userResponse.PolicyID = user.PolicyId
	policy, err := s.ar.GetUserPolicy(db, user.PolicyId)
	if err != nil {
		logrus.Error("get policy error", err)
	}
	userResponse.Permissions = policy
	return &userResponse, nil
}
