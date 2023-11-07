package model

type SignupRequest struct {
	Email string `json:"email" valid:""`
}

type ConfirmEmailRequest struct {
	Email     string `json:"email" valid:""`
	Timestamp string `json:"timestamp" valid:""`
	Signature string `json:"signature" valid:""`
}

type VerifyCardRequest struct {
	Token string `json:"token" valid:""`
	Email string `json:"email" valid:""`
	Name  string `json:"name" valid:""`
}

type CreateUserRequest struct {
	Username            string `json:"username" valid:""`
	FullName            string `json:"fullName" valid:""`
	Country             string `json:"country" valid:""`
	AddressLine         string `json:"addressLine" valid:""`
	AddressLine2        string `json:"addressLine2"`
	City                string `json:"city" valid:""`
	PostalCode          string `json:"postalCode" valid:""`
	State               string `json:"state" valid:""`
	PhoneNumber         string `json:"phoneNumber" valid:""`
	OrganizationName    string `json:"organizationName" valid:""`
	VAT                 string `json:"VAT" valid:""`
	OrganisationCountry string `json:"organisationCountry" valid:""`
	IsAgreementSigned   bool   `json:"isAgreementSigned" valid:""`
	Token               string `json:"token" valid:""`
	Timestamp           string `json:"timestamp" valid:""`
	Signature           string `json:"signature" valid:""`
	CustomerID          string `json:"customerID" valid:""`
	SourceID            string `json:"sourceID" valid:""`
	TenantCode          string `json:"tenantCode" valid:""`
}

// SignInRequest user sign_in request struct
type SignInRequest struct {
	Email   string `json:"email" valid:""`
	OobCode string `json:"oobCode" valid:""`
}

// SendSignInEmailRequest user send sign in email return struct
type SendSignInEmailRequest struct {
	Email string `json:"email" valid:""`
}

type UpdateUserLanguagePreferenceRequest struct {
	LanguagePreference string `json:"languagePreference" valid:""`
}

type UpdateUserThemePreferenceRequest struct {
	ThemePreference string `json:"themePreference" valid:""`
}

type UpdateUserProfileRequest struct {
	Email       string `json:"email" valid:""`
	FullName    string `json:"fullName" valid:""`
	PhoneNumber string `json:"phoneNumber" valid:""`
	Language    string `json:"language" valid:""`
	Theme       string `json:"theme" valid:""`
}

// OrganizationAccountRequest organization account request struct
type OrganizationAccountRequest struct {
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
