package repository

const (
	// QueryChildCount is a query for getting count of child rows
	QueryChildCount = `
		SELECT COUNT(*) as Count
		FROM user_group_lines 
		INNER JOIN users ON user_group_lines.user_id = users.id
		`

	// QueryChild is a query for getting child rows
	QueryChild = `
	SELECT CONCAT (user_group_lines.id, '-line') as id, 
	user_group_lines.user_id,
	users.full_name,
	users.email 
	FROM user_group_lines 
		INNER JOIN users ON user_group_lines.user_id = users.id`

	// QueryChildJoins is a query for getting child rows with joins
	QueryChildJoins = `
	INNER JOIN users ON user_group_lines.user_id = users.id
	`

	// QueryChildSuggestion is a query for getting child rows suggestion
	QueryChildSuggestion = `
	SELECT users.id AS user_id,
	users.full_name,
	users.email 
	FROM users where concat(id,full_name,email) like ? AND id not in (
		SELECT user_id FROM user_group_lines where parent_id = ?
	)
	`
)
