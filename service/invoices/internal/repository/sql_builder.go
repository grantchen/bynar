package repository

const (
	QueryCount = `SELECT COUNT(*) as Count FROM organizations`

	QuerySelect = `SELECT *,vat_number AS vat_no FROM organizations`

	QueryJoin = ``

	QueryPermissionFormat = ` AND user_group_int IN (SELECT parent_id FROM user_group_lines WHERE user_id = %d) `
)
