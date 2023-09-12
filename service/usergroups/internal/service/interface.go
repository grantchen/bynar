package service

import (
	"context"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

type UserGroupService interface {
	GetPageCount(ctx context.Context, treegrid *treegrid.Treegrid) (int64, error)
	GetPageData(ctx context.Context, tr *treegrid.Treegrid) ([]map[string]string, error)
	GetCellSuggestion(ctx context.Context, tr *treegrid.Treegrid) (*treegrid.PostResponse, error)
}
