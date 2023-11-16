package repository

const (
	// QueryParentCount is a query for getting count of parent rows
	QueryParentCount = `
	SELECT COUNT(payments.id) as Count
	FROM payments WHERE 1=1
	`
	// QueryParent is a query for getting parent rows
	QueryParent = `
	SELECT payments.*
	FROM payments
	`

	// QueryParentJoins is a query for getting parent rows with joins
	QueryParentJoins = `
	`
)
