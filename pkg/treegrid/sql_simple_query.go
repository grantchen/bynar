package treegrid

import (
	"bytes"
	"fmt"
)

func buildSimpleQueryCount(tableName string, fieldMapping map[string][]string) string {
	query := `select COUNT(*) FROM ` + tableName
	return query
}

func buildSimpleQueryCountWithGroupping(tableName string, fieldMapping map[string][]string, groupField string) string {
	return ""
}

func buildSimpleQuery(tableName string, fieldMapping map[string][]string) string {
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

func appendWhereCondition() {

}
