package service

import (
	"context"
	"encoding/json"
	"fmt"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/accounts/internal/model"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/errors"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/gip"
	"os"
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
		return errors.NewUnknownError("user no fund").WithInternal().WithCause(err)
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
	err := s.ar.CheckUserExists(email)
	if err == nil {
		return fmt.Errorf("account with email: %s has not signup", email)
	}
	exists, err := s.authProvider.IsUserExists(context.Background(), email)
	if err != nil {
		return err
	}
	if exists == false {
		return fmt.Errorf("account with email: %s has not signup", email)
	}
	return nil
}

// GetUserByUid after signin get user info
func (s *accountServiceHandler) GetUserByUid(uid string) (*model.GetUserResponse, error) {
	user, err := s.ar.SelectUserByUid(uid)
	if err != nil {
		return nil, err
	}
	return user, nil
}
