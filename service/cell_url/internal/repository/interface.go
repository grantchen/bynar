package repository

import (
	"context"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/cellurls/internal/model"
)

type DataGridRepository interface {
	// SearchWarehousesByValue retrieves a list of warehouses that match a given search value.
	// It performs a SQL query that selects specific columns from three tables, joins them together,
	// and filters the results based on a search value that is passed as a parameter.
	// The resulting list of warehouses is returned as a slice of WarehouseList structs.
	SearchWarehousesByValue(ctx context.Context, id string) ([]*model.WarehouseInfo, error)
}
