package service

import (
	"database/sql"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/payments/internal/models"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

type PaymentService interface {
	//GetTx(tx *sql.Tx, id interface{}) (*models.Payment, error)
	//Handle(tx *sql.Tx, m *models.Payment, moduleID int) error
	//HandleLine(tx *sql.Tx, payment *models.Payment, line *models.PaymentLine) (err error)

	GetTx(tx *sql.Tx, id interface{}) (*models.Payment, error)
	Handle(tx *sql.Tx, pr *models.Payment) error

	GetPageCount(treegrid *treegrid.Treegrid) (int64, error)
	GetPageData(tr *treegrid.Treegrid) ([]map[string]string, error)
}
