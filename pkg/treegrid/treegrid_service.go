package treegrid

import (
	"context"
	"database/sql"
)

type Service interface {
	GetPageCount(tr *Treegrid) (float64, error)
	GetPageData(tr *Treegrid) ([]map[string]string, error)
	Upload(req *PostRequest) (*PostResponse, error)
	GetCellData(ctx context.Context, req *Treegrid) (*PostResponse, error)
}

type ServiceFactoryFunc func(db *sql.DB, accountID int, organizationUuid string, permissionInfo *PermissionInfo, language string) Service
