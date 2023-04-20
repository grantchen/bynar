package repository

import (
	"database/sql"
	"log"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/logger"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/utils"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/transfers/internal/model"
	treegrid_model "git-codecommit.eu-central-1.amazonaws.com/v1/repos/transfers/internal/model/treegrid"
	sqlbuilder "git-codecommit.eu-central-1.amazonaws.com/v1/repos/transfers/internal/repository/sql_builder"
)

type transferRepository struct {
	db *sql.DB
}

// GetTransfersPageData implements TransferRepository
func (t *transferRepository) GetTransfersPageData(tg *treegrid_model.Treegrid) ([]map[string]string, error) {
	if tg.BodyParams.GetItemsRequest() {
		logger.Debug("get items request")

		query := sqlbuilder.QueryChild + " WHERE parent = " + tg.BodyParams.ID + tg.FilterWhere["child"] +
			sqlbuilder.OrderByQuery(tg.SortParams, model.TransferItemsFields)

		query = sqlbuilder.AddLimit(query)
		pos, _ := tg.BodyParams.IntPos()
		query = sqlbuilder.AddOffset(query, pos)

		logger.Debug("query", query, "args count", len(tg.FilterArgs["child"]))

		return t.getJSON(query, tg.FilterArgs["child"], tg)
	}

	// GROUP BY
	if tg.WithGroupBy() {
		logger.Debug("query with group by clause")

		return t.handleGroupBy(tg)
	}

	logger.Debug("get transfers without grouping")

	query := sqlbuilder.QueryParent + tg.FilterWhere["child"] + tg.FilterWhere["parent"] + sqlbuilder.OrderByQuery(tg.SortParams, nil)

	query = sqlbuilder.AddLimit(query)
	pos, _ := tg.BodyParams.IntPos()
	query = sqlbuilder.AddOffset(query, pos)
	mergedArgs := utils.MergeMaps(tg.FilterArgs["child"], tg.FilterArgs["parent"])

	logger.Debug("query", query)

	return t.getJSON(query, mergedArgs, tg)
}

// GetTransferCount implements TransferRepository
func (t *transferRepository) GetTransferCount(treegrid *treegrid_model.Treegrid) (int, error) {
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
	} else {
		if FilterWhere["child"] != "" {
			FilterWhere["child"] = " AND transfers.id IN (SELECT transfers_items.Parent from transfers_items " +
				sqlbuilder.QueryChildJoins +
				sqlbuilder.DummyWhere +
				FilterWhere["child"] + ") "
		}

		query = sqlbuilder.QueryParentCount + FilterWhere["child"] + FilterWhere["parent"]
	}

	mergedArgs := utils.MergeMaps(FilterArgs["child"], FilterArgs["parent"])
	rows, err := t.db.Query(query, mergedArgs...)
	if err != nil {
		log.Fatalln(err, "query", query, "colData", column)
	}

	return utils.CheckCount(rows), nil

}

func NewTransferRepository(db *sql.DB) TransferRepository {
	return &transferRepository{db: db}
}
