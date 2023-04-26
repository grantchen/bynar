package sqlbuilder

import (
	"fmt"
	"strconv"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

var pageSize = 100

func PageSizeString() string {
	return strconv.Itoa(pageSize)
}

func AddLimit(query string) string {
	return query + " LIMIT " + PageSizeString()
}

func AddOffset(query string, pos int) string {
	if pos == 0 {
		return query
	}

	return query + " OFFSET " + strconv.Itoa(pos*pageSize)
}

// OrderByQuery - making 'ORDER BY' query from SortParams and fieldsMapping
func OrderByQuery(s treegrid.SortParams, itemFields map[string]bool) (res string) {
	for k, v := range s {
		if itemFields == nil {
			res += fmt.Sprintf("%s %s, ", k, v)
			continue
		}

		if itemFields[k] {
			res += fmt.Sprintf("%s %s, ", k, v)
		}
	}

	if len(res) > 0 {
		res = " ORDER BY " + res[:len(res)-2]
	}

	return
}
