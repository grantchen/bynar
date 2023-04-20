package service

import (
	"math"

	treegrid_model "git-codecommit.eu-central-1.amazonaws.com/v1/repos/transfers/internal/model/treegrid"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/transfers/internal/repository"
)

const (
	pageSize int = 100
)

type transferService struct {
	transferRepository repository.TransferRepository
}

// GetTransfersPageData implements TransferService
func (t *transferService) GetTransfersPageData(tr *treegrid_model.Treegrid) ([]map[string]string, error) {
	return t.GetTransfersPageData(tr)
}

// GetPagesCount implements TransferService
func (t *transferService) GetPagesCount(tr *treegrid_model.Treegrid) float64 {
	rowsCount, _ := t.transferRepository.GetTransferCount(tr)

	return math.Ceil(float64(rowsCount) / float64(pageSize))
}

func NewTransferService(transferRepository repository.TransferRepository) TransferService {
	return &transferService{}
}
