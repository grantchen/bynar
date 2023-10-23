package repository

import (
	"database/sql"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/warehouses/internal/model"
)

type WarehousesRepository interface {
	IsContainCombination(tx *sql.Tx, status int, generalProductPostingGroupID int, generalBussinessPostingGroupID int) (bool, error)
	GetWarehouses(warehousesID int) (*model.Warehouses, error)
	GetWarehousesAsMap(warehousesID int) (map[string]interface{}, error)
}
