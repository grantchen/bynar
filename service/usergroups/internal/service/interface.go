package service

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

// UserGroupService is the interface for user group service
type UserGroupService interface {
	// GetPageCount returns the page count
	GetPageCount(tr *treegrid.Treegrid) (int64, error)
	// GetPageData returns the page data
	GetPageData(tr *treegrid.Treegrid) ([]map[string]string, error)
	// GetCellSuggestion returns the cell suggestion
	GetCellSuggestion(tr *treegrid.Treegrid) (*treegrid.PostResponse, error)
}
