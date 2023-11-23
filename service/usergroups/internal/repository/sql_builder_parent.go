package repository

const (
	// QueryParentCount is a query for getting count of parent rows
	QueryParentCount = `
		SELECT COUNT(user_groups.id) as Count 
		FROM user_groups 
		WHERE 1=1
		`

	// QueryParent is a query for getting parent rows
	// TODO remove JOIN users and delete user_group_lines when delete user
	QueryParent = `
		SELECT user_groups.id,
			   user_groups.code,
			   user_groups.description,
			   user_groups.status,
			   COALESCE(lines_t.Count, 0) AS Count
		FROM user_groups
				 LEFT JOIN (SELECT COUNT(user_group_lines.id) AS Count,
				                   user_group_lines.parent_id AS parent_id
							FROM user_group_lines
									 INNER JOIN users ON user_group_lines.user_id = users.id
							WHERE 2=2
							GROUP BY parent_id) lines_t
						   ON lines_t.parent_id = user_groups.id
		WHERE 1=1
		`

	// QueryParentJoins is a query for getting parent rows with joins
	QueryParentJoins = ``
)
