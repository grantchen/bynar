package service

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

// TransferService is the interface for transfer service
type TransferService interface {
	// GetPageCount returns the page count
	GetPageCount(tr *treegrid.Treegrid) (int64, error)
	// GetPageData returns the page data
	GetPageData(tr *treegrid.Treegrid) ([]map[string]string, error)
}
