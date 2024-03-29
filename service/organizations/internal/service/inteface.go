package service

import "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"

// OrganizationService is the interface for the organization service
type OrganizationService interface {
	GetPageCount(tr *treegrid.Treegrid) (int64, error)
	GetPageData(tr *treegrid.Treegrid) ([]map[string]string, error)
}
