package repositories

import (
	"database/sql"
	"fmt"
	"strings"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/payments/internal/models"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/logger"
)

var procurementFields = []string{
	"id",
	"document_id",
	"document_no",
	"transaction_no",
	"posting_date",
	"status",
	"discount_id",
	"currency_id",
	"currency_value",

	"total_discount",
	"total_exclusive_vat",
	"total_vat",
	"total_inclusive_vat",

	"subtotal_exclusive_vat",
	"subtotal_exclusive_vat_lcy",

	"total_discount_lcy",
	"total_exclusive_vat_lcy",
	"total_vat_lcy",
	"total_inclusive_vat_lcy",
	"paid",
	"remaining",
	"paid_status",
	// "entry_date",
	// "shipment_date",
	// "project_id",
	// "department_id",
	// "contract_id",
	// "user_group_id",
	// "store_id",
	// "document_date",
	// "budget_id",
	// "vendor_id",
	// "vendor_invoice_no",
	// "purchaser_id",
	// "responsibility_center_id",
	// "payment_terms_id",
	// "payment_method_id",
	// "transaction_type_id",
	// "payment_discount",
	// "shipment_method_id",
	// "payment_reference",
	// "creditor_no",
	// "on_hold",
	// "transaction_specification_id",
	// "transport_method_id",
	// "entry_point_id",
	// "campaign_id",
	// "area_id",
	// "vendor_shipment_no",
}

type procurementRepository struct {
	conn *sql.DB
}

func NewProcurementRepository(conn *sql.DB) ProcurementRepository {
	return &procurementRepository{conn: conn}
}

func (s *procurementRepository) GetDocID(id interface{}) (docID int, err error) {
	logger.Debug("get document id", id)
	query := `
	SELECT document_id 
	FROM procurements
	WHERE id = ?
	`

	err = s.conn.QueryRow(query, id).Scan(&docID)
	return
}

func (s *procurementRepository) GetStatus(id interface{}) (status int, err error) {
	query := `
	SELECT status 
	FROM procurements
	WHERE id = ?
	`

	err = s.conn.QueryRow(query, id).Scan(&status)
	return
}

func (p *procurementRepository) GetProcurement(tx *sql.Tx, id interface{}) (m *models.Procurement, err error) {
	query := `
	SELECT ` + strings.Join(procurementFields, ", ") + `
	FROM procurements
	WHERE id = ?
	`
	m = &models.Procurement{}
	err = tx.QueryRow(query, id).Scan(
		&m.ID,
		&m.DocumentID,
		&m.DocumentNo,
		&m.TransactionNo,
		&m.PostingDate,
		&m.Status,
		&m.InvoiceDiscountID,
		&m.CurrencyID,
		&m.CurrencyValue,
		&m.TotalDiscount,
		&m.TotalExclusiveVat,
		&m.TotalVat,
		&m.TotalInclusiveVat,
		&m.SubtotalExclusiveVat,
		&m.SubtotalExclusiveVatLcy,
		&m.TotalDiscountLcy,
		&m.TotalExclusiveVatLcy,
		&m.TotalVatLcy,
		&m.SubtotalExclusiveVatLcy,
		&m.Paid,
		&m.Remaining,
		&m.PaidStatus,
	)

	return
}

func (p *procurementRepository) SaveProcurement(tx *sql.Tx, m *models.Procurement) (err error) {
	query := `
	UPDATE procurements
	SET ` + strings.Join(procurementFields[1:], " = ?, ") + " = ? " + `
	WHERE id = ?
	`

	_, err = tx.Exec(query,
		m.DocumentID,
		m.DocumentNo,
		m.TransactionNo,
		m.PostingDate,
		m.Status,
		m.InvoiceDiscountID,
		m.CurrencyID,
		m.CurrencyValue,
		m.TotalDiscount,
		m.TotalExclusiveVat,
		m.TotalVat,
		m.TotalInclusiveVat,
		m.SubtotalExclusiveVat,
		m.SubtotalExclusiveVatLcy,
		m.TotalDiscountLcy,
		m.TotalExclusiveVatLcy,
		m.TotalVatLcy,
		m.SubtotalExclusiveVatLcy,
		m.Paid,
		m.Remaining,
		m.PaidStatus,
		m.ID,
	)

	return
}

var procurementLineFields = []string{
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
	// "tax_area_id",
	"vat_id",
	// "quantity_assign",
	// "quantity_assigned",
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

func (p *procurementRepository) GetProcurementLines(tx *sql.Tx, id interface{}) ([]*models.ProcurementLine, error) {
	query := `
	SELECT ` + strings.Join(procurementLineFields, ",") + `
	FROM procurement_lines
	WHERE parent_id = ?
	`

	res := make([]*models.ProcurementLine, 0, 10)
	rows, err := tx.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("do query: [%w], query: %s, id: %v", err, query, id)
	}
	defer rows.Close()

	for rows.Next() {
		prLine := &models.ProcurementLine{}
		if err := rows.Scan(
			&prLine.ID,
			&prLine.ParentID,
			&prLine.ItemType,
			&prLine.ItemID,
			&prLine.LocationID,
			&prLine.InputQuantity,
			&prLine.ItemUnitValue,
			&prLine.Quantity,
			&prLine.ItemUnitID,
			&prLine.DiscountID,
			&prLine.VatID,
			&prLine.SubtotalExclusiveVat,
			&prLine.TotalDiscount,
			&prLine.TotalExclusiveVat,
			&prLine.TotalVat,
			&prLine.TotalInclusiveVat,
			&prLine.SubtotalExclusiveVatLcy,
			&prLine.TotalDiscountLcy,
			&prLine.TotalExclusiveVatLcy,
			&prLine.TotalVatLcy,
			&prLine.TotalInclusiveVatLcy); err != nil {
			return res, fmt.Errorf("rows scan: [%w]", err)
		}

		res = append(res, prLine)
	}

	return res, nil
}

func (p *procurementRepository) SaveProcurementLine(tx *sql.Tx, prLine *models.ProcurementLine) (err error) {
	logger.Debug("save procurement line", prLine.ID, "unit value", prLine.ItemUnitValue)

	query := `
	UPDATE procurement_lines
	SET ` + strings.Join(procurementLineFields[1:], " = ?,") + " = ? " + `
	WHERE id = ?
	`

	_, err = tx.Exec(query,
		prLine.ParentID,
		prLine.ItemType,
		prLine.ItemID,
		prLine.LocationID,
		prLine.InputQuantity,
		prLine.ItemUnitValue,
		prLine.Quantity,
		prLine.ItemUnitID,
		prLine.DiscountID,
		prLine.VatID,
		prLine.SubtotalExclusiveVat,
		prLine.TotalDiscount,
		prLine.TotalExclusiveVat,
		prLine.TotalVat,
		prLine.TotalInclusiveVat,
		prLine.SubtotalExclusiveVatLcy,
		prLine.TotalDiscountLcy,
		prLine.TotalExclusiveVatLcy,
		prLine.TotalVatLcy,
		prLine.TotalInclusiveVatLcy,
		prLine.ID)

	return
}
