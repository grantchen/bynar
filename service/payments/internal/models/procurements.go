package models

type Procurement struct {
	ID                         int
	DocumentID                 int
	DocumentNo                 string
	TransactionNo              int
	StoreID                    int
	DocumentDate               string
	PostingDate                string
	EntryDate                  string
	ShipmentDate               string
	ProjectID                  int
	DepartmentID               int
	ContractID                 int
	UserGroupID                int
	Status                     int
	BudgetID                   int
	InvoiceDiscountID          int
	CurrencyID                 int
	CurrencyValue              float32
	VendorID                   int
	VendorInvoiceNo            float32
	PurchaserID                int
	ResponsibilityCenterID     int
	PaymentTermsID             int
	PaymentMethodID            int
	TransactionTypeID          int
	PaymentDiscount            float32
	ShipmentMethodID           int
	PaymentReference           int
	CreditorNo                 string
	OnHold                     string
	TransactionSpecificationID int
	TransportMethodID          int
	EntryPointID               int
	CampaignID                 int
	AreaID                     int
	VendorShipmentNo           string
	SubtotalExclusiveVat       float32
	TotalDiscount              float32
	TotalExclusiveVat          float32
	TotalVat                   float32
	TotalInclusiveVat          float32
	SubtotalExclusiveVatLcy    float32
	TotalDiscountLcy           float32
	TotalExclusiveVatLcy       float32
	TotalVatLcy                float32
	TotalInclusiveVatLcy       float32
	Paid                       float32
	Remaining                  float32
	PaidStatus                 uint8

	Discount DiscountVat

	Lines []*ProcurementLine
}

type ProcurementLine struct {
	ID                      int
	ParentID                int
	ItemType                int
	ItemID                  int
	LocationID              int
	InputQuantity           float32
	ItemUnitValue           float32
	Quantity                float32
	ItemUnitID              int
	DiscountID              int
	TaxAreaID               int
	VatID                   int
	QuantityAssign          float32
	QuantityAssigned        float32
	SubtotalExclusiveVat    float32
	TotalDiscount           float32
	TotalExclusiveVat       float32
	TotalVat                float32
	TotalInclusiveVat       float32
	SubtotalExclusiveVatLcy float32
	TotalDiscountLcy        float32
	TotalExclusiveVatLcy    float32
	TotalVatLcy             float32
	TotalInclusiveVatLcy    float32
}

type DiscountVat struct {
	Value      float32
	Percentage int
}

type Currency struct {
	ExchangeRate float32
}
