package treegrid

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/logger"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/utils"
)

// use for table with no child
type SimpleGridRowRepository interface {
	Add(tx *sql.Tx, gr GridRow) error
	Update(tx *sql.Tx, gr GridRow) error
	Delete(tx *sql.Tx, gr GridRow) error
	GetPageCount(tg *Treegrid) int64
	GetPageData(tg *Treegrid) ([]map[string]string, error)
	ValidateOnIntegrity(gr GridRow, validateFields []string) (bool, error)
}

type simpleGridRepository struct {
	db           *sql.DB
	tableName    string
	fieldMapping map[string][]string
	pageSize     int
}

// ValidateOnIntegrity implements SimpleGridRowRepository
func (s *simpleGridRepository) ValidateOnIntegrity(gr GridRow, validateFields []string) (bool, error) {
	query, args := gr.MakeValidateOnIntegrityQuery(s.tableName, s.fieldMapping, validateFields)
	fmt.Printf("ValidateOnIntegrity query: %s\n", query)

	var Count int
	err := s.db.QueryRow(query, args...).Scan(&Count)

	if err != nil {
		return false, fmt.Errorf("query row: [%w]", err)
	}
	return Count == 0, nil
}

// GetPageData implements SimpleGridRowRepository
func (s *simpleGridRepository) GetPageData(tg *Treegrid) ([]map[string]string, error) {
	b, _ := json.Marshal(tg)
	fmt.Printf("req data: %s\n", string(b))
	if tg.WithGroupBy() {
		return s.GetPageDataGroupBy(tg)
	}
	pos, _ := tg.BodyParams.IntPos()
	query := BuildSimpleQuery(s.tableName, s.fieldMapping)

	FilterWhere, FilterArgs := PrepQuerySimple(tg.FilterParams, s.fieldMapping)

	query = query + DummyWhere + FilterWhere
	query = AppendLimitToQuery(query, s.pageSize, pos)
	rows, err := s.db.Query(query, FilterArgs...)
	if err != nil {
		return nil, fmt.Errorf("do query: '%s': [%w]", query, err)
	}
	defer rows.Close()

	rowVals, err := utils.NewRowVals(rows)
	if err != nil {
		return nil, fmt.Errorf("new row vals: '%v': [%w]", rowVals, err)
	}

	tableData := make([]map[string]string, 0)
	for rows.Next() {
		if err := rowVals.Parse(rows); err != nil {
			return tableData, fmt.Errorf("parse rows: [%w]", err)
		}

		entry := rowVals.StringValues()
		if !tg.BodyParams.GetItemsRequest() {
			entry["Expanded"] = "0"
			entry["Count"] = "1"
		}

		tableData = append(tableData, entry)
	}
	// tableData := make([]map[string]string, 0, 0)

	return tableData, nil
}

func (s *simpleGridRepository) GetPageDataGroupBy(tg *Treegrid) ([]map[string]string, error) {
	FilterWhere, FilterArgs := PrepQuerySimple(tg.FilterParams, s.fieldMapping)
	query := BuildSimpleQueryGroupBy(s.tableName, s.fieldMapping, tg.GroupCols, FilterWhere)

	pos, _ := tg.BodyParams.IntPos()
	query = AppendLimitToQuery(query, s.pageSize, pos)
	fmt.Printf("query: %s\n", query)
	rows, err := s.db.Query(query, FilterArgs...)

	if err != nil {
		return nil, fmt.Errorf("do query: '%s': [%w]", query, err)
	}
	defer rows.Close()

	rowVals, err := utils.NewRowVals(rows)
	if err != nil {
		return nil, fmt.Errorf("new row vals: '%v': [%w]", rowVals, err)
	}

	tableData := make([]map[string]string, 0)
	for rows.Next() {
		if err := rowVals.Parse(rows); err != nil {
			return tableData, fmt.Errorf("parse rows: [%w]", err)
		}
		entry := rowVals.StringValues()
		entry["Def"] = "Group"
		entry["Expanded"] = "0"
		tableData = append(tableData, entry)
	}
	return tableData, nil
}

// GetPageCount implements SimpleGridRowRepository
func (s *simpleGridRepository) GetPageCount(tg *Treegrid) int64 {
	b, _ := json.Marshal(tg)
	fmt.Printf("req data: %s\n", string(b))
	var query string
	if !tg.WithGroupBy() {
		query = BuildSimpleQueryCount(s.tableName, s.fieldMapping)

	} else {
		query = BuildSimpleQueryGroupByCount(s.tableName, s.fieldMapping, tg.GroupCols)
	}

	rows, err := s.db.Query(query)
	if err != nil {
		fmt.Printf("parse rows: [%v]", err)
		return 0
	}

	return int64(math.Ceil(float64(utils.CheckCount(rows)) / float64(s.pageSize)))
}

func NewSimpleGridRowRepository(db *sql.DB, tableName string, fieldMapping map[string][]string, maxPage int) SimpleGridRowRepository {
	return &simpleGridRepository{
		db:           db,
		tableName:    tableName,
		fieldMapping: fieldMapping,
		pageSize:     maxPage,
	}
}

// Delete implements SimpleGridRowRepository
func (s *simpleGridRepository) Delete(tx *sql.Tx, gr GridRow) error {
	query, args := gr.MakeDeleteQuery(s.tableName)
	fmt.Printf("query: %s, %s\n", query, gr.GetID())
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
	args = append(args, gr.GetID())
	if _, err := tx.Exec(query, args...); err != nil {
		return fmt.Errorf("exec query: [%w], query: %s, args count: %d", err, query, len(args))
	}

	return nil
}
