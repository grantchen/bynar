package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/utils"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/checkout/models"
	errpkg "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/errors"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/gip"
)

// Signup is a service method which check the account is exist
func (s *accountServiceHandler) Signup(email string) *errpkg.Error {
	exists, err := s.authProvider.IsUserExists(context.Background(), email)
	if err != nil {
		return errpkg.NewUnknownError("check user exists fail", "").WithInternalCause(err)
	}
	if exists {
		return errpkg.NewUnknownError(fmt.Sprintf("email %s has already exist", email), errpkg.ErrCodeEmailAlreadyExists).WithInternalCause(err)
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

	err = gip.SendRegistrationEmail(email, continueUrl)
	if err != nil {
		return errpkg.NewUnknownError("send registration email fail", "").WithInternalCause(err)
	}
	return nil
}

// ConfirmEmail is a service method which confirms the email of new account
func (s *accountServiceHandler) ConfirmEmail(email, timestamp, signature string) (int, *errpkg.Error) {
	// Since unregistered users cannot verify oobcode in google identify platform, we decided to customize the verification
	// We sign mailboxes and timestamps to prevent data from being tampered with and to verify expiration dates
	// A character string to be signed in the following format
	needSignatureString := "email=%s&timestamp=%s"
	needSignatureString = fmt.Sprintf(needSignatureString, email, timestamp)
	// Sign data with a custom key
	nowSignature := utils.HmacSha1Signature(os.Getenv("SIGNUP_CUSTOM_VERIFICATION_KEY"), needSignatureString)

	// Verify whether the signatures are consistent
	if signature != nowSignature {
		return 0, errpkg.NewUnknownError("wrong signature", errpkg.ErrCodeSignatureInvalid)
	}

	// Verify that the timestamp is expired. The expiration time is 5 minutes
	timestampInt64, _ := strconv.ParseInt(timestamp, 10, 64)
	if (time.Now().UnixMilli()-timestampInt64)/1000 > 60*5 {
		return 0, errpkg.NewUnknownError("timestamp has expired", errpkg.ErrCodeTimestampExpired)
	}

	return 0, nil
}

// VerifyCard is a service method which verify card of new account
func (s *accountServiceHandler) VerifyCard(token, email, name string) (string, string, *errpkg.Error) {
	// Use checkout.com service to validate card
	resp, err := s.paymentProvider.ValidateCard(&models.ValidateCardRequest{Token: token, Email: email, Name: name})
	if err != nil {
		return "", "", errpkg.NewUnknownError(fmt.Sprintf("verify card failed: %s", err.Error()), "").WithInternal().WithCause(err)
	}
	return resp.Customer.ID, resp.Source.ID, nil
}

// CreateUser is a service method which handles the logic of new user registration
func (s *accountServiceHandler) CreateUser(email, _, _, _, fullName, country, addressLine, addressLine2, city, postalCode, state, phoneNumber, organizationName, vat, organisationCountry, customerID, sourceID, tenantCode string) (string, *errpkg.Error) {
	// recheck user exist
	exist, err := s.authProvider.IsUserExists(context.Background(), email)
	if err != nil {
		return "", errpkg.NewUnknownError("", "").WithInternalCause(err)
	}
	if exist {
		return "", errpkg.NewUnknownError(fmt.Sprintf("email %s has already exist", email), errpkg.ErrCodeEmailAlreadyExists).WithInternalCause(err)
	}
	// revalidate email
	// err = gip.VerificationEmail(signature)
	// if err != nil {
	// 	return "", err
	// }
	// check user exists in gip
	ok, err := s.authProvider.IsUserExists(context.TODO(), email)
	if err != nil && !errors.Is(err, gip.ErrUserNotFound) {
		return "", errpkg.NewUnknownError("check user exists failed", errpkg.ErrCodeNoUserFound).WithInternal().WithCause(err)
	}
	if ok || errors.Is(err, gip.ErrUserNotFound) {
		// If user exists in gip. delete it
		err = s.authProvider.DeleteUserByEmail(context.TODO(), email)
		if err != nil && !errors.Is(err, gip.ErrUserNotFound) {
			return "", errpkg.NewUnknownError("delete user failed", "").WithInternal().WithCause(err)
		}
	}
	// create use in gip
	uid, err := s.authProvider.CreateUser(context.TODO(), email, fullName, phoneNumber, false)
	if err != nil {
		return "", errpkg.NewUnknownError("create user failed: "+err.Error(), "").WithInternal().WithCause(err)
	}
	customClaims := map[string]interface{}{
		"country": organisationCountry,
		"status":  1,
	}
	// update custom user info in gip
	err = s.authProvider.UpdateUser(context.TODO(), uid, map[string]interface{}{"customClaims": customClaims})
	if err != nil {
		return "", errpkg.NewUnknownError("update user failed: "+err.Error(), "").WithInternal().WithCause(err)
	}
	// create user in db
	code, err := s.ar.CreateUser(uid, email, fullName, country, addressLine, addressLine2, city, postalCode, state, phoneNumber, organizationName, vat, organisationCountry, customerID, sourceID, tenantCode)
	if err != nil {
		if code == 0 {
			err = s.authProvider.DeleteUserByEmail(context.Background(), email)
			if err != nil {
				logrus.Error("delete user failed: ", err.Error())
			}
		}
		return "", errpkg.NewUnknownError("create user failed", "").WithInternal().WithCause(err)
	}
	// return idToken after created
	account, err := s.ar.SelectSignInColumns(uid)
	if err != nil {
		return "", errpkg.NewUnknownError("no user found", errpkg.ErrCodeNoUserFound).WithInternal().WithCause(err)
	}
	data, err := json.Marshal(&account)
	if err != nil {
		return "", errpkg.NewUnknownError("no user found", errpkg.ErrCode).WithInternal().WithCause(err)
	}
	claims := map[string]interface{}{}
	err = json.Unmarshal(data, &claims)
	if err != nil {
		return "", errpkg.NewUnknownError("no user found", errpkg.ErrCode).WithInternal().WithCause(err)
	}
	customToken, err := s.authProvider.CustomTokenWithClaims(context.Background(), uid, claims)
	if err != nil {
		return "", errpkg.NewUnknownError("set custom token with claims fail", "").WithInternal().WithCause(err)
	}
	return customToken, nil
}
