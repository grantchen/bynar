package service

import "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"

// SiteService interface
type SiteService interface {
	// GetPageCount get page count
	GetPageCount(tr *treegrid.Treegrid) (int64, error)
	// GetPageData get page data
	GetPageData(tr *treegrid.Treegrid) ([]map[string]string, error)
}
