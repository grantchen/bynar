package repository

import (
	"database/sql"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/models"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

type WorkflowRepository interface {
	GetWorkflowItem(accountID, documentID int) (models.WorkflowItem, error)
	CheckApprovalOrder(conn *sql.DB, tr *treegrid.MainRow, accountID int) (bool, error)
}

type UserRepository interface {
	AccountID() int
	ModuleID() int
	HasPermission(moduleID string, accountID string) (bool, error)
	GetUserID(accountID int) (int, error)
	GetUserGroupID(accountID int) (int, error)
}

type ProcurementRepository interface {
	GetDocID(id interface{}) (docID int, err error)
	GetStatus(id interface{}) (status int, err error)
	GetProcurement(tx *sql.Tx, id interface{}) (m *models.Procurement, err error)
	SaveProcurement(tx *sql.Tx, m *models.Procurement) (err error)
	GetProcurementLines(tx *sql.Tx, id interface{}) ([]*models.ProcurementLine, error)
	SaveProcurementLine(tx *sql.Tx, prLine *models.ProcurementLine) (err error)
}

type DocumentRepository interface {
	GetDocument(docID int) (m models.Document, err error)
	GetDocumentSeries(seriesID int) (m models.DocumentSeries, err error)
	GetDocumentSeriesItem(seriesID int) (m models.DocumentSeriesItem, err error)
	UpdateDocumentSeriesItem(tx *sql.Tx, item models.DocumentSeriesItem) (err error)
	UpdateDocNumber(tx *sql.Tx, id int, docNumber string) (err error)
}

type CurrencyRepository interface {
	GetDiscount(id int) (m models.DiscountVat, err error)
	GetVat(id int) (m models.DiscountVat, err error)
	GetCurrency(id int) (m models.Currency, err error)
	GetLedgerSetupCurrency() (curr float32, err error)
}

type CashReceiptRepository interface {
	GetDocID(id interface{}) (docID int, err error)
	GetStatus(id interface{}) (status int, err error)
	Get(tx *sql.Tx, id interface{}) (m *models.CashReceipt, err error)
	Save(tx *sql.Tx, m *models.CashReceipt) (err error)
	GetLines(tx *sql.Tx, parentID interface{}) ([]*models.CashReceiptLine, error)
	SaveLine(tx *sql.Tx, l *models.CashReceiptLine) (err error)
}

type CashManagementRepository interface {
	Get(bankID int) (m *models.CashManagement, err error)
	Update(tx *sql.Tx, m *models.CashManagement) (err error)
}

type PaymentRepository interface {
	GetDocID(id interface{}) (docID int, err error)
	GetStatus(id interface{}) (status int, err error)
	Get(tx *sql.Tx, id interface{}) (m *models.Payment, err error)
	Save(tx *sql.Tx, m *models.Payment) (err error)
	GetLines(tx *sql.Tx, parentID interface{}, applied int) ([]*models.PaymentLine, error)
	SaveLine(tx *sql.Tx, l *models.PaymentLine) (err error)
}

type UnitRepository interface {
	Get(id interface{}) (m models.Unit, err error)
}

type InventoryRepository interface {
	GetInventory(tx *sql.Tx, itemID int, locationID int) (m models.Inventory, err error)
	CreateInventory(tx *sql.Tx, itemID int, locationID int) (m models.Inventory, err error)
	UpdateInventory(tx *sql.Tx, inv models.Inventory) error
	AddValues(tx *sql.Tx, itemID, locationID int, quantity, val float32) (err error)
}

type BoundFlowRepository interface {
	SaveOutboundFlow(tx *sql.Tx, outFlow models.OutboundFlow) (err error)
	SaveInboundFlow(tx *sql.Tx, inFlow models.InboundFlow) (err error)
}

type AccountManagerRepository interface {
	// CheckPermission detect if user can access endpoints or not.
	CheckPermission(accountID int, organizationID int) (*PermissionInfo, bool, error)

	// CheckRole detect if user can access a specfic endpoint or not
	CheckRole(accountID int) (map[string]int, error)
}
