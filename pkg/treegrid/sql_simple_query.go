package treegrid

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
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
func BuildSimpleQueryGroupByCount(tableName string, fieldMapping map[string][]string, groupCol []string, where string, defaultQueryCount string) string {
	groupBy := make([]string, 0)
	for _, field := range groupCol {
		dbCol := fieldMapping[field][0]
		groupBy = append(groupBy, dbCol)
	}
	var query string
	if defaultQueryCount == "" {
		query = `select COUNT(*) FROM ` + where + tableName + " group by " + strings.Join(groupBy[:], ",")
	} else {
		query = defaultQueryCount + where + " group by " + strings.Join(groupBy[:], ",")
	}
	return query
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

	queryBuffer.WriteString(sel[level])
	queryBuffer.WriteString(", COUNT(*) AS Count FROM " + tableName)
	if innerJoin != "" {
		queryBuffer.WriteString(innerJoin)
	}
	queryBuffer.WriteString(whereCondition)
	queryBuffer.WriteString(" GROUP BY " + groupBy[level])

	return queryBuffer.String()
}

func AppendLimitToQuery(query string, pagesize int, pos int) string {
	return query + " LIMIT " + strconv.Itoa(pagesize) + " OFFSET " + strconv.Itoa(pos*pagesize)
}
