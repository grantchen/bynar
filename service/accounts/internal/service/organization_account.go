package service

import (
	"context"
	"database/sql"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/accounts/internal/model"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/i18n"
)

func (s *accountServiceHandler) GetOrganizationAccount(language string, accountID int, organizationUuid string) (*model.GetOrganizationAccountResponse, error) {
	account, err := s.ar.GetOrganizationAccount(language, accountID, organizationUuid)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (s *accountServiceHandler) UpdateOrganizationAccount(
	db *sql.DB,
	language string,
	accountID int,
	uid string,
	organizationUserId int,
	organizationUuid string,
	organizationAccount model.OrganizationAccountRequest,
) error {
	if organizationAccount.PhoneNumber[0] != '+' {
		organizationAccount.PhoneNumber = "+" + organizationAccount.PhoneNumber
	}

	err := s.authProvider.UpdateUser(context.Background(), uid, map[string]interface{}{
		"email":       organizationAccount.Email,
		"displayName": organizationAccount.FullName,
		"phoneNumber": organizationAccount.PhoneNumber,
	})
	if err != nil {
		return i18n.TranslationErrorToI18n(language, err)
	}

	err = s.ar.UpdateOrganizationAccount(db, language, accountID, organizationUserId, organizationUuid, organizationAccount)
	if err != nil {
		return i18n.TranslationErrorToI18n(language, err)
	}

	return nil
}

func (s *accountServiceHandler) DeleteOrganizationAccount(
	db *sql.DB,
	language string,
	tenantUuid string,
	organizationUuid string,
) error {
	// TODO delete cards from checkout.com

	// TODO delete existing file from Google cloud storage
	//if err = s.cloudStorageProvider.DeleteFiles(filePathPrefix); err != nil {
	//	return "", errors.NewUnknownError("upload file fail", "").WithInternalCause(err)
	//}

	err := s.ar.DeleteOrganizationAccount(db, language, tenantUuid, organizationUuid)
	if err != nil {
		return err
	}

	// TODO delete user from Google identify platform
	//err = s.authProvider.DeleteUser(context.Background(), uid)
	//if err != nil {
	//	return i18n.TranslationErrorToI18n(language, err)
	//}

	return nil
}
