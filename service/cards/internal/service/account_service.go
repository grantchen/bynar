package service

import (
	"database/sql"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/cards/internal/model"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/checkout/models"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/errors"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/utils"
)

func (r *cardServiceHandler) AddCard(req *models.ValidateCardRequest) error {
	total, err := r.cr.CountCard(req.ID)
	if err != nil {
		return err
	}
	if total >= 10 {
		return errors.NewUnknownError("user can't add more than 10 cards", "")
	}
	cardDetails, err := r.paymentProvider.ValidateCard(req)
	if err != nil {
		return errors.NewUnknownError("failed to validate card", "").WithInternal().WithCause(err)
	}
	if !cardDetails.Approved || cardDetails.Status != "Card Verified" {
		// TODO: delete card
		// err = r.paymentProvider.DeleteCard(cardDetails.Source.ID)
		// if err != nil {
		// 	return errors.NewUnknownError("failed to delete card", "").WithInternal().WithCause(err)
		// }
		return errors.NewUnknownError("failed to validate card", "")
	}
	err = r.cr.AddCard(req.ID, cardDetails.Customer.ID, cardDetails.Source.ID, total)
	if err != nil {
		errors.NewUnknownError("failed to add user card", "").WithInternal().WithCause(err)
	}
	return nil
}

func (r *cardServiceHandler) ListCards(accountID int) (model.ListCardsResponse, error) {
	// TODO:
	return model.ListCardsResponse{Name: "test", Email: "demo@demo.com", Instruments: []model.CardDetails{{ExpiryYear: 2023, ExpiryMonth: 11}}}, nil
}

func (r *cardServiceHandler) UpdateCard(accountID int, sourceID string) error {
	cardDetails, err := r.cr.FetchCardBySourceID(sourceID)
	if err != nil {
		return err
	}

	if cardDetails.UserID != accountID {
		return errors.NewUnknownError("user not authorized for this operation", "")
	}

	err = utils.WithTransaction(r.db, func(tx *sql.Tx) error {
		if err = r.cr.UpdateDefaultCard(tx, accountID, sourceID); err != nil {
			return err
		}
		// TODO: update card
		// customerInfo := models.UpdateCustomer{
		// 	Email:             cardDetails.Email,
		// 	Name:              cardDetails.FullName,
		// 	DefaultInstrument: sourceID,
		// }
		// if err = r.paymentProvider.UpdateCustomer(customerInfo, cardDetails.CustomerID); err != nil {
		// 	logrus.Errorf("DeleteCard: Error updating card from checkout %v", err)
		// 	return err
		// }
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *cardServiceHandler) DeleteCard(accountID int, sourceID string) error {
	cardDetails, err := r.cr.FetchCardBySourceID(sourceID)
	if err != nil {
		return err
	}

	if cardDetails.UserID != accountID {
		return errors.NewUnknownError("user not authorized for this operation", "")
	}

	if cardDetails.IsDefault {
		return errors.NewUnknownError("user cannot delete default card", "")
	}

	err = utils.WithTransaction(r.db, func(tx *sql.Tx) error {
		if err = r.cr.DeleteCard(tx, sourceID); err != nil {
			return errors.NewUnknownError("delete card failed from db", "").WithInternal().WithCause(err)
		}
		// TODO: deletecard
		// if err = r.paymentProvider.DeleteCard(cardDetails.SourceID); err != nil {
		// 	return errors.NewUnknownError("delete card failed from checkout", "").WithInternal().WithCause(err)
		// }
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
