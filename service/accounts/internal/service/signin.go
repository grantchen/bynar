package service

import (
	"context"
	"fmt"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/gip"
)

// Signin is a service method which handles the logic of user login
func (s *accountServiceHandler) Signin(email string) error {
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
	return gip.SendRegistrationEmail(email)
}
