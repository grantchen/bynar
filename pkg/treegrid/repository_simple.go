package treegrid

import (
	"database/sql"
	"fmt"
	"math"
	"strconv"

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
	MapSorted     map[string]bool
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

	return s.getPageData(tg, s.cfg.AdditionWhere)
}

func (s *simpleGridRepository) getPageData(tg *Treegrid, additionWhere string) ([]map[string]string, error) {
	pos, _ := tg.BodyParams.IntPos()
	query := BuildSimpleQuery(s.tableName, s.fieldMapping, s.cfg.QueryString)

	FilterWhere, FilterArgs := PrepQuerySimple(tg.FilterParams, s.fieldMapping)

	query = query + ParentDummyWhere + FilterWhere + " " + additionWhere + tg.OrderByChildQuery(s.cfg.MapSorted)
	query = AppendLimitToQuery(query, s.pageSize, pos)
	rows, err := s.db.Query(query, FilterArgs...)

	logger.Debug("query getPageData: ", query, "addition where: ", additionWhere)
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
		}

		tableData = append(tableData, entry)
	}
	return tableData, nil
}

func (s *simpleGridRepository) GetPageDataGroupBy(tg *Treegrid) ([]map[string]string, error) {
	level := tg.BodyParams.GetRowLevel()
	where := tg.BodyParams.GetRowWhere()

	// last level of group by, get data from table
	if level == len(tg.GroupCols) {
		where = where + s.cfg.AdditionWhere
		return s.getPageData(tg, where)
	}
	FilterWhere, FilterArgs := PrepQuerySimple(tg.FilterParams, s.fieldMapping)
	if level > 0 {
		FilterWhere = FilterWhere + tg.BodyParams.GetRowWhere() + s.cfg.AdditionWhere
	} else {
		FilterWhere = FilterWhere + s.cfg.AdditionWhere
	}

	groupWhere := ParentDummyWhere + FilterWhere
	query := BuildSimpleQueryGroupBy(s.tableName, s.fieldMapping, tg.GroupCols, groupWhere, level, s.cfg.QueryJoin)

	pos, _ := tg.BodyParams.IntPos()
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
	colData := NewColumn(tg.GroupCols[level], nil, s.fieldMapping)

	for rows.Next() {
		if err := rowVals.Parse(rows); err != nil {
			return tableData, fmt.Errorf("parse rows GetPageDataGroupBy: [%w]", err)
		}
		entry := rowVals.StringValues()
		entry["Def"] = "Group"
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
		if where != "" {
			where += " "
		}
		entry["Rows"] = strconv.Itoa(level+1) + where + "AND " + s.fieldMapping[tgCol][0] + " = '" + entry[tgCol] + "'"
		tableData = append(tableData, entry)
	}
	return tableData, nil
}

func createOrderMapping(fieldsMapping map[string][]string) map[string]bool {
	result := make(map[string]bool)
	for k, _ := range fieldsMapping {
		result[k] = true
	}
	return result
}

// GetPageCount implements SimpleGridRowRepository
func (s *simpleGridRepository) GetPageCount(tg *Treegrid) (int64, error) {
	var query string
	FilterWhere, FilterArgs := PrepQuerySimple(tg.FilterParams, s.fieldMapping)
	if !tg.WithGroupBy() {
		query = BuildSimpleQueryCount(s.tableName, s.fieldMapping, s.cfg.QueryCount)
		query = query + ParentDummyWhere + FilterWhere + s.cfg.AdditionWhere
	} else {
		where := ParentDummyWhere + FilterWhere + s.cfg.AdditionWhere
		query = BuildSimpleQueryGroupByCount(s.tableName, s.fieldMapping, tg.GroupCols, where)
	}

	logger.Debug("query GetPageCount: ", query)

	rows, err := s.db.Query(query, FilterArgs...)
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
		cfg:          &SimpleGridRepositoryCfg{MapSorted: createOrderMapping(fieldMapping)},
	}
}

func NewSimpleGridRowRepositoryWithCfg(db *sql.DB,
	tableName string,
	fieldMapping map[string][]string,
	maxPage int,
	cfg *SimpleGridRepositoryCfg) SimpleGridRowRepository {
	if cfg.MapSorted == nil {
		cfg.MapSorted = createOrderMapping(fieldMapping)
	}
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
	if len(args) == 0 {
		return fmt.Errorf("not field update detected")
	}
	args = append(args, gr.GetID())
	logrus.Info(query, "args", args)
	if _, err := tx.Exec(query, args...); err != nil {
		return fmt.Errorf("exec query: [%w], query: %s, args: %d", err, query, len(args))
	}

	return nil
}
