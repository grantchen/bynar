package repository

const (
	QueryCount = `SELECT COUNT(*) as Count FROM warehouses`

	QuerySelect = `SELECT * FROM warehouses`

	QueryJoin = ``

	// use only in this module
	QuerySelectWithoutJoin = `SELECT * FROM warehouses`
)
