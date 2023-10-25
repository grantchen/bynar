package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

func (t *transferRepository) validateAddTransferLine(tx *sql.Tx, item treegrid.GridRow) error {
	unitID, ok := item["item_unit_id"]
	if !ok {
		return fmt.Errorf("absent item_unit_id")
	}

	query := `SELECT value FROM units WHERE id = ?`
	var value float64
	if err := tx.QueryRow(query, unitID).Scan(&value); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("no units with item_unit_id: %s", unitID)
		}

		return fmt.Errorf("query row: [%w], query: %s", err, query)
	}

	unitVal := int(value)

	// check item unit val
	itemUnitVal, ok := item["item_unit_value"]
	if !ok {
		item["item_unit_value"] = unitVal
	} else {
		itemUnitValInt, _ := item.GetValInt("item_unit_value")

		if itemUnitValInt != unitVal {
			return fmt.Errorf("invalid item_unit_value: got '%d', want '%d'", itemUnitValInt, itemUnitVal)
		}
	}

	// check calculated quantity
	inputQuantity, _ := item.GetValInt("input_quantity")
	if inputQuantity == 0 {
		return fmt.Errorf("invalid input quantity: '%d'", inputQuantity)
	}

	item["quantity"] = inputQuantity * unitVal
	if _, ok := item["Parent"]; !ok {
		return fmt.Errorf("absent 'Parent' value ")
	}

	item["parent_id"] = item["Parent"]

	return nil
}

func (t *transferRepository) afterChangeTransferLine(tx *sql.Tx, item treegrid.GridRow) error {
	if _, ok := item["input_quantaty"]; ok {
		return t.updateTransferLineQuantityVals(tx, item)
	}

	if _, ok := item["item_unit_id"]; ok {
		return t.updateTransferLineQuantityVals(tx, item)

	}

	return nil
}

func (t *transferRepository) updateTransferLineQuantityVals(tx *sql.Tx, item treegrid.GridRow) error {
	query := `
UPDATE transfer_lines trl
SET item_unit_value = (SELECT value FROM units WHERE id = trl.item_unit_id),
quantity = input_quantity * (SELECT value FROM units WHERE id = trl.item_unit_id)
WHERE id = 1
	`
	_, err := tx.Exec(query, item.GetID())

	return err
}
