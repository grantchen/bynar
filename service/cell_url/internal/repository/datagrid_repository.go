package repository

import (
	"context"
	"database/sql"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/cellurls/internal/model"
)

type mysqlDataGridRepository struct {
	db *sql.DB
}

// SearchWarehousesByValue implements DataGridRepository
func (repository *mysqlDataGridRepository) SearchWarehousesByValue(_ context.Context, value string) ([]*model.WarehouseInfo, error) {
	var lsWareHouses []*model.WarehouseInfo

	stmt, err := repository.db.Prepare(`
		Select w.name as w_name, w.id, w.costing_methods, w.uuid, im.product_description,im.product_barcode,im.quantity,im.avg_cost
		From warehouses as w, inventory_managment as im, products as p
		Where w.uuid=im.warehouse_uuid and p.uuid=im.product_uuid and concat (w.name,'',w.costing_methods,'', w.uuid) like ?
	`)

	if err != nil {
		return nil, err
	}

	result, err := stmt.Query("%" + value + "%")
	lsWareHouses = make([]*model.WarehouseInfo, 0)
	for result.Next() {
		w := &model.WarehouseInfo{}
		err = result.Scan(&w.WName, &w.Id, &w.CostingMethods, &w.Uuid, &w.ProductDescription, &w.ProductBarcode, &w.ProductUuid, &w.AvgCost)
		lsWareHouses = append(lsWareHouses, w)
	}
	return lsWareHouses, nil
}

func NewMysqlDataGridRepository(db *sql.DB) DataGridRepository {
	return &mysqlDataGridRepository{db: db}
}
