package repository

import (
	"fmt"
	"strconv"
	"strings"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/logger"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/utils"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/transfers/internal/model"
	treegrid_model "git-codecommit.eu-central-1.amazonaws.com/v1/repos/transfers/internal/model/treegrid"
	sqlbuilder "git-codecommit.eu-central-1.amazonaws.com/v1/repos/transfers/internal/repository/sql_builder"
)

func (t *transferRepository) handleGroupBy(tg *treegrid_model.Treegrid) ([]map[string]string, error) {
	if tg.BodyParams.Rows != "" {
		logger.Debug("with rows:", tg.BodyParams.Rows)

		return t.getGroupData(tg.BodyParams.GetRowWhere(), tg)
	}

	parentBuild := sqlbuilder.QueryParent

	if tg.FilterWhere["parent"] != "" {
		logger.Debug("filters transfer", tg.FilterWhere["parent"])

		parentBuild += " " + tg.FilterWhere["parent"]
	}

	if tg.FilterWhere["child"] != "" {
		logger.Debug("filter items", tg.FilterWhere["child"])

		parentBuild = parentBuild + " AND id IN (" +
			"SELECT Parent FROM transfers_items " +
			"WHERE 1=1 " + tg.FilterWhere["child"] +
			sqlbuilder.OrderByQuery(tg.SortParams, model.TransferItemsFields) + ") "
	}

	mergedArgs := utils.MergeMaps(tg.FilterArgs["parent"], tg.FilterArgs["child"])

	logger.Debug("query for extraction", parentBuild)

	extWhereClause := utils.ExtractWhereClause(parentBuild, mergedArgs)

	logger.Debug("extracted where clause", extWhereClause)

	return t.getGroupData(extWhereClause, tg)
}

func (t *transferRepository) getGroupData(where string, tg *treegrid_model.Treegrid) ([]map[string]string, error) {
	logger.Debug("get group data")

	query, colData := t.prepareNameCountQuery(where, tg)

	logger.Debug("column", colData, "prepared query: \n", query)

	if colData.IsItem {
		return t.getChildData(tg.BodyParams.GetRowLevel(), tg.GroupCols, where, query, colData)
	}

	return t.getParentData(tg.BodyParams.GetRowLevel(), tg.GroupCols, where, query, colData)
}

// Prepare and Execute the query and return the results as JSON.
func (t *transferRepository) getJSON(sqlString string, mergedArgs []interface{}, tg *treegrid_model.Treegrid) ([]map[string]string, error) {
	stmt, err := t.db.Prepare(sqlString)
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
		}

		tableData = append(tableData, entry)
	}

	return tableData, nil
}

func (t *transferRepository) prepareNameCountQuery(where string, tg *treegrid_model.Treegrid) (query string, column model.Column) {
	// If both the level are equal then return the row

	level := tg.BodyParams.GetRowLevel()
	logger.Debug("Level", level, "len(groupCols)", len(tg.GroupCols))

	if level == len(tg.GroupCols) {

		logger.Debug("getting last level data")

		query = sqlbuilder.QueryParent
		query = strings.Replace(query, "WHERE 1=1", "", 1)
		query += where

		pos, _ := tg.BodyParams.IntPos()
		query = sqlbuilder.AddLimit(query)
		query = sqlbuilder.AddOffset(query, pos)

		return
	}

	column = model.NewColumn(tg.GroupCols[level])

	logger.Debug("getting data grouping by", tg.GroupCols[level])

	// multiple groupBy clauses
	if len(tg.GroupCols)-level > 1 {
		secColumn := model.NewColumn(tg.GroupCols[level+1])

		switch {
		case !column.IsItem && !secColumn.IsItem:
			return t.getCascadingGroupByParentParent(column, secColumn, where), column
		}
	}

	if column.IsItem {
		logger.Debug("grouping by item")

		var innerChildConditions, updatedWhere, whereCondition string

		if strings.Index(where, " AND id IN ( SELECT Parent FROM transfers_items WHERE ") > 0 {
			innerChildConditions = utils.GetStringBetween(where, " AND id IN ( SELECT Parent FROM transfers_items WHERE ", ")")
			subStr := utils.GetStringBetween(where, " AND id IN (", ")")
			updatedWhere = strings.Replace(where, subStr, "", 1)
			updatedWhere = strings.Replace(updatedWhere, "id IN ()", "1", 1)
		}

		if innerChildConditions != "" {
			innerChildConditions = " AND " + innerChildConditions
		}

		whereCondition = where
		if updatedWhere != "" {
			whereCondition = updatedWhere
		}

		if val, ok := model.ItemsFields[column.GridName]; ok {
			column.DBName = val
		}

		logger.Debug("innerchild conditions", innerChildConditions)

		query = `
		  SELECT ` + column.DBName + `, COUNT(*) Count ` +
			`   FROM transfers_items  
		  	    INNER JOIN items ON transfers_items.item_uuid = items.id
				INNER JOIN units ON transfers_items.item_unit_uuid = units.id 
				INNER JOIN item_types ON items.type_uuid = item_types.id 
		  	    WHERE Parent IN ( 
					SELECT id FROM transfers ` + whereCondition + " )" + innerChildConditions + " GROUP BY " + column.DBName

		return
	}

	query = "SELECT " + column.DBName + ", COUNT(*) Count FROM transfers " + sqlbuilder.QueryParentJoins + where + " GROUP BY " + column.DBName

	return
}

func (t *transferRepository) getCascadingGroupByParentParent(firstCol, secondCol model.Column, where string) string {
	query := `
	SELECT %s, COUNT(*) Count
	FROM (
		SELECT %s, %s 
		FROM transfers
		%s
		%s
		GROUP BY %s, %s
	) t
	GROUP BY %s
	`

	return fmt.Sprintf(query, firstCol.DBNameShort, firstCol.DBName, secondCol.DBName, sqlbuilder.QueryParentJoins, where, firstCol.DBName, secondCol.DBName, firstCol.DBNameShort)
}

func (t *transferRepository) getParentData(level int, group_cols []string, where string, query string, colData model.Column) ([]map[string]string, error) {
	logger.Debug("get parent data")

	rows, err := t.db.Query(query)
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

			tableData = append(tableData, tempObj)
			continue
		}

		tempObj["document_type"] = row.GetValue(colData.DBNameShort)

		// if grouping by document_type => DB, OPT, DT...
		// name = document_type
		// val = DB
		val := row.GetValue(colData.DBNameShort)
		where2 := " AND " + colData.DBName + "='" + val + "'"

		// Builds new attribute Rows for identification
		tempObj["Rows"] = strconv.Itoa(level+1) + where + where2

		// Query to get aggregated sum for "warehouse_destination_uuid" data for each group.
		calcQuery := "SELECT COALESCE(SUM(warehouseman_destination_approve),'') as value_sum, " +
			"COALESCE(MIN(document_date), '') AS min, " +
			"COALESCE(MAX(document_date), '') AS max " +
			// "FROM (SELECT warehouseman_destination_approve, document_date FROM transfers " + where + where2 + ") AS temp"
			"FROM (SELECT warehouseman_destination_approve, document_date FROM transfers " + sqlbuilder.QueryParentJoins + where + where2 + ") AS temp"

		calcRows, err := t.db.Query(calcQuery)
		if err != nil {
			return tableData, fmt.Errorf("do query: [%s]: [%w]", calcQuery, err)
		}

		var min, max, sum string

		for calcRows.Next() {
			err = calcRows.Scan(&min, &max, &sum)
			if err != nil {
				return tableData, fmt.Errorf("rows scan: [%w]", err)
			}
		}

		document_date := ""
		if min != "" && max != "" {
			document_date = min + "~" + max
		}

		tempObj["document_date"] = document_date
		tempObj["warehouse_destination_uuid"] = sum

		tableData = append(tableData, tempObj)

	}

	return tableData, nil
}

func (t *transferRepository) getChildData(level int, group_cols []string, where string, query string, columnData model.Column) ([]map[string]string, error) {
	logger.Debug("get child data")

	rows, err := t.db.Query(query)
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
		tempObj["document_type"] = row.GetValue(columnData.DBNameShort)

		if level == len(group_cols) {
			tableData = append(tableData, row.StringValues())
			continue
		}

		var value_sums, where2 string // sql.nullString

		// If grouped by child column
		if strings.Index(where, " AND "+columnData.DBName+"='") > 0 {
			// If the column is already added in the WHERE clause
			// Just update the query to adjust column value
			where2 = utils.ReplaceColumnValueInQuery(where, columnData.DBName, row.GetValue(columnData.DBNameShort))
		} else {
			where2 = t.addChildCondition(where, columnData.DBName, row.GetValue(columnData.DBNameShort))
		}
		tempObj["Rows"] = strconv.Itoa(level+1) + where2

		// Query to get aggregated sum for "item_quantity" data for each transfer_items group.
		calcQuery := `
SELECT COALESCE(sum(item_quantity),'') as value_sums 
FROM (
SELECT item_quantity 
FROM transfers_items ` + sqlbuilder.QueryChildJoins + ` 
WHERE Parent in (
	SELECT transfers.id FROM transfers ` + sqlbuilder.QueryParentJoins + where + `) AND ` + columnData.DBName + "='" + row.GetValue(columnData.DBNameShort) + "') AS temp;"
		calcRows, err := t.db.Query(calcQuery)
		if err != nil {
			return tableData, fmt.Errorf("do query: '%s': [%w]", calcQuery, err)
		}

		for calcRows.Next() {
			err = calcRows.Scan(&value_sums)
			if err != nil {
				return tableData, fmt.Errorf("rows scan: [%w]", err)
			}
		}
		// valueSums := calculations_rs->Get("value_sums")
		tempObj["item_quantity"] = value_sums

		tableData = append(tableData, tempObj)
	}

	return tableData, nil
}

func (t *transferRepository) addChildCondition(orignal_query string, col_name string, col_value string) string {
	if strings.Index(orignal_query, "id IN ( SELECT Parent FROM transfers_items") > 0 {
		// If already filtered on child column(s)
		orignal_query = strings.Replace(orignal_query, ")", " AND "+col_name+"='"+col_value+"' )", 1)
	} else {
		orignal_query += ` AND transfers.id IN ( SELECT Parent FROM transfers_items
		INNER JOIN items ON transfers_items.item_uuid = items.id  
		INNER JOIN units ON transfers_items.item_unit_uuid = units.id 
		INNER JOIN item_types ON items.type_uuid = item_types.id 
		WHERE ` + col_name + "='" + col_value + "' ) "
	}

	return orignal_query
}
