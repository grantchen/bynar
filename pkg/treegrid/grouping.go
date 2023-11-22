package treegrid

import "strings"

// GroupCols - columns for GROUP BY
type GroupCols []string

// ParseGroupCols parses group columns
func ParseGroupCols(params string) (grCols GroupCols) {
	grCols = strings.Split(params, ",")

	return
}

// ContainsAny checks if GroupCols contains any of fieldsMapping keys
func (gc GroupCols) ContainsAny(fieldsMapping map[string][]string) bool {
	for _, col := range gc {
		if _, ok := fieldsMapping[col]; ok {
			return true
		}
	}

	return false
}
