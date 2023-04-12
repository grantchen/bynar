package repositories

import (
	"database/sql"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/payments/internal/models"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/logger"
)

type currencyRepository struct {
	conn *sql.DB
}

func NewCurrencyRepository(conn *sql.DB) CurrencyRepository {
	return &currencyRepository{conn: conn}
}

func (d *currencyRepository) GetDiscount(id int) (m models.DiscountVat, err error) {
	query := `
	SELECT value, percentage
	FROM discounts
	WHERE id = ?
	`

	err = d.conn.QueryRow(query, id).Scan(&m.Value, &m.Percentage)
	return
}

func (d *currencyRepository) GetVat(id int) (m models.DiscountVat, err error) {
	query := `
	SELECT value, percentage
	FROM vats
	WHERE id = ?
	`

	err = d.conn.QueryRow(query, id).Scan(&m.Value, &m.Percentage)
	return
}

func (d *currencyRepository) GetCurrency(id int) (m models.Currency, err error) {
	query := `
	SELECT exchange_rate
	FROM currencies
	WHERE id = ?
	`

	err = d.conn.QueryRow(query, id).Scan(&m.ExchangeRate)
	return
}

func (d *currencyRepository) GetLedgerSetupCurrency() (curr float32, err error) {
	logger.Debug("get general ledger setup")

	query := `
	SELECT exchange_rate
	FROM currencies c
	INNER JOIN general_ledger_setup gls ON c.id = gls.local_currency_id
	WHERE gls.id = 1
	`

	err = d.conn.QueryRow(query).Scan(&curr)

	return
}
