package model

type ConfirmEmailResponse struct {
	AccountID int `json:"accountID"`
}

type CreateUserResponse struct {
	Token string `json:"token"`
}
