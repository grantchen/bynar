package svc

import (
	"database/sql"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/payments/internal/models"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/payments/internal/repositories"
)

type cashReceiptService struct {
	cashReceiptRep repositories.CashReceiptRepository
}

func NewCashReceiptService(cashReceiptRep repositories.CashReceiptRepository) CashReceiptService {
	return &cashReceiptService{
		cashReceiptRep: cashReceiptRep,
	}
}

func (s *cashReceiptService) GetPaymentTx(tx *sql.Tx, id interface{}) (*models.CashReceipt, error) {
	return s.cashReceiptRep.Get(tx, id)
}

func (s *cashReceiptService) Handle(tx *sql.Tx, m *models.CashReceipt, moduleID int) error {
	// todo: handle payment
	panic("not implemented yet")
}

func (s *cashReceiptService) HandleLine(tx *sql.Tx, pr *models.CashReceipt, l *models.CashReceiptLine) (err error) {
	// todo: handle payment line
	panic("not implemented yet")
}
