package treegrid

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

// ParentWhereTag is a tag for parent where in sql
const ParentWhereTag = "1=1"

// ChildWhereTag is a tag for child where in sql
const ChildWhereTag = "2=2"

// ConnectableSQL is a sql that can be concatenated with where
type ConnectableSQL struct {
	SQL  string        // sql
	Args []interface{} // sql args

	ParentQuery string // parent query
	ChildQuery  string // child query
}

// NewConnectableSQL creates a new ConnectableSQL
func NewConnectableSQL(sql string, args ...interface{}) *ConnectableSQL {
	return &ConnectableSQL{
		SQL:         sql,
		Args:        args,
		ParentQuery: ParentWhereTag,
		ChildQuery:  ChildWhereTag,
	}
}

// Set sets sql and args
func (s *ConnectableSQL) Set(sql string, args ...interface{}) {
	s.SQL = sql
	s.Args = args
}

// Append appends sql and args
func (s *ConnectableSQL) Append(sql string, args ...interface{}) {
	s.SQL += sql
	s.Args = append(s.Args, args...)
}

// ConcatParentWhere concatenates where to parent where in sql
func (s *ConnectableSQL) ConcatParentWhere(where string, args ...interface{}) {
	s.concatWhereToTag(true, where, args...)
}

// ConcatChildWhere concatenates where to child where in sql
func (s *ConnectableSQL) ConcatChildWhere(where string, args ...interface{}) {
	s.concatWhereToTag(false, where, args...)
}

// ParentWhere returns parent where
func (s *ConnectableSQL) ParentWhere() string {
	return " WHERE " + s.removeWhereTag(s.getWhereCondition(s.ParentQuery))
}

// ChildWhere returns child where
func (s *ConnectableSQL) ChildWhere() string {
	return " WHERE " + s.removeWhereTag(s.getWhereCondition(s.ChildQuery))
}

// AsSQL returns sql string
func (s *ConnectableSQL) AsSQL() string {
	return s.whereToSQL(s.SQL, s.Args)
}

// concatWhereToTag concatenates where to sql
func (s *ConnectableSQL) concatWhereToTag(isParent bool, where string, args ...interface{}) {
	if where == "" {
		return
	}

	where = s.getWhereCondition(where)
	concatWhere := where + " AND "

	tag := ""
	if isParent {
		tag = ParentWhereTag
	} else {
		tag = ChildWhereTag
	}

	// dump sql
	sqlAfterIndex := s.SQL
	sqlIndex := 0
	for {
		i := strings.Index(sqlAfterIndex, tag)
		if i == -1 {
			// save where to s.parentQueries or s.childQueries
			if sqlIndex != 0 {
				if isParent {
					//s.parentQueries = append(s.parentQueries, where)
					s.ParentQuery = strings.ReplaceAll(s.ParentQuery, tag, s.whereToSQL(concatWhere, args)+tag)
					s.ChildQuery = strings.ReplaceAll(s.ChildQuery, tag, s.ParentQuery)
				} else {
					//s.childQueries = append(s.childQueries, where)
					s.ChildQuery = strings.ReplaceAll(s.ChildQuery, tag, s.whereToSQL(concatWhere, args)+tag)
					s.ParentQuery = strings.ReplaceAll(s.ParentQuery, tag, s.ChildQuery)
				}
			}

			return
		}

		sqlIndex += i

		// insert where into s.sql
		s.SQL = s.SQL[:sqlIndex] + concatWhere + s.SQL[sqlIndex:]

		// args count until index in s.args
		argsCountUntilIndex := strings.Count(s.SQL[:sqlIndex], "?")
		// insert args into s.args
		s.Args = append(s.Args[:argsCountUntilIndex], append(args, s.Args[argsCountUntilIndex:]...)...)

		// set sqlAfterIndex to sqlAfterIndex after tag of index i
		sqlAfterIndex = sqlAfterIndex[i+len(tag):]
		sqlIndex += len(concatWhere) + len(tag)
	}
}

// getWhereCondition returns where condition
func (s *ConnectableSQL) getWhereCondition(where string) string {
	if where == "" {
		return ""
	}

	where = strings.TrimSpace(where)
	where = strings.TrimPrefix(where, "AND")
	where = strings.TrimPrefix(where, "WHERE")
	return where
}

// whereToSQL replaces placeholders in where with args
func (s *ConnectableSQL) whereToSQL(where string, args []interface{}) string {
	if where == "" {
		return ""
	}

	if len(args) == 0 {
		return where
	}

	// replace placeholders in where with args
	for _, arg := range args {
		where = strings.Replace(where, "?", fmt.Sprintf("'%v'", arg), 1)
	}
	return where
}

// removeWhereTag removes where tag
func (s *ConnectableSQL) removeWhereTag(sql string) string {
	sql = strings.ReplaceAll(sql, "AND "+ChildWhereTag, "")
	sql = strings.ReplaceAll(sql, "AND "+ParentWhereTag, "")
	return sql
}

const (
	// paramPlaceHolder is a placeholder for param in sql
	sqlParamPlaceHolder = "?"
)

// sqlNamedSearchHandle is a regexp for named search
var sqlNamedSearchHandle = regexp.MustCompile(`{{\S+?}}`)

// NamedSQL is used for expressing complex query sql
func NamedSQL(sql string, data map[string]string) (string, error) {
	length := len(data)
	if length == 0 {
		return sql, nil
	}

	var err error
	sqlStr := sqlNamedSearchHandle.ReplaceAllStringFunc(sql, func(paramName string) string {
		paramName = strings.TrimRight(strings.TrimLeft(paramName, "{"), "}")
		val, ok := data[paramName]
		if !ok {
			err = fmt.Errorf("%s not found", paramName)
			return ""
		}

		return val
	})

	if err != nil {
		return "", err
	}
	return sqlStr, nil
}

// NamedQuery is used for expressing complex query
func NamedQuery(sql string, data map[string]interface{}) (string, []interface{}, error) {
	length := len(data)
	if length == 0 {
		return sql, nil, nil
	}

	vals := make([]interface{}, 0, length)
	var err error
	cond := sqlNamedSearchHandle.ReplaceAllStringFunc(sql, func(paramName string) string {
		paramName = strings.TrimRight(strings.TrimLeft(paramName, "{"), "}")
		val, ok := data[paramName]
		if !ok {
			err = fmt.Errorf("%s not found", paramName)
			return ""
		}

		v := reflect.ValueOf(val)
		if v.Type().Kind() != reflect.Slice {
			vals = append(vals, val)
			return sqlParamPlaceHolder
		}

		length := v.Len()
		for i := 0; i < length; i++ {
			vals = append(vals, v.Index(i).Interface())
		}

		return createMultiPlaceholders(length)
	})

	if err != nil {
		return "", nil, err
	}
	return cond, vals, nil
}

// createMultiPlaceholders creates multi placeholders
func createMultiPlaceholders(num int) string {
	if 0 == num {
		return ""
	}
	length := (num << 1) | 1
	buff := make([]byte, length)
	buff[0], buff[length-1] = '(', ')'
	ll := length - 2
	for i := 1; i <= ll; i += 2 {
		buff[i] = '?'
	}
	ll = length - 3
	for i := 2; i <= ll; i += 2 {
		buff[i] = ','
	}
	return string(buff)
}
