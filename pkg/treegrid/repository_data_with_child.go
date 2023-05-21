package treegrid

import (
	"database/sql"
	"fmt"
	"log"
	"math"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/logger"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/utils"
)

type GridRowDataRepositoryWithChild interface {
	GetPageCount(tg *Treegrid) int64
	GetPageData(tg *Treegrid) ([]map[string]string, error)
}

type GridRowDataRepositoryWithChildCfg struct {
	MainCol                  string
	MapSorted                map[string]bool
	QueryParentJoins         string
	QueryParentCount         string
	QueryChildJoins          string
	QueryChildCount          string
	QueryChild               string
	QueryParent              string
	ChildJoinFieldWithParent string // example in user_group_lines: parent_id
	ParentIdField            string // example in user_groups: id
}

type gridRowDataRepositoryWithChild struct {
	db                 *sql.DB
	tableName          string
	lineTableName      string
	parentFieldMapping map[string][]string
	childFieldMapping  map[string][]string
	pageSize           int
	cfg                *GridRowDataRepositoryWithChildCfg
}

func NewGridRowDataRepositoryWithChild(
	db *sql.DB,
	tableName string,
	lineTableName string,
	parentFieldMapping map[string][]string,
	childFieldMapping map[string][]string,
	pageSize int,
	cfg *GridRowDataRepositoryWithChildCfg,
) GridRowDataRepositoryWithChild {
	return &gridRowDataRepositoryWithChild{
		db:                 db,
		tableName:          tableName,
		lineTableName:      lineTableName,
		parentFieldMapping: parentFieldMapping,
		childFieldMapping:  childFieldMapping,
		pageSize:           pageSize,
		cfg:                cfg,
	}
}

// GetPageCount implements GridRowDataRepositoryWithChild
func (g *gridRowDataRepositoryWithChild) GetPageCount(tg *Treegrid) int64 {
	var query string

	column := NewColumn(tg.GroupCols[0], g.childFieldMapping, g.parentFieldMapping)
	FilterWhere, FilterArgs := PrepQueryComplex(tg.FilterParams, g.parentFieldMapping, g.childFieldMapping)

	if column.IsItem {
		if FilterWhere["parent"] != "" {
			// FilterWhere["parent"] = " AND transfers_items.Parent IN (SELECT transfers.id from transfers " +
			FilterWhere["parent"] = fmt.Sprintf(" AND %s.%s IN (SELECT %s.%s FROM %s ",
				g.lineTableName,
				g.cfg.ChildJoinFieldWithParent,
				g.tableName,
				g.cfg.ParentIdField,
				g.tableName) +

				g.cfg.QueryParentJoins +
				DummyWhere +
				FilterWhere["parent"] + ") "
		}
		query = g.cfg.QueryChildCount + FilterWhere["child"] + FilterWhere["parent"]
		fmt.Printf("query count1: %s\n", query)
	} else {
		if FilterWhere["child"] != "" {
			// FilterWhere["child"] = " AND transfers.id IN (SELECT transfers_items.Parent from transfers_items " +
			FilterWhere["child"] = fmt.Sprintf(" AND %s.%s IN (SELECT %s.%s FROM %s ",
				g.tableName,
				g.cfg.ParentIdField,
				g.lineTableName,
				g.cfg.ChildJoinFieldWithParent,
				g.lineTableName) +

				g.cfg.QueryChildJoins +
				DummyWhere +
				FilterWhere["child"] + ") "

		}

		query = g.cfg.QueryParentCount + FilterWhere["child"] + FilterWhere["parent"]
		fmt.Printf("query count2: %s\n", query)
	}

	mergedArgs := utils.MergeMaps(FilterArgs["child"], FilterArgs["parent"])
	rows, err := g.db.Query(query, mergedArgs...)
	if err != nil {
		log.Println(err, "query", query, "colData", column)
	}

	return int64(math.Ceil(float64(utils.CheckCount(rows)) / float64(g.pageSize)))
	// return 0
}

func (g *gridRowDataRepositoryWithChild) toBooleanMapping(mapInput map[string][]string) map[string]bool {
	result := make(map[string]bool, 0)
	for _, v := range mapInput {
		result[v[0]] = true
	}
	return result
}

// GetPageData implements GridRowDataRepositoryWithChild
func (g *gridRowDataRepositoryWithChild) GetPageData(tg *Treegrid) ([]map[string]string, error) {

	PrepFilters(tg, g.parentFieldMapping, g.childFieldMapping)

	// items request
	if tg.BodyParams.GetItemsRequest() {
		logger.Debug("get items request")

		query := g.cfg.QueryChild +
			// " WHERE parent = " +
			fmt.Sprintf(" WHERE %s.%s = ", g.lineTableName, g.cfg.ChildJoinFieldWithParent) +
			tg.BodyParams.ID +
			tg.FilterWhere["child"] +
			tg.OrderByChildQuery(g.toBooleanMapping(g.childFieldMapping))
		pos, _ := tg.BodyParams.IntPos()
		query = AppendLimitToQuery(query, g.pageSize, pos)

		logger.Debug("query item request: ", query, "filter args: ", tg.FilterArgs["child"], " filter where: ", tg.FilterWhere["child"])

		return g.getJSON(query, tg.FilterArgs["child"], tg)
	}

	// GROUP BY
	if tg.WithGroupBy() {
		// logger.Debug("query with group by clause")

		// return t.handleGroupBy(tg)
	}

	logger.Debug("get transfers without grouping")

	query := g.cfg.QueryParent + tg.FilterWhere["parent"]
	if tg.FilterWhere["child"] != "" {
		// query += ` AND transfers.id IN ( SELECT Parent FROM transfers_items ` +
		query += fmt.Sprintf(" AND %s.%s IN ( SELECT %s FROM %s ", g.tableName, g.cfg.ParentIdField, g.cfg.ChildJoinFieldWithParent, g.lineTableName) +
			g.cfg.QueryChildJoins +
			tg.FilterWhere["child"] +
			`) `
	}

	query += tg.SortParams.OrderByQueryExludeChild2(g.toBooleanMapping(g.childFieldMapping), g.parentFieldMapping)

	pos, _ := tg.BodyParams.IntPos()
	query = AppendLimitToQuery(query, g.pageSize, pos)
	mergedArgs := utils.MergeMaps(tg.FilterArgs["parent"], tg.FilterArgs["child"])

	logger.Debug("query get page data", query, "args", mergedArgs)

	return g.getJSON(query, mergedArgs, tg)
}

func (g *gridRowDataRepositoryWithChild) getJSON(sqlString string, mergedArgs []interface{}, tg *Treegrid) ([]map[string]string, error) {
	stmt, err := g.db.Prepare(sqlString)
	if err != nil {
		return nil, fmt.Errorf("db prepare: [%w], sql string: [%s]", err, sqlString)
	}
	defer stmt.Close()

	rows, err := stmt.Query(mergedArgs...)
	if err != nil {
		return nil, fmt.Errorf("query: [%w], sql string: [%s]", err, sqlString)
	}
	defer rows.Close()

	rowVals, err := utils.NewRowVals(rows)
	if err != nil {
		return nil, fmt.Errorf("new row vals: [%w], row vals: [%v]", err, rowVals)
	}
	tableData := make([]map[string]string, 0, 100)

	for rows.Next() {
		if err := rowVals.Parse(rows); err != nil {
			return tableData, fmt.Errorf("parse rows: [%w]", err)
		}

		entry := rowVals.StringValues()
		if !tg.BodyParams.GetItemsRequest() {
			entry["Expanded"] = "0"
			entry["Count"] = "2"
			entry["has_child"] = "1"
		}

		tableData = append(tableData, entry)
	}

	return tableData, nil
}
