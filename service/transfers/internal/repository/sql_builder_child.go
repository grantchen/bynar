package repository

const (
	QueryChildCount = `
SELECT COUNT(*) as rowCount 
FROM transfer_lines 
	INNER JOIN items ON transfer_lines.item_id = items.id  
	INNER JOIN units ON transfer_lines.item_unit_id = units.id`

	QueryChild = `
SELECT 
	transfer_lines.id, 
	transfer_lines.parent_id, 
	items.no AS item_no, 
	items.description AS item_name, 
	units.code AS item_unit, 
	transfer_lines.input_quantity, 
	transfer_lines.item_unit_value, 
	transfer_lines.quantity, 
	transfer_lines.item_unit_id, 
	transfer_lines.shipment_date, 
	transfer_lines.receipt_date  
FROM transfer_lines 
	INNER JOIN items ON transfer_lines.item_id = items.id  
	INNER JOIN units ON transfer_lines.item_unit_id = units.id 
`

	QueryChildJoins = `
INNER JOIN items ON transfer_lines.item_id = items.id  
INNER JOIN units ON transfer_lines.item_unit_id = units.id `

	QueryChildSuggestion = ``
)
