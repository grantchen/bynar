package model

import "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/models"

// ConfirmEmailResponse confirm email api return struct
type ConfirmEmailResponse struct {
	AccountID int `json:"accountID"`
}

// VerifyCardResponse verify card api return struct
type VerifyCardResponse struct {
	CustomerID string `json:"customerID"`
	SourceID   string `json:"sourceID"`
}

// CreateUserResponse create user api return struct
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

	LanguagePreference string        `json:"languagePreference"`
	ThemePreference    string        `json:"themePreference"`
	ProfileURL         string        `json:"profileURL"`
	PolicyID           int           `json:"policyId"`
	Permissions        models.Policy `json:"permissions"`
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

	AccountId int `json:"account_id"`
}

// UserProfileResponse user profile api return struct
type UserProfileResponse struct {
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	FullName    string `json:"fullName"`
	Theme       string `json:"theme"`
	Language    string `json:"language"`
}

// GetOrganizationAccountResponse get organization account api return struct
type GetOrganizationAccountResponse struct {
	Email               string `json:"email" valid:""`
	FullName            string `json:"fullName" valid:""`
	Country             string `json:"country" valid:""`
	AddressLine         string `json:"addressLine" valid:""`
	AddressLine2        string `json:"addressLine2"`
	City                string `json:"city" valid:""`
	State               string `json:"state" valid:""`
	PostalCode          string `json:"postalCode" valid:""`
	PhoneNumber         string `json:"phoneNumber" valid:""`
	OrganizationName    string `json:"organizationName" valid:""`
	VAT                 string `json:"VAT" valid:""`
	OrganizationCountry string `json:"organizationCountry" valid:""`
}
