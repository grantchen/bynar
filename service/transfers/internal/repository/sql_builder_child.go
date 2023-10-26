package repository

const (
	QueryChildCount = `
SELECT COUNT(*) as Count 
FROM transfer_lines 
	INNER JOIN items ON transfer_lines.item_id = items.id  
	INNER JOIN units ON transfer_lines.item_unit_id = units.id`

	QueryChild = `
SELECT 
	transfer_lines.*,
	CONCAT (transfer_lines.id, '-line') as id
FROM transfer_lines 
	INNER JOIN items ON transfer_lines.item_id = items.id  
	INNER JOIN units ON transfer_lines.item_unit_id = units.id 
`

	QueryChildJoins = `
INNER JOIN items ON transfer_lines.item_id = items.id  
INNER JOIN units ON transfer_lines.item_unit_id = units.id `

	QueryChildSuggestion = ``
)
