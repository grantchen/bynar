package repository

import (
	"database/sql"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/payments/internal/models"
)

// PaymentRepository is an interface for payment repository
type PaymentRepository interface {
	//GetDocID get doc_id by id
	GetDocID(id interface{}) (docID int, err error)
	//GetStatus get status by id
	GetStatus(id interface{}) (status int, err error)
	//Get gets payment by id
	Get(tx *sql.Tx, id interface{}) (m *models.Payment, err error)
	//Save saves payment
	Save(tx *sql.Tx, m *models.Payment) (err error)
	// GetLines get payment lines by payment_id and applied
	GetLines(tx *sql.Tx, parentID interface{}, applied int) ([]*models.PaymentLine, error)
	// SaveLine saves payment line
	SaveLine(tx *sql.Tx, l *models.PaymentLine) (err error)
}
