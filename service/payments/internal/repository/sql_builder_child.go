package repository

const (
	QueryChildCount = `
		SELECT COUNT(*) as Count
		FROM payment_lines 
		INNER JOIN payments ON payment_lines.parent_id = payments.id
		`

	QueryChild = `
	SELECT CONCAT (payment_lines.id, '-line') as id, 
	payment_lines.parent_id,
	payment_lines.applies_document_type,
	payment_lines.applies_document_id,
	payment_lines.payment_type_id,
	payment_lines.amount,
	payment_lines.amount_lcy,
	payment_lines.applied
	FROM payment_lines 
		INNER JOIN payments ON payment_lines.parent_id = payments.id`

	QueryChildJoins = `
	INNER JOIN payments ON payment_lines.parent_id = payments.id
	`
)
