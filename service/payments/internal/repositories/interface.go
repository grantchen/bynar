package repository

import (
	"database/sql"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/payments/internal/models"
)

type PaymentRepository interface {
	GetDocID(id interface{}) (docID int, err error)
	GetStatus(id interface{}) (status int, err error)
	Get(tx *sql.Tx, id interface{}) (m *models.Payment, err error)
	Save(tx *sql.Tx, m *models.Payment) (err error)
	GetLines(tx *sql.Tx, parentID interface{}, applied int) ([]*models.PaymentLine, error)
	SaveLine(tx *sql.Tx, l *models.PaymentLine) (err error)
}
