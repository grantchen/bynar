package repository

const (
	// QueryCount is the query to get the count of organizations
	QueryCount = `SELECT COUNT(*) as Count FROM organizations`

	// QuerySelect is the query to get the organizations
	QuerySelect = `SELECT *,vat_number AS vat_no FROM organizations`

	// QueryPermissionFormat is the extra permission query to get the organizations
	QueryPermissionFormat = ` AND user_group_int IN (SELECT parent_id FROM user_group_lines WHERE user_id = %d) `
)
