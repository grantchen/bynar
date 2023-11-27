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

	stmt, err := g.db.Prepare(querySQL.SQL)
	if err != nil {
		return 0, fmt.Errorf("db prepare: [%w], query: [%s]", err, querySQL.SQL)
	}
	defer stmt.Close()

	rows, err := stmt.Query(querySQL.Args...)
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
	PrepRows(tg, g.parentFieldMapping, g.childFieldMapping)

	// items request(with or without group by)
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
		return g.getGroupData(tg)
	}

	// parent request without group by
	querySQL := NewConnectableSQL(g.cfg.QueryParent)
	querySQL.ConcatParentWhere(tg.FilterWhere["parent"], tg.FilterArgs["parent"]...)
	if tg.FilterWhere["child"] != "" {
		parentIn, err := g.childQueryToParentIn()
		if err != nil {
			return nil, err
		}
		querySQL.ConcatParentWhere(parentIn)
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
			entry["Def"] = "Node"
			// entry["MinLevels"] = "2"
		} else {
			entry["Def"] = "Data"
		}

		tableData = append(tableData, entry)
	}

	return tableData, nil
}

func (g *gridRowDataRepositoryWithChild) getGroupData(tg *Treegrid) ([]map[string]string, error) {
	querySQL, groupColumn := g.prepareNameCountQuery(tg)

	pos, _ := tg.BodyParams.IntPos()
	querySQL.SQL = AppendLimitToQuery(querySQL.SQL, g.pageSize, pos)

	if groupColumn.IsItem {
		return g.getChildData(tg, querySQL, groupColumn)
	}

	return g.getParentData(tg, querySQL, groupColumn)
}

// generateChildNameCountQuery generates query for child name count
func (g *gridRowDataRepositoryWithChild) generateChildNameCountQuery(tg *Treegrid, column Column) (querySQL *ConnectableSQL) {
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
		return
	}

	querySQL = NewConnectableSQL(queryCountSQL)
	querySQL.ConcatParentWhere(tg.RowsWhere["parent"], tg.RowsArgs["parent"]...)
	querySQL.ConcatChildWhere(tg.RowsWhere["child"], tg.RowsArgs["child"]...)
	querySQL.ConcatParentWhere(tg.FilterWhere["parent"], tg.FilterArgs["parent"]...)
	querySQL.ConcatChildWhere(tg.FilterWhere["child"], tg.FilterArgs["child"]...)
	return querySQL
}

func (g *gridRowDataRepositoryWithChild) prepareNameCountQuery(tg *Treegrid) (querySQL *ConnectableSQL, column Column) {
	// If both the level are equal then return the row
	level := tg.BodyParams.RowsLevel
	// query parent rows
	if level == len(tg.GroupCols) {
		querySQL = NewConnectableSQL(g.cfg.QueryParent)
		querySQL.ConcatParentWhere(tg.RowsWhere["parent"], tg.RowsArgs["parent"]...)
		querySQL.ConcatParentWhere(tg.FilterWhere["parent"], tg.FilterArgs["parent"]...)
		// query all child rows with filter, do not concat child RowsWhere
		if tg.FilterWhere["child"] != "" || tg.RowsWhere["child"] != "" {
			parentIn, err := g.childQueryToParentIn()
			if err != nil {
				return
			}

			parentInSQL := NewConnectableSQL(parentIn)
			parentInSQL.ConcatChildWhere(tg.RowsWhere["child"], tg.RowsArgs["child"]...) // only concat to parentInSQL!
			// concat parentInSQL to querySQL
			querySQL.ConcatParentWhere(parentInSQL.SQL, parentInSQL.Args...)
			querySQL.ConcatChildWhere(tg.FilterWhere["child"], tg.FilterArgs["child"]...) // concat to querySQL
		}
		return
	}

	column = NewColumn(tg.GroupCols[level], g.childFieldMapping, g.parentFieldMapping)

	// multiple groupBy clauses
	if len(tg.GroupCols)-level > 1 {
		secColumn := NewColumn(tg.GroupCols[level+1], g.childFieldMapping, g.parentFieldMapping)

		switch {
		case !column.IsItem && !secColumn.IsItem:
			return g.getCascadingGroupByParentParent(tg, column, secColumn), column
		case column.IsItem && !secColumn.IsItem, !column.IsItem && secColumn.IsItem:
			return g.getCascadingGroupByParentChild(tg, column, secColumn), column
		case column.IsItem && secColumn.IsItem:
			return g.getCascadingGroupByChildChild(tg, column, secColumn), column
		}
	}

	// last group is item(one or multiple groupBy clauses)
	if column.IsItem {
		if val, ok := g.childFieldMapping[column.GridName]; ok {
			column.DBName = val[0]
		}

		return g.generateChildNameCountQuery(tg, column), column
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
		return
	}

	querySQL = NewConnectableSQL(queryCountSQL)
	querySQL.ConcatParentWhere(tg.RowsWhere["parent"], tg.RowsArgs["parent"]...)
	querySQL.ConcatParentWhere(tg.FilterWhere["parent"], tg.FilterArgs["parent"]...)
	if tg.FilterWhere["child"] != "" {
		parentIn, err := g.childQueryToParentIn()
		if err != nil {
			return
		}
		querySQL.ConcatParentWhere(parentIn)
		querySQL.ConcatChildWhere(tg.FilterWhere["child"], tg.FilterArgs["child"]...)
	}

	return
}

func (g *gridRowDataRepositoryWithChild) getParentData(tg *Treegrid, querySQL *ConnectableSQL, groupColumn Column) ([]map[string]string, error) {
	stmt, err := g.db.Prepare(querySQL.SQL)
	if err != nil {
		return nil, fmt.Errorf("db prepare: [%w], query: [%s]", err, querySQL.SQL)
	}
	defer stmt.Close()

	rows, err := stmt.Query(querySQL.Args...)
	if err != nil {
		return nil, fmt.Errorf("do query: [%w], query: [%s]", err, querySQL.AsSQL())
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

		if tg.BodyParams.RowsLevel == len(tg.GroupCols) {
			for k := range row.StringValues() {
				tempObj[k] = row.StringValues()[k]
			}
			tempObj["Def"] = "Node"
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

		// Builds new attribute Rows for identification
		tempObj["Rows"], err = tg.BodyParams.GetNewRows(
			tg.BodyParams.RowsLevel+1,
			RowsFieldCond{GridName: groupColumn.GridName, Value: row.GetValue(groupColumn.DBNameShort)},
		)
		if err != nil {
			return nil, err
		}

		tableData = append(tableData, tempObj)
	}

	return tableData, nil
}

func (g *gridRowDataRepositoryWithChild) getChildData(tg *Treegrid, querySQL *ConnectableSQL, groupColumn Column) ([]map[string]string, error) {
	stmt, err := g.db.Prepare(querySQL.SQL)
	if err != nil {
		return nil, fmt.Errorf("db prepare: [%w], query: [%s]", err, querySQL.SQL)
	}
	defer stmt.Close()

	rows, err := stmt.Query(querySQL.Args...)
	if err != nil {
		return nil, fmt.Errorf("do query: [%w], query: [%s]", err, querySQL.AsSQL())
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
		if tg.BodyParams.RowsLevel == len(tg.GroupCols) {
			tableData = append(tableData, row.StringValues())
			continue
		}

		tempObj["Rows"], err = tg.BodyParams.GetNewRows(
			tg.BodyParams.RowsLevel+1,
			RowsFieldCond{GridName: groupColumn.GridName, Value: row.GetValue(groupColumn.DBNameShort)},
		)
		if err != nil {
			return nil, err
		}

		tableData = append(tableData, tempObj)
	}

	return tableData, nil
}

func (g *gridRowDataRepositoryWithChild) getCascadingGroupByParentParent(tg *Treegrid, firstCol, secondCol Column) (querySQL *ConnectableSQL) {
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
			return nil
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
		return nil
	}

	querySQL = NewConnectableSQL(query)
	querySQL.ConcatParentWhere(tg.RowsWhere["parent"], tg.RowsArgs["parent"]...)
	querySQL.ConcatParentWhere(tg.FilterWhere["parent"], tg.FilterArgs["parent"]...)
	if tg.FilterWhere["child"] != "" {
		parentIn, err := g.childQueryToParentIn()
		if err != nil {
			return nil
		}
		querySQL.ConcatParentWhere(parentIn)
		querySQL.ConcatChildWhere(tg.FilterWhere["child"], tg.FilterArgs["child"]...)
	}

	return querySQL
}

// when grouping by two parent (tansfer) columns
func (g *gridRowDataRepositoryWithChild) getCascadingGroupByParentChild(tg *Treegrid, firstCol, secondCol Column) (querySQL *ConnectableSQL) {
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
		return
	}

	querySQL = NewConnectableSQL(query)
	querySQL.ConcatParentWhere(tg.RowsWhere["parent"], tg.RowsArgs["parent"]...)
	querySQL.ConcatChildWhere(tg.RowsWhere["child"], tg.RowsArgs["child"]...)
	querySQL.ConcatParentWhere(tg.FilterWhere["parent"], tg.FilterArgs["parent"]...)
	querySQL.ConcatChildWhere(tg.FilterWhere["child"], tg.FilterArgs["child"]...)
	return querySQL
}

// when grouping by two child columns
func (g *gridRowDataRepositoryWithChild) getCascadingGroupByChildChild(tg *Treegrid, firstCol, secondCol Column) (querySQL *ConnectableSQL) {
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
		return
	}

	querySQL = NewConnectableSQL(query)
	querySQL.ConcatChildWhere(tg.RowsWhere["child"], tg.RowsArgs["child"]...)
	querySQL.ConcatParentWhere(tg.RowsWhere["parent"], tg.RowsArgs["parent"]...)
	querySQL.ConcatChildWhere(tg.FilterWhere["child"], tg.FilterArgs["child"]...)
	querySQL.ConcatParentWhere(tg.FilterWhere["parent"], tg.FilterArgs["parent"]...)
	return querySQL
}

// childQueryToParentIn generates parent id in query from child query
func (g *gridRowDataRepositoryWithChild) childQueryToParentIn() (string, error) {
	parentIn, err := NamedSQL(`
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
		return "", err
	}
	return parentIn, nil
}
