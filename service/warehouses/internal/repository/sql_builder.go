package repository

const (
	QueryCount = `SELECT COUNT(*) as Count FROM warehouses`

	QuerySelect = `SELECT * FROM warehouses`

	QueryJoin = ``

	// QuerySelectWithoutJoin use only in this module
	QuerySelectWithoutJoin = `SELECT * FROM warehouses`
)
