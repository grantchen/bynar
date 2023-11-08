package service

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

// SaleService is an interface for sale service
type SaleService interface {
	// GetPageCount returns page count
	GetPageCount(tr *treegrid.Treegrid) (int64, error)
	// GetPageData returns page data
	GetPageData(tr *treegrid.Treegrid) ([]map[string]string, error)
	// GetCellSuggestion returns cell suggestion
	GetCellSuggestion(tr *treegrid.Treegrid) (*treegrid.PostResponse, error)
}
