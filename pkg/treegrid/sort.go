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
func (s SortParams) OrderByChildQuery(childFieldMapping map[string][]string) (res string) {
	if childFieldMapping == nil {
		return ""
	}

	for k, v := range s {
		if f, ok := childFieldMapping[k]; ok {
			res += fmt.Sprintf("%s %s, ", f[0], v)
		}
	}

	if len(res) > 0 {
		res = " ORDER BY " + res[:len(res)-2]
	}

	return
}

// OrderByQueryExcludeChild - making 'ORDER BY' query from SortParams and fieldsMapping EXCLUDING child sort params
func (s SortParams) OrderByQueryExcludeChild(childFieldMapping map[string][]string, parentFieldMapping map[string][]string) (res string) {
	for k, v := range s {
		if childFieldMapping == nil || len(childFieldMapping[k]) > 0 {
			continue
		}

		if f, ok := parentFieldMapping[k]; ok {
			res += fmt.Sprintf("%s %s, ", f[0], v)
		}
	}

	if len(res) > 0 {
		res = " ORDER BY " + res[:len(res)-2]
	}

	return
}
