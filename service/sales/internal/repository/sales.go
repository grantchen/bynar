package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/logger"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/models"
)

type saleRepository struct {
	conn *sql.DB
}

func NewSaleRepository(conn *sql.DB) SaleRepository {
	return &saleRepository{conn: conn}
}

func (s *saleRepository) GetDocID(id interface{}) (docID int, err error) {

	logger.Debug("get document id", id)
	query := `
	SELECT document_id 
	FROM sales
	WHERE id = ?
	`

	err = s.conn.QueryRow(query, id).Scan(&docID)
	return
}

func (s *saleRepository) GetStatus(id interface{}) (status int, err error) {
	query := `
	SELECT status 
	FROM sales
	WHERE id = ?
	`

	err = s.conn.QueryRow(query, id).Scan(&status)
	return
}

var saleFields = []string{
	"id",
	"document_id",
	"document_no",
	"transaction_no",
	"store_id",
	"posting_date",
	"status",
	"currency_id",
	"currency_value",
	"discount_id",

	"subtotal_exclusive_vat",
	"total_discount",
	"total_exclusive_vat",
	"total_vat",
	"total_inclusive_vat",

	"subtotal_exclusive_vat_lcy",
	"total_discount_lcy",
	"total_exclusive_vat_lcy",
	"total_vat_lcy",
	"total_inclusive_vat_lcy",
}

func (s *saleRepository) GetSale(tx *sql.Tx, id interface{}) (m *models.Sale, err error) {
	query := `
	SELECT ` + strings.Join(saleFields, ", ") + `
	FROM sales
	WHERE id = ?
	`
	m = &models.Sale{}
	err = tx.QueryRow(query, id).Scan(&m.ID, &m.DocumentID, &m.DocumentNo, &m.TransactionNo, &m.StoreID, &m.PostingDate, &m.Status, &m.CurrencyID, &m.CurrencyValue, &m.DiscountID,
		&m.SubtotalExclusiveVat, &m.TotalDiscount, &m.TotalExclusiveVat, &m.TotalVat, &m.TotalInclusiveVat,
		&m.SubtotalExclusiveVatLcy, &m.TotalDiscountLcy, &m.TotalExclusiveVatLcy, &m.TotalVatLcy, &m.TotalInclusiveVatLcy)

	return
}

func (s *saleRepository) SaveSale(tx *sql.Tx, m *models.Sale) (err error) {
	query := `
	UPDATE sales
	SET ` + strings.Join(saleFields[1:], " = ?, ") + " = ? " + `
	WHERE id = ?
	`

	_, err = tx.Exec(query, m.DocumentID, m.DocumentNo, m.TransactionNo, m.StoreID, m.PostingDate, m.Status, m.CurrencyID, m.CurrencyValue, m.DiscountID,
		m.SubtotalExclusiveVat, m.TotalDiscount, m.TotalExclusiveVat, m.TotalVat, m.TotalInclusiveVat,
		m.SubtotalExclusiveVatLcy, m.TotalDiscountLcy, m.TotalExclusiveVatLcy, m.TotalVatLcy, m.TotalInclusiveVatLcy, m.ID)

	return
}

var saleLineFields = []string{
	"id",
	"parent_id",
	"item_type",
	"item_id",

	"location_id",
	"input_quantity",
	"item_unit_value",
	"quantity",

	"item_unit_id",
	"discount_id",
	"vat_id",
	"subtotal_exclusive_vat",

	"total_discount",
	"total_exclusive_vat",
	"total_vat",
	"total_inclusive_vat",

	"subtotal_exclusive_vat_lcy",
	"total_discount_lcy",
	"total_exclusive_vat_lcy",
	"total_vat_lcy",
	"total_inclusive_vat_lcy",
}

func (s *saleRepository) GetSaleLines(tx *sql.Tx, parentID interface{}) ([]*models.SaleLine, error) {
	query := `
	SELECT ` + strings.Join(saleLineFields, ",") + `
	FROM sale_lines
	WHERE parent_id = ?
	`

	res := make([]*models.SaleLine, 0, 10)
	rows, err := tx.Query(query, parentID)
	if err != nil {
		return nil, fmt.Errorf("do query: [%w], query: %s, id: %v", err, query, parentID)
	}
	defer rows.Close()

	for rows.Next() {
		line := &models.SaleLine{}
		if err := rows.Scan(
			&line.ID,
			&line.ParentID,
			&line.ItemType,
			&line.ItemID,

			&line.LocationID,
			&line.InputQuantity,
			&line.ItemUnitValue,
			&line.Quantity,

			&line.ItemUnitID,
			&line.DiscountID,
			&line.VatID,
			&line.SubtotalExclusiveVat,

			&line.TotalDiscount,
			&line.TotalExclusiveVat,
			&line.TotalVat,
			&line.TotalInclusiveVat,

			&line.SubtotalExclusiveVatLcy,
			&line.TotalDiscountLcy,
			&line.TotalExclusiveVatLcy,
			&line.TotalVatLcy,
			&line.TotalInclusiveVatLcy); err != nil {
			return res, fmt.Errorf("rows scan: [%w]", err)
		}

		res = append(res, line)
	}

	return res, nil
}

func (s *saleRepository) SaveSaleLine(tx *sql.Tx, l *models.SaleLine) (err error) {
	logger.Debug("save sale line", l.ID, "unit value", l.ItemUnitValue)

	query := `
	UPDATE sale_lines
	SET ` + strings.Join(saleLineFields[1:], " = ?,") + " = ? " + `
	WHERE id = ?
	`

	_, err = tx.Exec(query,
		l.ParentID,
		l.ItemType,
		l.ItemID,

		l.LocationID,
		l.InputQuantity,
		l.ItemUnitValue,
		l.Quantity,

		l.ItemUnitID,
		l.DiscountID,
		l.VatID,
		l.SubtotalExclusiveVat,

		l.TotalDiscount,
		l.TotalExclusiveVat,
		l.TotalVat,
		l.TotalInclusiveVat,

		l.SubtotalExclusiveVatLcy,
		l.TotalDiscountLcy,
		l.TotalExclusiveVatLcy,
		l.TotalVatLcy,
		l.TotalInclusiveVatLcy,
		l.ID,
	)

	return
}
