package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

type (
	inventory struct {
		ID        int
		Quantity  int
		Value     float32
		ValueFIFO float32
		IsOrigin  bool
	}

	item struct {
		LocationOriginID int
		LocationDestID   int
		PostingDate      string
		Quantity         int
		ItemID           int
	}
)

type inventoryRepository struct {
}

func NewInventoryRepository(db *sql.DB) InventoryRepository {
	return &inventoryRepository{}
}

func (ir *inventoryRepository) CheckQuantityAndValue(tx *sql.Tx, tr *treegrid.MainRow) (bool, error) {
	locOrigin, _, err := getLocations(tx, tr)
	if err != nil {
		return false, fmt.Errorf("get locations: [%w]", err)
	}

	query := `
	SELECT tl.id, tl.item_id, i.quantity
	FROM transfer_lines tl
	LEFT JOIN inventory i ON tl.item_id = i.id  
	WHERE 
		(tl.quantity > i.quantity OR i.id IS NULL)
		AND i.location_id = ?`

	args := make([]interface{}, 0, len(tr.Items))
	args = append(args, locOrigin)
	for k := range tr.Items {
		args = append(args, tr.Items[k].GetID())
	}

	exclaims := strings.Repeat("?,", len(tr.Items))
	if len(exclaims) > 0 {
		exclaims = strings.Trim(exclaims, ",")
		query += ` AND tl.id IN (%s)`
		query = fmt.Sprintf(query, exclaims)
	}

	rows, err := tx.Query(query, args...)
	if err != nil {
		return false, fmt.Errorf("[%w] query: %s", err, query)
	}

	var errStr string

	for rows.Next() {
		var id, itemId, invQuantity interface{}
		if err := rows.Scan(&id, itemId, invQuantity); err != nil {
			return false, fmt.Errorf("rows scan: [%w]", err)
		}
		errStr += fmt.Sprintf("Invalid quantity or item not found (id: %v,itemID: %v, inv quantity: %v)", id, itemId, invQuantity)
	}

	if errStr != "" {
		return false, errors.New(errStr)
	}

	return true, nil
}

func getLocations(tx *sql.Tx, tr *treegrid.MainRow) (locationOrigin, locationDest int, err error) {
	if tr.Fields.GetActionType() == treegrid.GridRowActionAdd {
		var ok bool
		if locationOrigin, ok = tr.Fields.GetValInt("location_origin_id"); !ok {
			return 0, 0, errors.New("location_origin_id not valid")
		}

		if locationDest, ok = tr.Fields.GetValInt("location_destination_id"); !ok {
			return 0, 0, errors.New("location_destination_id not valid")
		}
	} else {
		err = tx.QueryRow(`
			SELECT location_origin_id, location_destination_id
			FROM transfers
			WHERE id = ?`,
			tr.Fields.GetIDStr(),
		).Scan(&locationOrigin, &locationDest)
		if errors.Is(err, sql.ErrNoRows) {
			return 0, 0, errors.New("location_origin_id or location_destination_id not valid")
		}

		if locationOriginTmp, ok := tr.Fields.GetValInt("location_origin_id"); ok {
			locationOrigin = locationOriginTmp
		}

		if locationDestTmp, ok := tr.Fields.GetValInt("location_destination_id"); ok {
			locationDest = locationDestTmp
		}
	}

	return
}

func (ir *inventoryRepository) Save(tx *sql.Tx, tr *treegrid.MainRow) error {
	query := `
	SELECT  tl.quantity, tl.item_id, t.location_origin_id, t.location_destination_id, t.posting_date
	FROM transfers t
	INNER JOIN transfer_lines tl ON t.id = tl.parent_id
	WHERE tl.id = ?
	`
	rows, err := tx.Query(query, tr.Fields.GetID())
	if err != nil {
		return fmt.Errorf("do query: [%w], query: %s", err, query)
	}

	items := make([]item, 0, 10)
	for rows.Next() {
		var trItem item
		if err := rows.Scan(&trItem.ItemID, &trItem.Quantity, &trItem.LocationOriginID, &trItem.LocationDestID, &trItem.PostingDate); err != nil {
			return fmt.Errorf("rows scan: [%w]", err)
		}

		items = append(items, trItem)
	}

	for _, v := range items {
		if err := move(tx, v); err != nil {
			return fmt.Errorf("move item: [%w]", err)
		}
	}

	return nil
}

func move(tx *sql.Tx, trItem item) error {
	query := `
	SELECT id, quantity, value, value_fifo
	FROM inventory
	WHERE location_id = ? AND item_id = ?
	`
	var (
		invOrigin inventory
	)

	row := tx.QueryRow(query, trItem.LocationOriginID, trItem.ItemID)
	if err := row.Scan(&invOrigin.ID, &invOrigin.Quantity, &invOrigin.Value, &invOrigin.ValueFIFO); err != nil {
		return fmt.Errorf("row scan: [%w]", err)
	}

	inBoundFlow := &boundItem{
		ModuleID:         6, // TODO remove
		ItemID:           trItem.ItemID,
		LocationID:       trItem.LocationOriginID,
		PostingDate:      trItem.PostingDate,
		Quantity:         invOrigin.Quantity,
		Value:            invOrigin.Value,
		OutboundQuantity: trItem.Quantity,
	}

	invOrigin.IsOrigin = true
	if err := calcInventory(tx, invOrigin, trItem.Quantity, inBoundFlow); err != nil {
		return fmt.Errorf("calc inventory: [%w]", err)
	}

	if err := saveBoundFlow(tx, inBoundFlow); err != nil {
		return fmt.Errorf("save inboundflow: [%w]", err)
	}

	var invDest inventory
	row = tx.QueryRow(query, trItem.LocationDestID, trItem.ItemID)
	if err := row.Scan(&invDest.ID, &invDest.Quantity, &invDest.Value, &invDest.ValueFIFO); err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("row scan: [%w]", err)
	}

	// new dest location
	if invDest.ID == 0 {
		query = `
		INSERT INTO inventory (location_id, item_id, quantity, value, value_fifo)
		VALUES (?, ?, ?, ?, ?)
	`
		cost := invOrigin.Value / float32(invOrigin.Quantity)
		costFifo := invOrigin.ValueFIFO / float32(invOrigin.Quantity)
		newValue := cost * float32(trItem.Quantity)
		newValueFifo := costFifo * float32(trItem.Quantity)

		if _, err := tx.Exec(query, trItem.Quantity, newValue, newValueFifo); err != nil {
			return fmt.Errorf("tx exec: [%w], query: %s", err, query)
		}

		return nil
	}

	if err := calcInventory(tx, invDest, trItem.Quantity, nil); err != nil {
		return fmt.Errorf("calc inventory destination: [%w]", err)
	}

	return nil
}

func calcInventory(tx *sql.Tx, inv inventory, itemQuantity int, bItem *boundItem) error {
	if inv.Quantity == 0 {
		return fmt.Errorf("inventory quantity = 0, id: %v", inv.ID)
	}

	cost := inv.Value / float32(inv.Quantity)
	costFifo := inv.ValueFIFO / float32(inv.Quantity)

	mult := float32(1)
	if inv.IsOrigin {
		mult = -1
	}

	newValue := inv.Value - mult*cost*float32(itemQuantity)
	newValueFifo := inv.ValueFIFO - mult*costFifo*float32(itemQuantity)
	newQuantity := inv.Quantity - int(mult)*itemQuantity
	if inv.IsOrigin {
		bItem.OutboundValue = newValue - inv.Value
	}

	query := `
		UPDATE inventory
		SET quantity = ?, value = ?, value_fifo = ?
		WHERE id = ?
	`
	_, err := tx.Exec(query, newQuantity, newValue, newValueFifo, inv.ID)

	return err
}

type boundItem struct {
	ModuleID         int
	ItemID           int
	ParentID         interface{}
	LocationID       int
	PostingDate      string
	Quantity         int
	Value            float32
	OutboundQuantity int
	OutboundValue    float32
	Status           int
}

func (b *boundItem) CalsStatus() {
	if b.Quantity == b.OutboundQuantity && b.Value == b.OutboundValue {
		b.Status = 1

		return
	}

	b.Status = 0
}

func saveBoundFlow(tx *sql.Tx, bItem *boundItem) error {
	columnNames := []string{
		"module_id",
		"item_id",
		"parent_id",
		"location_id",
		"posting_date",
		"quantity",
		"value",
		"outbound_quantity",
		"outbound_value",
		"status",
	}
	columnNamesStr := strings.Join(columnNames, ",")
	columnVals := strings.Repeat("?,", len(columnNames))
	columnVals = columnVals[:len(columnVals)-1]

	query := `
	INSERT INTO inbound_flow (` + columnNamesStr + `)
	VALUES (` + columnVals + `)
	`
	if _, err := tx.Exec(query, bItem.ModuleID, bItem.ItemID, bItem.ParentID, bItem.LocationID, bItem.PostingDate,
		bItem.Quantity, bItem.Value, bItem.OutboundQuantity, bItem.OutboundValue, bItem.Status); err != nil {
		return fmt.Errorf("exec inbound flow: [%w], query: [%s]", err, query)
	}

	query = `
	INSERT INTO outbound_flow (module_id, item_it, parent_id, location_id, posting_date, quantity, value_avco, value_fifo)
	VALUES (?,?,?,?,?,?,?,?,)
	`

	if _, err := tx.Exec(query, bItem.ModuleID, bItem.ItemID, bItem.ParentID, bItem.LocationID,
		bItem.PostingDate, bItem.OutboundQuantity, bItem.OutboundValue, bItem.OutboundValue); err != nil {
		return fmt.Errorf("exec outbound flow: [%w], query: %s", err, query)
	}

	return nil
}
