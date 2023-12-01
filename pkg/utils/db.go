package utils

import (
	"database/sql"
	stderr "errors"
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

// CheckCount checks count value from sql.Rows
func CheckCount(rows *sql.Rows) (rowCount int) {
	if !rows.Next() {
		return 0
	}

	err := rows.Scan(&rowCount)
	if err != nil {
		logger.Debug(err)
	}

	return rowCount
}

// CheckCountWithError checks count value from sql.Rows
func CheckCountWithError(rows *sql.Rows) (rowCount int, err error) {
	if !rows.Next() {
		return 0, err
	}

	err = rows.Scan(&rowCount)
	if err != nil {
		return 0, err
	}

	return rowCount, nil
}

// CheckExist checks if row exists
func CheckExist(db *sql.DB, query string, args ...any) (exists bool, err error) {
	stmt, err := db.Prepare(query)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	var existFlag int
	if err = stmt.QueryRow(args...).Scan(&existFlag); err != nil {
		if stderr.Is(err, sql.ErrNoRows) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// CheckExistInTx checks if row exists in transaction
func CheckExistInTx(tx *sql.Tx, query string, args ...any) (exists bool, err error) {
	stmt, err := tx.Prepare(query)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	var existFlag int
	if err = stmt.QueryRow(args...).Scan(&existFlag); err != nil {
		if stderr.Is(err, sql.ErrNoRows) {
			return false, nil
		}

		return false, err
	}

	return true, nil
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
