/**
    @author: dongjs
    @date: 2023/9/12
    @description:
**/

package service

import (
	"errors"
	"fmt"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/gcs"
	"github.com/sirupsen/logrus"
	"mime/multipart"
	"path"
)

// UploadFileToGCS upload user's profile picture to gcs
func (s *accountServiceHandler) UploadFileToGCS(accountID, orgID int, body *multipart.Reader) (string, error) {
	part, err := body.NextPart()
	if err != nil {
		logrus.Errorf("UploadFileToGCS: NextPart err: %+v", err)
		return "", errors.New("upload file read error")
	}
	defer part.Close()
	logrus.Infof("Uploaded File: %+v\n", part.FileName())

	ext := path.Ext(part.FileName())
	filePath := fmt.Sprintf("%v/profile_picture/%v%v", orgID, accountID, ext)
	filePathPrefix := fmt.Sprintf("%v/profile_picture/%v", orgID, accountID)

	gcsClient, err := gcs.NewGCSClient()
	if err != nil {
		logrus.Errorf("UploadFileToGCS: NewGCSClient err: %+v", err)
		return "", errors.New("upload file read error")
	}
	defer gcsClient.StorageClient().Close()

	if err = gcsClient.DeleteFiles(filePathPrefix); err != nil {
		logrus.Errorf("UploadFileToGCS: DeleteFiles err: %+v", err)
		return "", errors.New("upload file error")
	}
	url, err := gcsClient.UploadFile(filePath, part)
	if err != nil {
		logrus.Errorf("UploadFileToGCS: SignedURL err: %+v", err)
		return "", errors.New("upload file error")
	}
	// update database profile_photo column in table users
	if err = s.ar.UpdateProfilePhotoOfUsers(accountID, url); err != nil {
		logrus.Errorf("UploadFileToGCS: UpdateProfilePhotoOfUsers err: %+v", err)
		return "", errors.New("upload file error")
	}
	return url, nil
}

// DeleteFileFromGCS delete user's profile picture from google cloud storage
func (s *accountServiceHandler) DeleteFileFromGCS(accountID, orgID int) error {
	filePathPrefix := fmt.Sprintf("%v/profile_picture/%v", orgID, accountID)
	gcsClient, err := gcs.NewGCSClient()
	if err != nil {
		logrus.Errorf("DeleteFileFromGCS: NewGCSClient err: %+v", err)
		return errors.New("delete file fail")
	}
	defer gcsClient.StorageClient().Close()
	err = gcsClient.DeleteFiles(filePathPrefix)
	if err != nil {
		logrus.Errorf("DeleteFileFromGCS: DeleteFiles err: %+v", err)
		return errors.New("delete file fail")
	}
	return nil

}
