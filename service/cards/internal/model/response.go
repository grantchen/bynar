package model

type AccountHolder struct {
	Phone struct{} `json:"phone"`
}

type CardDetails struct {
	ExpiryMonth   int           `json:"expiry_month"`
	ExpiryYear    int           `json:"expiry_year"`
	Scheme        string        `json:"scheme"`
	Last4         string        `json:"last4"`
	BIN           string        `json:"bin"`
	CardType      string        `json:"card_type"`
	CardCategory  string        `json:"card_category"`
	IssuerCountry string        `json:"issuer_country"`
	ProductID     string        `json:"product_id"`
	ProductType   string        `json:"product_type"`
	AccountHolder AccountHolder `json:"account_holder"`
	ID            string        `json:"id"`
	Type          string        `json:"type"`
	Fingerprint   string        `json:"fingerprint"`
}

type ListCardsResponse struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Email       string        `json:"email"`
	Metadata    struct{}      `json:"metadata"`
	Default     string        `json:"default"`
	Instruments []CardDetails `json:"instruments"`
}

type UserCard struct {
	ID         int
	CustomerID string
	UserID     int
	Status     bool
	IsDefault  bool
	SourceID   string
	Email      string
	FullName   string
}
