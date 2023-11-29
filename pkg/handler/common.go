package handler

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

type CallBackGetPageCount func(tr *treegrid.Treegrid) (float64, error)
type CallBackGetPageData func(tr *treegrid.Treegrid) ([]map[string]string, error)
type CallBackUploadData func(req *treegrid.PostRequest) (*treegrid.PostResponse, error)
type CallBackGetCellData func(req *treegrid.Treegrid) (*treegrid.PostResponse, error)
