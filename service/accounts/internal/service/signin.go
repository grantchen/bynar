package service

import (
	"context"
	"fmt"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/gip"
)

// SignIn is a service method which handles the logic of user login
func (s *accountServiceHandler) SignIn(email, oobCode string) (idToken string, err error) {
	if err := s.VerifyEmail(email); err != nil {
		return "", err
	}
	// todo call google identify platform api to check email and oobCode
	account, err := s.ar.SelectAccount(email)
	if err != nil || account == nil {
		return "", fmt.Errorf("SignIn: %s no user selected", email)
	}

	claims := map[string]interface{}{
		"uid":                  account.Uid,
		"organization_account": true,
		"organization_user_id": account.OrganizationUserId,
		"organization_status":  account.OrganizationStatus,
		"tenant_uuid":          account.TenantUuid,
		"organization_uuid":    account.OrganizationUuid,
	}
	token, err := s.authProvider.SignIn(context.Background(), account.Uid, claims)
	if err != nil {
		return "", fmt.Errorf("SignIn: %s generate idtoke err: %+v", email, err)
	}
	return token, nil

}

// SendSignInEmail send google identify platform with oobCode for sign in
func (s *accountServiceHandler) SendSignInEmail(email string) error {
	if err := s.VerifyEmail(email); err != nil {
		return err
	}
	return gip.SendRegistrationEmail(email)
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
