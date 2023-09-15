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
	ID                  int    `db:"id" json:"id"`
	Email               string `db:"email" json:"email"`
	FullName            string `json:"fullName"`
	Country             string `db:"country" json:"country"`
	AddressLine         string `db:"address" json:"addressLine"`
	AddressLine2        string `db:"address_2" json:"addressLine2"`
	City                string `db:"city" json:"city"`
	PostalCode          string `db:"postal_code" json:"postalCode"`
	State               string `db:"state" json:"state"`
	PhoneNumber         string `db:"phone" json:"phoneNumber"`
	UserGroups          string `db:"cognito_user_groups" json:"cognitoUserGroups"`
	OrganisationID      int    `db:"organization_id" json:"organisationID"`
	OrganisationAccount int    `db:"organization_account" json:"organisationAccount"`
	CanUpdate           bool   `db:"can_update" json:"canUpdate"`
	CanDelete           bool   `db:"can_delete" json:"canDelete"`
	Status              bool   `db:"account_confirm_status" json:"accountConfirmStatus"`

	LanguagePreference string `json:"languagePreference"`
	ProfileURL         string `json:"profileURL"`
	PolicyID           int    `json:"policyId"`
}
