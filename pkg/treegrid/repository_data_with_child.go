package treegrid

import (
	"database/sql"
	"fmt"
	"log"
	"math"
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
	column := NewColumn(tg.GroupCols[0], g.childFieldMapping, g.parentFieldMapping)
	filterWhere, filterArgs := PrepQueryComplex(tg.FilterParams, g.parentFieldMapping, g.childFieldMapping)

	querySQL := NewConnectableSQL("")

	if column.IsItem {
		if !tg.WithGroupBy() {
			querySQL.Set(g.cfg.QueryChildCount)
		} else {
			queryCountSQL, err := NamedSQL(`
				SELECT COUNT(*) FROM (SELECT {{groupColumn}} 
				                      FROM {{lineTableName}}
				                      {{queryChildJoins}}
				                      {{childWhere}}
				                      GROUP BY {{groupColumn}}) t;
				`,
				map[string]string{
					"groupColumn":     column.DBName,
					"lineTableName":   g.lineTableName,
					"queryChildJoins": g.cfg.QueryChildJoins,
					"childWhere":      ChildDummyWhere,
				})
			if err != nil {
				return 0, err
			}
			querySQL.Set(queryCountSQL)
		}

		if filterWhere["parent"] != "" {
			parentQuery, err := NamedSQL(`
				{{lineTableName}}.{{parentId}} IN (SELECT {{tableName}}.{{id}} 
										 FROM {{tableName}} 
										 {{queryParentJoins}}
										 {{parentWhere}})
				`,
				map[string]string{
					"lineTableName":    g.lineTableName,
					"parentId":         g.cfg.ChildJoinFieldWithParent,
					"tableName":        g.tableName,
					"id":               g.cfg.ParentIdField,
					"queryParentJoins": g.cfg.QueryParentJoins,
					"parentWhere":      ParentDummyWhere,
				})
			if err != nil {
				return 0, err
			}

			querySQL.ConcatChildWhere(parentQuery)
			querySQL.ConcatParentWhere(filterWhere["parent"], filterArgs["parent"]...)
		}

		querySQL.ConcatChildWhere(filterWhere["child"], filterArgs["child"]...)
	} else {
		if !tg.WithGroupBy() {
			querySQL.Set(g.cfg.QueryParentCount)
		} else {
			// if grouped columns contain any child column, join with line table
			var joinChildSQL string
			var err error
			if tg.GroupCols.ContainsAny(g.childFieldMapping) {
				joinChildSQL, err = NamedSQL(
					`INNER JOIN {{lineTableName}} ON {{lineTableName}}.{{parentId}} = {{tableName}}.{{id}}`,
					map[string]string{
						"lineTableName": g.lineTableName,
						"parentId":      g.cfg.ChildJoinFieldWithParent,
						"tableName":     g.tableName,
						"id":            g.cfg.ParentIdField,
					})
				if err != nil {
					return 0, err
				}
			}

			queryCountSQL, err := NamedSQL(`
				SELECT COUNT(*) FROM (SELECT {{groupColumn}} 
				                      FROM {{tableName}}
				                      {{queryParentJoins}}
				                      {{joinChildSQL}}
				                      {{parentWhere}}
				                      GROUP BY {{groupColumn}}) t;
				`,
				map[string]string{
					"groupColumn":      column.DBName,
					"tableName":        g.tableName,
					"queryParentJoins": g.cfg.QueryParentJoins,
					"joinChildSQL":     joinChildSQL,
					"parentWhere":      ParentDummyWhere,
				})
			if err != nil {
				return 0, err
			}
			querySQL.Set(queryCountSQL)
		}

		if filterWhere["child"] != "" {
			childQuery, err := NamedSQL(`
				{{tableName}}.{{id}} IN (SELECT {{lineTableName}}.{{parentId}} 
										 FROM {{lineTableName}} 
										 {{queryChildJoins}}
										 {{childWhere}})
				`,
				map[string]string{
					"tableName":       g.tableName,
					"id":              g.cfg.ParentIdField,
					"lineTableName":   g.lineTableName,
					"parentId":        g.cfg.ChildJoinFieldWithParent,
					"queryChildJoins": g.cfg.QueryChildJoins,
					"childWhere":      ChildDummyWhere,
				})
			if err != nil {
				return 0, err
			}

			querySQL.ConcatParentWhere(childQuery)
			querySQL.ConcatChildWhere(filterWhere["child"], filterArgs["child"]...)
		}

		querySQL.ConcatParentWhere(filterWhere["parent"], filterArgs["parent"]...)
	}

	rows, err := g.db.Query(querySQL.SQL, querySQL.Args...)
	if err != nil {
		log.Println(err, "query", querySQL.SQL, "colData", column)
		return 0, err
	}
	defer rows.Close()

	return int64(math.Ceil(float64(utils.CheckCount(rows)) / float64(g.pageSize))), nil
}

// GetPageData implements GridRowDataRepositoryWithChild
func (g *gridRowDataRepositoryWithChild) GetPageData(tg *Treegrid) ([]map[string]string, error) {
	PrepFilters(tg, g.parentFieldMapping, g.childFieldMapping)

	// items request
	if tg.BodyParams.GetItemsRequest() {
		querySQL := NewConnectableSQL(g.cfg.QueryChild)
		queryParentIdWhere, err := NamedSQL(`
			{{lineTableName}}.{{parentId}} = ?
			`,
			map[string]string{
				"lineTableName": g.lineTableName,
				"parentId":      g.cfg.ChildJoinFieldWithParent,
			})
		if err != nil {
			return nil, err
		}
		querySQL.ConcatChildWhere(queryParentIdWhere, tg.BodyParams.ID)
		querySQL.ConcatChildWhere(tg.FilterWhere["child"], tg.FilterArgs["child"]...)
		// order by
		querySQL.Append(tg.OrderByChildQuery(g.childFieldMapping, fmt.Sprintf("%s.id ASC", g.lineTableName)))
		// pagination
		pos, _ := tg.BodyParams.IntPos()
		querySQL.SQL = AppendLimitToQuery(querySQL.SQL, g.pageSize, pos)

		return g.getJSON(querySQL.SQL, querySQL.Args, tg)
	}

	// GROUP BY
	if tg.WithGroupBy() {
		return g.handleGroupBy(tg)
	}

	// parent request
	querySQL := NewConnectableSQL(g.cfg.QueryParent)
	querySQL.ConcatParentWhere(tg.FilterWhere["parent"], tg.FilterArgs["parent"]...)
	if tg.FilterWhere["child"] != "" {
		childQuery, err := NamedSQL(`
				{{tableName}}.{{id}} IN (SELECT {{lineTableName}}.{{parentId}} 
										 FROM {{lineTableName}} 
										 {{queryChildJoins}}
										 {{childWhere}})
				`,
			map[string]string{
				"tableName":       g.tableName,
				"id":              g.cfg.ParentIdField,
				"lineTableName":   g.lineTableName,
				"parentId":        g.cfg.ChildJoinFieldWithParent,
				"queryChildJoins": g.cfg.QueryChildJoins,
				"childWhere":      ChildDummyWhere,
			})
		if err != nil {
			return nil, err
		}

		querySQL.ConcatParentWhere(childQuery)
		querySQL.ConcatChildWhere(tg.FilterWhere["child"], tg.FilterArgs["child"]...)
	}

	// order by
	querySQL.Append(tg.SortParams.OrderByQueryExcludeChild(g.childFieldMapping, g.parentFieldMapping, fmt.Sprintf("%s.id ASC", g.tableName)))
	// pagination
	pos, _ := tg.BodyParams.IntPos()
	querySQL.SQL = AppendLimitToQuery(querySQL.SQL, g.pageSize, pos)

	return g.getJSON(querySQL.SQL, querySQL.Args, tg)
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
			// entry["MinLevels"] = "2"
		}

		tableData = append(tableData, entry)
	}

	return tableData, nil
}

func (g *gridRowDataRepositoryWithChild) handleGroupBy(tg *Treegrid) ([]map[string]string, error) {
	if tg.BodyParams.Rows != "" {
		return g.getGroupData(tg.BodyParams.GetRowParentWhere(), tg.BodyParams.GetRowChildWhere(), tg)
	}

	// first level grouped rowsu
	return g.getGroupData("WHERE TRUE", "WHERE TRUE", tg)
}

func (g *gridRowDataRepositoryWithChild) getGroupData(parentRowsWhere, childRowsWhere string, tg *Treegrid) ([]map[string]string, error) {
	query, groupColumn := g.prepareNameCountQuery(parentRowsWhere, childRowsWhere, tg)

	pos, _ := tg.BodyParams.IntPos()
	query = AppendLimitToQuery(query, g.pageSize, pos)

	if groupColumn.IsItem {
		return g.getChildData(tg.BodyParams.GetRowLevel(), tg.GroupCols, parentRowsWhere, childRowsWhere, query, groupColumn)
	}

	return g.getParentData(tg.BodyParams.GetRowLevel(), tg.GroupCols, parentRowsWhere, childRowsWhere, query, groupColumn)
}

// generateChildNameCountQuery generates query for child name count
func (g *gridRowDataRepositoryWithChild) generateChildNameCountQuery(tg *Treegrid, parentRowsWhere, childRowsWhere string, column Column) string {
	queryCountSQL, err := NamedSQL(`
		SELECT {{groupColumnShort}}, COUNT(*) AS Count
		FROM (SELECT {{groupColumn}} AS {{groupColumnShort}}, {{tableName}}.id
			  FROM {{lineTableName}}
			  {{queryChildJoins}}
			  INNER JOIN {{tableName}} ON {{tableName}}.{{id}} = {{lineTableName}}.{{parentId}}
			  {{parentWhere}} AND {{childWhere}}
			  GROUP BY {{groupColumn}}, {{tableName}}.id) t
		GROUP BY {{groupColumnShort}}
		`,
		map[string]string{
			"groupColumnShort": column.DBNameShort,
			"groupColumn":      column.DBName,
			"tableName":        g.tableName,
			"lineTableName":    g.lineTableName,
			"queryChildJoins":  g.cfg.QueryChildJoins,
			"id":               g.cfg.ParentIdField,
			"parentId":         g.cfg.ChildJoinFieldWithParent,
			"parentWhere":      ParentDummyWhere,
			"childWhere":       ChildWhereTag,
		})
	if err != nil {
		return ""
	}

	querySQL := NewConnectableSQL(queryCountSQL)
	querySQL.ConcatParentWhere(parentRowsWhere)
	querySQL.ConcatChildWhere(childRowsWhere)
	querySQL.ConcatParentWhere(tg.FilterWhere["parent"], tg.FilterArgs["parent"]...)
	querySQL.ConcatChildWhere(tg.FilterWhere["child"], tg.FilterArgs["child"]...)

	return querySQL.AsSQL()
}

func (g *gridRowDataRepositoryWithChild) prepareNameCountQuery(parentRowsWhere, childRowsWhere string, tg *Treegrid) (query string, column Column) {
	// If both the level are equal then return the row
	level := tg.BodyParams.GetRowLevel()
	// query parent rows
	if level == len(tg.GroupCols) {
		queryDataSQL := NewConnectableSQL(g.cfg.QueryParent)
		queryDataSQL.ConcatParentWhere(parentRowsWhere)
		// query all child rows with filter, do not concat childRowsWhere
		queryDataSQL.ConcatParentWhere(tg.FilterWhere["parent"], tg.FilterArgs["parent"]...)
		if tg.FilterWhere["child"] != "" {
			childQuery, err := NamedSQL(`
				{{tableName}}.{{id}} IN (SELECT {{lineTableName}}.{{parentId}} 
										 FROM {{lineTableName}} 
										 {{queryChildJoins}}
										 {{childWhere}})
				`,
				map[string]string{
					"tableName":       g.tableName,
					"id":              g.cfg.ParentIdField,
					"lineTableName":   g.lineTableName,
					"parentId":        g.cfg.ChildJoinFieldWithParent,
					"queryChildJoins": g.cfg.QueryChildJoins,
					"childWhere":      ChildDummyWhere,
				})
			if err != nil {
				return
			}

			queryDataSQL.ConcatParentWhere(childQuery)
			queryDataSQL.ConcatChildWhere(tg.FilterWhere["child"], tg.FilterArgs["child"]...)
		}
		query = queryDataSQL.AsSQL()
		return
	}

	column = NewColumn(tg.GroupCols[level], g.childFieldMapping, g.parentFieldMapping)

	// multiple groupBy clauses
	if len(tg.GroupCols)-level > 1 {
		secColumn := NewColumn(tg.GroupCols[level+1], g.childFieldMapping, g.parentFieldMapping)

		switch {
		case !column.IsItem && !secColumn.IsItem:
			return g.getCascadingGroupByParentParent(tg, column, secColumn, parentRowsWhere, childRowsWhere), column
		case column.IsItem && !secColumn.IsItem, !column.IsItem && secColumn.IsItem:
			return g.getCascadingGroupByParentChild(tg, column, secColumn, parentRowsWhere, childRowsWhere), column
		case column.IsItem && secColumn.IsItem:
			return g.getCascadingGroupByChildChild(tg, column, secColumn, parentRowsWhere, childRowsWhere), column
		}
	}

	if column.IsItem {
		if val, ok := g.childFieldMapping[column.GridName]; ok {
			column.DBName = val[0]
		}

		query = g.generateChildNameCountQuery(tg, parentRowsWhere, childRowsWhere, column)
		return query, column
	}

	// only one groupBy clause, query parent rows
	queryCountSQL, err := NamedSQL(`
		SELECT {{groupColumn}}, COUNT(*) AS Count
		FROM {{tableName}}
		{{queryParentJoins}}
		{{parentWhere}}
		GROUP BY {{groupColumn}}
		`,
		map[string]string{
			"groupColumn":      column.DBName,
			"tableName":        g.tableName,
			"queryParentJoins": g.cfg.QueryParentJoins,
			"parentWhere":      ParentDummyWhere,
		})
	if err != nil {
		return "", column
	}

	querySQL := NewConnectableSQL(queryCountSQL)
	querySQL.ConcatParentWhere(parentRowsWhere)
	querySQL.ConcatParentWhere(tg.FilterWhere["parent"], tg.FilterArgs["parent"]...)
	if tg.FilterWhere["child"] != "" {
		childQuery, err := NamedSQL(`
				{{tableName}}.{{id}} IN (SELECT {{lineTableName}}.{{parentId}} 
										 FROM {{lineTableName}} 
										 {{queryChildJoins}}
										 {{childWhere}})
				`,
			map[string]string{
				"tableName":       g.tableName,
				"id":              g.cfg.ParentIdField,
				"lineTableName":   g.lineTableName,
				"parentId":        g.cfg.ChildJoinFieldWithParent,
				"queryChildJoins": g.cfg.QueryChildJoins,
				"childWhere":      ChildDummyWhere,
			})
		if err != nil {
			return
		}

		querySQL.ConcatParentWhere(childQuery)
		querySQL.ConcatChildWhere(tg.FilterWhere["child"], tg.FilterArgs["child"]...)
	}
	query = querySQL.AsSQL()

	return
}

func (g *gridRowDataRepositoryWithChild) getParentData(level int, groupCols []string, parentWhere, childWhere string, query string, groupColumn Column) ([]map[string]string, error) {
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

		if level == len(groupCols) {
			for k := range row.StringValues() {
				tempObj[k] = row.StringValues()[k]
			}
			tableData = append(tableData, tempObj)
			continue
		}

		// remove Def for editable parent rows
		tempObj["Def"] = "Group"

		docType := row.GetValue(groupColumn.DBNameShort)
		if docType == "" {
			docType = row.GetValue(groupColumn.GridName)
		}

		tempObj[g.cfg.MainCol] = docType
		tempObj[g.cfg.MainCol+"Type"] = groupColumn.Type()
		// TODO get from treegrid config
		if groupColumn.IsDate {
			tempObj[g.cfg.MainCol+"Format"] = "yyyy-MM-dd"
		}

		val := docType
		where2 := " AND " + groupColumn.WhereSQL(val)

		// Builds new attribute Rows for identification
		tempObj["Rows"] = SetBodyParamRows(level+1, parentWhere+where2, childWhere)
		tableData = append(tableData, tempObj)
	}

	return tableData, nil
}

func (g *gridRowDataRepositoryWithChild) getChildData(level int, groupCols []string, parentWhere, childWhere string, query string, groupColumn Column) ([]map[string]string, error) {
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
		tempObj[g.cfg.MainCol] = row.GetValue(groupColumn.DBNameShort)
		if level == len(groupCols) {
			tableData = append(tableData, row.StringValues())
			continue
		}

		// rows where condition
		var rowsParentWhere string
		var rowsChildWhere string

		// If grouped by child column
		if strings.Index(parentWhere, " AND "+groupColumn.DBName+"='") > 0 {
			// If the column is already added in the WHERE clause
			// Just update the query to adjust column value
			rowsParentWhere = utils.ReplaceColumnValueInQuery(parentWhere, groupColumn.DBName, row.GetValue(groupColumn.DBNameShort))
		} else {
			rowsParentWhere = g.addChildCondition(parentWhere, groupColumn.DBName, row.GetValue(groupColumn.DBNameShort))
		}
		rowsChildWhere = childWhere + " AND " + groupColumn.DBName + "='" + row.GetValue(groupColumn.DBNameShort) + "' "
		tempObj["Rows"] = SetBodyParamRows(level+1, rowsParentWhere, rowsChildWhere)

		tableData = append(tableData, tempObj)
	}

	return tableData, nil
}

func (g *gridRowDataRepositoryWithChild) addChildCondition(parentWhere, colName string, colValue string) string {
	tmp := fmt.Sprintf("%s IN ( SELECT %s FROM %s", g.cfg.ParentIdField, g.cfg.ChildJoinFieldWithParent, g.lineTableName)
	if strings.Index(parentWhere, tmp) > 0 {
		// If already filtered on child column(s)
		parentWhere = strings.Replace(parentWhere, ")", " AND "+colName+"='"+colValue+"' )", 1)
	} else {
		// parentWhere += ` AND transfers.id IN ( SELECT Parent FROM transfers_items
		// INNER JOIN items ON transfers_items.item_uuid = items.id
		// INNER JOIN units ON transfers_items.item_unit_uuid = units.id
		// INNER JOIN item_types ON items.type_uuid = item_types.id
		// WHERE ` + colName + "='" + colValue + "' ) "

		parentWhere += fmt.Sprintf(` AND %s.%s IN ( SELECT %s FROM %s `, g.tableName, g.cfg.ParentIdField, g.cfg.ChildJoinFieldWithParent, g.lineTableName) +
			g.cfg.QueryChildJoins +
			" WHERE " + colName + "='" + colValue + "' ) "
	}

	return parentWhere
}

func (g *gridRowDataRepositoryWithChild) getCascadingGroupByParentParent(tg *Treegrid, firstCol, secondCol Column, parentRowsWhere, childRowsWhere string) string {
	// if grouped columns contain any child column, join with line table
	var joinChildSQL string
	var err error
	if tg.GroupCols.ContainsAny(g.childFieldMapping) {
		joinChildSQL, err = NamedSQL(
			`INNER JOIN {{lineTableName}} ON {{lineTableName}}.{{parentId}} = {{tableName}}.{{id}}`,
			map[string]string{
				"lineTableName": g.lineTableName,
				"parentId":      g.cfg.ChildJoinFieldWithParent,
				"tableName":     g.tableName,
				"id":            g.cfg.ParentIdField,
			})
		if err != nil {
			return ""
		}
	}

	query, err := NamedSQL(`
		SELECT {{firstColDBNameShort}}, COUNT(*) Count
		FROM (
			SELECT {{firstColDBName}}, {{secondColDBName}}
			FROM {{tableName}}
			{{queryParentJoins}}
			{{joinChildSQL}}
			{{parentWhere}}
			GROUP BY {{firstColDBName}}, {{secondColDBName}}
		) t
		GROUP BY {{firstColDBNameShort}}
	`,
		map[string]string{
			"firstColDBNameShort": firstCol.DBNameShort,
			"firstColDBName":      firstCol.DBName,
			"secondColDBName":     secondCol.DBName,
			"tableName":           g.tableName,
			"queryParentJoins":    g.cfg.QueryParentJoins,
			"joinChildSQL":        joinChildSQL,
			"parentWhere":         ParentDummyWhere,
		})
	if err != nil {
		return ""
	}

	querySQL := NewConnectableSQL(query)
	querySQL.ConcatParentWhere(parentRowsWhere)
	querySQL.ConcatParentWhere(tg.FilterWhere["parent"], tg.FilterArgs["parent"]...)
	if tg.FilterWhere["child"] != "" {
		childQuery, err := NamedSQL(`
				{{tableName}}.{{id}} IN (SELECT {{lineTableName}}.{{parentId}} 
										 FROM {{lineTableName}} 
										 {{queryChildJoins}}
										 {{childWhere}})
				`,
			map[string]string{
				"tableName":       g.tableName,
				"id":              g.cfg.ParentIdField,
				"lineTableName":   g.lineTableName,
				"parentId":        g.cfg.ChildJoinFieldWithParent,
				"queryChildJoins": g.cfg.QueryChildJoins,
				"childWhere":      ChildDummyWhere,
			})
		if err != nil {
			return ""
		}

		querySQL.ConcatParentWhere(childQuery)
		querySQL.ConcatChildWhere(tg.FilterWhere["child"], tg.FilterArgs["child"]...)
	}

	return querySQL.AsSQL()
}

// when grouping by two parent (tansfer) columns
func (g *gridRowDataRepositoryWithChild) getCascadingGroupByParentChild(tg *Treegrid, firstCol, secondCol Column, parentRowsWhere, childRowsWhere string) string {
	query, err := NamedSQL(`
		SELECT {{firstColDBNameShort}}, COUNT(*) Count FROM (
			SELECT
				{{firstColDBName}}, {{secondColDBName}}, COUNT(*) Count
			FROM {{lineTableName}}
				{{queryChildJoins}}
				INNER JOIN {{tableName}} ON {{tableName}}.{{id}} = {{lineTableName}}.{{parentId}}
				{{parentWhere}} AND {{childWhere}}
			GROUP BY {{firstColDBName}}, {{secondColDBName}}) t
		GROUP BY {{firstColDBNameShort}}
	`,
		map[string]string{
			"firstColDBNameShort": firstCol.DBNameShort,
			"firstColDBName":      firstCol.DBName,
			"secondColDBName":     secondCol.DBName,
			"tableName":           g.tableName,
			"lineTableName":       g.lineTableName,
			"queryChildJoins":     g.cfg.QueryChildJoins,
			"id":                  g.cfg.ParentIdField,
			"parentId":            g.cfg.ChildJoinFieldWithParent,
			"parentWhere":         ParentDummyWhere,
			"childWhere":          ChildWhereTag,
		})
	if err != nil {
		return ""
	}

	querySQL := NewConnectableSQL(query)
	querySQL.ConcatParentWhere(parentRowsWhere)
	querySQL.ConcatChildWhere(childRowsWhere)
	querySQL.ConcatParentWhere(tg.FilterWhere["parent"], tg.FilterArgs["parent"]...)
	querySQL.ConcatChildWhere(tg.FilterWhere["child"], tg.FilterArgs["child"]...)
	return querySQL.AsSQL()
}

// when grouping by two child columns
func (g *gridRowDataRepositoryWithChild) getCascadingGroupByChildChild(tg *Treegrid, firstCol, secondCol Column, parentRowsWhere, childRowsWhere string) string {
	query, err := NamedSQL(`
		SELECT {{firstColDBNameShort}}, COUNT(*) AS Count
		FROM (
			SELECT
				{{firstColDBName}}, {{secondColDBName}}, COUNT(*) AS Count
			FROM {{lineTableName}}
				{{queryChildJoins}}
				INNER JOIN {{tableName}} ON {{tableName}}.{{id}} = {{lineTableName}}.{{parentId}}
				{{childWhere}} AND {{parentWhere}}
			GROUP BY {{firstColDBName}}, {{secondColDBName}}) t
		GROUP BY {{firstColDBNameShort}}
	`,
		map[string]string{
			"firstColDBNameShort": firstCol.DBNameShort,
			"firstColDBName":      firstCol.DBName,
			"secondColDBName":     secondCol.DBName,
			"lineTableName":       g.lineTableName,
			"tableName":           g.tableName,
			"queryChildJoins":     g.cfg.QueryChildJoins,
			"id":                  g.cfg.ParentIdField,
			"parentId":            g.cfg.ChildJoinFieldWithParent,
			"childWhere":          ChildDummyWhere,
			"parentWhere":         ParentWhereTag,
		})
	if err != nil {
		return ""
	}

	querySQL := NewConnectableSQL(query)
	querySQL.ConcatChildWhere(childRowsWhere)
	querySQL.ConcatParentWhere(parentRowsWhere)
	querySQL.ConcatChildWhere(tg.FilterWhere["child"], tg.FilterArgs["child"]...)
	querySQL.ConcatParentWhere(tg.FilterWhere["parent"], tg.FilterArgs["parent"]...)
	return querySQL.AsSQL()
}
