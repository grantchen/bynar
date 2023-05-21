package treegrid

import (
	"strconv"
	"strings"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/logger"
)

func PrepFilters(tr *Treegrid, fieldAliasesParent map[string][]string, fieldAliasesChild map[string][]string) {
	tr.FilterWhere, tr.FilterArgs = PrepQueryComplex(tr.FilterParams, fieldAliasesParent, fieldAliasesChild)
}

// with parent and child
func PrepQueryComplex(f FilterParams, fieldAliasesParent map[string][]string, fieldAliasesChild map[string][]string) (map[string]string, map[string][]interface{}) {
	FilterWhere := map[string]string{}
	FilterArgs := map[string][]interface{}{}

	// Filter process
	var curField, curFieldValue, curOperation, curMarker string

	for key, el := range f.Filters() {

		if key == "id" || strings.Contains(key, "Filter") {
			continue
		}

		if list, ok := fieldAliasesParent[key]; ok && list[0] != "" {
			curMarker = "parent"
			curField = list[0]
		} else if list, ok := fieldAliasesChild[key]; ok && list[0] != "" {
			// cur column is a parent column
			curMarker = "child"
			curField = list[0]
		} else {
			break
		}

		curOperation = f.Filters()[key+"Filter"].(string)
		curFieldValue = el.(string)

		operatorInt, _ := strconv.Atoi(curOperation)

		if curField == "" {
			logger.Debug("undefined filter column: ", key)

			continue
		}

		if operatorInt < 7 {
			filterWhere, filterArgs := PrepareFilters(operatorInt, curField, curFieldValue)
			FilterWhere[curMarker] += " AND " + filterWhere
			FilterArgs[curMarker] = append(FilterArgs[curMarker], filterArgs...)
			continue
		}

		// LIKE operators for text search
		filterWhere, filterArgs := PrepareTextFilters(operatorInt, curField, curFieldValue)
		FilterWhere[curMarker] += " AND " + filterWhere
		FilterArgs[curMarker] = append(FilterArgs[curMarker], filterArgs...)
	}

	return FilterWhere, FilterArgs
}

// only without child
func PrepQuerySimple(f FilterParams, fieldAliases map[string][]string) (string, []interface{}) {
	dummyMap := make(map[string][]string, 0)
	FilterWhere, FilterArgs := PrepQueryComplex(f, fieldAliases, dummyMap)
	return FilterWhere["parent"], FilterArgs["parent"]
}

func PrepareTextFilters(operator int, colName, val string) (whereSql string, args []interface{}) {
	var valPrefix, valSuffix, sqlOperator, sqlOperatorJoin string

	switch operator {

	// 7-12 text operators
	case 7:
		valSuffix = "%"
		sqlOperator = " LIKE "
		sqlOperatorJoin = " OR "
	case 8:
		valSuffix = "%"
		sqlOperator = " NOT LIKE "
		sqlOperatorJoin = " AND "
	case 9:
		valPrefix = "%"
		sqlOperator = " LIKE "
		sqlOperatorJoin = " OR "
	case 10:
		valPrefix = "%"
		sqlOperator = " NOT LIKE "
		sqlOperatorJoin = " AND "
	case 11:
		valPrefix = "%"
		valSuffix = "%"
		sqlOperator = " LIKE "
		sqlOperatorJoin = " OR "
	case 12:
		valPrefix = "%"
		valSuffix = "%"
		sqlOperator = " NOT LIKE "
		sqlOperatorJoin = " AND "
	}

	sqlParts := strings.Split(val, ";")
	lenSqlParts := len(sqlParts)
	args = make([]interface{}, 0, lenSqlParts)

	for k := range sqlParts {
		whereSql += colName + sqlOperator + " ? "
		if k != (lenSqlParts - 1) {
			whereSql += sqlOperatorJoin
		}

		args = append(args, valPrefix+sqlParts[k]+valSuffix)
	}

	return "(" + whereSql + ")", args
}

func PrepareFilters(operator int, colName, val string) (whereSql string, args []interface{}) {
	var sqlOperator, sqlOperatorJoin, rangeOperator string

	logger.Debug("operator", operator, "colName", colName, "val", val)

	switch operator {
	case 1:
		sqlOperator = " = "
		sqlOperatorJoin = " OR "
		rangeOperator = " BETWEEN "
	case 2:
		sqlOperator = " != "
		sqlOperatorJoin = " AND "
		rangeOperator = " NOT BETWEEN "
	case 3:
		whereSql = " (" + colName + " < ? )"
		args = append(args, val)

		return
	case 4:
		whereSql = " (" + colName + " <= ? )"
		args = append(args, val)

		return
	case 5:
		whereSql = " (" + colName + " > ? )"
		args = append(args, val)

		return
	case 6:
		whereSql = " (" + colName + " >= ? )"
		args = append(args, val)

		return
	}

	sqlParts := strings.Split(val, ";")
	lenSqlParts := len(sqlParts)
	args = make([]interface{}, 0, lenSqlParts)

	for k, val := range sqlParts {
		rangeVals := strings.Split(val, "~")

		// value without range condition "1;2;3;6"
		if len(rangeVals) == 1 {
			if !IsDateCol(colName) {
				whereSql += colName + sqlOperator + " ? "
			} else {
				whereSql += colName + sqlOperator + " STR_TO_DATE(?,'%m/%d/%Y') "
			}

			if k != (lenSqlParts - 1) {
				whereSql += sqlOperatorJoin
			}

			args = append(args, sqlParts[k])

			continue
		}

		// condition with range "1~4"
		if len(rangeVals) != 2 {
			logger.Debug("invalid range condition", val)

			continue
		}
		if !IsDateCol(colName) {
			whereSql += " ( " + colName + rangeOperator + " ? AND ? )"
		} else {
			whereSql += " ( " + colName + rangeOperator + " STR_TO_DATE(?,'%m/%d/%Y') AND STR_TO_DATE(?,'%m/%d/%Y') )"
		}

		if k != (lenSqlParts - 1) {
			whereSql += sqlOperatorJoin
		}

		args = append(args, rangeVals[0], rangeVals[1])
	}

	return " (" + whereSql + ") ", args
}

func IsDateCol(colName string) bool {
	return strings.Contains(colName, "_date")
}
