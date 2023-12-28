package repository

import (
	"database/sql"
	"fmt"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/utils"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/warehouses/internal/model"
)

type warehousesRepository struct {
	db *sql.DB
}

// GetWarehousesAsMap implements WarehousesRepository
func (g *warehousesRepository) GetWarehousesAsMap(warehousesID int) (map[string]interface{}, error) {
	rows, err := g.getWarehouses(warehousesID)
	if err != nil {
		return nil, err
	}
	rowVals, err := utils.NewRowVals(rows)

	if err != nil {
		return nil, fmt.Errorf("parse new row error: [%w]", err)
	}

	for rows.Next() {
		if err := rowVals.Parse(rows); err != nil {
			return nil, fmt.Errorf("parse rows: [%w]", err)
		}

		return rowVals.Values(), nil
	}

	return nil, fmt.Errorf("not found general posting setup with id: [%d]", warehousesID)
}

// GetWarehouses implements UploadRepository
func (g *warehousesRepository) GetWarehouses(warehousesID int) (*model.Warehouses, error) {
	rows, err := g.getWarehouses(warehousesID)
	if err != nil {
		return nil, err
	}
	rowVals, err := utils.NewRowVals(rows)

	if err != nil {
		return nil, fmt.Errorf("parse new row error: [%w]", err)
	}

	for rows.Next() {
		if err := rowVals.Parse(rows); err != nil {
			return nil, fmt.Errorf("parse rows: [%w]", err)
		}

		data := rowVals.StringValues()
		warehouses, err := model.ParseFromMapStr(data)
		if err != nil {
			return nil, fmt.Errorf("parse json: [%w]", err)
		}
		return warehouses, nil
	}

	return nil, fmt.Errorf("not found general posting setup with id: [%d]", warehousesID)
}

func (g *warehousesRepository) getWarehouses(warehousesID int) (*sql.Rows, error) {
	query := QuerySelectWithoutJoin + " WHERE id = ?"
	stmt, err := g.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("db prepare: [%w], sql string: [%s]", err, query)
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)

	rows, err := stmt.Query(warehousesID)
	if err != nil {
		return nil, fmt.Errorf("query error: [%w], sql string: [%s]", err, query)
	}
	return rows, nil
}

// IsContainCombination implements UploadRepository
func (*warehousesRepository) IsContainCombination(*sql.Tx, int, int, int) (bool, error) {
	panic("unimplemented")
}

func NewPostingSetupRepository(db *sql.DB) WarehousesRepository {
	return &warehousesRepository{db: db}
}
