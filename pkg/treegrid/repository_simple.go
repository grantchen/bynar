package treegrid

import (
	"database/sql"
	"fmt"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/logger"
)

// use for table with no child
type SimpleGridRowRepository interface {
	Add(tx *sql.Tx, gr GridRow) error
	Update(tx *sql.Tx, gr GridRow) error
	Delete(tx *sql.Tx, gr GridRow) error
}

type simpleGridRepository struct {
	db           *sql.DB
	tableName    string
	fieldMapping map[string][]string
}

func NewSimpleGridRowRepository(db *sql.DB, tableName string, fieldMapping map[string][]string) SimpleGridRowRepository {
	return &simpleGridRepository{
		db:           db,
		tableName:    tableName,
		fieldMapping: fieldMapping,
	}
}

// Delete implements SimpleGridRowRepository
func (s *simpleGridRepository) Delete(tx *sql.Tx, gr GridRow) error {
	logger.Debug("delete parent")

	query, args := gr.MakeDeleteQuery(s.tableName)
	args = append(args, gr.GetID())

	if _, err := tx.Exec(query, args...); err != nil {
		return fmt.Errorf("exec query: [%w], query: %s, args count: %d", err, query, len(args))
	}

	return nil
}

// Add implements SimpleGridRowRepository
func (s *simpleGridRepository) Add(tx *sql.Tx, gr GridRow) error {
	query, args := gr.MakeInsertQuery(s.tableName, s.fieldMapping)
	logger.Debug(query, "args", args)
	res, err := tx.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("exec query: [%w], query: %s", err, query)
	}
	newID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("last inserted id: [%w]", err)
	}

	// update id for row and child items
	gr["NewId"] = newID
	return nil
}

// Update implements SimpleGridRowRepository
func (s *simpleGridRepository) Update(tx *sql.Tx, gr GridRow) error {
	query, args := gr.MakeUpdateQuery(s.tableName, s.fieldMapping)
	if _, err := tx.Exec(query, args...); err != nil {
		return fmt.Errorf("exec query: [%w], query: %s, args count: %d", err, query, len(args))
	}

	return nil
}
