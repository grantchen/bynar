package repository

import "database/sql"

// siteRepository implements SiteRepository
type siteRepository struct {
	db *sql.DB
}

// NewSiteRepository create new instance of SiteRepository
func NewSiteRepository(db *sql.DB) SiteRepository {
	return &siteRepository{db: db}
}
