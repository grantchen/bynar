package treegrid

import "strings"

// ParentWhereTag is a tag for parent where in sql
const ParentWhereTag = "1=1"

// ChildWhereTag is a tag for child where in sql
const ChildWhereTag = "2=2"

// ConcatParentWhereToSQL concatenates where to sql
func ConcatParentWhereToSQL(sql, where string) string {
	if where != "" {
		where = strings.TrimSpace(where)
		where = strings.TrimPrefix(where, "AND")
		where = strings.TrimPrefix(where, "WHERE")
		sql = strings.ReplaceAll(sql, ParentWhereTag, where+" AND "+ParentWhereTag)
	}
	return sql
}

// ConcatChildWhereToSQL concatenates where to sql
func ConcatChildWhereToSQL(sql, where string) string {
	if where != "" {
		where = strings.TrimSpace(where)
		where = strings.TrimPrefix(where, "AND")
		where = strings.TrimPrefix(where, "WHERE")
		sql = strings.ReplaceAll(sql, ChildWhereTag, where+" AND "+ChildWhereTag)
	}
	return sql
}
