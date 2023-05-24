package repository

import "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"

type UserGroupsRepository interface {
	GetPageCount(tg *treegrid.Treegrid) int64
	GetPageData(tg *treegrid.Treegrid) ([]map[string]string, error)
}
