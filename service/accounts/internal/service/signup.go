package service

import (
	"errors"
	"github.com/sirupsen/logrus"
)

// CreateUser is a service method which handles the logic of new user registration
func (s *accountServiceHandler) CreateUser(email string) (string, error) {
	if len(email) == 0 {
		logrus.Errorf("email doesn't met required validation criteria")
		return "", errors.New("email doesn't met required validation criteria")
	}
	return "", s.ar.CreateUser(email)

}

// ConfirmEmail is a service method which confirms the email of new account
func (s *accountServiceHandler) ConfirmEmail(email, code string) (int, error) {
	return 0, nil
}

// ResendVerificationCode is a service method which resend the verification code for email verification
func (s *accountServiceHandler) ResendVerificationCode(email string) error {
	return nil
}

// AddUserDetails is a service method which add contract infomation of new account
func (s *accountServiceHandler) AddUserDetails(fullName, country, address, address2, city, postalCode, state, phone string) error {
	return nil
}

// AddTaxDetails is a service method which add tax infomation of new account
func (s *accountServiceHandler) AddTaxDetails(organization, number, country string) error {
	return nil
}

// AddCreditCard is a service method which validating the user with help of Checkout.com API
func (s *accountServiceHandler) AddCreditCard(number, date, code string) error {
	return nil
}
