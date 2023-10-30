package repository

const (
	QueryParentCount = `
	SELECT COUNT(payments.id) as Count 
	FROM payments
	`

	QueryParent = `
	SELECT payments.*
	FROM payments
	`

	// empty
	QueryParentJoins = `
	`
)
