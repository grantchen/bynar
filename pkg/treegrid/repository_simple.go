package treegrid

import (
	"database/sql"
	"fmt"
	"math"

	"github.com/sirupsen/logrus"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/logger"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/utils"
)

// use for table with no child
type SimpleGridRowRepository interface {
	Add(tx *sql.Tx, gr GridRow) error
	Update(tx *sql.Tx, gr GridRow) error
	Delete(tx *sql.Tx, gr GridRow) error
	GetPageCount(tg *Treegrid) (int64, error)
	GetPageData(tg *Treegrid) ([]map[string]string, error)
	ValidateOnIntegrity(tx *sql.Tx, gr GridRow, validateFields []string) (bool, error)
}

type SimpleGridRepositoryCfg struct {
	MainCol       string
	QueryString   string
	QueryJoin     string
	QueryCount    string
	AdditionWhere string // using addtion, example for permission
}

type simpleGridRepository struct {
	db           *sql.DB
	tableName    string
	fieldMapping map[string][]string
	pageSize     int
	cfg          *SimpleGridRepositoryCfg
}

// ValidateOnIntegrity implements SimpleGridRowRepository
func (s *simpleGridRepository) ValidateOnIntegrity(tx *sql.Tx, gr GridRow, validateFields []string) (bool, error) {
	query, args := gr.MakeValidateOnIntegrityQuery(s.tableName, s.fieldMapping, validateFields)
	fmt.Printf("ValidateOnIntegrity query: %s, %s\n", query, args)

	var Count int
	err := tx.QueryRow(query, args...).Scan(&Count)

	if err != nil {
		return false, fmt.Errorf("query row: [%w]", err)
	}
	return Count == 0, nil
}

// GetPageData implements SimpleGridRowRepository
func (s *simpleGridRepository) GetPageData(tg *Treegrid) ([]map[string]string, error) {
	if tg.WithGroupBy() {
		return s.GetPageDataGroupBy(tg)
	}

	return s.getPageData(tg, s.cfg.AdditionWhere, nil)
}

func (s *simpleGridRepository) getPageData(tg *Treegrid, additionWhere string, additionWhereArgs []interface{}) ([]map[string]string, error) {
	pos, _ := tg.BodyParams.IntPos()
	query := BuildSimpleQuery(s.tableName, s.fieldMapping, s.cfg.QueryString)

	FilterWhere, FilterArgs := PrepQuerySimple(tg.FilterParams, s.fieldMapping)

	query = JoinSQLs(query, ParentDummyWhere, FilterWhere, additionWhere,
		tg.OrderByChildQuery(s.fieldMapping, fmt.Sprintf("%s.id ASC", s.tableName)))
	query = AppendLimitToQuery(query, s.pageSize, pos)
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("prepare query: '%s': [%w]", query, err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(append(FilterArgs, additionWhereArgs...)...)
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
			return tableData, fmt.Errorf("parse rows getPageData: [%w]", err)
		}

		entry := rowVals.StringValues()
		if !tg.BodyParams.GetItemsRequest() {
			entry["Expanded"] = "0"
			entry["Count"] = "0"
			entry["Def"] = "Node"
			entry["CDef"] = "Data"
		}

		tableData = append(tableData, entry)
	}
	return tableData, nil
}

func (s *simpleGridRepository) GetPageDataGroupBy(tg *Treegrid) ([]map[string]string, error) {
	level := tg.BodyParams.RowsLevel
	rowsWhere, rowsArgs := PrepRowsSimple(tg.BodyParams, s.fieldMapping)
	// last level of group by, get data from table
	if level == len(tg.GroupCols) {
		return s.getPageData(tg, JoinSQLs(rowsWhere, s.cfg.AdditionWhere), rowsArgs)
	}

	var queryArgs []interface{}
	FilterWhere, FilterArgs := PrepQuerySimple(tg.FilterParams, s.fieldMapping)
	if level > 0 {
		FilterWhere = JoinSQLs(FilterWhere, rowsWhere, s.cfg.AdditionWhere)
		queryArgs = append(FilterArgs, rowsArgs...)
	} else {
		FilterWhere = JoinSQLs(FilterWhere, s.cfg.AdditionWhere)
		queryArgs = FilterArgs
	}

	groupWhere := JoinSQLs(ParentDummyWhere, FilterWhere)
	query := BuildSimpleQueryGroupBy(s.tableName, s.fieldMapping, tg.GroupCols, groupWhere, level, s.cfg.QueryJoin)

	pos, _ := tg.BodyParams.IntPos()
	query = AppendLimitToQuery(query, s.pageSize, pos)
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("prepare query: '%s': [%w]", query, err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(queryArgs...)
	if err != nil {
		return nil, fmt.Errorf("query rows: [%w]", err)
	}
	defer rows.Close()

	rowVals, err := utils.NewRowVals(rows)
	if err != nil {
		return nil, fmt.Errorf("new row vals: '%v': [%w]", rowVals, err)
	}

	tableData := make([]map[string]string, 0)
	colData := NewColumn(tg.GroupCols[level], nil, s.fieldMapping)

	for rows.Next() {
		if err := rowVals.Parse(rows); err != nil {
			return tableData, fmt.Errorf("parse rows GetPageDataGroupBy: [%w]", err)
		}
		entry := rowVals.StringValues()
		entry["Def"] = "Group"
		entry["CDef"] = "Node"
		entry["Expanded"] = "0"
		tgCol := tg.GroupCols[level]
		if s.cfg != nil && s.cfg.MainCol != "" {
			entry[s.cfg.MainCol] = entry[tgCol]
			entry[s.cfg.MainCol+"Type"] = colData.Type()
			// TODO get from treegrid config
			if colData.IsDate {
				entry[s.cfg.MainCol+"Format"] = "yyyy-MM-dd"
			}
		}

		rowsConds, err := tg.BodyParams.GetNewRows(level+1, RowsFieldCond{GridName: tgCol, Value: entry[tgCol]})
		if err != nil {
			return nil, err
		}
		entry["Rows"] = rowsConds
		tableData = append(tableData, entry)
	}
	return tableData, nil
}

// GetPageCount implements SimpleGridRowRepository
func (s *simpleGridRepository) GetPageCount(tg *Treegrid) (int64, error) {
	var query string
	FilterWhere, FilterArgs := PrepQuerySimple(tg.FilterParams, s.fieldMapping)
	if !tg.WithGroupBy() {
		query = BuildSimpleQueryCount(s.tableName, s.fieldMapping, s.cfg.QueryCount)
		query = JoinSQLs(query, ParentDummyWhere, FilterWhere, s.cfg.AdditionWhere)
	} else {
		where := JoinSQLs(ParentDummyWhere, FilterWhere, s.cfg.AdditionWhere)
		query = BuildSimpleQueryGroupByCount(s.tableName, s.fieldMapping, tg.GroupCols, where)
	}

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return 0, fmt.Errorf("prepare query: '%s': [%w]", query, err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(FilterArgs...)
	if err != nil {
		fmt.Printf("parse rows: [%v]", err)
		return 0, err
	}
	defer rows.Close()

	return int64(math.Ceil(float64(utils.CheckCount(rows)) / float64(s.pageSize))), nil
}

func NewSimpleGridRowRepository(db *sql.DB, tableName string, fieldMapping map[string][]string, maxPage int) SimpleGridRowRepository {
	return &simpleGridRepository{
		db:           db,
		tableName:    tableName,
		fieldMapping: fieldMapping,
		pageSize:     maxPage,
		cfg:          &SimpleGridRepositoryCfg{},
	}
}

func NewSimpleGridRowRepositoryWithCfg(db *sql.DB,
	tableName string,
	fieldMapping map[string][]string,
	maxPage int,
	cfg *SimpleGridRepositoryCfg) SimpleGridRowRepository {
	return &simpleGridRepository{
		db:           db,
		tableName:    tableName,
		fieldMapping: fieldMapping,
		pageSize:     maxPage,
		cfg:          cfg,
	}
}

// Delete implements SimpleGridRowRepository
func (s *simpleGridRepository) Delete(tx *sql.Tx, gr GridRow) error {
	changedRow := GenGridRowChangeError(gr)

	query, args := gr.MakeDeleteQuery(s.tableName)
	fmt.Printf("query: %s, %s\n", query, gr.GetID())
	args = append(args, gr.GetID())

	if _, err := tx.Exec(query, args...); err != nil {
		return fmt.Errorf("exec query: [%w], query: %s, args count: %d", err, query, len(args))
	}

	changedRow.Color = ChangedSuccessColor
	changedRow.Deleted = 1
	SetGridRowChangedResult(gr, changedRow)

	return nil
}

// Add implements SimpleGridRowRepository
func (s *simpleGridRepository) Add(tx *sql.Tx, gr GridRow) error {
	changedRow := GenGridRowChangeError(gr)

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
	changedRow.Color = ChangedSuccessColor
	changedRow.Added = 1
	changedRow.NewId = fmt.Sprintf("%v$%d", changedRow.Parent, newID) // full id
	SetGridRowChangedResult(gr, changedRow)

	return nil
}

// Update implements SimpleGridRowRepository
func (s *simpleGridRepository) Update(tx *sql.Tx, gr GridRow) error {
	changedRow := GenGridRowChangeError(gr)

	query, args := gr.MakeUpdateQuery(s.tableName, s.fieldMapping)
	if len(args) == 0 {
		return fmt.Errorf("not field update detected")
	}
	args = append(args, gr.GetID())
	logrus.Info(query, "args", args)
	if _, err := tx.Exec(query, args...); err != nil {
		return fmt.Errorf("exec query: [%w], query: %s, args: %d", err, query, len(args))
	}

	changedRow.Color = ChangedSuccessColor
	changedRow.Changed = 1
	SetGridRowChangedResult(gr, changedRow)

	return nil
}
