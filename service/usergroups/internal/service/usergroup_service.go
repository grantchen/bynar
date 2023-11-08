package service

import (
	"database/sql"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

// UserGroupService implements UserGroupService
type userGroupService struct {
	gridRowDataRepositoryWithChild treegrid.GridRowDataRepositoryWithChild
	db                             *sql.DB
}

// GetCellSuggestion implements UserGroupsService to get cell suggestion
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
	resp.Changes = append(resp.Changes, treegrid.CreateSuggestionResult(tr.BodyParams.Col, suggestion, tr))
	return resp, nil
}

// GetPageCount implements UserGroupsService to get page count
func (u *userGroupService) GetPageCount(tr *treegrid.Treegrid) (int64, error) {
	return u.gridRowDataRepositoryWithChild.GetPageCount(tr)
}

// GetPageData implements UserGroupsService to get page data
func (u *userGroupService) GetPageData(tr *treegrid.Treegrid) ([]map[string]string, error) {
	return u.gridRowDataRepositoryWithChild.GetPageData(tr)
}

// NewUserGroupService returns a new UserGroupService
func NewUserGroupService(
	db *sql.DB,
	gridRowDataRepositoryWithChild treegrid.GridRowDataRepositoryWithChild,
) UserGroupService {
	return &userGroupService{
		db:                             db,
		gridRowDataRepositoryWithChild: gridRowDataRepositoryWithChild,
	}
}
