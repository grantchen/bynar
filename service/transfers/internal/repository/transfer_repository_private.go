package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/i18n"
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

// validate transfer line params and calculate quantity
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
	var unitValue float64
	if err := tx.QueryRow(query, unitID).Scan(&unitValue); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("unit not found with item_unit_id: %s", unitID)
		}

		return fmt.Errorf("query row: [%w], query: %s", err, query)
	}

	// check item unit val
	itemUnitVal, ok := item["item_unit_value"]
	if !ok {
		item["item_unit_value"] = unitValue
	} else {
		itemUnitValFloat, _ := item.GetValFloat64("item_unit_value")

		if itemUnitValFloat != unitValue {
			return fmt.Errorf("invalid item_unit_value: got '%v', want '%d'", itemUnitValFloat, itemUnitVal)
		}
	}

	// check calculated quantity
	inputQuantity, _ := item.GetValFloat64("input_quantity")
	if inputQuantity == 0 {
		return fmt.Errorf("invalid input quantity: '%v'", inputQuantity)
	}

	item["quantity"] = inputQuantity * unitValue
	if _, ok := item["Parent"]; !ok {
		return fmt.Errorf("absent 'Parent' value ")
	}

	item["parent_id"] = item["Parent"]

	return nil
}

// afterChangeTransferLine updates quantity and item_unit_value
func (t *transferRepository) afterChangeTransferLine(tx *sql.Tx, item treegrid.GridRow) error {
	if _, ok := item["input_quantity"]; ok {
		return t.updateTransferLineQuantityVals(tx, item)
	}

	if _, ok := item["item_unit_id"]; ok {
		return t.updateTransferLineQuantityVals(tx, item)

	}

	return nil
}

// updateTransferLineQuantityVals updates quantity and item_unit_value
func (t *transferRepository) updateTransferLineQuantityVals(tx *sql.Tx, item treegrid.GridRow) error {
	query := `
UPDATE transfer_lines trl
SET item_unit_value = (SELECT value FROM units WHERE id = trl.item_unit_id),
quantity = input_quantity * (SELECT value FROM units WHERE id = trl.item_unit_id)
WHERE id = ?
	`
	_, err := tx.Exec(query, item.GetLineID())

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
			return i18n.TranslationI18n(t.language, "DocumentNotExist", map[string]string{
				"DocumentId": fmt.Sprint(id),
			})
		}

		return i18n.TranslationErrorToI18n(t.language, err)
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
			return i18n.TranslationI18n(t.language, "StoreNotExist", map[string]string{
				"StoreId": fmt.Sprint(id),
			})
		}

		return i18n.TranslationErrorToI18n(t.language, err)
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
			return i18n.TranslationI18n(t.language, "ItemNotExist", map[string]string{
				"ItemId": fmt.Sprint(id),
			})
		}

		return i18n.TranslationErrorToI18n(t.language, err)
	}

	return nil
}

// validate item_unit_id
func (t *transferRepository) validateItemUnitID(tx *sql.Tx, item treegrid.GridRow) error {
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
			return i18n.TranslationI18n(t.language, "UnitNotExist", map[string]string{
				"ItemUnitId": fmt.Sprint(id),
			})
		}

		return i18n.TranslationErrorToI18n(t.language, err)
	}

	return nil
}
