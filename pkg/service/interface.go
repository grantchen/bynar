package service

import (
	"database/sql"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/middleware"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/models"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

type PaymentService interface {
	GetTx(tx *sql.Tx, id interface{}) (*models.Payment, error)
	Handle(tx *sql.Tx, m *models.Payment, moduleID int) error
	HandleLine(tx *sql.Tx, payment *models.Payment, line *models.PaymentLine) (err error)
}

type DocumentService interface {
	Handle(tx *sql.Tx, modelID, docID int, docNo string) error
}

type CashReceiptService interface {
	GetPaymentTx(tx *sql.Tx, id interface{}) (*models.CashReceipt, error)
	Handle(tx *sql.Tx, m *models.CashReceipt, moduleID int) error
	HandleLine(tx *sql.Tx, pr *models.CashReceipt, l *models.CashReceiptLine) (err error)
}

type ApprovalService interface {
	Check(tr *treegrid.MainRow, moduleID, accountID int) (bool, error)
}

type ApprovalCashPaymentService interface {
	Check(tr *treegrid.MainRow, moduleID, accountID int) (bool, error)
}

type AccountManagerService interface {
	CheckPermission(claims *middleware.IdTokenClaims) (*repository.PermissionInfo, bool, error)
	GetNewStringConnection(tenantUuid, organizationUid string, permission *repository.PermissionInfo) (string, error)
	GetRole(uid int) (map[string]int, error)
}
