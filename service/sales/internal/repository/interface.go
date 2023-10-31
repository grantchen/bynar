package repository

import (
	"database/sql"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/models"
)

type SaleRepository interface {
	GetDocID(id interface{}) (docID int, err error)
	GetStatus(id interface{}) (status int, err error)
	GetSale(tx *sql.Tx, id interface{}) (m *models.Sale, err error)
	SaveSale(tx *sql.Tx, m *models.Sale) (err error)
	GetSaleLines(tx *sql.Tx, parentID interface{}) ([]*models.SaleLine, error)
	SaveSaleLine(tx *sql.Tx, l *models.SaleLine) (err error)
}
