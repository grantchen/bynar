package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/payments/internal/models"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/logger"
)

type paymentRepository struct {
	conn          *sql.DB
	tableName     string
	lineTableName string
}

func NewPayment(conn *sql.DB, tableName, lineTableName string) PaymentRepository {
	return &paymentRepository{
		conn:          conn,
		tableName:     tableName,
		lineTableName: lineTableName,
	}
}

func (s *paymentRepository) GetDocID(id interface{}) (docID int, err error) {
	logger.Debug("get document id", id)
	query := `
	SELECT document_id 
	FROM ` + s.tableName + `
	WHERE id = ?
	`

	err = s.conn.QueryRow(query, id).Scan(&docID)
	return
}

func (s *paymentRepository) GetStatus(id interface{}) (status int, err error) {
	query := `
	SELECT status 
	FROM ` + s.tableName + `
	WHERE id = ?
	`

	err = s.conn.QueryRow(query, id).Scan(&status)
	return
}

var paymentFields = []string{
	"id",
	"batch_id",
	"document_id",
	"document_no",
	"external_document_no",
	"transaction_no",
	"store_id",
	"document_date",
	"posting_date",
	"entry_date",
	"account_type",
	"account_id",
	"recipient_bank_account_id",
	"balance_account_type",
	"balance_account_id",
	"amount",
	"amount_lcy",
	"currency_id",
	"currency_value",
	"user_group_id",
	"payment_method_id",
	"payment_reference",
	"creditor_no",
	"bank_payment_type_id",
	"bank_id",
	"paid",
	"remaining",
	"paid_status",
	"status",
}

func (s *paymentRepository) Get(tx *sql.Tx, id interface{}) (m *models.Payment, err error) {
	query := `
	SELECT ` + strings.Join(paymentFields, ", ") + `
	FROM ` + s.tableName + `
	WHERE id = ?
	`
	m = &models.Payment{}
	err = tx.QueryRow(query, id).Scan(
		&m.ID,
		&m.BatchID,
		&m.DocumentID,
		&m.DocumentNo,
		&m.ExternalDocumentNo,
		&m.TransactionNo,
		&m.StoreID,
		&m.DocumentDate,
		&m.PostingDate,
		&m.EntryDate,
		&m.AccountType,
		&m.AccountID,
		&m.RecipientBankAccountID,
		&m.BalanceAccountType,
		&m.BalanceAccountID,
		&m.Amount,
		&m.AmountLcy,
		&m.CurrencyID,
		&m.CurrencyValue,
		&m.UserGroupID,
		&m.PaymentMethodID,
		&m.PaymentReference,
		&m.CreditorNo,
		&m.BankPaymentTypeID,
		&m.BankID,
		&m.Paid,
		&m.Remaining,
		&m.PaidStatus,
		&m.Status,
	)

	return
}

func (s *paymentRepository) Save(tx *sql.Tx, m *models.Payment) (err error) {
	query := `
	UPDATE ` + s.tableName + `
	SET ` + strings.Join(paymentFields[1:], " = ?, ") + " = ? " + `
	WHERE id = ?
	`

	_, err = tx.Exec(query,
		&m.BatchID,
		&m.DocumentID,
		&m.DocumentNo,
		&m.ExternalDocumentNo,
		&m.TransactionNo,
		&m.StoreID,
		&m.DocumentDate,
		&m.PostingDate,
		&m.EntryDate,
		&m.AccountType,
		&m.AccountID,
		&m.RecipientBankAccountID,
		&m.BalanceAccountType,
		&m.BalanceAccountID,
		&m.Amount,
		&m.AmountLcy,
		&m.CurrencyID,
		&m.CurrencyValue,
		&m.UserGroupID,
		&m.PaymentMethodID,
		&m.PaymentReference,
		&m.CreditorNo,
		&m.BankPaymentTypeID,
		&m.BankID,
		&m.Paid,
		&m.Remaining,
		&m.PaidStatus,
		&m.Status,
		&m.ID,
	)

	return
}

var paymentLineFields = []string{
	"id",
	"parent_id",
	"applies_document_type",
	"applies_document_id",
	"payment_type_id",
	"amount",
	"amount_lcy",
	"applied",
}

func (s *paymentRepository) GetLines(tx *sql.Tx, parentID interface{}, applied int) ([]*models.PaymentLine, error) {
	query := `
	SELECT ` + strings.Join(paymentLineFields, ",") + `
	FROM ` + s.lineTableName + `
	WHERE parent_id = ? AND applied = ?
	`

	res := make([]*models.PaymentLine, 0, 10)
	rows, err := tx.Query(query, parentID, applied)
	if err != nil {
		return nil, fmt.Errorf("do query: [%w], query: %s, id: %v", err, query, parentID)
	}
	defer rows.Close()

	for rows.Next() {
		line := &models.PaymentLine{}
		if err := rows.Scan(
			&line.ID,
			&line.ParentID,
			&line.AppliesDocumentType,
			&line.AppliesDocumentID,
			&line.PaymentTypeID,
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

func (s *paymentRepository) SaveLine(tx *sql.Tx, l *models.PaymentLine) (err error) {
	logger.Debug("save line", l.ID)

	query := `
	UPDATE ` + s.lineTableName + `
	SET ` + strings.Join(paymentLineFields[1:], " = ?,") + " = ? " + `
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
