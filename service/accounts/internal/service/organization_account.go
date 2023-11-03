package service

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/accounts/internal/model"
)

func (s *accountServiceHandler) GetOrganizationAccount(language string, accountID int, organizationUuid string) (*model.GetOrganizationAccountResponse, error) {
	account, err := s.ar.GetOrganizationAccount(language, accountID, organizationUuid)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (s *accountServiceHandler) UpdateOrganizationAccount(language string, accountID int, organizationUuid string, organizationAccount model.OrganizationAccountRequest) error {
	err := s.ar.UpdateOrganizationAccount(language, accountID, organizationUuid, organizationAccount)
	if err != nil {
		return err
	}

	//s.authProvider.UpdateUser()

	return nil
}

func (s *accountServiceHandler) DeleteOrganizationAccount(language string, accountID int, organizationUuid string) error {
	err := s.ar.DeleteOrganizationAccount(language, accountID, organizationUuid)
	if err != nil {
		return err
	}

	// remove users from gip
	//s.authProvider.DeleteUser()

	return nil
}
