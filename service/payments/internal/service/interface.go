package service

import (
	"database/sql"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/payments/internal/models"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

// PaymentService is the interface for payment service
type PaymentService interface {
	GetTx(tx *sql.Tx, id interface{}) (*models.Payment, error)
	// Handle handles payment calculation
	Handle(tx *sql.Tx, pr *models.Payment) error
	// HandleLine handles payment line calculation
	HandleLine(tx *sql.Tx, payment *models.Payment, line *models.PaymentLine) (err error)
	// GetPageCount returns the page count
	GetPageCount(treegrid *treegrid.Treegrid) (int64, error)
	// GetPageData returns the page data
	GetPageData(tr *treegrid.Treegrid) ([]map[string]string, error)
}
