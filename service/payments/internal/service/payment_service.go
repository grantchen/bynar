package service

import (
	"database/sql"
	"fmt"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/payments/internal/models"
	repository "git-codecommit.eu-central-1.amazonaws.com/v1/repos/payments/internal/repositories"
	pkg_repository "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/repository"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

type paymentService struct {
	gridRowDataRepositoryWithChild treegrid.GridRowDataRepositoryWithChild
	db                             *sql.DB
	paymentRepository              repository.PaymentRepository
	procRepository                 pkg_repository.ProcurementRepository
	currRepository                 pkg_repository.CurrencyRepository
	cashManageRepository           pkg_repository.CashManagementRepository
}

// GetPageCount implements PaymentsService
func (u *paymentService) GetPageCount(tr *treegrid.Treegrid) (int64, error) {
	return u.gridRowDataRepositoryWithChild.GetPageCount(tr)
}

// GetPageData implements PaymentsService
func (u *paymentService) GetPageData(tr *treegrid.Treegrid) ([]map[string]string, error) {
	return u.gridRowDataRepositoryWithChild.GetPageData(tr)
}

func (u *paymentService) Handle(tx *sql.Tx, m *models.Payment, moduleID int) error {
	if err := u.handlePayment(tx, m); err != nil {
		return fmt.Errorf("handle payment: [%w]", err)
	}
	// exchange rate of currency table
	localCurrency, err := u.currRepository.GetLedgerSetupCurrency()
	if err != nil {
		return fmt.Errorf("get ledger setup currency: [%w]", err)
	}

	// get payment by id = m.id and applied = To apply
	lines, err := u.paymentRepository.GetLines(tx, m.ID, ToApply)
	if err != nil {
		return fmt.Errorf("get lines: [%w]", err)
	}

	for _, l := range lines {
		l.AmountLcy = l.Amount * localCurrency
		if err := u.HandleLine(tx, m, l); err != nil {
			return fmt.Errorf("handle line: [%w]; id: %d, parentID: %d", err, l.ID, l.ParentID)
		}
	}

	return nil
}
func (u *paymentService) GetTx(tx *sql.Tx, id interface{}) (*models.Payment, error) {
	return u.paymentRepository.Get(tx, id)
}

func (u *paymentService) handlePayment(tx *sql.Tx, m *models.Payment) error {
	if m.PaidStatus == 1 {
		return nil
	}

	curr, err := u.currRepository.GetCurrency(m.CurrencyID)
	if err != nil {
		return fmt.Errorf("get currency: [%w]", err)
	}

	m.CurrencyValue = curr.ExchangeRate
	m.AmountLcy = m.Amount * m.CurrencyValue

	cashManagement, err := u.cashManageRepository.Get(m.BankID)
	if err != nil {
		return fmt.Errorf("get cash management: [%w]", err)
	}

	cashManagement.Amount -= m.Amount

	if err := u.cashManageRepository.Update(tx, cashManagement); err != nil {
		return fmt.Errorf("update cash management: [%w]", err)
	}
	m.PaidStatus = 1

	if err := u.paymentRepository.Save(tx, m); err != nil {
		return fmt.Errorf("save payment: [%w]", err)
	}

	return nil
}

func (u *paymentService) HandleLine(tx *sql.Tx, payment *models.Payment, line *models.PaymentLine) (err error) {
	pr, err := u.procRepository.GetProcurement(tx, line.AppliesDocumentID)
	if err != nil {
		return fmt.Errorf("get procurement: [%w], id: %d", err, line.AppliesDocumentID)
	}

	if err := u.paymentRepository.SaveLine(tx, line); err != nil {
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

	if err := u.procRepository.SaveProcurement(tx, pr); err != nil {
		return fmt.Errorf("save procurement: [%w]", err)
	}

	return nil
}

func NewPaymentService(
	db *sql.DB,
	gridRowDataRepositoryWithChild treegrid.GridRowDataRepositoryWithChild,
	paymentRepository repository.PaymentRepository,
	procRepository pkg_repository.ProcurementRepository,
	currRepository pkg_repository.CurrencyRepository,
	cashManageRepository pkg_repository.CashManagementRepository,
) PaymentService {
	return &paymentService{
		db:                             db,
		gridRowDataRepositoryWithChild: gridRowDataRepositoryWithChild,
		paymentRepository:              paymentRepository,
		procRepository:                 procRepository,
		currRepository:                 currRepository,
		cashManageRepository:           cashManageRepository,
	}
}
