package repository

const (
	QueryParentCount = `
	SELECT COUNT(user_groups.id) as rowCount 
	FROM user_groups 
	WHERE 1=1 `

	QueryParent = `
	SELECT user_groups.id,
	user_groups.code,
	user_groups.description,
	user_groups.status
	FROM user_groups
	WHERE 1=1 
	`

	// empty
	QueryParentJoins = `
	`

	QueryParentBuild = `
	SELECT * FROM user_groups
	`
)
