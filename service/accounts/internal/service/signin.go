package service

import (
	"errors"
	"github.com/sirupsen/logrus"
)

// Signin is a service method which handles the logic of user login
func (s *accountServiceHandler) Signin(email string) error {
	if len(email) == 0 {
		logrus.Errorf("email doesn't met required validation criteria")
		return errors.New("email doesn't met required validation criteria")
	}
	return nil

}
