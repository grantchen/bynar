package model

type SignupRequest struct {
	Email string `json:"email" valid:""`
}

type ConfirmEmailRequest struct {
	Email string `json:"email" valid:""`
	Code  string `json:"code" valid:""`
}

type VerifyCardRequest struct {
	ID    int    `json:"id" valid:""`
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
}

// SignInRequest user sign_in request struct
type SignInRequest struct {
	Email string `json:"email" valid:""`
}
