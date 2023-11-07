package service

import (
	"context"
	"database/sql"
	"fmt"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/accounts/internal/model"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/i18n"
	"log"
)

// Maximum number of gip users allowed to batch delete at a time.
const gipMaxDeleteAccountsBatchSize = 1000

// GetOrganizationAccount returns the organization account.
func (s *accountServiceHandler) GetOrganizationAccount(language string, accountID int, organizationUuid string) (*model.GetOrganizationAccountResponse, error) {
	account, err := s.ar.GetOrganizationAccount(language, accountID, organizationUuid)
	if err != nil {
		return nil, err
	}

	return account, nil
}

// UpdateOrganizationAccount updates the organization account.
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

	// check if organization vat is duplicated
	err := s.ar.IsOrganizationVATDuplicated(language, organizationUuid, organizationAccount.VAT)
	if err != nil {
		return err
	}

	err = s.authProvider.UpdateUser(context.Background(), uid, map[string]interface{}{
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

// DeleteOrganizationAccount deletes the organization account.
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

	// delete users from GIP in batches
	err = s.userUidsInBatches(userUids, gipMaxDeleteAccountsBatchSize, func(dividedIds []string) error {
		_, err = s.authProvider.DeleteUsers(context.Background(), dividedIds)
		if err != nil {
			log.Print(err)
			// return nil for continue deleting other users
			return nil
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// userUidsInBatches divides the given user IDs into batches and calls the given function for each batch.
func (s *accountServiceHandler) userUidsInBatches(ids []string, chunkSize int, batchFunc func(dividedIds []string) error) error {
	idsLength := len(ids)
	for i := 0; i < idsLength; i += chunkSize {
		end := i + chunkSize
		if end > idsLength {
			end = idsLength
		}
		err := batchFunc(ids[i:end])
		if err != nil {
			return err
		}
	}

	return nil
}
