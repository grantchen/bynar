package repository

import (
	"database/sql"
	"fmt"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/logger"
	"strings"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/models"
)

// saleRepository implements SaleRepository
type saleRepository struct {
	conn *sql.DB
}

// NewSaleRepository returns new SaleRepository
func NewSaleRepository(conn *sql.DB) SaleRepository {
	return &saleRepository{conn: conn}
}

// GetDocID returns document id by sale id
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

// GetStatus returns status by sale id
func (s *saleRepository) GetStatus(id interface{}) (status int, err error) {
	query := `
	SELECT status 
	FROM sales
	WHERE id = ?
	`

	err = s.conn.QueryRow(query, id).Scan(&status)
	return
}

// saleFields is a slice of sale field names
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

// GetSale returns sale by id
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

// SaveSale saves sale
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

// SaleLineFieldNames is a slice of sale line field names
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

// GetSaleLines returns sale lines by parent id
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

// SaveSaleLine saves sale line
func (s *saleRepository) SaveSaleLine(tx *sql.Tx, l *models.SaleLine) (err error) {
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
