package service

import (
	"context"
	"errors"
	"fmt"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/utils"
	"os"
	"strconv"
	"time"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/checkout/models"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/gip"
)

// Signup is a service method which check the account is exist
func (s *accountServiceHandler) Signup(email string) error {
	err := s.ar.CheckUserExists(email)
	if err != nil {
		return err
	}

	// Since unregistered users cannot verify oobcode in google identify platform, we decided to customize the verification
	// We sign mailboxes and timestamps to prevent data from being tampered with and to verify expiration dates
	// A character string to be signed in the following format
	needSignatureString := "email=%s&timestamp=%s"
	needSignatureString = fmt.Sprintf(needSignatureString, email, strconv.FormatInt(time.Now().UnixMilli(), 10))
	// Sign data with a custom key
	signature := utils.HmacSha1Signature(os.Getenv("SIGNUP_CUSTOM_VERIFICATION_KEY"), needSignatureString)
	// Append the signature to continueUrl
	continueUrl := "%s?%s&signature=%s"
	continueUrl = fmt.Sprintf(continueUrl, os.Getenv("SIGNUP_REDIRECT_URL"), needSignatureString, signature)

	return gip.SendRegistrationEmail(email, continueUrl)
}

// ConfirmEmail is a service method which confirms the email of new account
func (s *accountServiceHandler) ConfirmEmail(email, timestamp, signature string) (int, error) {
	// Since unregistered users cannot verify oobcode in google identify platform, we decided to customize the verification
	// We sign mailboxes and timestamps to prevent data from being tampered with and to verify expiration dates
	// A character string to be signed in the following format
	needSignatureString := "email=%s&timestamp=%s"
	needSignatureString = fmt.Sprintf(needSignatureString, email, timestamp)
	// Sign data with a custom key
	nowSignature := utils.HmacSha1Signature(os.Getenv("SIGNUP_CUSTOM_VERIFICATION_KEY"), needSignatureString)

	// Verify whether the signatures are consistent
	if signature != nowSignature {
		return 0, errors.New("wrong signature")
	}

	// Verify that the timestamp is expired. The expiration time is 5 minutes
	timestampInt64, _ := strconv.ParseInt(timestamp, 10, 64)
	if (time.Now().UnixMilli()-timestampInt64)/1000 > 60*5 {
		return 0, errors.New("the timestamp has expired")
	}

	return 0, nil
}

// VerifyCard is a service method which verify card of new account
func (s *accountServiceHandler) VerifyCard(token, email, name string) (string, string, error) {
	resp, err := s.paymentProvider.ValidateCard(&models.ValidateCardRequest{Token: token, Email: email, Name: name})
	if err != nil {
		return "", "", err
	}
	return resp.Customer.ID, resp.Source.ID, nil
}

// CreateUser is a service method which handles the logic of new user registration
func (s *accountServiceHandler) CreateUser(email, code, sign, token, fullName, country, addressLine, addressLine2, city, postalCode, state, phoneNumber, organizationName, vat, organisationCountry, customerID, sourceID string) (string, error) {
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
	return uid, s.ar.CreateUser(uid, email, fullName, country, addressLine, addressLine2, city, postalCode, state, phoneNumber, organizationName, vat, organisationCountry, customerID, sourceID)
}
