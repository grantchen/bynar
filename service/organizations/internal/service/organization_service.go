package service

import (
	"database/sql"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

// organizationService is the implementation of OrganizationService
type organizationService struct {
	db                           *sql.DB
	simpleOrganizationRepository treegrid.SimpleGridRowRepository
}

// GetPageCount implements OrganizationService
func (o *organizationService) GetPageCount(tr *treegrid.Treegrid) (int64, error) {
	return o.simpleOrganizationRepository.GetPageCount(tr)
}

// GetPageData implements OrganizationService
func (o *organizationService) GetPageData(tr *treegrid.Treegrid) ([]map[string]string, error) {
	return o.simpleOrganizationRepository.GetPageData(tr)
}

// NewOrganizationService create new instance of OrganizationService
func NewOrganizationService(db *sql.DB, simpleOrganizationService treegrid.SimpleGridRowRepository) OrganizationService {
	return &organizationService{db: db, simpleOrganizationRepository: simpleOrganizationService}
}
