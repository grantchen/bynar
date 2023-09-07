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
