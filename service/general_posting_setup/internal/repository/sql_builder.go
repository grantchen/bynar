package repository

const (
	QueryCount = `SELECT COUNT(*) as Count FROM general_posting_setup`

	QuerySelect = `SELECT * FROM general_posting_setup`

	QueryJoin = ``

	// QuerySelectWithoutJoin use only in this module
	QuerySelectWithoutJoin = `SELECT * FROM general_posting_setup`
)
