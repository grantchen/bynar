package repository

import (
	"database/sql"
	"fmt"
	"log"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/logger"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/utils"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/transfers/internal/model"
	sqlbuilder "git-codecommit.eu-central-1.amazonaws.com/v1/repos/transfers/internal/repository/sql_builder"
)

type transferRepository struct {
	gridTreeRepository treegrid.GridRowRepositoryWithChild
	db                 *sql.DB
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

// SaveDocumentID implements TransferRepository
func (*transferRepository) SaveDocumentID(tx *sql.Tx, tr *treegrid.MainRow, docID string) error {
	return nil
}

// SaveTransfer implements TransferRepository
func (t *transferRepository) SaveTransfer(tx *sql.Tx, tr *treegrid.MainRow) error {
	return t.gridTreeRepository.SaveMainRow(tx, tr)
}

// SaveTransferLines implements TransferRepository
func (t *transferRepository) SaveTransferLines(tx *sql.Tx, tr *treegrid.MainRow) error {
	for _, item := range tr.Items {
		switch item.GetActionType() {
		case treegrid.GridRowActionAdd:
			if err := t.validateAddTransferLine(tx, item); err != nil {
				return fmt.Errorf("validate TransferLine: [%w]", err)
			}
			err := t.gridTreeRepository.SaveLineAdd(tx, item)
			if err != nil {
				return err
			}

			continue
		case treegrid.GridRowActionChanged:
			err := t.gridTreeRepository.SaveLineUpdate(tx, item)
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

		query := sqlbuilder.QueryChild + " WHERE parent = " + tg.BodyParams.ID + tg.OrderByChildQuery(model.TransferItemsFields)

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

	logger.Debug("get transfers without grouping")

	query := sqlbuilder.QueryParent + tg.FilterWhere["parent"]
	if tg.FilterWhere["child"] != "" {
		query += ` AND transfers.id IN ( SELECT Parent FROM transfers_items ` + sqlbuilder.QueryChildJoins + tg.FilterWhere["child"] + `) `
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
				sqlbuilder.QueryParentJoins +
				sqlbuilder.DummyWhere +
				FilterWhere["parent"] + ") "
		}
		query = sqlbuilder.QueryChildCount + FilterWhere["child"] + FilterWhere["parent"]
		fmt.Printf("query count1: %s\n", query)
	} else {
		if FilterWhere["child"] != "" {
			FilterWhere["child"] = " AND transfers.id IN (SELECT transfers_items.Parent from transfers_items " +
				sqlbuilder.QueryChildJoins +
				sqlbuilder.DummyWhere +
				FilterWhere["child"] + ") "
		}

		query = sqlbuilder.QueryParentCount + FilterWhere["child"] + FilterWhere["parent"]
		fmt.Printf("query count2: %s\n", query)
	}

	mergedArgs := utils.MergeMaps(FilterArgs["child"], FilterArgs["parent"])

	rows, err := t.db.Query(query, mergedArgs...)
	if err != nil {
		log.Fatalln(err, "query", query, "colData", column)
	}

	return utils.CheckCount(rows), nil

}

func NewTransferRepository(db *sql.DB) TransferRepository {
	grRepository := treegrid.NewGridRepository(db,
		"transfers",
		"transfer_lines",
		TransferFieldNames,
		TransferLineFieldNames,
	)
	return &transferRepository{
		db:                 db,
		gridTreeRepository: grRepository,
	}
}
