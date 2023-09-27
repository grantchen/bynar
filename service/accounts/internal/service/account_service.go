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
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/accounts/internal/model"
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

// UpdateUserProfile update user profile
func (s *accountServiceHandler) UpdateUserProfile(db *sql.DB, userId int, uid string, userProfile model.UpdateUserProfileRequest) error {
	prevDetail, err := s.ar.GetUserDetail(db, userId)
	if err != nil {
		return errors.NewUnknownError("no user found").WithInternalCause(err)
	}
	// if person changes email then only validate email from abstract api
	var (
		needUpdate       = false
		needUpdateClaims = false
	)
	gipUpdateParam := map[string]interface{}{}
	if prevDetail.Email != userProfile.Email {
		//todo verify email
		gipUpdateParam["email"] = userProfile.Email
		needUpdate = true
	}
	phoneNumber := userProfile.PhoneNumber
	if phoneNumber[0] != '+' {
		phoneNumber = "+" + phoneNumber
	}
	// if person changes phone number then only validate phone number from abstract api
	if prevDetail.Phone != phoneNumber {
		//todo verify phoneNumber
		gipUpdateParam["phoneNumber"] = phoneNumber
		needUpdate = true
	}
	if prevDetail.FullName != userProfile.FullName {
		gipUpdateParam["displayName"] = userProfile.FullName
		needUpdate = true
	}
	if prevDetail.Theme != userProfile.Theme || prevDetail.LanguagePreference != userProfile.Language {
		needUpdate = true
		needUpdateClaims = true
	}
	if false == needUpdate {
		return nil
	}
	//update gip user info
	err = s.authProvider.UpdateUser(context.Background(), uid, gipUpdateParam)
	if err != nil {
		return errors.NewUnknownError("update user profile fail").WithInternalCause(err)
	}
	// update database user profile
	err = s.ar.UpdateUserProfile(db, userId, uid, userProfile)
	if err != nil {
		return errors.NewUnknownError("update user profile fail").WithInternalCause(err)
	}
	//update gip custom claims
	if needUpdateClaims {
		err = s.UpdateGipCustomClaims(uid)
		if err != nil {
			return errors.NewUnknownError("update user profile fail").WithInternalCause(err)
		}
	}
	return nil
}

// UpdateGipCustomClaims update custom claims
func (s *accountServiceHandler) UpdateGipCustomClaims(uid string) error {
	account, err := s.ar.SelectSignInColumns(uid)
	if err != nil || account == nil {
		return errors.NewUnknownError("query claims data fail").WithInternalCause(err)
	}

	claims, err := convertSignInToClaims(account)
	if err != nil {
		return errors.NewUnknownError("convert to claims struct fail").WithInternalCause(err)
	}
	err = s.authProvider.SetCustomUserClaims(context.Background(), uid, claims)
	if err != nil {
		return errors.NewUnknownError("set custom claims fail").WithInternalCause(err)
	}
	return nil
}

// GetUserProfileById get user profile by userId
func (s *accountServiceHandler) GetUserProfileById(db *sql.DB, userId int) (*model.UserProfileResponse, error) {
	detail, err := s.ar.GetUserDetail(db, userId)
	if err != nil {
		return nil, errors.NewUnknownError("no record found").WithInternalCause(err)
	}
	return &model.UserProfileResponse{
		Email:       detail.Email,
		PhoneNumber: detail.Phone,
		FullName:    detail.FullName,
		Theme:       detail.Theme,
		Language:    detail.LanguagePreference,
	}, nil
}
