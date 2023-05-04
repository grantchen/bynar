package repository

import (
	"database/sql"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/models"
)

type boundFlowRepository struct {
}

func NewBoundFlows() BoundFlowRepository {
	return &boundFlowRepository{}
}

func (b *boundFlowRepository) SaveOutboundFlow(tx *sql.Tx, outFlow models.OutboundFlow) (err error) {
	query := `
	INSERT INTO outbound_flows (module_id, location_id, item_id, parent_id, transaction_id, posting_date, quantity, value_avco, value_fifo)
		VALUES(?,?,?,?,?,?,?,?,?)`

	_, err = tx.Exec(query,
		outFlow.ModuleID, outFlow.LocationID, outFlow.ItemID, outFlow.ParentID, outFlow.TransactionNo,
		outFlow.PostingDate, outFlow.Quantity, outFlow.ValueAvco, outFlow.ValueFifo)

	return
}

func (b *boundFlowRepository) SaveInboundFlow(tx *sql.Tx, inFlow models.InboundFlow) (err error) {
	query := `
	INSERT INTO inbound_flows (module_id, location_id, item_id, parent_id, posting_date, quantity, value,  outbound_quantity, outbound_value, status)
		VALUES(?,?,?,?,?,?,?,?,?, ?)`

	_, err = tx.Exec(query,
		inFlow.ModuleID, inFlow.LocationID, inFlow.ItemID, inFlow.ParentID,
		inFlow.PostingDate, inFlow.Quantity, inFlow.Value, inFlow.OutboundQuantity, inFlow.OutboundValue, inFlow.Status)

	return
}
