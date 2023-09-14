package service

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

type GeneralPostingSetupService interface {
	GetPageCount(treegrid *treegrid.Treegrid) (int64, error)
	GetPageData(tr *treegrid.Treegrid) ([]map[string]string, error)
}

type UploadService interface {
	Handle(req *treegrid.PostRequest) (*treegrid.PostResponse, error)
}
