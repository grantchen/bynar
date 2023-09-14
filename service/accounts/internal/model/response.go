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
	LanguagePreference string `json:"languagePreference"`
	Email              string `json:"email"`
	FullName           string `json:"fullName"`
	ProfileUrl         string `json:"profileUrl"`
	//todo add other column
}
