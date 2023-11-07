package repository

import (
	"database/sql"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/cards/internal/model"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/checkout/models"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/errors"
)

func (r *cardRepositoryHandler) CountCard(id int) (int, error) {
	var total int
	err := r.db.QueryRow(`SELECT count(*) FROM accounts_cards WHERE user_id = ?`, id).Scan(&total)
	if err != nil {
		return 0, errors.NewUnknownError("failed to count user cards", "").WithInternal().WithCause(err)
	}
	return total, nil
}

func (r *cardRepositoryHandler) AddCard(userID int, customerID, sourceID string, total int) error {
	isDefault := 1
	if total > 0 {
		isDefault = 0
	}
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM accounts_cards WHERE user_payment_gateway_id=? AND source_id=?`, customerID, sourceID).Scan(&count)
	if err != nil {
		return errors.NewUnknownError("select card failed", "").WithInternal().WithCause(err)
	}
	if count > 0 {
		return errors.NewUnknownError("card exists", "")
	}
	_, err = r.db.Exec(`INSERT INTO accounts_cards
				(user_payment_gateway_id, source_id, user_id, is_default, account_id)
				VALUE
				(?, ?, ?, ?, ?)`,
		customerID,
		sourceID,
		userID,
		isDefault,
		userID)
	if err != nil {
		return errors.NewUnknownError("add card failed", "").WithInternal().WithCause(err)
	}
	return nil
}

func (r *cardRepositoryHandler) FetchCardBySourceID(sourceID string) (cardDetails model.UserCard, err error) {
	err = r.db.QueryRow(`SELECT ac.id,
		user_payment_gateway_id,
		user_id,
		ac.status,
		is_default,
		source_id,
		email,
		full_name
	FROM accounts_cards ac
	JOIN accounts a on ac.user_id = a.id
	WHERE source_id = ?`, sourceID).
		Scan(&cardDetails.ID, &cardDetails.CustomerID, &cardDetails.UserID,
			&cardDetails.Status, &cardDetails.IsDefault, &cardDetails.SourceID,
			&cardDetails.Email, &cardDetails.FullName)
	if err != nil {
		return cardDetails, errors.NewUnknownError("fetch card failed", "").WithInternal().WithCause(err)
	}
	return cardDetails, nil
}

func (r *cardRepositoryHandler) ListCards(accountID int) (model.ListCardsResponse, map[string]bool, error) {
	resp := model.ListCardsResponse{Instruments: make([]models.CardDetails, 0)}
	ins := map[string]bool{}
	err := r.db.QueryRow(`SELECT user_payment_gateway_id, source_id, email, full_name FROM accounts_cards ac JOIN accounts a on ac.user_id = a.id WHERE ac.user_id = ? AND is_default = ?`, accountID, true).Scan(&resp.ID, &resp.Default, &resp.Email, &resp.Name)
	if err != nil {
		return resp, ins, errors.NewUnknownError("list card failed", "").WithInternal().WithCause(err)
	}
	rows, err := r.db.Query(`SELECT source_id FROM accounts_cards WHERE user_id = ?`, accountID)
	if err != nil {
		return resp, ins, errors.NewUnknownError("list card failed", "").WithInternal().WithCause(err)
	}
	for rows.Next() {
		var inst string
		rows.Scan(&inst)
		ins[inst] = true
	}
	defer rows.Close()
	return resp, ins, nil
}

func (r *cardRepositoryHandler) UpdateDefaultCard(tx *sql.Tx, accountID int, sourceID string) error {
	_, err := tx.Exec(`UPDATE accounts_cards
					SET is_default = IF(source_id = ?, 1, 0)
					WHERE user_id = ?`, sourceID, accountID)
	if err != nil {
		return errors.NewUnknownError("failed to update card", "").WithInternal().WithCause(err)
	}
	return nil
}

func (r *cardRepositoryHandler) DeleteCard(tx *sql.Tx, sourceID string) error {
	_, err := tx.Exec(`DELETE
					FROM accounts_cards
					WHERE source_id = ?`, sourceID)
	if err != nil {
		return errors.NewUnknownError("failed to delete card", "").WithInternal().WithCause(err)
	}
	return nil
}
