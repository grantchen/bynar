/**
    @author: dongjs
    @date: 2023/9/12
    @description:
**/

package service

import (
	"context"
	"fmt"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/accounts/internal/model"
	sql_db "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/errors"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/utils"
	"github.com/sirupsen/logrus"
	"mime/multipart"
	"os"
	"path"
	"strings"
)

var profilePictureType = []string{"png", "jpg", "jpeg"}

// UploadFileToGCS upload user's profile picture to gcs
func (s *accountServiceHandler) UploadFileToGCS(tenantId, organizationUuid, email string, body *multipart.Reader) (string, error) {
	part, err := body.NextPart()
	if err != nil {
		return "", errors.NewUnknownError("file read error").WithInternal().WithCause(err)
	}
	defer part.Close()

	ext := path.Ext(part.FileName())
	// check file type is jpg or png
	if !utils.IsStringArrayInclude(profilePictureType, strings.ToLower(strings.Split(ext, ".")[1])) {
		return "", errors.NewUnknownError("profile picture type is not png or jpg")
	}
	organization, err := s.ar.GetOrganizationDetail(organizationUuid)
	if err != nil || organization == nil {
		return "", errors.NewUnknownError("organization not found").WithInternal().WithCause(err)
	}
	user, err := s.ar.GetUserAccountDetail(email)
	if err != nil || organization == nil {
		return "", errors.NewUnknownError("user not found").WithInternal().WithCause(err)
	}
	filePath := fmt.Sprintf("%v/profile_picture/%v%v", organization.ID, user.ID, ext)
	filePathPrefix := fmt.Sprintf("%v/profile_picture/%v", organization.ID, user.ID)

	//delete existing file from Google cloud storage
	if err = s.cloudStorageProvider.DeleteFiles(filePathPrefix); err != nil {
		return "", errors.NewUnknownError("upload file error").WithInternal().WithCause(err)
	}
	// upload profile picture to Google cloud storage
	url, err := s.cloudStorageProvider.UploadFile(filePath, part)
	if err != nil {
		return "", errors.NewUnknownError("upload file error").WithInternal().WithCause(err)
	}
	// update database profile_photo column in table users
	if err = s.UpdateProfilePhotoOfUsers(tenantId, organizationUuid, email, url); err != nil {
		return "", errors.NewUnknownError("upload file error").WithInternal().WithCause(err)
	}
	return url, nil
}

// DeleteFileFromGCS delete user's profile picture from google cloud storage
func (s *accountServiceHandler) DeleteFileFromGCS(tenantId, organizationUuid, email string) error {
	organization, err := s.ar.GetOrganizationDetail(organizationUuid)
	if err != nil || organization == nil {
		return errors.NewUnknownError("organization not found").WithInternal().WithCause(err)
	}
	user, err := s.ar.GetUserAccountDetail(email)
	if err != nil || organization == nil {
		return errors.NewUnknownError("user not found").WithInternal().WithCause(err)
	}
	filePathPrefix := fmt.Sprintf("%v/profile_picture/%v", organization.ID, user.ID)
	err = s.cloudStorageProvider.DeleteFiles(filePathPrefix)
	if err != nil {
		return errors.NewUnknownError("delete file fail").WithInternal().WithCause(err)
	}
	// update database profile_photo column in table users
	if err = s.UpdateProfilePhotoOfUsers(tenantId, organizationUuid, email, ""); err != nil {
		return errors.NewUnknownError("delete file fail").WithInternal().WithCause(err)
	}
	return nil
}

// GetUserDetail get user details from organization_schema(uuid)
func (s *accountServiceHandler) GetUserDetail(tenantUuid, organizationUuid, email string) (*model.User, error) {
	if len(os.Getenv(tenantUuid)) == 0 {
		return nil, errors.NewUnknownError("no mysql conn environment of " + tenantUuid)
	}
	envs := strings.Split(os.Getenv(tenantUuid), "/")
	connStr := envs[0] + "/" + organizationUuid
	if len(envs) > 1 {
		connStr += envs[1]
	}
	db, err := sql_db.InitializeConnection(connStr)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	var querySql = `select a.id,
       a.email,
       coalesce(a.full_name,''),
       coalesce(a.phone,''),
       a.status,
       coalesce(a.language_preference,''),
       coalesce(a.policy_id,0),
       coalesce(a.theme,''),
       coalesce(a.profile_photo,'')
		from users a 
		where a.email = ? and status = ? limit 1`
	var user = model.User{}
	err = db.QueryRow(querySql, email, true).Scan(
		&user.ID, &user.Email, &user.FullName, &user.Phone, &user.Status,
		&user.LanguagePreference, &user.PolicyId, &user.Theme, &user.ProfilePhoto)
	if err != nil {
		return nil, fmt.Errorf("query row: [%w]", err)
	}
	return &user, nil
}

// UpdateProfilePhotoOfUsers update column profile_photo in table users of organization_schema(uuid)
func (s *accountServiceHandler) UpdateProfilePhotoOfUsers(tenantUuid, organizationUuid string, email string, profilePhoto string) error {
	if len(os.Getenv(tenantUuid)) == 0 {
		return errors.NewUnknownError("no mysql conn environment of " + tenantUuid)
	}
	envs := strings.Split(os.Getenv(tenantUuid), "/")
	connStr := envs[0] + "/" + organizationUuid
	if len(envs) > 1 {
		connStr += envs[1]
	}
	db, err := sql_db.InitializeConnection(connStr)
	logrus.Info("init db ", connStr)
	if err != nil {
		return err
	}
	defer db.Close()
	if _, err = db.Exec(`UPDATE users SET profile_photo = ? WHERE email = ?`, profilePhoto, email); err != nil {
		return err
	}
	return nil
}

// Update user language preference
func (s *accountServiceHandler) UpdateUserLanguagePreference(tenantId, organizationUuid, email, languagePreference string) error {
	// Update the language_preference field in the users table
	if err := s.ar.UpdateUserLanguagePreference(email, languagePreference); err != nil {
		return err
	}

	account, err := s.ar.SelectSignInColumns(email)
	if err != nil || account == nil {
		return err
	}
	claims, err := convertSignInToClaims(account)
	if err != nil {
		return err
	}
	err = s.authProvider.SetCustomUserClaims(context.Background(), account.Uid, claims)
	if err != nil {
		return err
	}

	return nil
}
