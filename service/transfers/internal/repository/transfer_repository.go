package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/i18n"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

// transferRepository is a repository for transfer
type transferRepository struct {
	gridTreeRepository treegrid.GridRowRepositoryWithChild
	db                 *sql.DB
	language           string
}

// Save implements TransferRepository
func (t *transferRepository) Save(tx *sql.Tx, tr *treegrid.MainRow) error {
	if err := t.SaveTransfer(tx, tr); err != nil {
		templateData := map[string]string{
			"Message": err.Error(),
		}
		return i18n.TranslationI18n(t.language, "SaveTransfer", templateData)
	}

	if err := t.SaveTransferLines(tx, tr); err != nil {
		templateData := map[string]string{
			"Message": err.Error(),
		}
		return i18n.TranslationI18n(t.language, "SaveTransferLine", templateData)
	}

	return nil
}

// SaveDocumentID implements TransferRepository
func (*transferRepository) SaveDocumentID(*sql.Tx, *treegrid.MainRow, string) error {
	return nil
}

// SaveTransfer implements TransferRepository
func (t *transferRepository) SaveTransfer(tx *sql.Tx, tr *treegrid.MainRow) error {
	requiredFieldsMapping := tr.Fields.FilterFieldsMapping(
		TransferFieldNames,
		[]string{
			"document_id",
			"transaction_no",
			"store_id",
			"location_origin_id",
			"location_destination_id",
		})
	positiveFieldsMapping := tr.Fields.FilterFieldsMapping(
		TransferFieldNames,
		[]string{
			"document_id",
			//"item_id",
			//"item_unit_id",
			//"project_id",
			//"area_id",
			//"department_id",
			//"in_transit_id",
			//"shipment_method_id",
			//"shipping_agent_id",
			//"shipping_agent_service_id",
			//"transaction_type_id",
			//"transaction_specification_id",
			//"user_group_id",
			"store_id",
			"location_origin_id",
			"location_destination_id",
		})

	switch tr.Fields.GetActionType() {
	case treegrid.GridRowActionAdd:
		err := tr.Fields.ValidateOnRequiredAll(requiredFieldsMapping, t.language)
		if err != nil {
			return err
		}

		err = tr.Fields.ValidateOnPositiveNumber(positiveFieldsMapping, t.language)
		if err != nil {
			return err
		}

		err = t.validateTransferParams(tx, tr)
		if err != nil {
			return err
		}
	case treegrid.GridRowActionChanged:
		err := tr.Fields.ValidateOnRequired(requiredFieldsMapping, t.language)
		if err != nil {
			return err
		}

		err = tr.Fields.ValidateOnPositiveNumber(positiveFieldsMapping, t.language)
		if err != nil {
			return err
		}

		err = t.validateTransferParams(tx, tr)
		if err != nil {
			return err
		}
	case treegrid.GridRowActionDeleted:
		// ignore id start with CR
		idStr := tr.Fields.GetIDStr()
		if !strings.HasPrefix(idStr, "CR") {
			stmt, err := tx.Prepare("DELETE FROM transfer_lines WHERE parent_id = ?")
			if err != nil {
				return err
			}

			defer func(stmt *sql.Stmt) {
				_ = stmt.Close()
			}(stmt)

			_, err = stmt.Exec(idStr)
			if err != nil {
				return err
			}
		}
	}

	return t.gridTreeRepository.SaveMainRow(tx, tr)
}

// SaveTransferLines implements TransferRepository
func (t *transferRepository) SaveTransferLines(tx *sql.Tx, tr *treegrid.MainRow) error {
	requiredFieldsMapping := tr.Fields.FilterFieldsMapping(
		TransferLineFieldNames,
		[]string{
			"item_id",
			"item_unit_id",
		})
	positiveFieldsMapping := tr.Fields.FilterFieldsMapping(
		TransferLineFieldNames,
		[]string{
			"item_id",
			"input_quantity",
			"item_unit_id",
		})

	for _, item := range tr.Items {
		switch item.GetActionType() {
		case treegrid.GridRowActionAdd:
			err := item.ValidateOnRequiredAll(requiredFieldsMapping, t.language)
			if err != nil {
				return err
			}

			err = item.ValidateOnPositiveNumber(positiveFieldsMapping, t.language)
			if err != nil {
				return err
			}

			if err = t.validateAddTransferLine(tx, item); err != nil {
				return err
			}

			err = t.gridTreeRepository.SaveLineAdd(tx, item)
			if err != nil {
				return err
			}

			continue
		case treegrid.GridRowActionChanged:
			err := item.ValidateOnRequired(requiredFieldsMapping, t.language)
			if err != nil {
				return err
			}

			err = item.ValidateOnPositiveNumber(positiveFieldsMapping, t.language)
			if err != nil {
				return err
			}

			// check item_id
			if err = t.validateItemID(tx, item); err != nil {
				return err
			}

			// check item_unit_id
			if err = t.validateItemUnitID(tx, item); err != nil {
				return err
			}

			err = t.gridTreeRepository.SaveLineUpdate(tx, item)
			if err != nil {
				return err
			}

			if err := t.afterChangeTransferLine(tx, item); err != nil {
				return fmt.Errorf("afterChangeTransferLine: [%w]", err)
			}

			continue
		case treegrid.GridRowActionDeleted:
			err := t.gridTreeRepository.SaveLineDelete(tx, item)

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

// UpdateStatus implements TransferRepository
func (*transferRepository) UpdateStatus(_ *sql.Tx, _ int) error {
	return nil
}

// NewTransferRepository returns a new TransferRepository
func NewTransferRepository(db *sql.DB, language string) TransferRepository {
	grRepository := treegrid.NewGridRepository(db,
		"transfers",
		"transfer_lines",
		TransferFieldNames,
		TransferLineFieldNames,
	)
	return &transferRepository{
		db:                 db,
		gridTreeRepository: grRepository,
		language:           language,
	}
}
