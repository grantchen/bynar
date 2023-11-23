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
			   COALESCE(lines_t.Count, 0) AS Count
		FROM payments
				 LEFT JOIN (SELECT COUNT(payment_lines.id) AS Count,
				                   payment_lines.parent_id AS parent_id
							FROM payment_lines
							WHERE 2=2
							GROUP BY parent_id) lines_t
						   ON lines_t.parent_id = payments.id
		WHERE 1=1
		`

	// QueryParentJoins is a query for getting parent rows with joins
	QueryParentJoins = ``
)
