package service

import (
	"database/sql"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/models"
)

type (
	ProcurementService interface {
		GetProcurementTx(tx *sql.Tx, id interface{}) (*models.Procurement, error)
		Handle(tx *sql.Tx, pr *models.Procurement, moduleID int) error
	}
)
