package repository

const (
	// QueryChildCount is a query for getting child count
	QueryChildCount = `
		SELECT COUNT(*) as Count 
		FROM sale_lines 
			INNER JOIN items ON sale_lines.item_id = items.id  
			INNER JOIN units ON sale_lines.item_unit_id = units.id
		WHERE 2=2
		`

	// QueryChild is a query for getting child
	QueryChild = `
		SELECT 
			sale_lines.*,
			CONCAT (sale_lines.id, '-line') as id
		FROM sale_lines 
			INNER JOIN items ON sale_lines.item_id = items.id  
			INNER JOIN units ON sale_lines.item_unit_id = units.id
		WHERE 2=2
		`

	// QueryChildJoins is a query for getting child joins
	QueryChildJoins = `
INNER JOIN items ON sale_lines.item_id = items.id  
INNER JOIN units ON sale_lines.item_unit_id = units.id `

	// QueryChildSuggestion is a query for getting child suggestion
	QueryChildSuggestion = ``
)
