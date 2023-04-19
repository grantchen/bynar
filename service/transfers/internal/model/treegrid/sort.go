package treegrid_model

import (
	"errors"
	"strings"
)

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
