package service

import (
	"database/sql"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

type siteService struct {
	db                   *sql.DB
	simpleSiteRepository treegrid.SimpleGridRowRepository
}

// GetPageCount implements SiteService
func (o *siteService) GetPageCount(tr *treegrid.Treegrid) (int64, error) {
	return o.simpleSiteRepository.GetPageCount(tr)
}

// GetPageData implements SiteService
func (o *siteService) GetPageData(tr *treegrid.Treegrid) ([]map[string]string, error) {
	return o.simpleSiteRepository.GetPageData(tr)
}

func NewSiteService(db *sql.DB, simpleSiteService treegrid.SimpleGridRowRepository) SiteService {
	return &siteService{db: db, simpleSiteRepository: simpleSiteService}
}
