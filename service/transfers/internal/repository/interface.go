package repository

import (
	"database/sql"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

// TransferRepository is a repository for transfer
type TransferRepository interface {
	// Save saves transfer
	Save(tx *sql.Tx, tr *treegrid.MainRow) error
	// SaveDocumentID saves document id
	SaveDocumentID(tx *sql.Tx, tr *treegrid.MainRow, docID string) error
	// UpdateStatus updates status
	UpdateStatus(tx *sql.Tx, status int) error
}

// InventoryRepository is a repository for inventory
type InventoryRepository interface {
	// CheckQuantityAndValue checks quantity and value
	CheckQuantityAndValue(tx *sql.Tx, tr *treegrid.MainRow) (bool, error)
	// Save saves inventory
	Save(tx *sql.Tx, tr *treegrid.MainRow) error
}

// DocumentRepository is a repository for document
type DocumentRepository interface {
	// IsAuto returns true if document is auto
	IsAuto(tx *sql.Tx, tr *treegrid.MainRow) (bool, error)
	// Generate generates document
	Generate(tx *sql.Tx, tr *treegrid.MainRow) (string, error)
}

// UserRepository is a repository for user
type UserRepository interface {
}
