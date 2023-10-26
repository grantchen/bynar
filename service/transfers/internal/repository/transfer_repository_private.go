package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

// validate transfer params
func (t *transferRepository) validateTransferParams(tx *sql.Tx, tr *treegrid.MainRow) error {
	// check document_id
	if err := t.validateDocumentID(tx, tr); err != nil {
		return err
	}

	// check store_id
	if err := t.validateStoreID(tx, tr); err != nil {
		return err
	}

	return nil
}

func (t *transferRepository) validateAddTransferLine(tx *sql.Tx, item treegrid.GridRow) error {
	// check item_id
	if err := t.validateItemID(tx, item); err != nil {
		return err
	}

	unitID, ok := item["item_unit_id"]
	if !ok {
		return fmt.Errorf("absent item_unit_id")
	}

	query := `SELECT value FROM units WHERE id = ?`
	var value float64
	if err := tx.QueryRow(query, unitID).Scan(&value); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("unit not found with item_unit_id: %s", unitID)
		}

		return fmt.Errorf("query row: [%w], query: %s", err, query)
	}

	unitIntValue := int(value)

	// check item unit val
	itemUnitVal, ok := item["item_unit_value"]
	if !ok {
		item["item_unit_value"] = unitIntValue
	} else {
		itemUnitValInt, _ := item.GetValInt("item_unit_value")

		if itemUnitValInt != unitIntValue {
			return fmt.Errorf("invalid item_unit_value: got '%d', want '%d'", itemUnitValInt, itemUnitVal)
		}
	}

	// check calculated quantity
	inputQuantity, _ := item.GetValFloat64("input_quantity")
	if inputQuantity == 0 {
		return fmt.Errorf("invalid input quantity: '%v'", inputQuantity)
	}

	item["quantity"] = inputQuantity * float64(unitIntValue)
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

// validate document_id
func (t *transferRepository) validateDocumentID(tx *sql.Tx, tr *treegrid.MainRow) error {
	id, ok := tr.Fields["document_id"]
	if !ok {
		return nil
	}

	query := `SELECT 1 FROM documents WHERE id = ?`
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var existFlag int
	if err = stmt.QueryRow(id).Scan(&existFlag); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("document not found with document_id: %s", id)
		}

		return err
	}

	return nil
}

// validate store_id
func (t *transferRepository) validateStoreID(tx *sql.Tx, tr *treegrid.MainRow) error {
	id, ok := tr.Fields["store_id"]
	if !ok {
		return nil
	}

	query := `SELECT 1 FROM stores WHERE id = ?`
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var existFlag int
	if err = stmt.QueryRow(id).Scan(&existFlag); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("store not found with store_id: %s", id)
		}

		return err
	}

	return nil
}

// validate item_id
func (t *transferRepository) validateItemID(tx *sql.Tx, item treegrid.GridRow) error {
	id, ok := item["item_id"]
	if !ok {
		return nil
	}

	query := `SELECT 1 FROM items WHERE id = ?`
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var existFlag int
	if err = stmt.QueryRow(id).Scan(&existFlag); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("item not found with item_id: %s", id)
		}

		return err
	}

	return nil
}

// validate item_unit_id
func (t *transferRepository) validateItemUintID(tx *sql.Tx, item treegrid.GridRow) error {
	id, ok := item["item_unit_id"]
	if !ok {
		return nil
	}

	query := `SELECT 1 FROM units WHERE id = ?`
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var existFlag int
	if err = stmt.QueryRow(id).Scan(&existFlag); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("uint not found with item_unit_id: %s", id)
		}

		return err
	}

	return nil
}
