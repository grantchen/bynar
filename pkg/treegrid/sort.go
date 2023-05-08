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
type SortParams map[string]SortType

// ParseSortParams converts sortValsStr="col1,col2", sortTypesStr= "0,1" to SortParams
func ParseSortParams(sortValsStr string, sortTypesStr string) (SortParams, error) {
	params := make(map[string]SortType)
	sortVals := strings.Split(sortValsStr, ",")
	sortTypes := strings.Split(sortTypesStr, ",")

	if sortVals[0] == "" {
		return params, nil
	}

	if len(sortTypes) != len(sortVals) {
		return nil, errors.New("invalid sort params, " + sortValsStr + ", " + sortTypesStr)
	}

	for i := 0; i < len(sortVals); i++ {
		params[sortVals[i]] = SortType(sortTypes[i])
	}

	return params, nil
}

// OrderByChildQuery - making 'ORDER BY' query from SortParams and fieldsMapping
func (s SortParams) OrderByChildQuery(itemFields map[string]bool) (res string) {
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

// OrderByQueryExludeChild - making 'ORDER BY' query from SortParams and fieldsMapping EXCLUDING child sort params
func (s SortParams) OrderByQueryExludeChild(childFields map[string]bool, fieldAlias map[string]string) (res string) {
	for k, v := range s {
		if childFields == nil || childFields[k] {
			continue
		}

		if f, ok := fieldAlias[k]; ok {
			// res += fmt.Sprintf("%s %s, ", k, v)

			res += fmt.Sprintf("%s %s, ", f, v)
		}
	}

	if len(res) > 0 {
		res = " ORDER BY " + res[:len(res)-2]
	}

	return
}
