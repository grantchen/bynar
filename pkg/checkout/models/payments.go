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
	AcquirerTransactionID            string `json:"acquirer_transaction_id"`
	RetrievalReferenceNumber         string `json:"retrieval_reference_number"`
	MerchantCategoryCode             string `json:"merchant_category_code"`
	SchemeMerchantID                 string `json:"scheme_merchant_id"`
	Aft                              bool   `json:"aft"`
	RecommendationCode               string `json:"recommendation_code"`
	PartnerOrderId                   string `json:"partner_order_id"`
	PartnerSessionId                 string `json:"partner_session_id"`
	PartnerClientToken               string `json:"partner_client_token"`
	PartnerPaymentId                 string `json:"partner_payment_id"`
	PartnerStatus                    string `json:"partner_status"`
	PartnerTransactionId             string `json:"partner_transaction_id"`
	PartnerErrorCodes                string `json:"partner_error_codes"`
	PartnerErrorMessage              string `json:"partner_error_message"`
	PartnerAuthorizationCode         string `json:"partner_authorization_code"`
	PartnerAuthorizationResponseCode string `json:"partner_authorization_response_code"`
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

// CheckOutErrorResponse checkout return error struct
type CheckOutErrorResponse struct {
	RequestId  string   `json:"request_id"`
	ErrorType  string   `json:"error_type"`
	ErrorCodes []string `json:"error_codes"`
}

type UpdateCustomer struct {
	Email             string `json:"email"`
	Name              string `json:"name"`
	DefaultInstrument string `json:"default"`
}

type PaymentSource struct {
	Type string `json:"type"`
}
type Destination struct {
	Type string `json:"type"`
}

type ThreeDS struct {
	Downgraded                 bool   `json:"downgraded"`
	Enrolled                   string `json:"enrolled"`
	SignatureValid             string `json:"signature_valid"`
	AuthenticationResponse     string `json:"authentication_response"`
	AuthenticationStatusReason string `json:"authentication_status_reason"`
	Cryptogram                 string `json:"cryptogram"`
	XID                        string `json:"xid"`
	Version                    string `json:"version"`
	Exemption                  string `json:"exemption"`
	ExemptionApplied           string `json:"exemption_applied"`
	Challenged                 bool   `json:"challenged"`
	UpgradeReason              string `json:"upgrade_reason"`
}

type BillingDescriptor struct {
	Name      string `json:"name"`
	City      string `json:"city"`
	Reference string `json:"reference"`
}

type Shipping struct {
	Address Address       `json:"address"`
	Phone   ShippingPhone `json:"phone"`
}

type Address struct {
	AddressLine1 string `json:"address_line1"`
	AddressLine2 string `json:"address_line2"`
	City         string `json:"city"`
	State        string `json:"state"`
	Zip          string `json:"zip"`
	Country      string `json:"country"`
}

type ShippingPhone struct {
	CountryCode string `json:"country_code"`
	Number      string `json:"number"`
}

type Sender struct {
	Type      string `json:"type"`
	Reference string `json:"reference"`
}

type Marketplace struct {
	SubEntityID string      `json:"sub_entity_id"`
	SubEntities []SubEntity `json:"sub_entities"`
}

type SubEntity struct {
	ID         string     `json:"id"`
	Amount     int        `json:"amount"`
	Reference  string     `json:"reference"`
	Commission Commission `json:"commission"`
}

type Commission struct {
	Amount     int     `json:"amount"`
	Percentage float64 `json:"percentage"`
}
type AmountAllocation struct {
	ID         string     `json:"id"`
	Amount     int        `json:"amount"`
	Reference  string     `json:"reference"`
	Commission Commission `json:"commission"`
}

type Recipient struct {
	DOB           string  `json:"dob"`
	AccountNumber string  `json:"account_number"`
	Address       Address `json:"address"`
	Zip           string  `json:"zip"`
	FirstName     string  `json:"first_name"`
	LastName      string  `json:"last_name"`
}

type ProcessingPayments struct {
	PreferredScheme              string                  `json:"preferred_scheme"`
	AppID                        string                  `json:"app_id"`
	PartnerCustomerID            string                  `json:"partner_customer_id"`
	PartnerPaymentID             string                  `json:"partner_payment_id"`
	TaxAmount                    int                     `json:"tax_amount"`
	PurchaseCountry              string                  `json:"purchase_country"`
	Locale                       string                  `json:"locale"`
	RetrievalReferenceNumber     string                  `json:"retrieval_reference_number"`
	PartnerOrderID               string                  `json:"partner_order_id"`
	PartnerStatus                string                  `json:"partner_status"`
	PartnerTransactionID         string                  `json:"partner_transaction_id"`
	PartnerErrorCodes            []string                `json:"partner_error_codes"`
	PartnerErrorMessage          string                  `json:"partner_error_message"`
	PartnerAuthorizationCode     string                  `json:"partner_authorization_code"`
	PartnerAuthorizationResponse string                  `json:"partner_authorization_response_code"`
	FraudStatus                  string                  `json:"fraud_status"`
	ProviderAuthorizedPayment    AuthorizedPaymentMethod `json:"provider_authorized_payment_method"`
	CustomPaymentMethodIDs       []string                `json:"custom_payment_method_ids"`
	AFT                          bool                    `json:"aft"`
	MerchantCategoryCode         string                  `json:"merchant_category_code"`
	SchemeMerchantID             string                  `json:"scheme_merchant_id"`
}

type AuthorizedPaymentMethod struct {
	Type                 string `json:"type"`
	Description          string `json:"description"`
	NumberOfInstallments int    `json:"number_of_installments"`
	NumberOfDays         int    `json:"number_of_days"`
}

type PaymentDetailsItem struct {
	Name           string `json:"name"`
	Quantity       int    `json:"quantity"`
	UnitPrice      int    `json:"unit_price"`
	Reference      string `json:"reference"`
	CommodityCode  string `json:"commodity_code"`
	UnitOfMeasure  string `json:"unit_of_measure"`
	TotalAmount    int    `json:"total_amount"`
	TaxAmount      int    `json:"tax_amount"`
	DiscountAmount int    `json:"discount_amount"`
	WxpayGoodsID   int    `json:"wxpay_goods_id"`
	URL            string `json:"url"`
	ImageURL       string `json:"image_url"`
}

type Metadata struct {
	CouponCode string `json:"coupon_code"`
	PartnerID  int    `json:"partner_id"`
}

type Action struct {
	ID              string `json:"id"`
	Type            string `json:"type"`
	ResponseCode    string `json:"response_code"`
	ResponseSummary string `json:"response_summary"`
}

type PaymentLinks struct {
	Self      Link `json:"self"`
	Actions   Link `json:"actions"`
	Authorize Link `json:"authorize"`
	Refund    Link `json:"refund"`
}

type FetchPaymentDetails struct {
	ID                string               `json:"id"`
	RequestedOn       string               `json:"requested_on"`
	Source            PaymentSource        `json:"source"`
	Destination       Destination          `json:"destination"`
	Amount            int                  `json:"amount"`
	Currency          string               `json:"currency"`
	PaymentType       string               `json:"payment_type"`
	Reference         string               `json:"reference"`
	Description       string               `json:"description"`
	Approved          bool                 `json:"approved"`
	ExpiresOn         string               `json:"expires_on"`
	Status            string               `json:"status"`
	Balances          Balances             `json:"balances"`
	ThreeDS           ThreeDS              `json:"3ds"`
	Risk              Risk                 `json:"risk"`
	Customer          Customer             `json:"customer"`
	BillingDescriptor BillingDescriptor    `json:"billing_descriptor"`
	Shipping          Shipping             `json:"shipping"`
	PaymentIP         string               `json:"payment_ip"`
	Sender            Sender               `json:"sender"`
	Marketplace       Marketplace          `json:"marketplace"`
	AmountAllocations []AmountAllocation   `json:"amount_allocations"`
	Recipient         Recipient            `json:"recipient"`
	Processing        ProcessingPayments   `json:"processing"`
	Items             []PaymentDetailsItem `json:"items"`
	Metadata          Metadata             `json:"metadata"`
	ECI               string               `json:"eci"`
	SchemeID          string               `json:"scheme_id"`
	Actions           []Action             `json:"actions"`
	Links             PaymentLinks         `json:"_links"`
}

type CustomerResponse struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Email       string        `json:"email"`
	Metadata    struct{}      `json:"metadata"`
	Default     string        `json:"default"`
	Instruments []CardDetails `json:"instruments"`
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

type AccountHolder struct {
	Phone struct{} `json:"phone"`
}
