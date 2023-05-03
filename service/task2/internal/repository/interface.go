package repository

import (
	"database/sql"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

type TransferRepository interface {
	Save(tx *sql.Tx, tr *treegrid.MainRow) error
	SaveTransfer(tx *sql.Tx, tr *treegrid.MainRow) error
	SaveTransferLines(tx *sql.Tx, tr *treegrid.MainRow) error
	SaveDocumentID(tx *sql.Tx, tr *treegrid.MainRow, docID string) error
	UpdateStatus(tx *sql.Tx, status int) error
}

type InventoryRepository interface {
	CheckQuantityAndValue(tx *sql.Tx, tr *treegrid.MainRow) (bool, error)
	Save(tx *sql.Tx, tr *treegrid.MainRow) error
}

type DocumentRepository interface {
	IsAuto(tx *sql.Tx, tr *treegrid.MainRow) (bool, error)
	Generate(tx *sql.Tx, tr *treegrid.MainRow) (string, error)
	Save(tx *sql.Tx, tr *treegrid.MainRow) error
}

type WorkflowRepository interface {
	GetModuleID() (int, error)
	CheckApprovalOrder(accountID, status int) (bool, error)
}

type UserRepository interface {
	GetUserID(accountID int) (int, error)
}
