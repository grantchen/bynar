package repository

const (
	QueryParentCount = `
		SELECT COUNT(id) as Count
		FROM procurements
		WHERE 1=1
		`

	QueryParent = `
		SELECT procurements.id,
			   document_id,
			   document_no,
			   transaction_no,
			   store_id,
			   document_date,
			   posting_date,
			   entry_date,
			   shipment_date,
			   project_id,
			   department_id,
			   contract_id,
			   user_group_id,
			   status,
			   budget_id,
			   currency_id,
			   currency_value,
			   vendor_id,
			   vendor_invoice_no,
			   purchaser_id,
			   responsibility_center_id,
			   payment_terms_id,
			   payment_method_id,
			   transaction_type_id,
			   payment_discount,
			   shipment_method_id,
			   payment_reference,
			   creditor_no,
			   on_hold,
			   transaction_specification_id,
			   transport_method_id,
			   entry_point_id,
			   campaign_id,
			   area_id,
			   vendor_shipment_no,
			   procurements.subtotal_exclusive_vat,
			   procurements.total_discount,
			   procurements.total_exclusive_vat,
			   procurements.total_vat,
			   procurements.total_inclusive_vat,
			   procurements.subtotal_exclusive_vat_lcy,
			   procurements.total_discount_lcy,
			   procurements.total_exclusive_vat_lcy,
			   procurements.total_vat_lcy,
			   procurements.total_inclusive_vat_lcy,
			   COALESCE(lines_t.Count, 0) AS Count
		FROM procurements
				 LEFT JOIN (SELECT COUNT(procurement_lines.id) AS Count,
								   procurement_lines.parent_id AS parent_id
							FROM procurement_lines
							WHERE 2=2
							GROUP BY parent_id) lines_t
						   ON lines_t.parent_id = procurements.id
		WHERE 1=1
		`

	// empty
	QueryParentJoins = ``
)
