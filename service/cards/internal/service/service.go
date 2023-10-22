package service

import (
	"database/sql"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/cards/internal/model"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/cards/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/checkout"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/checkout/models"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/gip"
)

// CardService is a interface which provide helper methods to access account related operations
type CardService interface {
	AddCard(req *models.ValidateCardRequest) error
	ListCards(accountID int) (model.ListCardsResponse, error)
	UpdateCard(accountID int, sourceID string) error
	DeleteCard(accountID int, sourceID string) error
}

type cardServiceHandler struct {
	db              *sql.DB
	cr              repository.CardRepository
	authProvider    gip.AuthProvider
	paymentProvider checkout.PaymentClient
}

// NewCardService initiates the account service object
func NewCardService(db *sql.DB, authProvider gip.AuthProvider, paymentProvider checkout.PaymentClient) CardService {
	cr := repository.NewCardRepository(db)
	return &cardServiceHandler{db, cr, authProvider, paymentProvider}
}
