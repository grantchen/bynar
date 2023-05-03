package repository

import (
	"database/sql"
	"fmt"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

type transferRepository struct {
	grRepository treegrid.GridRowRepository
}

func NewTransferRepository(db *sql.DB) TransferRepository {
	grRepository := treegrid.NewGridRepository(db,
		"transfer",
		"transfer_line",
		TransferFieldNames,
		TransferLineFieldNames,
	)
	return &transferRepository{grRepository: grRepository}
}

// Save implements TransferRepository
func (t *transferRepository) Save(tx *sql.Tx, tr *treegrid.MainRow) error {
	if err := t.SaveTransfer(tx, tr); err != nil {
		return fmt.Errorf("save transfer: [%w]", err)
	}

	if err := t.SaveTransferLines(tx, tr); err != nil {
		return fmt.Errorf("save transfer line: [%w]", err)
	}

	return nil
}

// SaveTransfer implements TransferRepository
func (t *transferRepository) SaveTransfer(tx *sql.Tx, tr *treegrid.MainRow) error {
	return t.grRepository.SaveMainRow(tx, tr)
}

// SaveTransferLines implements TransferRepository
func (t *transferRepository) SaveTransferLines(tx *sql.Tx, tr *treegrid.MainRow) error {

	for _, item := range tr.Items {
		switch item.GetActionType() {
		case treegrid.GridRowActionAdd:
			if err := t.validateAddTransferLine(tx, item); err != nil {
				return fmt.Errorf("validate TransferLine: [%w]", err)
			}
			err := t.grRepository.SaveLineAdd(tx, item)
			if err != nil {
				return err
			}

			continue
		case treegrid.GridRowActionChanged:
			err := t.grRepository.SaveLineUpdate(tx, item)
			if err != nil {
				return err
			}

			if err := t.afterChangeTransferLine(tx, item); err != nil {
				return fmt.Errorf("afterChangeTransferLine: [%w]", err)
			}

			continue
		case treegrid.GridRowActionDeleted:
			err := t.grRepository.SaveLineDelete(tx, item)

			if err != nil {
				return err
			}
			continue
		default:
			return fmt.Errorf("undefined row type: %s", item.GetActionType())
		}

	}

	return nil

}

func (t *transferRepository) validateAddTransferLine(tx *sql.Tx, item treegrid.GridRow) error {
	unitID, ok := item["item_unit_id"]
	if !ok {
		return fmt.Errorf("absent item_unit_id")
	}

	query := `SELECT value FROM units WHERE id = ?`
	var unitVal int
	if err := tx.QueryRow(query, unitID).Scan(&unitVal); err != nil {
		return fmt.Errorf("query row: [%w], query: %s", err, query)
	}

	// check item unit val
	itemUnitVal, ok := item["item_unit_value"]
	if !ok {
		item["item_unit_value"] = unitVal
	} else {
		itemUnitValInt, _ := item.GetStrInt("item_unit_value")

		if itemUnitValInt != unitVal {
			return fmt.Errorf("invalid item_unit_value: got '%d', want '%d'", itemUnitValInt, itemUnitVal)
		}
	}

	// check calculated quantity
	inputQuantity, _ := item.GetStrInt("input_quantity")
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

func (t *transferRepository) SaveDocumentID(tx *sql.Tx, tr *treegrid.MainRow, docID string) error {
	return nil
}

func (t *transferRepository) UpdateStatus(tx *sql.Tx, status int) error {
	return nil
}
