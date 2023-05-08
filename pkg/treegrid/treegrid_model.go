package treegrid

import "fmt"

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

func NewTreegrid(req *Request) (*Treegrid, error) {
	sortParams, err := ParseSortParams(req.Cfg.SortCols, req.Cfg.SortTypes)
	if err != nil {
		return nil, fmt.Errorf("parse sort params: [%w]", err)
	}

	grCols := ParseGroupCols(req.Cfg.GroupCols)

	filterParams, err := ParseFilterParams(req.Filters)
	if err != nil {
		return nil, fmt.Errorf("parse filter params: [%w]", err)
	}

	tr := &Treegrid{
		RawSortCols:  req.Cfg.SortCols,
		RawSortTypes: req.Cfg.SortTypes,
		RawGroupBy:   req.Cfg.GroupCols,
		SortParams:   sortParams,
		GroupCols:    grCols,
		FilterParams: filterParams,
		FilterWhere:  map[string]string{},
		FilterArgs:   map[string][]interface{}{},
	}

	if len(req.Body) > 0 {
		tr.BodyParams = req.Body[0]
	}

	return tr, nil
}

func (tr *Treegrid) WithGroupBy() bool {
	return len(tr.RawGroupBy) != 0
}

// Refresh sets new filter params
func (tr *Treegrid) SetFilters(filterParams []map[string]interface{}) {
	tr.FilterParams = filterParams
}
