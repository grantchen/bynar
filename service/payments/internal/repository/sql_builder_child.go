package repository

const (
	// QueryChildCount is a query for getting count of child rows
	QueryChildCount = `
		SELECT COUNT(*) as Count
		FROM payment_lines 
		`
	// QueryChild is a query for getting child rows
	QueryChild = `
	SELECT
	payment_lines.*,
	CONCAT (payment_lines.id, '-line') as id
	FROM payment_lines `

	// QueryChildJoins is a query for getting child rows with joins
	QueryChildJoins = ` INNER JOIN payments ON payment_lines.parent_id = payments.id `

	// QueryChildSuggestion is a query for getting child rows suggestion
	QueryChildSuggestion = ``
)
