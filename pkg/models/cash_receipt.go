package models

type CashReceipt struct {
	ID                 int
	BatchID            int
	DocumentID         int
	DocumentNo         string
	TransactionNo      int
	StoreID            int
	DocumentDate       string
	PostingDate        string
	EntryDate          string
	AccountType        int
	AccountID          int
	BalanceAccountType int
	BalanceAccountID   int
	Amount             float32
	AmountLcy          float32
	CurrencyValue      float32
	UserGroupID        int
	Status             int
	BankID             int
}

type CashReceiptLine struct {
	ID                  int
	ParentID            int
	AppliesDocumentType int
	AppliesDocumentID   int
	Amount              float32
	AmountLcy           float32
	Applied             int
}

type CashManagement struct {
	Id          int
	Type        int
	Bank_Id     int
	Currency_Id int
	Amount      float32
}
