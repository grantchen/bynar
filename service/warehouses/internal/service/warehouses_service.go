package service

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

type warehousesService struct {
	warehousesSimpleRepository treegrid.SimpleGridRowRepository
}

func NewWarehousesService(warehousesSimpleRepository treegrid.SimpleGridRowRepository) WarehousesService {
	return &warehousesService{
		warehousesSimpleRepository: warehousesSimpleRepository,
	}
}

// GetPageCount implements WarehousesService
func (g *warehousesService) GetPageCount(tr *treegrid.Treegrid) (int64, error) {
	return g.warehousesSimpleRepository.GetPageCount(tr)
}

// GetPageData implements WarehousesService
func (g *warehousesService) GetPageData(tr *treegrid.Treegrid) ([]map[string]string, error) {
	return g.warehousesSimpleRepository.GetPageData(tr)
}
