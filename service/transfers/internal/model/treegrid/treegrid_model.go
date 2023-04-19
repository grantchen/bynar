package treegrid_model

type Treegrid struct {
	RawSortCols  string
	RawSortTypes string
	RawGroupBy   string
	BodyParams   BodyParam

	SortParams
	GroupCols
	FilterParams
	FilterWhere map[string]string
	FilterArgs  map[string][]interface{}
}
