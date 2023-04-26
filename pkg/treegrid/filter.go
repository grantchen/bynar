package treegrid

// FilterParams stores filter's data
type FilterParams []map[string]interface{}

// ParseFilterParams maps input data to FilterParams
func ParseFilterParams(filterParams []map[string]interface{}) (f FilterParams, err error) {
	f = filterParams

	return
}

// MainFilter returns main filter params
func (f FilterParams) MainFilter() map[string]interface{} {
	if len(f) == 0 {
		return map[string]interface{}{}
	}

	return f[0]
}
