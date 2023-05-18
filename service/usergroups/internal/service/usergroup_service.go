package service

import (
	"database/sql"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

type userGroupService struct {
	gridRowDataRepositoryWithChild treegrid.GridRowDataRepositoryWithChild
	db                             *sql.DB
}

// GetTransferCount implements UserGroupsService
func (u *userGroupService) GetPageCount(tr *treegrid.Treegrid) int64 {
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
