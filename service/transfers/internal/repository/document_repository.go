package repository

import (
	"database/sql"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

// DocumentRepository is a repository for document
type documentRepository struct {
}

// NewDocumentRepository returns a new DocumentRepository
func NewDocumentRepository(db *sql.DB) DocumentRepository {
	return &documentRepository{}
}

// IsAuto returns true if document is auto
func (dr *documentRepository) IsAuto(tx *sql.Tx, tr *treegrid.MainRow) (bool, error) {
	return false, nil
}

// Generate generates document
func (dr *documentRepository) Generate(tx *sql.Tx, tr *treegrid.MainRow) (string, error) {
	return "", nil
}
