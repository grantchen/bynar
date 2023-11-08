package repository

import "database/sql"

// OrganizationRepository is implementation of OrganizationRepository
type organizationRepository struct {
	db *sql.DB
}

// NewOrganizationRepository create new instance of OrganizationRepository
func NewOrganizationRepository(db *sql.DB) OrganizationRepository {
	return &organizationRepository{db: db}
}
