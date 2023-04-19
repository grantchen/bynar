package repository

import treegrid_model "git-codecommit.eu-central-1.amazonaws.com/v1/repos/transfers/internal/model/treegrid"

type TransferRepository interface {
	GetTransferCount(treegrid *treegrid_model.Treegrid) (int, error)
}
