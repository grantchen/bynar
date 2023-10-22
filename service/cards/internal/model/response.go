package model

import "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/checkout/models"

type AccountHolder struct {
	Phone struct{} `json:"phone"`
}

type ListCardsResponse struct {
	ID          string               `json:"id"`
	Name        string               `json:"name"`
	Email       string               `json:"email"`
	Metadata    struct{}             `json:"metadata"`
	Default     string               `json:"default"`
	Instruments []models.CardDetails `json:"instruments"`
}

type UserCard struct {
	ID         int
	CustomerID string
	UserID     int
	Status     bool
	IsDefault  bool
	SourceID   string
	Email      string
	FullName   string
}
