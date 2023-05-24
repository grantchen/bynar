package repository

const (
	QueryChildCount = `
		SELECT COUNT(*) as Count
		FROM user_group_lines 
		INNER JOIN users ON user_group_lines.user_id = users.id
		`

	QueryChild = `
	SELECT CONCAT (user_group_lines.id, '-line') as id, 
	user_group_lines.user_id,
	users.full_name,
	users.email 
	FROM user_group_lines 
		INNER JOIN users ON user_group_lines.user_id = users.id`

	QueryChildJoins = `
	INNER JOIN users ON user_group_lines.user_id = users.id
	`
)
