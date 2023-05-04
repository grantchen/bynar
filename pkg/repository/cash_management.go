package repository

import (
	"database/sql"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/models"
)

// CashManagement
type cashManagementRepository struct {
	conn *sql.DB
}

func NewCashManagementRepository(conn *sql.DB) CashManagementRepository {
	return &cashManagementRepository{conn: conn}
}

func (c *cashManagementRepository) Get(bankID int) (m *models.CashManagement, err error) {
	query := `
	SELECT id, type, bank_id, amount, currency_id
	FROM cash_managements
	WHERE bank_id = ?
	`
	m = &models.CashManagement{}
	err = c.conn.QueryRow(query, bankID).Scan(&m.Id, &m.Type, &m.Bank_Id, &m.Amount, &m.Currency_Id)

	return
}

func (c *cashManagementRepository) Update(tx *sql.Tx, m *models.CashManagement) (err error) {
	query := `
	UPDATE cash_managements
	SET amount = ?
	WHERE id = ?
	`

	_, err = tx.Exec(query, m.Amount, m.Id)

	return
}
