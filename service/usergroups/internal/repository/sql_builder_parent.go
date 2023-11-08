package repository

const (
	// QueryParentCount is a query for getting count of parent rows
	QueryParentCount = `
	SELECT COUNT(user_groups.id) as Count 
	FROM user_groups 
	WHERE 1=1 `

	// QueryParent is a query for getting parent rows
	QueryParent = `
	SELECT user_groups.id,
	user_groups.code,
	user_groups.description,
	user_groups.status
	FROM user_groups
	WHERE 1=1 
	`

	// QueryParentJoins is a query for getting parent rows with joins
	QueryParentJoins = `
	`
)
