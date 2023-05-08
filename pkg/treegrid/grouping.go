package treegrid

import "strings"

// GroupCols - columns for GROUP BY
type GroupCols []string

func ParseGroupCols(params string) (grCols GroupCols) {
	grCols = strings.Split(params, ",")

	return
}
