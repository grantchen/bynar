package repository

const (
	QueryChildCount = `
SELECT COUNT(*) as Count 
FROM sale_lines 
	INNER JOIN items ON sale_lines.item_id = items.id  
	INNER JOIN units ON sale_lines.item_unit_id = units.id`

	QueryChild = `
SELECT 
	sale_lines.*,
	CONCAT (sale_lines.id, '-line') as id
FROM sale_lines 
	INNER JOIN items ON sale_lines.item_id = items.id  
	INNER JOIN units ON sale_lines.item_unit_id = units.id 
`

	QueryChildJoins = `
INNER JOIN items ON sale_lines.item_id = items.id  
INNER JOIN units ON sale_lines.item_unit_id = units.id `

	QueryChildSuggestion = ``
)
