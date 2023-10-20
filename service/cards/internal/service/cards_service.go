package service

import (
	"database/sql"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/cards/internal/model"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/checkout/models"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/errors"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/utils"
	"github.com/sirupsen/logrus"
)

func (r *cardServiceHandler) AddCard(req *models.ValidateCardRequest) error {
	total, err := r.cr.CountCard(req.ID)
	if err != nil {
		return err
	}
	if total >= 10 {
		return errors.NewUnknownError("user can't add more than 10 cards", "")
	}
	logrus.Infof("----- add card request  %#v\n", req)
	cardDetails, err := r.paymentProvider.ValidateCard(req)
	if err != nil {
		return errors.NewUnknownError("failed to validate card", "").WithInternal().WithCause(err)
	}
	// TODO: delete log after tested
	logrus.Infof("----- add card details %#v\n", cardDetails)
	var rsp models.CustomerResponse
	if cardDetails.Source.ID == "" {
		rsp, err = r.paymentProvider.FetchCustomerDetails(cardDetails.Customer.ID)
		if err != nil {
			return errors.NewUnknownError("failed to fetch card", "").WithInternal().WithCause(err)
		}
		logrus.Infof("--- fetch cards %#v\n", rsp.Instruments)
		cardDetails.Source.ID = rsp.Instruments[len(rsp.Instruments)-1].ID
	}
	if !cardDetails.Approved || cardDetails.Status != "Card Verified" {
		err = r.paymentProvider.DeleteCard(cardDetails.Source.ID)
		if err != nil {
			return errors.NewUnknownError("failed to delete card", "").WithInternal().WithCause(err)
		}
		return errors.NewUnknownError("card not approved or not verified", "")
	}
	err = r.cr.AddCard(req.ID, cardDetails.Customer.ID, cardDetails.Source.ID, total)
	if err != nil {
		derr := r.paymentProvider.DeleteCard(cardDetails.Source.ID)
		if derr != nil {
			return errors.NewUnknownError("failed to delete card", "").WithInternal().WithCause(derr)
		}
		return errors.NewUnknownError("failed to add user card", "").WithInternal().WithCause(err)
	}
	return nil
}

func (r *cardServiceHandler) ListCards(accountID int) (model.ListCardsResponse, error) {
	card, err := r.cr.ListCards(accountID)
	if err != nil {
		return card, err
	}
	info, err := r.paymentProvider.FetchCustomerDetails(card.ID)
	if err != nil {
		return card, err
	}
	card.Name = info.Name
	card.Email = info.Email
	card.Default = info.Default
	card.Instruments = info.Instruments
	return card, nil
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
		customerInfo := models.UpdateCustomer{
			Email:             cardDetails.Email,
			Name:              cardDetails.FullName,
			DefaultInstrument: sourceID,
		}
		if err = r.paymentProvider.UpdateCustomer(customerInfo, cardDetails.CustomerID); err != nil {
			return errors.NewUnknownError("update card failed", "").WithInternal().WithCause(err)
		}
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
		if err = r.paymentProvider.DeleteCard(sourceID); err != nil {
			return errors.NewUnknownError("delete card failed from checkout", "").WithInternal().WithCause(err)
		}
		return nil
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
		if err = r.paymentProvider.DeleteCard(cardDetails.SourceID); err != nil {
			return errors.NewUnknownError("delete card failed from checkout", "").WithInternal().WithCause(err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
