package utils

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/logger"
)

type RowValues struct {
	values         map[string]interface{}
	columnPointers []interface{}
	columnNames    []string
}

func NewRowVals(rows *sql.Rows) (*RowValues, error) {
	columnNames, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("get columns: [%w]", err)
	}

	r := &RowValues{
		values:         make(map[string]interface{}),
		columnPointers: make([]interface{}, len(columnNames)),
		columnNames:    columnNames,
	}

	columns := make([]interface{}, len(columnNames))
	for i := range columns {
		r.columnPointers[i] = &columns[i]
	}

	return r, nil
}

func (r *RowValues) Values() map[string]interface{} {
	return r.values
}

func (r *RowValues) StringValues() map[string]string {
	res := make(map[string]string, len(r.columnNames))
	for _, v := range r.columnNames {
		res[v] = r.GetValue(v)
	}

	return res
}

func (r *RowValues) Parse(rows *sql.Rows) error {
	r.values = make(map[string]interface{})

	if err := rows.Scan(r.columnPointers...); err != nil {
		return fmt.Errorf("rows scan: [%w]", err)
	}

	for i, colName := range r.columnNames {
		val := r.columnPointers[i].(*interface{})
		r.values[colName] = *val
	}

	return nil
}

// func(r *RowValues)
func (r *RowValues) GetValue(columnName string) string {
	// return fmt.Sprintf("%v", r.values[columnName])
	if val, ok := r.values[columnName].([]byte); ok {
		return string(val)
	}

	val := fmt.Sprintf("%v", r.values[columnName])
	if val == "<nil>" {
		return ""
	}

	return val
}

func CheckCount(rows *sql.Rows) (rowCount int) {
	for rows.Next() {
		err := rows.Scan(&rowCount)
		if err != nil {
			logger.Debug(err)
		}
	}

	return rowCount
}

func CheckCoutWithError(rows *sql.Rows) (rowCount int, err error) {
	for rows.Next() {
		err := rows.Scan(&rowCount)
		if err != nil {
			return 0, err
		}
	}

	return rowCount, nil
}

func ExtractWhereClause(query string, params []interface{}) string {
	paramValueCount := strings.Count(query, "?")
	for count := 0; count < paramValueCount; count++ {
		query = strings.Replace(query, "?", fmt.Sprintf("'%s'", params[count]), 1)
	}

	pos := strings.Index(query, "WHERE")
	if pos == -1 {
		return "WHERE 1=1"
	}

	return query[pos:]
}

// Returns the replaced column in the query.
func ReplaceColumnValueInQuery(orignal_query string, col_name string, col_value string) string {
	updated_query := ""
	if strings.Index(orignal_query, " AND "+col_name+"=''") > 0 {
		// If the column value was empty
		updated_query = strings.Replace(orignal_query, " AND "+col_name+"=''", " AND "+col_name+"='"+col_value+"'", 1)
	} else {
		// If the column value was not empty
		sub_string := GetStringBetween(orignal_query, " AND "+col_name+"='", "'")
		updated_query = strings.Replace(orignal_query, sub_string, col_value, 1)
	}
	return updated_query
}

// TODO
// Generate organization connection
func GenerateOrganizationConnection(tenantUuid, organizationUuid string) string {
	envs := strings.Split(os.Getenv(tenantUuid), "/")
	connStr := envs[0] + "/" + organizationUuid
	if len(envs) > 1 {
		connStr += envs[1]
	}

	return connStr
}

// A TxFn is a function that will be called with an initialized `Transaction` object
// that can be used for executing statements and queries against a database.
type TxFn func(tx *sql.Tx) error

// WithTransaction creates a new transaction and handles rollback/commit based on the
// error object returned by the `TxFn`
// nolint: gocritic // no need to lint
func WithTransaction(db *sql.DB, fn TxFn) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return
	}

	defer func() {
		if p := recover(); p != nil {
			// a panic occurred, rollback and re-panic
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			// something went wrong, rollback
			_ = tx.Rollback()
		} else {
			// all good, commit
			err = tx.Commit()
		}
	}()

	err = fn(tx)
	return err
}
