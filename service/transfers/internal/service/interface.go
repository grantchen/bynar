package service

import "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"

type TransferService interface {
	GetPagesCount(tr *treegrid.Treegrid) float64
	GetTransfersPageData(tr *treegrid.Treegrid) ([]map[string]string, error)
}
