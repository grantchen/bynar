package repository

var (
	// UserGroupFieldNames is a map for user groups fields
	UserGroupFieldNames = map[string][]string{
		"id":          {"user_groups.id"},
		"code":        {"user_groups.code"},
		"description": {"user_groups.description"},
		"status":      {"user_groups.status"},
	}

	// UserGroupLineFieldNames is a map for user group lines fields
	UserGroupLineFieldNames = map[string][]string{
		"id-line":   {"user_group_lines.id"},
		"Parent":    {"user_group_lines.parent_id"},
		"user_id":   {"user_group_lines.user_id"},
		"email":     {"users.email"},
		"full_name": {"users.full_name"},
	}

	// UserGroupLineFieldUploadNames is a map for user group lines fields for upload
	UserGroupLineFieldUploadNames = map[string][]string{
		"id":      {"user_group_lines.id"},
		"Parent":  {"user_group_lines.parent_id"}, // for easy parse
		"user_id": {"user_group_lines.user_id"},
	}

	// UserUploadNames is a map for users fields for upload
	UserUploadNames = map[string][]string{
		"email":     {"users.email"},
		"full_name": {"users.full_name"},
	}
)
