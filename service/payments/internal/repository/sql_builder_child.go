package repository

const (
	QueryChildCount = `
		SELECT COUNT(*) as Count
		FROM payment_lines 
		`

	QueryChild = `
	SELECT
	payment_lines.*,
	CONCAT (payment_lines.id, '-line') as id
	FROM payment_lines `

	QueryChildJoins      = ` INNER JOIN payments ON payment_lines.parent_id = payments.id `
	QueryChildSuggestion = ``
)
