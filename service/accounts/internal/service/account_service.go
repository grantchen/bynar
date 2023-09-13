/**
    @author: dongjs
    @date: 2023/9/12
    @description:
**/

package service

import (
	"errors"
	"fmt"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/accounts/internal/model"
	sql_db "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/utils"
	"github.com/sirupsen/logrus"
	"mime/multipart"
	"os"
	"path"
	"strings"
)

var profilePictureType = []string{"png", "jpg"}

// UploadFileToGCS upload user's profile picture to gcs
func (s *accountServiceHandler) UploadFileToGCS(tenantId, organizationUuid, email string, body *multipart.Reader) (string, error) {
	part, err := body.NextPart()
	if err != nil {
		logrus.Errorf("UploadFileToGCS: NextPart err: %+v", err)
		return "", errors.New("upload file read error")
	}
	defer part.Close()
	logrus.Infof("Uploaded File: %+v\n", part.FileName())
	ext := path.Ext(part.FileName())
	if !utils.IsStringArrayInclude(profilePictureType, strings.ToLower(ext)) {
		return "", errors.New("profile picture type is not png or jpg")
	}
	organization, err := s.ar.GetOrganizationDetail(organizationUuid)
	if err != nil || organization == nil {
		logrus.Errorf("UploadFileToGCS: GetOrganizationDetail err: %+v", err)
		return "", errors.New("organization not found")
	}
	user, err := s.GetUserDetail(tenantId, organizationUuid, email)
	if err != nil || organization == nil {
		logrus.Errorf("UploadFileToGCS: GetUserDetail err: %+v", err)
		return "", errors.New("user not found")
	}
	filePath := fmt.Sprintf("%v/profile_picture/%v%v", organization.ID, user.ID, ext)
	filePathPrefix := fmt.Sprintf("%v/profile_picture/%v", organization.ID, user.ID)

	if err = s.cloudStorageProvider.DeleteFiles(filePathPrefix); err != nil {
		logrus.Errorf("UploadFileToGCS: DeleteFiles err: %+v", err)
		return "", errors.New("upload file error")
	}
	url, err := s.cloudStorageProvider.UploadFile(filePath, part)
	if err != nil {
		logrus.Errorf("UploadFileToGCS: SignedURL err: %+v", err)
		return "", errors.New("upload file error")
	}
	// update database profile_photo column in table users
	if err = s.UpdateProfilePhotoOfUsers(tenantId, organizationUuid, user.ID, url); err != nil {
		logrus.Errorf("UploadFileToGCS: UpdateProfilePhotoOfUsers err: %+v", err)
		return "", errors.New("upload file error")
	}
	return url, nil
}

// DeleteFileFromGCS delete user's profile picture from google cloud storage
func (s *accountServiceHandler) DeleteFileFromGCS(tenantId, organizationUuid, email string) error {
	organization, err := s.ar.GetOrganizationDetail(organizationUuid)
	if err != nil || organization == nil {
		logrus.Errorf("DeleteFileFromGCS: GetOrganizationDetail err: %+v", err)
		return errors.New("organization not found")
	}
	user, err := s.GetUserDetail(tenantId, organizationUuid, email)
	if err != nil || organization == nil {
		logrus.Errorf("DeleteFileFromGCS: GetUserDetail err: %+v", err)
		return errors.New("user not found")
	}
	filePathPrefix := fmt.Sprintf("%v/profile_picture/%v", organization.ID, user.ID)
	err = s.cloudStorageProvider.DeleteFiles(filePathPrefix)
	if err != nil {
		logrus.Errorf("DeleteFileFromGCS: DeleteFiles err: %+v", err)
		return errors.New("delete file fail")
	}
	// update database profile_photo column in table users
	if err = s.UpdateProfilePhotoOfUsers(tenantId, organizationUuid, user.ID, ""); err != nil {
		logrus.Errorf("DeleteFileFromGCS: UpdateProfilePhotoOfUsers err: %+v", err)
		return errors.New("delete file fail")
	}
	return nil
}

// GetUserDetail get user details from organization_schema(uuid)
func (s *accountServiceHandler) GetUserDetail(tenantUuid, organizationUuid, email string) (*model.User, error) {
	connStr := os.Getenv(tenantUuid) + organizationUuid
	db, err := sql_db.InitializeConnection(connStr)
	logrus.Info("init db ", connStr)
	if err != nil {
		logrus.Errorf("GetUserDetail: init db connection url: %s error:%+v", connStr, err)
		return nil, err
	}
	defer db.Close()
	var querySql = `select a.id,a.eamil,a.full_name,a.phone,a.status,a.language_preference,a.policy_id,a.theme from users a where a.eamil = ? limit 1`
	var user = model.User{}
	err = db.QueryRow(querySql, email).Scan(
		&user.ID, &user.Email, &user.FullName, &user.Phone, &user.Status, &user.LanguagePreference, &user.PolicyId, &user.Theme)
	if err != nil {
		return nil, fmt.Errorf("query row: [%w]", err)
	}
	return &user, nil
}

// UpdateProfilePhotoOfUsers update column profile_photo in table users of organization_schema(uuid)
func (s *accountServiceHandler) UpdateProfilePhotoOfUsers(tenantUuid, organizationUuid string, accountID int, profilePhoto string) error {
	connStr := os.Getenv(tenantUuid) + organizationUuid
	db, err := sql_db.InitializeConnection(connStr)
	logrus.Info("init db ", connStr)
	if err != nil {
		logrus.Errorf("UpdateProfilePhotoOfUsers: init db connection url: %s error:%+v", connStr, err)
		return err
	}
	defer db.Close()
	if _, err = db.Exec(`UPDATE users SET profile_photo = ? WHERE id = ?`, profilePhoto, accountID); err != nil {
		return err
	}
	return nil
}
