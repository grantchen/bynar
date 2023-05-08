package models

import "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/models"

type Payment struct {
	ID                     int
	BatchID                int
	DocumentID             int
	DocumentNo             string
	ExternalDocumentNo     string
	TransactionNo          int
	StoreID                int
	DocumentDate           string
	PostingDate            string
	EntryDate              string
	AccountType            int
	AccountID              int
	RecipientBankAccountID int
	BalanceAccountType     int
	BalanceAccountID       int
	Amount                 float32
	AmountLcy              float32
	CurrencyID             int
	CurrencyValue          float32
	UserGroupID            int
	Status                 int
	PaymentMethodID        int
	PaymentReference       string
	CreditorNo             string
	BankPaymentTypeID      int
	BankID                 int
	Paid                   float32
	Remaining              float32
	PaidStatus             int

	CashManagement *models.CashManagement
}

type PaymentLine struct {
	ID                  int
	ParentID            int
	AppliesDocumentType int
	AppliesDocumentID   int
	PaymentTypeID       int
	Amount              float32
	AmountLcy           float32
	Applied             int
}
