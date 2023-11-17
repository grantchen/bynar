package repository

const (
	QueryParentCount = `
		SELECT COUNT(id) as Count
		FROM procurements
		WHERE 1=1
		`

	QueryParent = `
		SELECT id,
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
			   subtotal_exclusive_vat,
			   total_discount,
			   total_exclusive_vat,
			   total_vat,
			   total_inclusive_vat,
			   subtotal_exclusive_vat_lcy,
			   total_discount_lcy,
			   total_exclusive_vat_lcy,
			   total_vat_lcy,
			   total_inclusive_vat_lcy,
			   COUNT(procurement_lines.id) AS Count
		FROM procurements
				 LEFT JOIN procurement_lines ON procurement_lines.parent_id = procurements.id
		WHERE 1=1
		GROUP BY procurements.id
		`

	// empty
	QueryParentJoins = ``
)
