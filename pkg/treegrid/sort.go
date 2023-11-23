package treegrid

import (
	"errors"
	"fmt"
	"strings"
)

// SortType - type of sorting for ORDER BY
type SortType string

const (
	// sort by ASC
	SortASC SortType = "0"
	// sort by DESC
	SortDESC SortType = "1"
)

// String - stringer interface implementation
func (s SortType) String() string {
	if s == SortDESC {
		return "DESC"
	}

	return "ASC"
}

// SortParams - storage for managing column's sorting
type SortParams []struct {
	Col  string   // column name
	Type SortType // sorting type
}

// ParseSortParams converts sortValsStr="col1,col2", sortTypesStr= "0,1" to SortParams
func ParseSortParams(sortValsStr string, sortTypesStr string) (SortParams, error) {
	sortVals := strings.Split(sortValsStr, ",")
	sortTypes := strings.Split(sortTypesStr, ",")

	if sortVals[0] == "" {
		return nil, nil
	}

	if len(sortTypes) != len(sortVals) {
		return nil, errors.New("invalid sort params, " + sortValsStr + ", " + sortTypesStr)
	}

	params := make(SortParams, len(sortVals))
	for i := 0; i < len(sortVals); i++ {
		params[i].Col = sortVals[i]
		params[i].Type = SortType(sortTypes[i])
	}

	return params, nil
}

// OrderByChildQuery - making 'ORDER BY' query from SortParams and fieldsMapping
func (s SortParams) OrderByChildQuery(childFieldMapping map[string][]string) (res string) {
	if childFieldMapping == nil {
		return ""
	}

	for _, sort := range s {
		if f, ok := childFieldMapping[sort.Col]; ok {
			res += fmt.Sprintf("%s %s, ", f[0], sort.Type)
		}
	}

	if len(res) > 0 {
		res = " ORDER BY " + res[:len(res)-2]
	}

	return
}

// OrderByQueryExcludeChild - making 'ORDER BY' query from SortParams and fieldsMapping EXCLUDING child sort params
func (s SortParams) OrderByQueryExcludeChild(childFieldMapping map[string][]string, parentFieldMapping map[string][]string) (res string) {
	for _, sort := range s {
		if childFieldMapping == nil || len(childFieldMapping[sort.Col]) > 0 {
			continue
		}

		if f, ok := parentFieldMapping[sort.Col]; ok {
			res += fmt.Sprintf("%s %s, ", f[0], sort.Type)
		}
	}

	if len(res) > 0 {
		res = " ORDER BY " + res[:len(res)-2]
	}

	return
}
