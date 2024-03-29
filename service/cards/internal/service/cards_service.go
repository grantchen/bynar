package service

import (
	"database/sql"
	"strings"

	"github.com/sirupsen/logrus"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/cards/internal/model"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/checkout/models"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/errors"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/utils"
)

// AddCard add card through checkout and db
func (r *cardServiceHandler) AddCard(req *models.ValidateCardRequest) error {
	// get user card num in db
	total, err := r.cr.CountCard(req.ID)
	if err != nil {
		return err
	}
	logrus.Infof("----- add card request  %#v\n", req)
	cardDetails, err := r.paymentProvider.ValidateCard(req)
	if err != nil {
		return errors.NewUnknownError("failed to validate card", "failed-validate-card").WithInternal().WithCause(err)
	}
	logrus.Infof("----- add card details %#v\n", cardDetails)
	// get sourceid from checkout if response empty
	var rsp models.CustomerResponse
	if cardDetails.Source.ID == "" {
		rsp, err = r.paymentProvider.FetchCustomerDetails(cardDetails.Customer.ID)
		if err != nil {
			return errors.NewUnknownError("failed to fetch card", "failed-fetch-card").WithInternal().WithCause(err)
		}
		logrus.Infof("--- fetch cards %#v\n", rsp.Instruments)
		cardDetails.Source.ID = rsp.Instruments[len(rsp.Instruments)-1].ID
	}
	// check card status
	if !cardDetails.Approved || cardDetails.Status != "Card Verified" {
		err = r.paymentProvider.DeleteCard(cardDetails.Source.ID)
		if err != nil {
			return errors.NewUnknownError("failed to delete card", "failed-delete-card").WithInternal().WithCause(err)
		}
		return errors.NewUnknownError("card not approved or not verified", "failed-validate-card")
	}
	err = r.cr.AddCard(req.ID, cardDetails.Customer.ID, cardDetails.Source.ID, total)
	if err != nil {
		if !strings.Contains(err.Error(), "card exists") {
			derr := r.paymentProvider.DeleteCard(cardDetails.Source.ID)
			if derr != nil {
				return errors.NewUnknownError("failed to delete card", "failed-delete-card").WithInternal().WithCause(derr)
			}
		}
		return errors.NewUnknownError("failed to add user card", "failed-add-card").WithInternal().WithCause(err)
	}
	return nil
}

// ListCards list user's cards from checkout and db
func (r *cardServiceHandler) ListCards(accountID int) (model.ListCardsResponse, error) {
	card, ins, err := r.cr.ListCards(accountID)
	if err != nil {
		return card, errors.NewUnknownError("failed to fetch card", "failed-fetch-card").WithInternal().WithCause(err)
	}
	info, err := r.paymentProvider.FetchCustomerDetails(card.ID)
	if err != nil {
		return card, errors.NewUnknownError("failed to fetch card", "failed-fetch-card").WithInternal().WithCause(err)
	}
	card.Name = info.Name
	card.Email = info.Email
	card.Default = info.Default
	instruments := make([]models.CardDetails, 0)
	for i := range info.Instruments {
		if ins[info.Instruments[i].ID] {
			instruments = append(instruments, info.Instruments[i])
		}
	}
	card.Instruments = instruments
	return card, nil
}

// UpdateCard set default card of user in checkout and db
func (r *cardServiceHandler) UpdateCard(accountID int, sourceID string) error {
	// get card from checkout
	cardDetails, err := r.cr.FetchCardBySourceID(sourceID)
	if err != nil {
		return errors.NewUnknownError("failed to fetch card", "failed-fetch-card").WithInternal().WithCause(err)
	}

	if cardDetails.UserID != accountID {
		return errors.NewUnknownError("user not authorized for this operation", "user-not-allowed")
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
			return errors.NewUnknownError("update card failed", "failed-update-card").WithInternal().WithCause(err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// DeleteCard delete card in checkout and db
func (r *cardServiceHandler) DeleteCard(accountID int, sourceID string) error {
	// get card from checkout
	cardDetails, err := r.cr.FetchCardBySourceID(sourceID)
	if err != nil {
		if err = r.paymentProvider.DeleteCard(sourceID); err != nil {
			return errors.NewUnknownError("delete card failed from checkout", "failed-delete-card").WithInternal().WithCause(err)
		}
		return nil
	}
	//  check permission
	if cardDetails.UserID != accountID {
		return errors.NewUnknownError("user not authorized for this operation", "user-not-allowed")
	}
	// check default
	if cardDetails.IsDefault {
		return errors.NewUnknownError("user cannot delete default card", "failed-delete-card")
	}

	err = utils.WithTransaction(r.db, func(tx *sql.Tx) error {
		if err = r.cr.DeleteCard(tx, sourceID); err != nil {
			return errors.NewUnknownError("delete card failed from db", "failed-delete-card").WithInternal().WithCause(err)
		}
		if err = r.paymentProvider.DeleteCard(cardDetails.SourceID); err != nil {
			return errors.NewUnknownError("delete card failed from checkout", "failed-delete-card").WithInternal().WithCause(err)
		}
		return nil
	})
	if err != nil {
		return errors.NewUnknownError("delete card failed", "failed-delete-card").WithInternal().WithCause(err)
	}
	return nil
}
