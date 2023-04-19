package repository

import (
	"database/sql"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/transfers/internal/model"
	treegrid_model "git-codecommit.eu-central-1.amazonaws.com/v1/repos/transfers/internal/model/treegrid"
)

type transferRepository struct {
	db *sql.DB
}

// GetTransferCount implements TransferRepository
func (*transferRepository) GetTransferCount(treegrid *treegrid_model.Treegrid) (int, error) {
	var query string

	column := model.NewColumn(treegrid.GroupCols[0])

	FilterWhere, FilterArgs := prepQuery(treegrid.FilterParams)

	if column.IsItem {
		if FilterWhere["parent"] != "" {
			FilterWhere["parent"] = " AND transfers_items.Parent IN (SELECT transfers.id from transfers " + QueryParentJoins + DummyWhere + FilterWhere["parent"] + ") "
		}
		query = QueryChildCount + FilterWhere["child"] + FilterWhere["parent"]
	} else {
		if FilterWhere["child"] != "" {
			FilterWhere["child"] = " AND transfers.id IN (SELECT transfers_items.Parent from transfers_items " + QueryChildJoins + DummyWhere + FilterWhere["child"] + ") "
		}

		query = QueryParentCount + FilterWhere["child"] + FilterWhere["parent"]
	}

}

func NewTransferRepository(db *sql.DB) TransferRepository {
	return &transferRepository{db: db}
}
