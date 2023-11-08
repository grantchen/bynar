package repository

const (
	// QueryChildCount is a query for child count
	QueryChildCount = `
SELECT COUNT(*) as Count 
FROM transfer_lines 
	INNER JOIN items ON transfer_lines.item_id = items.id  
	INNER JOIN units ON transfer_lines.item_unit_id = units.id`

	// QueryChild is a query for child
	QueryChild = `
SELECT 
	transfer_lines.*,
	CONCAT (transfer_lines.id, '-line') as id
FROM transfer_lines 
	INNER JOIN items ON transfer_lines.item_id = items.id  
	INNER JOIN units ON transfer_lines.item_unit_id = units.id 
`

	// QueryChildJoins is a query for child joins
	QueryChildJoins = `
INNER JOIN items ON transfer_lines.item_id = items.id  
INNER JOIN units ON transfer_lines.item_unit_id = units.id `

	// QueryChildSuggestion is a query for child suggestion
	QueryChildSuggestion = ``
)
