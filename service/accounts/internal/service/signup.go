package service

import (
	"context"
	"errors"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/checkout/models"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/gip"
)

// Signup is a service method which check the account is exist
func (s *accountServiceHandler) Signup(email string) error {
	err := s.ar.CheckUserExists(email)
	if err != nil {
		return err
	}
	return gip.SendRegistrationEmail(email)
}

// ConfirmEmail is a service method which confirms the email of new account
func (s *accountServiceHandler) ConfirmEmail(email, code string) (int, error) {
	err := gip.VerificationEmail(code)
	if err != nil {
		return 0, err
	}
	return 0, nil
}

// VerifyCard is a service method which verify card of new account
func (s *accountServiceHandler) VerifyCard(token, email, name string) error {
	_, err := s.paymentProvider.ValidateCard(&models.ValidateCardRequest{Token: token, Email: email, Name: name})
	if err != nil {
		return err
	}
	return nil
}

// // ResendVerificationCode is a service method which resend the verification code for email verification
// func (s *accountServiceHandler) ResendVerificationCode(email string) error {
// 	return nil
// }

// // AddUserDetails is a service method which add contract infomation of new account
// func (s *accountServiceHandler) AddUserDetails(fullName, country, address, address2, city, postalCode, state, phone string) error {
// 	return nil
// }

// // AddTaxDetails is a service method which add tax infomation of new account
// func (s *accountServiceHandler) AddTaxDetails(organization, number, country string) error {
// 	return nil
// }

// // AddCreditCard is a service method which validating the user with help of Checkout.com API
// func (s *accountServiceHandler) AddCreditCard(number, date, code string) error {
// 	return nil
// }

// CreateUser is a service method which handles the logic of new user registration
func (s *accountServiceHandler) CreateUser(email, code, sign, token, fullName, country, addressLine, addressLine2, city, postalCode, state, phoneNumber, organizationName, vat, organisationCountry string) (string, error) {
	// recheck user exist
	err := s.ar.CheckUserExists(email)
	if err != nil {
		return "", err
	}
	// revalidate email
	err = gip.VerificationEmail(code)
	if err != nil {
		return "", err
	}
	// revalidate card
	cardResp, err := s.paymentProvider.ValidateCard(&models.ValidateCardRequest{Token: token, Email: email, Name: fullName})
	if err != nil {
		return "", err
	}
	// create user in gip
	ok, err := s.authProvider.IsUserExists(context.TODO(), email)
	if err != nil && !errors.Is(err, gip.ErrUserNotFound) {
		return "", err
	}
	if ok || errors.Is(err, gip.ErrUserNotFound) {
		err = s.authProvider.DeleteUserByEmail(context.TODO(), email)
		if err != nil && !errors.Is(err, gip.ErrUserNotFound) {
			return "", err
		}
	}
	uid, err := s.authProvider.CreateUser(context.TODO(), email, fullName, phoneNumber)
	if err != nil {
		return uid, err
	}
	customClaims := map[string]interface{}{
		"country": organisationCountry,
	}
	err = s.authProvider.UpdateUser(context.TODO(), uid, map[string]interface{}{"customClaims": customClaims})
	if err != nil {
		return uid, err
	}
	// create user in db
	return uid, s.ar.CreateUser(uid, email, fullName, country, addressLine, addressLine2, city, postalCode, state, phoneNumber, organizationName, vat, organisationCountry, cardResp.Customer.ID, cardResp.Source.ID)
}
