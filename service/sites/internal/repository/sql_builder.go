package repository

const (
	// QueryCount is the query to get the count of sites
	QueryCount = `SELECT COUNT(*) as Count FROM sites`

	// QuerySelect is the query to get all sites
	QuerySelect = `SELECT * FROM sites`

	// QueryPermissionFormat is the extra permission query to get the sites
	QueryPermissionFormat = ``
)
