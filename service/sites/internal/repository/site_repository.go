package repository

import "database/sql"

type siteRepository struct {
	db *sql.DB
}

func NewSiteRepository(db *sql.DB) SiteRepository {
	return &siteRepository{db: db}
}
