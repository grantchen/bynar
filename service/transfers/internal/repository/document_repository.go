package repository

import (
	"database/sql"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

// DocumentRepository is a repository for document
type documentRepository struct {
}

// NewDocumentRepository returns a new DocumentRepository
func NewDocumentRepository(_ *sql.DB) DocumentRepository {
	return &documentRepository{}
}

// IsAuto returns true if document is auto
func (dr *documentRepository) IsAuto(*sql.Tx, *treegrid.MainRow) (bool, error) {
	return false, nil
}

// Generate generates document
func (dr *documentRepository) Generate(*sql.Tx, *treegrid.MainRow) (string, error) {
	return "", nil
}
