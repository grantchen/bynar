package service

import (
	"database/sql"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/models"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/repository"
)

type cashReceiptService struct {
	cashReceiptRep repository.CashReceiptRepository
}

// return new cashReceiptService
func _(cashReceiptRep repository.CashReceiptRepository) CashReceiptService {
	return &cashReceiptService{
		cashReceiptRep: cashReceiptRep,
	}
}

func (s *cashReceiptService) GetPaymentTx(tx *sql.Tx, id interface{}) (*models.CashReceipt, error) {
	return s.cashReceiptRep.Get(tx, id)
}

func (s *cashReceiptService) Handle(_ *sql.Tx, _ *models.CashReceipt) error {
	// todo: handle payment
	panic("not implemented yet")
}

func (s *cashReceiptService) HandleLine(_ *sql.Tx, _ *models.CashReceipt, _ *models.CashReceiptLine) (err error) {
	// todo: handle payment line
	panic("not implemented yet")
}
