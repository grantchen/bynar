package treegrid

import (
	"bytes"
	"fmt"
	"strconv"
)

var (
	// Dummy where clause used before AND clauses
	DummyWhere = " where 1=1 "

	// WHERE Parent = ''

)

func BuildSimpleQueryCount(tableName string, fieldMapping map[string][]string, defaultQueryCount string) string {
	if defaultQueryCount != "" {
		return defaultQueryCount
	}
	query := `select COUNT(*) FROM ` + tableName
	return query
}

// BuildSimpleQueryGroupByCount build query group by, note that WHERE key word is NOT added in here
func BuildSimpleQueryGroupByCount(tableName string, fieldMapping map[string][]string, groupCol []string, where string) string {
	groupBy := make([]string, 0)
	for _, field := range groupCol {
		dbCol := fieldMapping[field][0]
		groupBy = append(groupBy, dbCol)
	}

	return fmt.Sprintf(`SELECT COUNT(*) FROM (SELECT %s FROM %s%s GROUP BY %s) t;`, groupBy[0], tableName, where, groupBy[0])
}

func BuildSimpleQuery(tableName string, fieldMapping map[string][]string, defaultQuery string) string {
	if defaultQuery != "" {
		return defaultQuery
	}
	var queryBuffer bytes.Buffer
	queryBuffer.WriteString(`select `)
	totalField := len(fieldMapping)
	idx := 0

	for tgCol, dbCol := range fieldMapping {
		queryBuffer.WriteString(fmt.Sprintf("%s AS %s", dbCol[0], tgCol))
		// fmt.Println(idx)
		if idx < totalField-1 {
			queryBuffer.WriteString(",")
		}
		idx += 1
	}
	queryBuffer.WriteString(fmt.Sprintf(" FROM %s", tableName))
	return queryBuffer.String()
}

// BuildSimpleQueryGroupBy build query group by, note that WHERE key word is added in here
func BuildSimpleQueryGroupBy(tableName string, fieldMapping map[string][]string, groupCol []string, whereCondition string, level int, innerJoin string) string {
	var queryBuffer bytes.Buffer
	queryBuffer.WriteString(`select `)
	groupBy := make([]string, 0)
	sel := make([]string, 0)
	for _, field := range groupCol {
		dbCol := fieldMapping[field][0]
		sel = append(sel, fmt.Sprintf("%s AS %s", dbCol, field))
		groupBy = append(groupBy, dbCol)
	}

	queryFrom := ""
	if level == len(groupCol)-1 {
		// get the last level parent data with children count
		queryFrom = tableName
		queryBuffer.WriteString(sel[level])
		queryBuffer.WriteString(", COUNT(*) AS Count FROM " + queryFrom)
		if innerJoin != "" {
			queryBuffer.WriteString(innerJoin)
		}
		queryBuffer.WriteString(whereCondition)
		queryBuffer.WriteString(" GROUP BY " + groupBy[level])
	} else {
		// get the parent data with children count not the last level
		queryFrom = fmt.Sprintf(`
			(SELECT %s, %s
			FROM %s
			%s
			%s
			GROUP BY %s, %s) t
			`,
			groupBy[level],
			groupBy[level+1],
			tableName,
			innerJoin,
			whereCondition,
			groupBy[level],
			groupBy[level+1],
		)

		queryBuffer.WriteString(sel[level])
		queryBuffer.WriteString(", COUNT(*) AS Count FROM " + queryFrom)
		queryBuffer.WriteString(" GROUP BY " + groupBy[level])
	}

	return queryBuffer.String()
}

func AppendLimitToQuery(query string, pagesize int, pos int) string {
	return query + " LIMIT " + strconv.Itoa(pagesize) + " OFFSET " + strconv.Itoa(pos*pagesize)
}
