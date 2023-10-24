package repository

const (
	QueryChildCount = `
SELECT COUNT(*) as rowCount 
FROM transfer_lines 
	INNER JOIN items ON transfer_lines.item_id = items.id  
	INNER JOIN units ON transfer_lines.item_unit_id = units.id 
WHERE 1=1 `

	// TODO
	QueryChild = `
SELECT 
	transfer_lines.id, 
	transfer_lines.parent_id, 
	item_types.code AS item_type, 
	items.no AS item_no, 
	items.description AS item_name, 
	units.code AS item_unit, 
	transfer_lines.input_quantity, 
	transfer_lines.item_quantity_unit, 
	transfer_lines.item_quantity, 
	transfer_lines.item_tempory, 
	transfer_lines.item_uuid, 
	transfer_lines.item_unit_uuid  
FROM transfer_lines 
	INNER JOIN items ON transfer_lines.item_uuid = items.id  
	INNER JOIN units ON transfer_lines.item_unit_uuid = units.id 
`

	QueryChildJoins = `
INNER JOIN items ON transfer_lines.item_id = items.id  
INNER JOIN units ON transfer_lines.item_unit_id = units.id `

	QueryChildSuggestion = ``
)
