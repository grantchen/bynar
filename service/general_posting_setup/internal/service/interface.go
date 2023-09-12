package service

import (
	"context"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

type GeneralPostingSetupService interface {
	GetPageCount(ctx context.Context, treegrid *treegrid.Treegrid) (int64, error)
	GetPageData(ctx context.Context, tr *treegrid.Treegrid) ([]map[string]string, error)
}

type UploadService interface {
	Handle(ctx context.Context, req *treegrid.PostRequest) (*treegrid.PostResponse, error)
}
