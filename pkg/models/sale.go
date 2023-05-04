package models

type Sale struct {
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
	CurrencyID                 int
	CurrencyValue              float32
	CustomerID                 int
	SalespersonID              int
	ResponsibilityCenterID     int
	PaymentTermsID             int
	PaymentMethodID            int
	TransactionTypeID          int
	PaymentDiscount            float32
	ShipmentMethodID           int
	PaymentReference           int
	TransactionSpecificationID int
	TransportMethodID          int
	ExitPointID                int
	CampaignID                 int
	AreaID                     int
	PackageTrackingNo          string
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
	DirectDebitMandateID       int
	AgentID                    int
	AgentServiceID             int
	DiscountID                 int

	Discount DiscountVat
}

type SaleLine struct {
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
