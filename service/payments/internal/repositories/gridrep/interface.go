package gridrep

import (
	"database/sql"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/payments/internal/treegrid"
)

type GridRowReppository interface {
	IsChild(gr treegrid.GridRow) bool
	GetParentID(gr treegrid.GridRow) (parentID interface{}, err error)
	GetStatus(id interface{}) (status interface{}, err error)
	Save(tx *sql.Tx, tr *treegrid.MainRow) error
}
