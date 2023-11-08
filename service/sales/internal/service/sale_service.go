package service

import (
	"database/sql"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

// SaleService is implementation of SaleService
type saleService struct {
	gridRowDataRepositoryWithChild treegrid.GridRowDataRepositoryWithChild
	db                             *sql.DB
}

// GetCellSuggestion implements SaleService
func (u *saleService) GetCellSuggestion(tr *treegrid.Treegrid) (*treegrid.PostResponse, error) {
	data, err := u.gridRowDataRepositoryWithChild.GetChildSuggestion(tr)

	resp := &treegrid.PostResponse{}

	if err != nil {
		resp.IO.Result = -1
		resp.IO.Message += err.Error() + "\n"
		return resp, err
	}

	suggestion := &treegrid.Suggestion{
		Items: data,
	}
	resp.Changes = append(resp.Changes, treegrid.CreateSuggestionResult(tr.BodyParams.Col, suggestion, tr))
	return resp, nil
}

// GetPageCount implements SalesService
func (u *saleService) GetPageCount(tr *treegrid.Treegrid) (int64, error) {
	return u.gridRowDataRepositoryWithChild.GetPageCount(tr)
}

// GetPageData implements SalesService
func (u *saleService) GetPageData(tr *treegrid.Treegrid) ([]map[string]string, error) {
	return u.gridRowDataRepositoryWithChild.GetPageData(tr)
}

// NewSaleService returns new instance of SaleService
func NewSaleService(
	db *sql.DB,
	gridRowDataRepositoryWithChild treegrid.GridRowDataRepositoryWithChild,
) SaleService {
	return &saleService{
		db:                             db,
		gridRowDataRepositoryWithChild: gridRowDataRepositoryWithChild,
	}
}
