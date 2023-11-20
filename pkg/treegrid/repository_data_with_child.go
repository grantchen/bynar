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
	GetChildSuggestion(tg *Treegrid) ([]map[string]interface{}, error)
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
	QueryChildSuggestion     string
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

// GetChildSuggestion implements GridRowDataRepositoryWithChild
func (g *gridRowDataRepositoryWithChild) GetChildSuggestion(tg *Treegrid) ([]map[string]interface{}, error) {
	var query string
	id := tg.BodyParams.ID
	origin := tg.BodyParams.TreegridOriginID
	split := strings.Split(origin, "$")

	if len(split) > 1 {
		id = split[len(split)-2]
	}
	searchValue := "%" + tg.BodyParams.Val + "%"
	query = g.cfg.QueryChildSuggestion
	AppendLimitToQuery(query, g.pageSize, 0)
	params := []interface{}{searchValue, id}
	logger.Debug("query: ", query, "param: ", params)
	data, err := g.getJSON(query, params, tg)
	if err != nil {
		return nil, err
	}

	mapResult := make([]map[string]interface{}, 0)

	for _, element := range data {
		newElement := make(map[string]interface{})

		for k, v := range element {
			newElement[k] = v
		}
		mapResult = append(mapResult, newElement)
	}

	return mapResult, nil
}

// GetPageCount implements GridRowDataRepositoryWithChild
func (g *gridRowDataRepositoryWithChild) GetPageCount(tg *Treegrid) (int64, error) {
	var query string

	column := NewColumn(tg.GroupCols[0], g.childFieldMapping, g.parentFieldMapping)
	filterWhere, filterArgs := PrepQueryComplex(tg.FilterParams, g.parentFieldMapping, g.childFieldMapping)

	sqlQuery := ""

	if column.IsItem {
		if filterWhere["parent"] != "" {
			// sqlQuery = " AND transfers_items.Parent IN (SELECT transfers.id from transfers " +
			sqlQuery = fmt.Sprintf(" AND %s.%s IN (SELECT %s.%s FROM %s ",
				g.lineTableName,
				g.cfg.ChildJoinFieldWithParent,
				g.tableName,
				g.cfg.ParentIdField,
				g.tableName) +
				g.cfg.QueryParentJoins +
				ConcatParentWhereToSQL(ParentDummyWhere, filterWhere["parent"]) +
				") "
		}

		if !tg.WithGroupBy() {
			query = g.cfg.QueryChildCount
			query = ConcatChildWhereToSQL(query, filterWhere["child"])
			query = ConcatParentWhereToSQL(query, filterWhere["parent"])
			query = ConcatChildWhereToSQL(query, sqlQuery)
		} else {
			where := ChildDummyWhere
			where = ConcatChildWhereToSQL(where, filterWhere["child"])
			where = ConcatParentWhereToSQL(where, filterWhere["parent"])
			where = ConcatChildWhereToSQL(where, sqlQuery)
			query = fmt.Sprintf(`SELECT COUNT(*) FROM (SELECT %s FROM %s%s%s GROUP BY %s) t;`,
				column.DBName,
				g.lineTableName,
				g.cfg.QueryChildJoins,
				where,
				column.DBName,
			)
		}
	} else {
		if filterWhere["child"] != "" {
			// sqlQuery = " AND transfers.id IN (SELECT transfers_items.Parent from transfers_items " +
			sqlQuery = fmt.Sprintf(" AND %s.%s IN (SELECT %s.%s FROM %s ",
				g.tableName,
				g.cfg.ParentIdField,
				g.lineTableName,
				g.cfg.ChildJoinFieldWithParent,
				g.lineTableName) +
				g.cfg.QueryChildJoins +
				ParentDummyWhere +
				filterWhere["child"] + ") "

		}

		if !tg.WithGroupBy() {
			query = g.cfg.QueryParentCount
			query = ConcatChildWhereToSQL(query, filterWhere["child"])
			query = ConcatParentWhereToSQL(query, filterWhere["parent"])
			query = ConcatParentWhereToSQL(query, sqlQuery)
		} else {
			where := ParentDummyWhere
			where = ConcatChildWhereToSQL(where, filterWhere["child"])
			where = ConcatParentWhereToSQL(where, filterWhere["parent"])
			query = ConcatParentWhereToSQL(query, sqlQuery)
			query = fmt.Sprintf(`SELECT COUNT(*) FROM (SELECT %s FROM %s%s%s GROUP BY %s) t;`,
				column.DBName,
				g.tableName,
				g.cfg.QueryParentJoins,
				where,
				column.DBName,
			)
		}
	}

	// TODO args is not matched query ？use @args
	mergedArgs := utils.MergeMaps(filterArgs["child"], filterArgs["parent"])
	rows, err := g.db.Query(query, mergedArgs...)
	if err != nil {
		log.Println(err, "query", query, "colData", column)
		return 0, err
	}
	defer rows.Close()

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
		query := ConcatChildWhereToSQL(g.cfg.QueryChild,
			fmt.Sprintf(" WHERE %s.%s = ", g.lineTableName, g.cfg.ChildJoinFieldWithParent)+
				tg.BodyParams.ID+
				tg.FilterWhere["child"]+
				tg.OrderByChildQuery(g.toBooleanMapping(g.childFieldMapping)))
		pos, _ := tg.BodyParams.IntPos()
		query = AppendLimitToQuery(query, g.pageSize, pos)
		return g.getJSON(query, tg.FilterArgs["child"], tg)
	}

	// GROUP BY
	if tg.WithGroupBy() {
		return g.handleGroupBy(tg)
	}

	query := ConcatParentWhereToSQL(g.cfg.QueryParent, tg.FilterWhere["parent"])
	if tg.FilterWhere["child"] != "" {
		query = ConcatParentWhereToSQL(query, fmt.Sprintf(" AND %s.%s IN ( SELECT %s FROM %s ", g.tableName, g.cfg.ParentIdField, g.cfg.ChildJoinFieldWithParent, g.lineTableName)+
			g.cfg.QueryChildJoins+
			tg.FilterWhere["child"]+
			`) `)
	}

	query += tg.SortParams.OrderByQueryExludeChild2(g.toBooleanMapping(g.childFieldMapping), g.parentFieldMapping)

	pos, _ := tg.BodyParams.IntPos()
	query = AppendLimitToQuery(query, g.pageSize, pos)
	mergedArgs := utils.MergeMaps(tg.FilterArgs["parent"], tg.FilterArgs["child"])

	// TODO
	return g.getJSON(query, mergedArgs, tg)
}

func (g *gridRowDataRepositoryWithChild) getJSON(sqlString string, mergedArgs []interface{}, tg *Treegrid) ([]map[string]string, error) {
	stmt, err := g.db.Prepare(sqlString)
	if err != nil {
		return nil, fmt.Errorf("db prepare: [%w], sql string: [%s]", err, sqlString)
	} else {
		fmt.Println(sqlString)
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
			// entry["MinLevels"] = "2"
			entry["has_child"] = "1"
		}

		tableData = append(tableData, entry)
	}

	return tableData, nil
}

func (g *gridRowDataRepositoryWithChild) handleGroupBy(tg *Treegrid) ([]map[string]string, error) {
	if tg.BodyParams.Rows != "" {
		return g.getGroupData(tg.BodyParams.GetRowWhere(), tg)
	}

	// parentBuild := QueryParent
	parentBuild := g.cfg.QueryParent

	if tg.FilterWhere["parent"] != "" {
		parentBuild = ConcatParentWhereToSQL(parentBuild, tg.FilterWhere["parent"])
	}

	if tg.FilterWhere["child"] != "" {
		parentBuild = ConcatParentWhereToSQL(parentBuild,
			fmt.Sprintf(" AND %s.%s IN (SELECT %s FROM %s", g.tableName, g.cfg.ParentIdField, g.cfg.ChildJoinFieldWithParent, g.lineTableName)+
				g.cfg.QueryChildJoins+
				ParentDummyWhere+
				tg.FilterWhere["child"]+
				tg.OrderByChildQuery(g.toBooleanMapping(g.childFieldMapping))+") ")
	}

	mergedArgs := utils.MergeMaps(tg.FilterArgs["parent"], tg.FilterArgs["child"])

	extWhereClause := utils.ExtractWhereClause(parentBuild, mergedArgs, ParentWhereTag)
	return g.getGroupData(extWhereClause, tg)
}

func (g *gridRowDataRepositoryWithChild) getGroupData(where string, tg *Treegrid) ([]map[string]string, error) {
	// TODO here 加 where 条件到 left join 中
	query, colData := g.prepareNameCountQuery(where, tg)

	pos, _ := tg.BodyParams.IntPos()
	query = AppendLimitToQuery(query, g.pageSize, pos)

	if colData.IsItem {
		return g.getChildData(tg.BodyParams.GetRowLevel(), tg.GroupCols, where, query, colData)
	}

	return g.getParentData(tg.BodyParams.GetRowLevel(), tg.GroupCols, where, query, colData)
}

func (g *gridRowDataRepositoryWithChild) generateChildNameCountQuery(where string, column Column, tg *Treegrid) string {
	additionChildCondition := utils.ExtractWhereClause("WHERE "+tg.FilterWhere["child"], tg.FilterArgs["child"], ChildWhereTag)
	additionChildCondition = additionChildCondition[len("WHERE"):]
	where = ConcatChildWhereToSQL(where, additionChildCondition)

	s := fmt.Sprintf(`
SELECT %s, COUNT(*) AS Count
FROM (SELECT %s AS %s, %s.id
      FROM %s
               %s
               INNER JOIN %s ON %s.%s = %s.%s
      %s
      GROUP BY %s, %s.id) t
GROUP BY %s
		`,
		column.DBNameShort,
		column.DBName, column.DBNameShort, g.tableName,
		g.lineTableName,
		g.cfg.QueryChildJoins,
		g.tableName, g.tableName, g.cfg.ParentIdField, g.lineTableName, g.cfg.ChildJoinFieldWithParent,
		where,
		column.DBName, g.tableName,
		column.DBNameShort,
	)

	return s
}

func (g *gridRowDataRepositoryWithChild) prepareNameCountQuery(where string, tg *Treegrid) (query string, column Column) {
	// If both the level are equal then return the row
	level := tg.BodyParams.GetRowLevel()
	if level == len(tg.GroupCols) {
		query = ConcatParentWhereToSQL(g.cfg.QueryParent, where)
		return
	}

	column = NewColumn(tg.GroupCols[level], g.childFieldMapping, g.parentFieldMapping)

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
		if val, ok := g.childFieldMapping[column.GridName]; ok {
			column.DBName = val[0]
		}

		query = g.generateChildNameCountQuery(where, column, tg)
		return query, column
	}

	query = "SELECT " + column.DBName + ", COUNT(*) Count FROM " + g.tableName + " " + g.cfg.QueryParentJoins + where + " GROUP BY " + column.DBName

	return
}

func (g *gridRowDataRepositoryWithChild) getParentData(level int, group_cols []string, where string, query string, colData Column) ([]map[string]string, error) {
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

			// for expendable parent rows
			tempObj["has_child"] = "1"

			tableData = append(tableData, tempObj)
			continue
		}

		// remove Def for editable parent rows
		tempObj["Def"] = "Group"

		docType := row.GetValue(colData.DBNameShort)
		if docType == "" {
			docType = row.GetValue(colData.GridName)
		}

		tempObj[g.cfg.MainCol] = docType
		tempObj[g.cfg.MainCol+"Type"] = colData.Type()
		// TODO get from treegrid config
		if colData.IsDate {
			tempObj[g.cfg.MainCol+"Format"] = "yyyy-MM-dd"
		}

		val := docType
		where2 := " AND " + colData.WhereSQL(val)

		// Builds new attribute Rows for identification
		tempObj["Rows"] = strconv.Itoa(level+1) + where + where2
		tableData = append(tableData, tempObj)
	}

	return tableData, nil
}

func (g *gridRowDataRepositoryWithChild) getChildData(level int, groupCols []string, where string, query string, columnData Column) ([]map[string]string, error) {
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

		orignal_query += fmt.Sprintf(` AND %s.%s IN ( SELECT %s FROM %s `, g.tableName, g.cfg.ParentIdField, g.cfg.ChildJoinFieldWithParent, g.lineTableName) +
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
