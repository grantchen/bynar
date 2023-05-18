package treegrid

import (
	"database/sql"
	"fmt"
	"log"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/utils"
)

type GridRowDataRepositoryWithChild interface {
	GetPageCount(tg *Treegrid) int64
	GetPageData(tg *Treegrid) ([]map[string]string, error)
}

type GridRowDataRepositoryWithChildCfg struct {
	MainCol                  string
	MapSorted                map[string]bool
	QueryParentJoins         string
	QueryParentCount         string
	QueryChildJoins          string
	QueryChildCount          string
	ChildJoinFieldWithParent string // example in user_group_lines: parent_id
	ParentIdField            string // example in user_groups: id
}

type gridRowDataRepositoryWithChild struct {
	db                 *sql.DB
	tableName          string
	lineTableName      string
	parentFieldMapping map[string][]string
	childFieldMapping  map[string][]string
	pageSize           int
	cfg                *GridRowDataRepositoryWithChildCfg
}

func NewGridRowDataRepositoryWithChild(
	db *sql.DB,
	tableName string,
	lineTableName string,
	parentFieldMapping map[string][]string,
	childFieldMapping map[string][]string,
	pageSize int,
	cfg *GridRowDataRepositoryWithChildCfg,
) GridRowDataRepositoryWithChild {
	return &gridRowDataRepositoryWithChild{
		db:                 db,
		tableName:          tableName,
		lineTableName:      lineTableName,
		parentFieldMapping: parentFieldMapping,
		childFieldMapping:  childFieldMapping,
		pageSize:           pageSize,
		cfg:                cfg,
	}
}

// GetPageCount implements GridRowDataRepositoryWithChild
func (g *gridRowDataRepositoryWithChild) GetPageCount(tg *Treegrid) int64 {
	var query string

	column := NewColumn(tg.GroupCols[0], g.childFieldMapping, g.parentFieldMapping)
	FilterWhere, FilterArgs := PrepQueryComplex(tg.FilterParams, g.childFieldMapping, g.parentFieldMapping)

	if column.IsItem {
		if FilterWhere["parent"] != "" {
			// FilterWhere["parent"] = " AND transfers_items.Parent IN (SELECT transfers.id from transfers " +
			FilterWhere["parent"] = fmt.Sprintf(" AND %s.%s IN (SELECT %s.%s FROM %s ",
				g.lineTableName,
				g.cfg.ChildJoinFieldWithParent,
				g.tableName,
				g.cfg.ParentIdField,
				g.tableName) +

				g.cfg.QueryParentJoins +
				DummyWhere +
				FilterWhere["parent"] + ") "
		}
		query = g.cfg.QueryChildCount + FilterWhere["child"] + FilterWhere["parent"]
		fmt.Printf("query count1: %s\n", query)
	} else {
		if FilterWhere["child"] != "" {
			// FilterWhere["child"] = " AND transfers.id IN (SELECT transfers_items.Parent from transfers_items " +
			FilterWhere["child"] = fmt.Sprintf(" AND %s.%s IN (SELECT %s.%s FROM %s ",
				g.tableName,
				g.cfg.ParentIdField,
				g.lineTableName,
				g.cfg.ChildJoinFieldWithParent,
				g.lineTableName) +

				g.cfg.QueryChildJoins +
				DummyWhere +
				FilterWhere["child"] + ") "
		}

		query = g.cfg.QueryParentCount + FilterWhere["child"] + FilterWhere["parent"]
		fmt.Printf("query count2: %s\n", query)
	}

	mergedArgs := utils.MergeMaps(FilterArgs["child"], FilterArgs["parent"])
	rows, err := g.db.Query(query, mergedArgs...)
	if err != nil {
		log.Fatalln(err, "query", query, "colData", column)
	}

	return int64(utils.CheckCount(rows))
	// return 0
}

// GetPageData implements GridRowDataRepositoryWithChild
func (*gridRowDataRepositoryWithChild) GetPageData(tg *Treegrid) ([]map[string]string, error) {
	result := make([]map[string]string, 0)
	return result, nil
}
