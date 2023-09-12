package handler

import (
	"context"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/aws/scope"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

type CallBackGetPageCount func(ctx context.Context, tr *treegrid.Treegrid) (float64, error)
type CallBackGetPageData func(ctx context.Context, tr *treegrid.Treegrid) ([]map[string]string, error)
type CallBackUploadData func(ctx context.Context, req *treegrid.PostRequest) (*treegrid.PostResponse, error)
type CallBackGetCellData func(ctx context.Context, req *treegrid.Treegrid) (*treegrid.PostResponse, error)

type CallBackLambdaUploadData func(scope *scope.RequestScope, req *treegrid.PostRequest) (*treegrid.PostResponse, error)
