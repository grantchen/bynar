package repository

import (
	"database/sql"
	"fmt"
	"log"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/errors"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/i18n"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/logger"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/utils"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/transfers/internal/model"
	sqlbuilder "git-codecommit.eu-central-1.amazonaws.com/v1/repos/transfers/internal/repository/sql_builder"
)

type transferRepository struct {
	gridTreeRepository treegrid.GridRowRepositoryWithChild
	db                 *sql.DB
	language           string
}

// Save implements TransferRepository
func (t *transferRepository) Save(tx *sql.Tx, tr *treegrid.MainRow) error {
	if err := t.SaveTransfer(tx, tr); err != nil {
		return fmt.Errorf("%s %s: [%w]",
			i18n.Localize(t.language, errors.ErrCodeSave),
			i18n.Localize(t.language, errors.ErrCodeTransfer),
			i18n.ErrMsgToI18n(err, t.language))
	}

	if err := t.SaveTransferLines(tx, tr); err != nil {
		return fmt.Errorf("%s %s: [%w]",
			i18n.Localize(t.language, errors.ErrCodeSave),
			i18n.Localize(t.language, errors.ErrCodeTransferLine),
			i18n.ErrMsgToI18n(err, t.language))
	}

	return nil
}

// SaveDocumentID implements TransferRepository
func (*transferRepository) SaveDocumentID(tx *sql.Tx, tr *treegrid.MainRow, docID string) error {
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
			//"location_origin_id",
			//"location_destination_id",
		})

	switch tr.Fields.GetActionType() {
	case treegrid.GridRowActionAdd:
		err := tr.Fields.ValidateOnRequiredAll(requiredFieldsMapping)
		if err != nil {
			return err
		}

		err = tr.Fields.ValidateOnPositiveNumber(positiveFieldsMapping)
		if err != nil {
			return fmt.Errorf(i18n.Localize(t.language, "", err.Error()))
		}
	case treegrid.GridRowActionChanged:
		err := tr.Fields.ValidateOnRequired(requiredFieldsMapping)
		if err != nil {
			return err
		}

		err = tr.Fields.ValidateOnPositiveNumber(positiveFieldsMapping)
		if err != nil {
			return fmt.Errorf(i18n.Localize(t.language, "", err.Error()))
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
			"item_unit_id",
		})

	for _, item := range tr.Items {
		switch item.GetActionType() {
		case treegrid.GridRowActionAdd:
			err := item.ValidateOnRequiredAll(requiredFieldsMapping)
			if err != nil {
				return err
			}

			err = item.ValidateOnPositiveNumber(positiveFieldsMapping)
			if err != nil {
				return fmt.Errorf(i18n.Localize(t.language, "", err.Error()))
			}

			if err = t.validateAddTransferLine(tx, item); err != nil {
				return fmt.Errorf("validate TransferLine: [%w]", err)
			}
			err = t.gridTreeRepository.SaveLineAdd(tx, item)
			if err != nil {
				return err
			}

			continue
		case treegrid.GridRowActionChanged:
			err := item.ValidateOnRequired(requiredFieldsMapping)
			if err != nil {
				return err
			}

			err = item.ValidateOnPositiveNumber(positiveFieldsMapping)
			if err != nil {
				return fmt.Errorf(i18n.Localize(t.language, "", err.Error()))
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
func (*transferRepository) UpdateStatus(tx *sql.Tx, status int) error {
	return nil
}

// GetTransfersPageData implements TransferRepository
func (t *transferRepository) GetTransfersPageData(tg *treegrid.Treegrid) ([]map[string]string, error) {
	// Prepare filter for WHERE condition with args
	sqlbuilder.PrepFilters(tg)

	// items request
	if tg.BodyParams.GetItemsRequest() {
		logger.Debug("get items request")

		query := QueryChild + " WHERE parent = " + tg.BodyParams.ID + tg.OrderByChildQuery(model.TransferItemsFields)

		query = sqlbuilder.AddLimit(query)
		pos, _ := tg.BodyParams.IntPos()
		query = sqlbuilder.AddOffset(query, pos)

		logger.Debug("query", query)

		return t.getJSON(query, tg.FilterArgs["child"], tg)
	}

	// GROUP BY
	if tg.WithGroupBy() {
		logger.Debug("query with group by clause")

		return t.handleGroupBy(tg)
	}

	logger.Debug("get without grouping")

	query := QueryParent + tg.FilterWhere["parent"]
	if tg.FilterWhere["child"] != "" {
		query += ` AND transfers.id IN ( SELECT Parent FROM transfers_items ` + QueryChildJoins + tg.FilterWhere["child"] + `) `
	}

	query += tg.SortParams.OrderByQueryExludeChild(model.TransferItemsFields, model.FieldAliases)

	query = sqlbuilder.AddLimit(query)
	pos, _ := tg.BodyParams.IntPos()
	query = sqlbuilder.AddOffset(query, pos)
	mergedArgs := utils.MergeMaps(tg.FilterArgs["parent"], tg.FilterArgs["child"])

	logger.Debug("query", query, "args", mergedArgs)

	return t.getJSON(query, mergedArgs, tg)
}

// GetTransferCount implements TransferRepository
func (t *transferRepository) GetTransferCount(treegrid *treegrid.Treegrid) (int, error) {
	var query string

	column := model.NewColumn(treegrid.GroupCols[0])

	FilterWhere, FilterArgs := sqlbuilder.PrepQuery(treegrid.FilterParams)

	if column.IsItem {
		if FilterWhere["parent"] != "" {
			FilterWhere["parent"] = " AND transfers_items.Parent IN (SELECT transfers.id from transfers " +
				QueryParentJoins +
				sqlbuilder.DummyWhere +
				FilterWhere["parent"] + ") "
		}
		query = QueryChildCount + FilterWhere["child"] + FilterWhere["parent"]
		fmt.Printf("query count1: %s\n", query)
	} else {
		if FilterWhere["child"] != "" {
			FilterWhere["child"] = " AND transfers.id IN (SELECT transfers_items.Parent from transfers_items " +
				QueryChildJoins +
				sqlbuilder.DummyWhere +
				FilterWhere["child"] + ") "
		}

		query = QueryParentCount + FilterWhere["child"] + FilterWhere["parent"]
		fmt.Printf("filter where[child]: %s\n", FilterWhere["child"])
		fmt.Printf("query count2: %s\n", query)
	}

	mergedArgs := utils.MergeMaps(FilterArgs["child"], FilterArgs["parent"])

	rows, err := t.db.Query(query, mergedArgs...)
	if err != nil {
		log.Println(err, "query", query, "colData", column)
		return 0, err
	}

	return utils.CheckCount(rows), nil

}

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
