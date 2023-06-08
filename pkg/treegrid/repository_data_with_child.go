package treegrid

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/logger"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/utils"
)

type GridRowDataRepositoryWithChild interface {
	GetPageCount(tg *Treegrid) (int64, error)
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
func (g *gridRowDataRepositoryWithChild) GetPageCount(tg *Treegrid) (int64, error) {
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
		logger.Debug("query count1:", query)
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
		logger.Debug("query count2: ", query)
	}

	mergedArgs := utils.MergeMaps(FilterArgs["child"], FilterArgs["parent"])
	rows, err := g.db.Query(query, mergedArgs...)
	if err != nil {
		log.Println(err, "query", query, "colData", column)
		return 0, err
	}

	return int64(math.Ceil(float64(utils.CheckCount(rows)) / float64(g.pageSize))), nil
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
		logger.Debug("get items request, body id: ", tg.BodyParams.ID)

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
		logger.Debug("query with group by clause")

		return g.handleGroupBy(tg)
	}

	logger.Debug("get without grouping")

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
			// entry["MinLevels"] = "2"
			entry["has_child"] = "1"
		}

		tableData = append(tableData, entry)
	}

	return tableData, nil
}

func (g *gridRowDataRepositoryWithChild) handleGroupBy(tg *Treegrid) ([]map[string]string, error) {
	if tg.BodyParams.Rows != "" {
		logger.Debug("with rows:", tg.BodyParams.Rows)

		return g.getGroupData(tg.BodyParams.GetRowWhere(), tg)
	}

	// parentBuild := QueryParent
	parentBuild := g.cfg.QueryParent

	if tg.FilterWhere["parent"] != "" {
		logger.Debug("filters table", tg.FilterWhere["parent"])

		parentBuild += " " + tg.FilterWhere["parent"]
	}

	if tg.FilterWhere["child"] != "" {
		logger.Debug("filter items", tg.FilterWhere["child"])

		// parentBuild = parentBuild + " AND id IN (" +
		// 	"SELECT Parent FROM transfers_items " +
		// 	"WHERE 1=1 " + tg.FilterWhere["child"] +
		// 	tg.OrderByChildQuery(model.TransferItemsFields) + ") "

		parentBuild = parentBuild +
			fmt.Sprintf("AND %s.%s IN (SELECT %s FROM %s", g.tableName, g.cfg.ParentIdField, g.cfg.ChildJoinFieldWithParent, g.lineTableName) +
			g.cfg.QueryChildJoins +
			DummyWhere +
			tg.FilterWhere["child"] +
			tg.OrderByChildQuery(g.toBooleanMapping(g.childFieldMapping)) + ") "

		logger.Debug("parentBuild ", parentBuild)
	}

	mergedArgs := utils.MergeMaps(tg.FilterArgs["parent"], tg.FilterArgs["child"])

	logger.Debug("query for extraction", parentBuild)

	extWhereClause := utils.ExtractWhereClause(parentBuild, mergedArgs)

	logger.Debug("extracted where clause", extWhereClause)

	return g.getGroupData(extWhereClause, tg)
}

func (g *gridRowDataRepositoryWithChild) getGroupData(where string, tg *Treegrid) ([]map[string]string, error) {
	logger.Debug("get group data")

	query, colData := g.prepareNameCountQuery(where, tg)

	pos, _ := tg.BodyParams.IntPos()
	query = AppendLimitToQuery(query, g.pageSize, pos)

	logger.Debug("column", colData, "prepared query: \n", query)

	if colData.IsItem {
		return g.getChildData(tg.BodyParams.GetRowLevel(), tg.GroupCols, where, query, colData)
	}

	return g.getParentData(tg.BodyParams.GetRowLevel(), tg.GroupCols, where, query, colData)
}

func (g *gridRowDataRepositoryWithChild) generateNameCountQuery(where string, column Column, tg *Treegrid) string {
	// query := `
	// SELECT %s, COUNT(*) Count
	// FROM transfers_items
	// INNER JOIN items ON transfers_items.item_uuid = items.id
	// INNER JOIN units ON transfers_items.item_unit_uuid = units.id
	// INNER JOIN item_types ON items.type_uuid = item_types.id
	// INNER JOIN transfers ON transfers_items.Parent = transfers.id
	// INNER JOIN documents ON transfers.document_type_uuid = documents.id
	// INNER JOIN stores ON transfers.store_origin_uuid = stores.id
	// INNER JOIN stores ss ON transfers.store_destination_uuid = ss.id
	// INNER JOIN warehouses wh_origin ON transfers.warehouse_origin_uuid = wh_origin.id
	// INNER JOIN warehouses wh_destination ON transfers.warehouse_destination_uuid = wh_destination.id
	// INNER JOIN responsibility_center ON transfers.responsibility_center_uuid = responsibility_center.id
	// %s
	// GROUP BY %s
	// `
	var query string = g.cfg.QueryChildCount

	// join with parent table
	// if strings.Index(where, g.tableName) > 0 {
	query = query + fmt.Sprintf(`
		INNER JOIN %s ON %s.%s = %s.%s
		`, g.tableName, g.tableName, g.cfg.ParentIdField, g.lineTableName, g.cfg.ChildJoinFieldWithParent)

	query = strings.Replace(query, "SELECT ", `SELECT %s, `, 1) +
		`%s %s GROUP BY %s`
	// }

	// fix here
	additionChildCondition := utils.ExtractWhereClause("WHERE "+tg.FilterWhere["child"], tg.FilterArgs["child"])
	additionChildCondition = additionChildCondition[len("WHERE"):]
	logger.Debug("additionChildCondition: ", additionChildCondition, "filter where: ", tg.FilterWhere["child"], "filter args", tg.FilterArgs["child"])
	return fmt.Sprintf(query, column.DBName, where, additionChildCondition, column.DBName)
}

func (g *gridRowDataRepositoryWithChild) prepareNameCountQuery(where string, tg *Treegrid) (query string, column Column) {
	// If both the level are equal then return the row

	level := tg.BodyParams.GetRowLevel()
	logger.Debug("Level", level, "len(groupCols)", len(tg.GroupCols))

	if level == len(tg.GroupCols) {
		logger.Debug("getting last level data")

		query = g.cfg.QueryParent
		query = strings.Replace(query, "WHERE 1=1", "", 1)
		query += where

		return
	}

	column = NewColumn(tg.GroupCols[level], g.childFieldMapping, g.parentFieldMapping)

	logger.Debug("getting data grouping by", tg.GroupCols[level])

	// multiple groupBy clauses
	if len(tg.GroupCols)-level > 1 {
		secColumn := NewColumn(tg.GroupCols[level+1], g.childFieldMapping, g.parentFieldMapping)

		switch {
		case !column.IsItem && !secColumn.IsItem:
			return g.getCascadingGroupByParentParent(column, secColumn, where), column
		case column.IsItem && !secColumn.IsItem, !column.IsItem && secColumn.IsItem:
			return g.getCascadingGroupByParentChild(column, secColumn, where), column
		}
	}

	if column.IsItem {
		logger.Debug("grouping by item")

		if val, ok := g.childFieldMapping[column.GridName]; ok {
			column.DBName = val[0]
		}

		query = g.generateNameCountQuery(where, column, tg)
		logger.Debug("query check: ", query, "where:", where)
		return query, column
	}

	query = "SELECT " + column.DBName + ", COUNT(*) Count FROM " + g.tableName + " " + g.cfg.QueryParentJoins + where + " GROUP BY " + column.DBName

	return
}

func (g *gridRowDataRepositoryWithChild) getParentData(level int, group_cols []string, where string, query string, colData Column) ([]map[string]string, error) {
	logger.Debug("get parent data")

	rows, err := g.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("do query: '%s': [%w]", query, err)
	}
	defer rows.Close()

	row, err := utils.NewRowVals(rows)
	if err != nil {
		return nil, fmt.Errorf("new row vals: '%v': [%w]", row, err)
	}

	tableData := make([]map[string]string, 0, 100)

	// Loop through the rows and form the row content.
	for rows.Next() {
		tempObj := make(map[string]string)
		tempObj["Def"] = "Group"
		tempObj["Expanded"] = "0"

		err := row.Parse(rows)
		if err != nil {
			return tableData, fmt.Errorf("parse rows: [%w]", err)
		}
		tempObj["Count"] = row.GetValue("Count")

		if level == len(group_cols) {
			for k := range row.StringValues() {
				tempObj[k] = row.StringValues()[k]
			}

			tempObj["Count"] = "1"

			tableData = append(tableData, tempObj)
			continue
		}

		docType := row.GetValue(colData.DBNameShort)
		if docType == "" {
			docType = row.GetValue(colData.GridName)
		}

		tempObj[g.cfg.MainCol] = docType

		// if grouping by document_type => DB, OPT, DT...
		// name = document_type
		// val = DB
		val := docType
		where2 := " AND " + colData.DBName + "='" + val + "'"

		// Builds new attribute Rows for identification
		tempObj["Rows"] = strconv.Itoa(level+1) + where + where2

		// Query to get aggregated sum for "warehouse_destination_uuid" data for each group.
		// calcQuery := "SELECT COALESCE(SUM(warehouseman_destination_approve),'') as value_sum, " +
		// 	"COALESCE(MIN(document_date), '') AS min, " +
		// 	"COALESCE(MAX(document_date), '') AS max " +
		// 	// "FROM (SELECT warehouseman_destination_approve, document_date FROM transfers " + where + where2 + ") AS temp"
		// 	"FROM (SELECT warehouseman_destination_approve, document_date FROM transfers " + sqlbuilder.QueryParentJoins + where + where2 + ") AS temp"

		// calcRows, err := t.db.Query(calcQuery)
		// if err != nil {
		// 	return tableData, fmt.Errorf("do query: [%s]: [%w]", calcQuery, err)
		// }

		// var min, max, sum string

		// for calcRows.Next() {
		// 	err = calcRows.Scan(&min, &max, &sum)
		// 	if err != nil {
		// 		return tableData, fmt.Errorf("rows scan: [%w]", err)
		// 	}
		// }

		// document_date := ""
		// if min != "" && max != "" {
		// 	document_date = min + "~" + max
		// }

		// tempObj["document_date"] = document_date
		// tempObj["warehouse_destination_uuid"] = sum

		tableData = append(tableData, tempObj)

	}

	return tableData, nil
}

func (g *gridRowDataRepositoryWithChild) getChildData(level int, groupCols []string, where string, query string, columnData Column) ([]map[string]string, error) {
	logger.Debug("get child data")

	rows, err := g.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("do query: [%w], query: [%s]", err, query)
	}
	defer rows.Close()

	row, err := utils.NewRowVals(rows)
	if err != nil {
		return nil, fmt.Errorf("new row values: [%w], values: [%v]", err, row)
	}

	tableData := make([]map[string]string, 0, 100)
	for rows.Next() {
		tempObj := make(map[string]string)
		tempObj["Def"] = "Group"
		tempObj["Expanded"] = "0"

		err := row.Parse(rows)
		if err != nil {
			return tableData, fmt.Errorf("parse row: [%w]", err)
		}

		tempObj["Count"] = row.GetValue("Count")
		tempObj[g.cfg.MainCol] = row.GetValue(columnData.DBNameShort)
		if level == len(groupCols) {
			tableData = append(tableData, row.StringValues())
			continue
		}

		// var value_sums, where2 string // sql.nullString
		var where2 string // sql.nullString

		// If grouped by child column
		if strings.Index(where, " AND "+columnData.DBName+"='") > 0 {
			// If the column is already added in the WHERE clause
			// Just update the query to adjust column value
			where2 = utils.ReplaceColumnValueInQuery(where, columnData.DBName, row.GetValue(columnData.DBNameShort))
		} else {
			where2 = g.addChildCondition(where, columnData.DBName, row.GetValue(columnData.DBNameShort))
		}
		tempObj["Rows"] = strconv.Itoa(level+1) + where2

		// Query to get aggregated sum for "item_quantity" data for each transfer_items group.
		// 		calcQuery := `
		// SELECT COALESCE(sum(item_quantity),'') as value_sums
		// FROM (
		// SELECT item_quantity
		// FROM transfers_items ` + sqlbuilder.QueryChildJoins + `
		// WHERE Parent in (
		// 	SELECT transfers.id FROM transfers ` + sqlbuilder.QueryParentJoins + where + `) AND ` + columnData.DBName + "='" + row.GetValue(columnData.DBNameShort) + "') AS temp;"
		// 		calcRows, err := t.db.Query(calcQuery)
		// 		if err != nil {
		// 			return tableData, fmt.Errorf("do query: '%s': [%w]", calcQuery, err)
		// 		}

		// 		for calcRows.Next() {
		// 			err = calcRows.Scan(&value_sums)
		// 			if err != nil {
		// 				return tableData, fmt.Errorf("rows scan: [%w]", err)
		// 			}
		// 		}
		// 		// valueSums := calculations_rs->Get("value_sums")
		// 		tempObj["item_quantity"] = value_sums

		tableData = append(tableData, tempObj)
	}

	return tableData, nil
}

func (g *gridRowDataRepositoryWithChild) addChildCondition(orignal_query string, col_name string, col_value string) string {
	tmp := fmt.Sprintf("%s IN ( SELECT %s FROM %s", g.cfg.ParentIdField, g.cfg.ChildJoinFieldWithParent, g.lineTableName)
	if strings.Index(orignal_query, tmp) > 0 {
		// If already filtered on child column(s)
		orignal_query = strings.Replace(orignal_query, ")", " AND "+col_name+"='"+col_value+"' )", 1)
	} else {
		// orignal_query += ` AND transfers.id IN ( SELECT Parent FROM transfers_items
		// INNER JOIN items ON transfers_items.item_uuid = items.id
		// INNER JOIN units ON transfers_items.item_unit_uuid = units.id
		// INNER JOIN item_types ON items.type_uuid = item_types.id
		// WHERE ` + col_name + "='" + col_value + "' ) "

		orignal_query += fmt.Sprintf(`AND %s.%s IN ( SELECT %s FROM %s `, g.tableName, g.cfg.ParentIdField, g.cfg.ChildJoinFieldWithParent, g.lineTableName) +
			g.cfg.QueryChildJoins +
			" WHERE " + col_name + "='" + col_value + "' ) "
	}

	return orignal_query
}

func (g *gridRowDataRepositoryWithChild) getCascadingGroupByParentParent(firstCol, secondCol Column, where string) string {
	query := `
	SELECT %s, COUNT(*) Count
	FROM (
		SELECT %s, %s 
		FROM %s
		%s
		%s
		GROUP BY %s, %s
	) t
	GROUP BY %s
	`

	return fmt.Sprintf(query, firstCol.DBNameShort, firstCol.DBName, secondCol.DBName, g.tableName, g.cfg.QueryParentJoins, where, firstCol.DBName, secondCol.DBName, firstCol.DBNameShort)
}

// when grouping by two parent (tansfer) columns
func (g *gridRowDataRepositoryWithChild) getCascadingGroupByParentChild(firstCol, secondCol Column, where string) string {
	// query := `
	// SELECT %s, COUNT(*) Count FROM (
	// 	SELECT
	// 		%s, %s, COUNT(*) Count
	// 	FROM transfers_items
	// 		INNER JOIN items ON transfers_items.item_uuid = items.id
	// 		INNER JOIN units ON transfers_items.item_unit_uuid = units.id
	// 		INNER JOIN item_types ON items.type_uuid = item_types.id
	// 		INNER JOIN transfers ON transfers_items.Parent = transfers.id
	// 		INNER JOIN documents ON transfers.document_type_uuid = documents.id
	// 		INNER JOIN stores ON transfers.store_origin_uuid = stores.id
	// 		INNER JOIN stores ss ON transfers.store_destination_uuid = ss.id
	// 		INNER JOIN warehouses wh_origin ON transfers.warehouse_origin_uuid = wh_origin.id
	// 		INNER JOIN warehouses wh_destination ON transfers.warehouse_destination_uuid = wh_destination.id
	// 		INNER JOIN responsibility_center ON transfers.responsibility_center_uuid = responsibility_center.id
	// 		%s
	// 	GROUP BY %s, %s) t
	// GROUP BY %s
	// `

	query := `
	SELECT %s, COUNT(*) Count FROM (
		SELECT
			%s, %s, COUNT(*) Count
		FROM %s
			%s
			INNER JOIN %s ON %s.%s = %s.%s 
			%s 
		GROUP BY %s, %s) t
	GROUP BY %s
	`

	return fmt.Sprintf(query,
		firstCol.DBNameShort,
		firstCol.DBName, secondCol.DBName,
		g.lineTableName,
		g.cfg.QueryChildJoins,
		g.tableName, g.tableName, g.cfg.ParentIdField, g.lineTableName, g.cfg.ChildJoinFieldWithParent,
		where, firstCol.DBName, secondCol.DBName, firstCol.DBNameShort)
}
