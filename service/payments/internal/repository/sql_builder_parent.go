package repository

const (
	// QueryParentCount is a query for getting count of parent rows
	QueryParentCount = `
		SELECT COUNT(payments.id) as Count
		FROM payments
		WHERE 1=1
		`

	// QueryParent is a query for getting parent rows
	QueryParent = `
		SELECT payments.*,
			   COUNT(payment_lines.id) AS Count
		FROM payments
				 LEFT JOIN payment_lines ON payment_lines.parent_id = payments.id
		WHERE 1=1
		GROUP BY payments.id
		`

	// QueryParentJoins is a query for getting parent rows with joins
	QueryParentJoins = ``
)
