package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	"firebase.google.com/go/v4/auth"

	"github.com/sirupsen/logrus"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/accounts/internal/model"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/errors"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/gip"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/models"
)

// SignIn is a service method which handles the logic of user login
func (s *accountServiceHandler) SignIn(email, oobCode string) (string, *errors.Error) {
	var exists = false
	var err error
	if exists, err = s.authProvider.IsUserExists(context.Background(), email); err != nil {
		return "", errors.NewUnknownError("sign in fail", "").WithInternalCause(err)
	}
	if !exists {
		return "", errors.NewUnknownError("sign in fail: email not sign up", errors.ErrCodeNoUserFound)
	}
	account, err := s.ar.SelectSignInColumns(email)
	if err != nil || account == nil {
		return "", errors.NewUnknownError("sign in fail", errors.ErrCodeNoUserFound).WithInternalCause(err)
	}
	err = gip.SignInWithEmailLink(email, oobCode)
	if err != nil {
		return "", errors.NewUnknownError("sign in fail", "").WithInternalCause(err)
	}
	claims, err := convertSignInToClaims(account)
	if err != nil {
		return "", errors.NewUnknownError("sign in fail", "").WithInternalCause(err)
	}
	token, err := s.authProvider.SignIn(context.Background(), account.Uid, claims)
	if err != nil {
		return "", errors.NewUnknownError("sign in fail", "").WithInternalCause(err)
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
func (s *accountServiceHandler) SendSignInEmail(email string) *errors.Error {
	var (
		userRecord *auth.UserRecord = nil
		err        error            = nil
	)
	if userRecord, err = s.authProvider.GetUserByEmail(context.Background(), email); err != nil || userRecord == nil {
		return errors.NewUnknownError("email not signed up", "").WithInternalCause(err)
	}

	account, err := s.ar.SelectSignInColumns(userRecord.UserInfo.UID)
	if err != nil || account == nil {
		return errors.NewUnknownError("no user found", errors.ErrCodeNoUserFound).WithMetadata(map[string]string{"uid": userRecord.UserInfo.UID}).WithInternalCause(err)
	}
	claims, err := convertSignInToClaims(account)
	if err != nil {
		return errors.NewUnknownError("send email fail", errors.ErrCode).WithInternalCause(err)
	}
	err = s.authProvider.SetCustomUserClaims(context.Background(), account.Uid, claims)
	if err != nil {
		return errors.NewUnknownError("set custom user claims fail", "").WithInternalCause(err)
	}
	if err = gip.SendRegistrationEmail(email, fmt.Sprintf("%s?email=%s", os.Getenv("SIGNIN_REDIRECT_URL"), email)); err != nil {
		return errors.NewUnknownError("send email fail", "").WithInternalCause(err)
	}
	return nil
}

// VerifyEmail check email is stored in db and google identify platform
func (s *accountServiceHandler) VerifyEmail(email string) *errors.Error {
	exists, err := s.authProvider.IsUserExists(context.Background(), email)
	if err != nil {
		return errors.NewUnknownError("send email fail", "").WithInternalCause(err)
	}
	if !exists {
		return errors.NewUnknownError(fmt.Sprintf("account with email: %s has not signup", email), errors.ErrCodeNoUserFound).WithInternalCause(err)
	}
	return nil
}

// GetUserDetails after signing get user info
func (s *accountServiceHandler) GetUserDetails(db *sql.DB, uid string, userId int) (*model.GetUserResponse, *errors.Error) {
	account, err := s.ar.GetUserAccountDetail(uid)
	var userResponse = model.GetUserResponse{}
	if err == nil && account != nil {
		userResponse.Email = account.Email.String
		userResponse.FullName = account.FullName.String
		userResponse.Country = account.Country.String
		userResponse.AddressLine = account.Address.String
		userResponse.AddressLine2 = account.Address2.String
		userResponse.City = account.City.String
		userResponse.PostalCode = account.PostalCode.String
		userResponse.State = account.State.String
		userResponse.PhoneNumber = account.Phone.String
	}
	user, err := s.ar.GetUserDetail(db, userId)
	if err != nil {
		return nil, errors.NewUnknownError("user not found", errors.ErrCodeNoUserFound).WithInternalCause(err)
	}
	userResponse.ID = user.ID
	userResponse.Email = user.Email
	userResponse.FullName = user.FullName
	userResponse.LanguagePreference = user.LanguagePreference
	userResponse.ThemePreference = user.Theme
	userResponse.ProfileURL = user.ProfilePhoto
	policy := models.Policy{Services: make([]models.ServicePolicy, 0)}
	err = json.Unmarshal([]byte(user.Policies), &policy)
	if err != nil {
		logrus.Error("get policy error", err)
	}
	userResponse.Permissions = policy
	return &userResponse, nil
}
