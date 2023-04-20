package service

import treegrid_model "git-codecommit.eu-central-1.amazonaws.com/v1/repos/transfers/internal/model/treegrid"

type TransferService interface {
	GetPagesCount(tr *treegrid_model.Treegrid) float64
	GetTransfersPageData(tr *treegrid_model.Treegrid) ([]map[string]string, error)
}
