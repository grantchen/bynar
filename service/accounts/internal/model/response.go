package model

type ConfirmEmailResponse struct {
	AccountID int `json:"accountID"`
}

type VerifyCardResponse struct {
	CustomerID string `json:"customerID"`
	SourceID   string `json:"sourceID"`
}

type CreateUserResponse struct {
	Token string `json:"token"`
}

// SignInResponse sign in api return struct
type SignInResponse struct {
	IdToke string `json:"idToke"`
}

// GetUserResponse get user api return struct
type GetUserResponse struct {
	ID           int    `db:"id" json:"id"`
	Email        string `db:"email" json:"email"`
	FullName     string `json:"fullName"`
	Country      string `db:"country" json:"country"`
	AddressLine  string `db:"address" json:"addressLine"`
	AddressLine2 string `db:"address_2" json:"addressLine2"`
	City         string `db:"city" json:"city"`
	PostalCode   string `db:"postal_code" json:"postalCode"`
	State        string `db:"state" json:"state"`
	PhoneNumber  string `db:"phone" json:"phoneNumber"`

	LanguagePreference string         `json:"languagePreference"`
	ThemePreference    string         `json:"themePreference"`
	ProfileURL         string         `json:"profileURL"`
	PolicyID           int            `json:"policyId"`
	Permissions        map[string]int `json:"permissions"`
}

// SignIn generate token columns in database when user sign in
type SignIn struct {
	Uid                  string `json:"uid"`
	OrganizationUuid     string `json:"organization_uuid"`
	OrganizationUserId   int    `json:"organization_user_id"`
	OrganizationStatus   bool   `json:"organization_status"`
	OrganizationAccount  bool   `json:"organization_account"`
	OrganizationVerified bool   `json:"organization_verified"`
	TenantUuid           string `json:"tenant_uuid"`
	TenantStatus         bool   `json:"tenant_status"`
	TenantSuspended      bool   `json:"tenant_suspended"`

	Language string `json:"language"`
	Theme    string `json:"theme"`
}
