/**
    @author: dongjs
    @date: 2023/9/6
    @description:
**/

package models

import "time"

// ValidateCardRequest validate card request struct
// Token is Frames card details tokenized
type ValidateCardRequest struct {
	ID    int    `json:"id"`
	Token string `json:"token"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// ValidateCard validate card response struct
type ValidateCard struct {
	ID              string     `json:"id"`
	ActionID        string     `json:"action_id"`
	Amount          int        `json:"amount"`
	Currency        string     `json:"currency"`
	Approved        bool       `json:"approved"`
	Status          string     `json:"status"`
	AuthCode        string     `json:"auth_code"`
	ResponseCode    string     `json:"response_code"`
	ResponseSummary string     `json:"response_summary"`
	Balances        Balances   `json:"balances"`
	Risk            Risk       `json:"risk"`
	Source          Source     `json:"source"`
	Customer        Customer   `json:"customer"`
	ProcessedOn     time.Time  `json:"processed_on"`
	SchemeID        string     `json:"scheme_id"`
	Processing      Processing `json:"processing"`
	ExpiresOn       time.Time  `json:"expires_on"`
	Links           Links      `json:"_links"`
}

// Balances validate card balances response
type Balances struct {
	TotalAuthorized    int `json:"total_authorized"`
	TotalVoided        int `json:"total_voided"`
	AvailableToVoid    int `json:"available_to_void"`
	TotalCaptured      int `json:"total_captured"`
	AvailableToCapture int `json:"available_to_capture"`
	TotalRefunded      int `json:"total_refunded"`
	AvailableToRefund  int `json:"available_to_refund"`
}

// Risk validate card risk response
type Risk struct {
	Flagged bool    `json:"flagged"`
	Score   float64 `json:"score"`
}

// Source validate card source response
type Source struct {
	ID                      string `json:"id"`
	Type                    string `json:"type"`
	Phone                   Phone  `json:"phone"`
	ExpiryMonth             int    `json:"expiry_month"`
	ExpiryYear              int    `json:"expiry_year"`
	Scheme                  string `json:"scheme"`
	Last4                   string `json:"last4"`
	Fingerprint             string `json:"fingerprint"`
	Bin                     string `json:"bin"`
	CardType                string `json:"card_type"`
	CardCategory            string `json:"card_category"`
	IssuerCountry           string `json:"issuer_country"`
	ProductID               string `json:"product_id"`
	ProductType             string `json:"product_type"`
	AvsCheck                string `json:"avs_check"`
	CvvCheck                string `json:"cvv_check"`
	PaymentAccountReference string `json:"payment_account_reference"`
}

// Phone validate card pone response
type Phone struct {
	Number string `json:"number"`
}

// Customer validate card customer response
type Customer struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// Processing validate card processing response
type Processing struct {
	AcquirerTransactionID    string `json:"acquirer_transaction_id"`
	RetrievalReferenceNumber string `json:"retrieval_reference_number"`
	MerchantCategoryCode     string `json:"merchant_category_code"`
	SchemeMerchantID         string `json:"scheme_merchant_id"`
	Aft                      bool   `json:"aft"`
}

// Links validate card links response
type Links struct {
	Self    Link `json:"self"`
	Actions Link `json:"actions"`
}

// Link validate card link response
type Link struct {
	Href string `json:"href"`
}

// CardValidationPayload validate card request struct
type CardValidationPayload struct {
	Source              TokenSource `json:"source"`
	Currency            string      `json:"currency"`
	Customer            NewCustomer `json:"customer"`
	ProcessingChannelID string      `json:"processing_channel_id"`
}

// TokenSource validate card request struct
type TokenSource struct {
	Type  string `json:"type"`
	Token string `json:"token"`
}

// NewCustomer validate card request struct
type NewCustomer struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}
