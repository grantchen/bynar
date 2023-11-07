package service

import (
	"context"
	"database/sql"
	"fmt"
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
	err := s.ar.IsCanDeleteOrganizationAccount(language, organizationUuid)
	if err != nil {
		return err
	}

	organizationID, err := s.ar.GetOrganizationIdByUuid(language, organizationUuid)
	if err != nil {
		return err
	}

	// delete all customers of organization from checkout.com
	err = s.deleteCustomersFromCheckout(language, organizationID)
	if err != nil {
		return err
	}

	// delete existing file from Google cloud storage
	if err = s.cloudStorageProvider.DeleteFiles(fmt.Sprintf("%d/", organizationID)); err != nil {
		return err
	}

	// delete all users of organization from GIP
	err = s.deleteUsersFromGIP(language, organizationID)
	if err != nil {
		return err
	}

	err = s.ar.DeleteOrganizationAccount(db, language, tenantUuid, organizationUuid)
	if err != nil {
		return err
	}

	return nil
}

// delete all customers of organization from checkout.com
func (s *accountServiceHandler) deleteCustomersFromCheckout(language string, organizationID int) error {
	customerIDs, err := s.ar.GetCustomerIDsByOrganizationID(language, organizationID)
	if err != nil {
		return err
	}

	for _, customerID := range customerIDs {
		err = s.paymentProvider.DeleteCustomer(customerID)
		if err != nil {
			continue
		}
	}

	return nil
}

// delete all users of organization from GIP
func (s *accountServiceHandler) deleteUsersFromGIP(language string, organizationID int) error {
	userUids, err := s.ar.GetGipUserUidsByOrganizationID(language, organizationID)
	if err != nil {
		return err
	}

	if len(userUids) == 0 {
		return nil
	}

	_, err = s.authProvider.DeleteUsers(context.Background(), userUids)
	if err != nil {
		return err
	}

	return nil
}
