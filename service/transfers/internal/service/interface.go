package service

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

type TransferService interface {
	GetPageCount(tr *treegrid.Treegrid) (int64, error)
	GetPageData(tr *treegrid.Treegrid) ([]map[string]string, error)
}
