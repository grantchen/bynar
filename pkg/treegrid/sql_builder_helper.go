package treegrid

import "strings"

// WhereTag is a tag for where in sql
const WhereTag = "1=1"

// ConcatWhereToSQL concatenates where to sql
func ConcatWhereToSQL(sql, where string) string {
	if where != "" {
		where = strings.TrimSpace(where)
		where = strings.TrimPrefix(where, "AND")
		where = strings.TrimPrefix(where, "WHERE")
		sql = strings.ReplaceAll(sql, WhereTag, where+" AND "+WhereTag)
	}
	return sql

}
