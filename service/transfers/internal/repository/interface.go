package repository

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

type TransferRepository interface {
	GetTransferCount(treegrid *treegrid.Treegrid) (int, error)
	GetTransfersPageData(tr *treegrid.Treegrid) ([]map[string]string, error)
}
