package service

import (
	"database/sql"
	stderr "errors"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/transfers/internal/repository"
)

var (
	ErrInvalidQuantity = stderr.New("invalid quantity")
)

type transferService struct {
	db                             *sql.DB
	language                       string
	gridRowDataRepositoryWithChild treegrid.GridRowDataRepositoryWithChild
	transferRepository             repository.TransferRepository
}

// GetPageData implements TransferService
func (t *transferService) GetPageData(tr *treegrid.Treegrid) ([]map[string]string, error) {
	return t.gridRowDataRepositoryWithChild.GetPageData(tr)
}

// GetPageCount implements TransferService
func (t *transferService) GetPageCount(tr *treegrid.Treegrid) (int64, error) {
	return t.gridRowDataRepositoryWithChild.GetPageCount(tr)
}

func NewTransferService(
	db *sql.DB,
	transferRepository repository.TransferRepository,
	gridRowDataRepositoryWithChild treegrid.GridRowDataRepositoryWithChild,
) TransferService {
	return &transferService{
		db:                             db,
		transferRepository:             transferRepository,
		gridRowDataRepositoryWithChild: gridRowDataRepositoryWithChild,
	}
}
