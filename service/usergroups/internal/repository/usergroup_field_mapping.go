package repository

var (
	UserGroupFieldNames = map[string][]string{
		"id":          {"user_groups.id"},
		"code":        {"user_groups.code"},
		"description": {"user_groups.description"},
		"status":      {"user_groups.status"},
	}

	UserGroupLineFieldNames = map[string][]string{
		"id":        {"user_group_lines.id"},
		"parent":    {"user_group_lines.parent"},
		"user_id":   {"user_group_lines.user_id"},
		"email":     {"users.email"},
		"full_name": {"users.full_name"},
	}
)
