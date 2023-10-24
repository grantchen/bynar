package service

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

type TransferService interface {
	GetPageCount(tr *treegrid.Treegrid) (float64, error)
	GetPageData(tr *treegrid.Treegrid) ([]map[string]string, error)
	HandleUpload(req *treegrid.PostRequest, accountID int) (*treegrid.PostResponse, error)
}
