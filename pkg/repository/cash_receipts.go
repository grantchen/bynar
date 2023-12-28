package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/logger"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/models"
)

// CashReceipts
type cashReceiptRepository struct {
	conn          *sql.DB
	tableName     string
	lineTableName string
}

// return new cashReceiptRepository
func _(conn *sql.DB, tableName, lineTableName string) CashReceiptRepository {
	return &cashReceiptRepository{
		conn:          conn,
		tableName:     tableName,
		lineTableName: lineTableName,
	}
}

func (s *cashReceiptRepository) GetDocID(id interface{}) (docID int, err error) {
	logger.Debug("get document id", id)

	query := `
	SELECT document_id 
	FROM ` + s.tableName + `
	WHERE id = ?
	`

	err = s.conn.QueryRow(query, id).Scan(&docID)
	return
}

func (s *cashReceiptRepository) GetStatus(id interface{}) (status int, err error) {
	logger.Debug("get status", id)

	query := `
	SELECT status 
	FROM ` + s.tableName + `
	WHERE id = ?
	`

	err = s.conn.QueryRow(query, id).Scan(&status)
	return
}

var cashReceiptFields = []string{
	"id",
	"batch_id",
	"document_id",
	"document_no",
	"transaction_no",
	"store_id",
	"document_date",
	"posting_date",
	"entry_date",
	"account_type",
	"account_id",
	"balance_account_type",
	"balance_account_id",
	"amount",
	"amount_lcy",
	"currency_value",
	"user_group_id",
	"bank_id",
	"status",
}

func (s *cashReceiptRepository) Get(tx *sql.Tx, id interface{}) (m *models.CashReceipt, err error) {
	logger.Debug("get model", id)

	query := `
	SELECT ` + strings.Join(cashReceiptFields, ", ") + `
	FROM ` + s.tableName + `
	WHERE id = ?
	`
	m = &models.CashReceipt{}
	err = tx.QueryRow(query, id).Scan(
		&m.ID,
		&m.BatchID,
		&m.DocumentID,
		&m.DocumentNo,
		&m.TransactionNo,
		&m.StoreID,
		&m.DocumentDate,
		&m.PostingDate,
		&m.EntryDate,
		&m.AccountType,
		&m.AccountID,
		&m.BalanceAccountType,
		&m.BalanceAccountID,
		&m.Amount,
		&m.AmountLcy,
		&m.CurrencyValue,
		&m.UserGroupID,
		&m.BankID,
		&m.Status,
	)

	return
}

func (s *cashReceiptRepository) Save(tx *sql.Tx, m *models.CashReceipt) (err error) {
	logger.Debug("save model", m.ID)

	query := `
	UPDATE ` + s.tableName + `
	SET ` + strings.Join(cashReceiptFields[1:], " = ?, ") + " = ? " + `
	WHERE id = ?
	`

	_, err = tx.Exec(query,
		&m.BatchID,
		&m.DocumentID,
		&m.DocumentNo,
		&m.TransactionNo,
		&m.StoreID,
		&m.DocumentDate,
		&m.PostingDate,
		&m.EntryDate,
		&m.AccountType,
		&m.AccountID,
		&m.BalanceAccountType,
		&m.BalanceAccountID,
		&m.Amount,
		&m.AmountLcy,
		&m.CurrencyValue,
		&m.UserGroupID,
		&m.BankID,
		&m.Status,
		&m.ID,
	)

	return
}

var cashReceiptLineFields = []string{
	"id",
	"parent_id",
	"applies_document_type",
	"applies_document_id",
	"amount",
	"amount_lcy",
	"applied",
}

func (s *cashReceiptRepository) GetLines(tx *sql.Tx, parentID interface{}) ([]*models.CashReceiptLine, error) {
	query := `
	SELECT ` + strings.Join(cashReceiptLineFields, ",") + `
	FROM ` + s.lineTableName + `
	WHERE parent_id = ?
	`

	res := make([]*models.CashReceiptLine, 0, 10)
	rows, err := tx.Query(query, parentID)
	if err != nil {
		return nil, fmt.Errorf("do query: [%w], query: %s, id: %v", err, query, parentID)
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	for rows.Next() {
		line := &models.CashReceiptLine{}
		if err := rows.Scan(
			&line.ID,
			&line.ParentID,
			&line.AppliesDocumentType,
			&line.AppliesDocumentID,
			&line.Amount,
			&line.AmountLcy,
			&line.Applied,
		); err != nil {
			return res, fmt.Errorf("rows scan: [%w]", err)
		}

		res = append(res, line)
	}

	return res, nil
}

func (s *cashReceiptRepository) SaveLine(tx *sql.Tx, l *models.CashReceiptLine) (err error) {
	logger.Debug("save line", l.ID)

	query := `
	UPDATE ` + s.lineTableName + `
	SET ` + strings.Join(cashReceiptLineFields[1:], " = ?,") + " = ? " + `
	WHERE id = ?
	`

	_, err = tx.Exec(query,
		&l.ParentID,
		&l.AppliesDocumentType,
		&l.AppliesDocumentID,
		&l.Amount,
		&l.AmountLcy,
		&l.Applied,
		&l.ID,
	)

	return
}
