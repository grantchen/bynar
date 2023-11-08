package service

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

type LanguageService interface {
	GetPageCount(tr *treegrid.Treegrid) (int64, error)
	GetPageData(tr *treegrid.Treegrid) ([]map[string]string, error)
	GetCellSuggestion(tr *treegrid.Treegrid) (*treegrid.PostResponse, error)
}
