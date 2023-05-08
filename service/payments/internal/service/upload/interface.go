package svc_upload

import (
	"database/sql"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/payments/internal/models"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

type UploadService interface {
	Handle(req *treegrid.PostRequest) (*treegrid.PostResponse, error)
}
type PaymentService interface {
	GetTx(tx *sql.Tx, id interface{}) (*models.Payment, error)
	Handle(tx *sql.Tx, pr *models.Payment, moduleID int) error
}
