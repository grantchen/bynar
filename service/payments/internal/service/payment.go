package service

import (
	"database/sql"
	"fmt"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/payments/internal/models"
	repository "git-codecommit.eu-central-1.amazonaws.com/v1/repos/payments/internal/repositories"
	pkg_repository "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/repository"
)

const (
	ToApply        = 1
	SuccessApplied = 2
)

type paymentService struct {
	paymentRepository   repository.PaymentRepository
	procRepository      pkg_repository.ProcurementRepository
	currRepository      pkg_repository.CurrencyRepository
	cashManagRepository pkg_repository.CashManagementRepository
}

func NewPaymentService(
	paymentRepository repository.PaymentRepository,
	procRepository pkg_repository.ProcurementRepository,
	currRepository pkg_repository.CurrencyRepository,
	cashManagRepository pkg_repository.CashManagementRepository) *paymentService {
	return &paymentService{
		paymentRepository:   paymentRepository,
		procRepository:      procRepository,
		currRepository:      currRepository,
		cashManagRepository: cashManagRepository,
	}
}

func (s *paymentService) GetTx(tx *sql.Tx, id interface{}) (*models.Payment, error) {
	return s.paymentRepository.Get(tx, id)
}

func (s *paymentService) Handle(tx *sql.Tx, m *models.Payment, moduleID int) error {
	if err := s.handlePayment(tx, m); err != nil {
		return fmt.Errorf("handle payment: [%w]", err)
	}
	// exchange rate of currency table
	localCurrency, err := s.currRepository.GetLedgerSetupCurrency()
	if err != nil {
		return fmt.Errorf("get ledger setup currency: [%w]", err)
	}

	// get payment by id = m.id and applied = To apply
	lines, err := s.paymentRepository.GetLines(tx, m.ID, ToApply)
	if err != nil {
		return fmt.Errorf("get lines: [%w]", err)
	}

	for _, l := range lines {
		l.AmountLcy = l.Amount * localCurrency
		if err := s.HandleLine(tx, m, l); err != nil {
			return fmt.Errorf("handle line: [%w]; id: %d, parentID: %d", err, l.ID, l.ParentID)
		}
	}

	return nil
}

func (s *paymentService) HandleLine(tx *sql.Tx, payment *models.Payment, line *models.PaymentLine) (err error) {
	pr, err := s.procRepository.GetProcurement(tx, line.AppliesDocumentID)
	if err != nil {
		return fmt.Errorf("get procurement: [%w], id: %d", err, line.AppliesDocumentID)
	}

	if err := s.paymentRepository.SaveLine(tx, line); err != nil {
		return fmt.Errorf("save payment line: [%w]", err)
	}

	pr.Paid = line.Amount
	pr.Remaining = pr.TotalInclusiveVat - pr.Paid
	switch {
	case pr.TotalInclusiveVat == pr.Paid:
		pr.PaidStatus = 1
	case pr.TotalInclusiveVat > pr.Paid:
		pr.PaidStatus = 0
	case pr.TotalInclusiveVat < pr.Paid:
		pr.PaidStatus = 2
	}

	if err := s.procRepository.SaveProcurement(tx, pr); err != nil {
		return fmt.Errorf("save procurement: [%w]", err)
	}

	return nil
}

func (s *paymentService) handlePayment(tx *sql.Tx, m *models.Payment) error {
	if m.PaidStatus == 1 {
		return nil
	}

	curr, err := s.currRepository.GetCurrency(m.CurrencyID)
	if err != nil {
		return fmt.Errorf("get currency: [%w]", err)
	}

	m.CurrencyValue = curr.ExchangeRate
	m.AmountLcy = m.Amount * m.CurrencyValue

	cashManagement, err := s.cashManagRepository.Get(m.BankID)
	if err != nil {
		return fmt.Errorf("get cash management: [%w]", err)
	}

	cashManagement.Amount -= m.Amount

	if err := s.cashManagRepository.Update(tx, cashManagement); err != nil {
		return fmt.Errorf("update cash management: [%w]", err)
	}
	m.PaidStatus = 1

	if err := s.paymentRepository.Save(tx, m); err != nil {
		return fmt.Errorf("save payment: [%w]", err)
	}

	return nil
}
