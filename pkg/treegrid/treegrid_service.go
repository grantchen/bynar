package treegrid

import (
	"context"
	"database/sql"
)

type TreeGridService interface {
	GetPageCount(tr *Treegrid) float64
	GetPageData(tr *Treegrid) ([]map[string]string, error)
	Upload(req *PostRequest) (*PostResponse, error)
	GetCellData(ctx context.Context, req *Treegrid) (*PostResponse, error)
}

type TreeGridServiceFactoryFunc func(db *sql.DB, AccountID int, permissionInfo *PermissionInfo) TreeGridService
