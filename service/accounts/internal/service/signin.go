package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/accounts/internal/model"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/gip"
	"github.com/sirupsen/logrus"
	"os"
)

// SignIn is a service method which handles the logic of user login
func (s *accountServiceHandler) SignIn(email, oobCode string) (idToken string, err error) {
	if err = s.VerifyEmail(email); err != nil {
		logrus.Errorf("SignIn: verify email error: %+v", err)
		return "", fmt.Errorf("email is not signed up")
	}
	account, err := s.ar.SelectSignInColumns(email)
	if err != nil || account == nil {
		logrus.Errorf("SignIn: %s no user selected", email)
		return "", fmt.Errorf("sign in failed")
	}
	err = gip.SignInWithEmailLink(email, oobCode)
	if err != nil {
		logrus.Errorf("SignIn: verification oobCode err: %+v", err)
		return "", fmt.Errorf("sign in failed")
	}
	claims, err := convertSignInToClaims(account)
	if err != nil {
		logrus.Errorf("SignIn: convert sign in claims err: %+v", err)
		return "", fmt.Errorf("sign in failed")
	}
	token, err := s.authProvider.SignIn(context.Background(), account.Uid, claims)
	if err != nil {
		logrus.Errorf("SignIn: %s generate idtoke err: %+v", email, err)
		return "", fmt.Errorf("sign in failed")
	}
	return token, nil

}

// covert signIn struct to claims map
func convertSignInToClaims(signIn *model.SignIn) (map[string]interface{}, error) {
	data, err := json.Marshal(&signIn)
	if err != nil {
		return nil, fmt.Errorf("convertSignInToClaims: marshal singin to []byte err: %+v", err)
	}
	claims := map[string]interface{}{}
	err = json.Unmarshal(data, &claims)
	if err != nil {
		return nil, fmt.Errorf("convertSignInToClaims: unmarshal []byte to claims map err: %+v", err)
	}
	return claims, nil
}

// SendSignInEmail send google identify platform with oobCode for sign in
func (s *accountServiceHandler) SendSignInEmail(email string) error {
	if err := s.VerifyEmail(email); err != nil {
		logrus.Errorf("SendSignInEmail: verify email error: %+v", err)
		return errors.New("email is not signed up")
	}
	if err := gip.SendRegistrationEmail(email, os.Getenv("SIGNIN_REDIRECT_URL")); err != nil {
		logrus.Errorf("SendSignInEmail: send registration email error: %+v", err)
		return errors.New("email sending failed")
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
