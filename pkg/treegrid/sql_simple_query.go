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

func BuildSimpleQueryCount(tableName string, fieldMapping map[string][]string) string {
	query := `select COUNT(*) FROM ` + tableName
	return query
}

func BuildSimpleQueryGroupByCount(tableName string, fieldMapping map[string][]string, groupCol []string) string {
	groupBy := make([]string, 0)
	for _, field := range groupCol {
		dbCol := fieldMapping[field][0]
		groupBy = append(groupBy, dbCol)
	}
	query := `select COUNT(*) FROM ` + tableName + " group by " + strings.Join(groupBy[:], ",")
	return query
}

func BuildSimpleQueryCountWithGroupping(tableName string, fieldMapping map[string][]string, groupField string) string {
	return ""
}

func BuildSimpleQuery(tableName string, fieldMapping map[string][]string) string {
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

func BuildSimpleQueryGroupBy(tableName string, fieldMapping map[string][]string, groupCol []string, whereCondition string) string {
	var queryBuffer bytes.Buffer
	queryBuffer.WriteString(`select `)
	groupBy := make([]string, 0)
	sel := make([]string, 0)
	for _, field := range groupCol {
		dbCol := fieldMapping[field][0]
		sel = append(sel, fmt.Sprintf("%s AS %s", dbCol, field))
		groupBy = append(groupBy, dbCol)
	}

	queryBuffer.WriteString(strings.Join(sel[:], ","))
	queryBuffer.WriteString(", COUNT(*) AS Count FROM " + tableName)
	queryBuffer.WriteString(DummyWhere + whereCondition)
	queryBuffer.WriteString(" GROUP BY " + strings.Join(groupBy[:], ","))

	return queryBuffer.String()
}

func AppendLimitToQuery(query string, pagesize int, pos int) string {
	return query + " LIMIT " + strconv.Itoa(pagesize) + " OFFSET " + strconv.Itoa(pos*pagesize)
}

func AppendWhereCondition() {

}
