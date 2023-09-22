/**
    @author: dongjs
    @date: 2023/9/12
    @description:
**/

package service

import (
	"context"
	"database/sql"
	"fmt"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/errors"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/utils"
	"mime/multipart"
	"path"
	"strings"
)

var profilePictureType = []string{"png", "jpg", "jpeg"}

// UploadFileToGCS upload user's profile picture to gcs
func (s *accountServiceHandler) UploadFileToGCS(db *sql.DB, organizationUuid string, userId int, body *multipart.Reader) (string, error) {
	part, err := body.NextPart()
	if err != nil {
		return "", errors.NewUnknownError("file read fail").WithInternalCause(err)
	}
	defer part.Close()

	ext := path.Ext(part.FileName())
	// check file type is jpg or png
	if !utils.IsStringArrayInclude(profilePictureType, strings.ToLower(strings.Split(ext, ".")[1])) {
		return "", errors.NewUnknownError("profile picture type is not png or jpg")
	}
	organization, err := s.ar.GetOrganizationDetail(organizationUuid)
	if err != nil || organization == nil {
		return "", errors.NewUnknownError("organization not found").WithInternalCause(err)
	}

	filePath := fmt.Sprintf("%v/profile_picture/%v%v", organization.ID, userId, ext)
	filePathPrefix := fmt.Sprintf("%v/profile_picture/%v", organization.ID, userId)

	//delete existing file from Google cloud storage
	if err = s.cloudStorageProvider.DeleteFiles(filePathPrefix); err != nil {
		return "", errors.NewUnknownError("upload file fail").WithInternalCause(err)
	}
	// upload profile picture to Google cloud storage
	url, err := s.cloudStorageProvider.UploadFile(filePath, part)
	if err != nil {
		return "", errors.NewUnknownError("upload file fail").WithInternalCause(err)
	}
	// update database profile_photo column in table users
	if err = s.ar.UpdateProfilePhotoOfUsers(db, userId, url); err != nil {
		return "", errors.NewUnknownError("upload file fail").WithInternalCause(err)
	}
	return url, nil
}

// DeleteFileFromGCS delete user's profile picture from google cloud storage
func (s *accountServiceHandler) DeleteFileFromGCS(db *sql.DB, organizationUuid string, useId int) error {
	organization, err := s.ar.GetOrganizationDetail(organizationUuid)
	if err != nil || organization == nil {
		return errors.NewUnknownError("organization not found").WithInternalCause(err)
	}

	filePathPrefix := fmt.Sprintf("%v/profile_picture/%v", organization.ID, useId)
	err = s.cloudStorageProvider.DeleteFiles(filePathPrefix)
	if err != nil {
		return errors.NewUnknownError("delete file fail").WithInternalCause(err)
	}
	// update database profile_photo column in table users
	if err = s.ar.UpdateProfilePhotoOfUsers(db, useId, ""); err != nil {
		return errors.NewUnknownError("delete file fail").WithInternalCause(err)
	}
	return nil
}

// Update user language preference
func (s *accountServiceHandler) UpdateUserLanguagePreference(db *sql.DB, uid string, userId int, languagePreference string) error {
	// Update the language_preference field in the users table
	if err := s.ar.UpdateUserLanguagePreference(db, userId, languagePreference); err != nil {
		return errors.NewUnknownError("update user language preference fail").WithInternalCause(err)
	}

	// Set custom user claims
	account, err := s.ar.SelectSignInColumns(uid)
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

// Update user theme preference
func (s *accountServiceHandler) UpdateUserThemePreference(db *sql.DB, userId int, themePreference string) error {
	// Update the theme field in the users table
	if err := s.ar.UpdateUserThemePreference(db, userId, themePreference); err != nil {
		return errors.NewUnknownError("update user theme preference fail").WithInternalCause(err)
	}

	return nil
}
