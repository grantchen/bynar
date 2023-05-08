package service

import (
	"context"
	"log"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/cellurls/internal/model"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/cellurls/internal/repository"
)

type dataGridService struct {
	dataGridRepository repository.DataGridRepository
}

// GetGridDataByValue retrieves a grid data object containing a list of warehouses that match a given search value.
// If there is an error during the search, it returns an error.
// Otherwise, it creates a new grid data object containing the list of warehouses and returns it.
func (service *dataGridService) GetGridDataByValue(ctx context.Context, value string, id string) (*model.Grid, error) {
	lsWarehouses, err := service.dataGridRepository.SearchWarehousesByValue(ctx, value)
	if err != nil {
		log.Printf("Error getting warehouses for search value %s: %s", value, err.Error())
		return nil, err
	}

	lsWarehouses = append(lsWarehouses, service.getTitleOfGrid())
	for i, j := 0, len(lsWarehouses)-1; i < j; i, j = i+1, j-1 {
		lsWarehouses[i], lsWarehouses[j] = lsWarehouses[j], lsWarehouses[i]
	}

	changes := []*model.Change{{ID: id, Data: &model.DataItem{Items: lsWarehouses}}}
	gridData := &model.Grid{Changes: changes}
	return gridData, nil
}

func (service *dataGridService) getTitleOfGrid() *model.WarehouseInfo {
	return &model.WarehouseInfo{
		WName:              "Warehouse Description",
		CostingMethods:     "Warehouse Costing Methods",
		Uuid:               "Warehouse UUID",
		ProductDescription: "Product Description",
		ProductBarcode:     "Product Barcode",
		Quantity:           "Product Quantity",
		AvgCost:            "Product Average Cost",
		ProductUuid:        "Product UUID",
	}
}

func NewDataGridService(dataGridRepository repository.DataGridRepository) DataGridService {
	return &dataGridService{
		dataGridRepository: dataGridRepository,
	}
}
