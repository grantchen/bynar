package repository

const (
	// QueryParentCount is a query for getting count of parent rows
	QueryParentCount = `
		SELECT COUNT(user_groups.id) as Count 
		FROM user_groups 
		WHERE 1=1
		`

	// QueryParent is a query for getting parent rows
	QueryParent = `
		SELECT user_groups.id,
			   user_groups.code,
			   user_groups.description,
			   user_groups.status,
			   COUNT(user_group_lines.id) AS Count
		FROM user_groups
				 LEFT JOIN user_group_lines ON user_group_lines.parent_id = user_groups.id
		WHERE 1=1
		GROUP BY user_groups.id
		`

	// QueryParentJoins is a query for getting parent rows with joins
	QueryParentJoins = ``
)
