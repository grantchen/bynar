package service

import "database/sql"

type organizationService struct {
	db *sql.DB
}

func NewOrganizationService(db *sql.DB) OrganizationService {
	return &organizationService{db: db}
}
