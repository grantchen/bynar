package repository

const (
	QueryChildCount = `
		SELECT COUNT(*) as Count
		FROM procurement_lines WHERE 1=1
		`

	QueryChild = `
	SELECT CONCAT (id, '-line') as id,
	parent_id,
	item_type,
	item_id,
	location_id,
	input_quantity,
	item_unit_value,
	quantity,
	item_unit_id,
	discount_id,
	tax_area_id,
	vat_id,
	quantity_assign,
	quantity_assigned,
	subtotal_exclusive_vat,
	total_discount,
	total_exclusive_vat,
	total_vat,
	total_inclusive_vat,
	subtotal_exclusive_vat_lcy,
	total_discount_lcy,
	total_exclusive_vat_lcy,
	total_vat_lcy,
	total_inclusive_vat_lcy
	FROM procurement_lines WHERE 1=1
	`

	QueryChildJoins = ` INNER JOIN procurements ON parent_id = procurements.id `
)
