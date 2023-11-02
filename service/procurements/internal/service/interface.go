package service

import (
	"database/sql"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/models"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

type ProcurementsService interface {
	GetTx(tx *sql.Tx, id interface{}) (*models.Procurement, error)
	Handle(tx *sql.Tx, pr *models.Procurement) error
	HandleLine(tx *sql.Tx, pr *models.Procurement, line *models.ProcurementLine) (err error)
	GetPageCount(treegrid *treegrid.Treegrid) (int64, error)
	GetPageData(tr *treegrid.Treegrid) ([]map[string]string, error)
}
