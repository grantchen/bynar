package repository

import (
	"database/sql"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/models"
)

// SaleRepository is an interface for sale repository
type SaleRepository interface {
	// GetDocID returns document id by sale id
	GetDocID(id interface{}) (docID int, err error)
	// GetStatus returns status by sale id
	GetStatus(id interface{}) (status int, err error)
	// GetSale returns sale by id
	GetSale(tx *sql.Tx, id interface{}) (m *models.Sale, err error)
	// SaveSale saves sale
	SaveSale(tx *sql.Tx, m *models.Sale) (err error)
	// GetSaleLines returns sale lines by parent id
	GetSaleLines(tx *sql.Tx, parentID interface{}) ([]*models.SaleLine, error)
	// SaveSaleLine saves sale line
	SaveSaleLine(tx *sql.Tx, l *models.SaleLine) (err error)
}
