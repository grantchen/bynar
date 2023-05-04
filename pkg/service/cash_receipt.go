package service

import (
	"database/sql"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/models"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/repository"
)

type cashReceiptService struct {
	cashReceiptRep repository.CashReceiptRepository
}

func NewCashReceiptService(cashReceiptRep repository.CashReceiptRepository) CashReceiptService {
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
