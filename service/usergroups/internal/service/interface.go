package service

import "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"

type UserGroupService interface {
	GetPageCount(treegrid *treegrid.Treegrid) int64
	GetPageData(tr *treegrid.Treegrid) ([]map[string]string, error)
}
