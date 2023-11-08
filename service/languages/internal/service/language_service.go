package service

import (
	"database/sql"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

type languageService struct {
	gridRowDataRepositoryWithChild treegrid.GridRowDataRepositoryWithChild
	db                             *sql.DB
}

// GetCellSuggestion implements languageservice
func (u *languageService) GetCellSuggestion(tr *treegrid.Treegrid) (*treegrid.PostResponse, error) {
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

// GetTransferCount implements languagesService
func (u *languageService) GetPageCount(tr *treegrid.Treegrid) (int64, error) {
	return u.gridRowDataRepositoryWithChild.GetPageCount(tr)
}

// GetTransfersPageData implements languagesService
func (u *languageService) GetPageData(tr *treegrid.Treegrid) ([]map[string]string, error) {
	return u.gridRowDataRepositoryWithChild.GetPageData(tr)
}

func Newlanguageservice(
	db *sql.DB,
	gridRowDataRepositoryWithChild treegrid.GridRowDataRepositoryWithChild,
) LanguageService {
	return &languageService{
		db:                             db,
		gridRowDataRepositoryWithChild: gridRowDataRepositoryWithChild,
	}
}
