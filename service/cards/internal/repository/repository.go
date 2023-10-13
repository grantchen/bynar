package repository

import (
	"database/sql"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/cards/internal/model"
)

// CardRepository provides a interface on db level for cards
type CardRepository interface {
	CountCard(id int) (int, error)
	AddCard(userID int, customerID, sourceID string, total int) error
	FetchCardBySourceID(sourceID string) (model.UserCard, error)
	ListCards(accountID int) (model.ListCardsResponse, error)
	UpdateDefaultCard(tx *sql.Tx, accountID int, sourceID string) error
	DeleteCard(tx *sql.Tx, sourceID string) error
}

type cardRepositoryHandler struct {
	db *sql.DB
}

func NewCardRepository(db *sql.DB) CardRepository {
	return &cardRepositoryHandler{db}
}
