package service

import (
	"database/sql"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/logger"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

type userGroupService struct {
	gridRowDataRepositoryWithChild treegrid.GridRowDataRepositoryWithChild
	db                             *sql.DB
}

// GetCellSuggestion implements UserGroupService
func (u *userGroupService) GetCellSuggestion(tr *treegrid.Treegrid) (*treegrid.PostResponse, error) {
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
	logger.Debug("data: ", suggestion)
	resp.Changes = append(resp.Changes, treegrid.CreateSuggestionResult(suggestion, tr))
	return resp, nil
}

// GetTransferCount implements UserGroupsService
func (u *userGroupService) GetPageCount(tr *treegrid.Treegrid) (int64, error) {
	return u.gridRowDataRepositoryWithChild.GetPageCount(tr)
}

// GetTransfersPageData implements UserGroupsService
func (u *userGroupService) GetPageData(tr *treegrid.Treegrid) ([]map[string]string, error) {
	return u.gridRowDataRepositoryWithChild.GetPageData(tr)
}

func NewUserGroupService(
	db *sql.DB,
	gridRowDataRepositoryWithChild treegrid.GridRowDataRepositoryWithChild,
) UserGroupService {
	return &userGroupService{
		db:                             db,
		gridRowDataRepositoryWithChild: gridRowDataRepositoryWithChild,
	}
}
