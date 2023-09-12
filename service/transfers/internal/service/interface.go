package service

import (
	"context"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

type TransferService interface {
	GetPagesCount(ctx context.Context, tr *treegrid.Treegrid) (float64, error)
	GetTransfersPageData(ctx context.Context, tr *treegrid.Treegrid) ([]map[string]string, error)
	HandleUpload(ctx context.Context, req *treegrid.PostRequest, accountID int) (*treegrid.PostResponse, error)
}
