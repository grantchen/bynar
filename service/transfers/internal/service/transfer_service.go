package service

import (
	"database/sql"
	stderr "errors"
	"math"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/transfers/internal/repository"
)

const (
	pageSize int = 100
)

var (
	ErrForbiddenStatus = stderr.New("forbidden transfer status")
	ErrInvalidQuantity = stderr.New("invalid quantity")
)

type transferService struct {
	db                             *sql.DB
	language                       string
	gridRowDataRepositoryWithChild treegrid.GridRowDataRepositoryWithChild
	//userRepository      repository.UserRepository
	//workflowRepository  repository.WorkflowRepository
	transferRepository repository.TransferRepository
	//inventoryRepository repository.InventoryRepository
	//documentRepository  repository.DocumentRepository
}

// GetPageData implements TransferService
func (t *transferService) GetPageData(tr *treegrid.Treegrid) ([]map[string]string, error) {
	return t.transferRepository.GetTransfersPageData(tr)
}

// GetPageCount implements TransferService
func (t *transferService) GetPageCount(tr *treegrid.Treegrid) (int64, error) {
	rowsCount, _ := t.transferRepository.GetTransferCount(tr)

	return int64(math.Ceil(float64(rowsCount) / float64(pageSize))), nil
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

//
//func (t *transferService) HandleUpload(req *treegrid.PostRequest, accountID int) (*treegrid.PostResponse, error) {
//	trList, err := treegrid.ParseRequestUpload2(req)
//	if err != nil {
//		return nil, fmt.Errorf("parse requst: [%w]", err)
//	}
//
//	for _, tr := range trList.MainRows() {
//		ok, err := t.workflowRepository.CheckApprovalOrder(accountID, tr.Status())
//		if err != nil {
//			return nil, fmt.Errorf("check order: [%w], transfer id: %s", err, tr.IDString())
//		}
//
//		if !ok {
//			return nil, fmt.Errorf("%w: status: %d", ErrForbiddenStatus, tr.Status())
//		}
//
//		if err := t.handleUpload(tr); err != nil {
//			return nil, fmt.Errorf("handle transfer: [%w], id: %s", err, tr.IDString())
//		}
//	}
//
//	return &treegrid.PostResponse{}, nil
//}
//
//func (t *transferService) handleUpload(tr *treegrid.MainRow) error {
//	tx, err := t.conn.BeginTx(context.Background(), nil)
//	if err != nil {
//		return fmt.Errorf("begin transaction: [%w]", err)
//	}
//	defer tx.Rollback()
//
//	switch tr.Status() {
//	// update/add
//	case 0:
//		if err := t.transferRepository.Save(tx, tr); err != nil {
//			return fmt.Errorf("transfer svc save '%s': [%w]", tr.IDString(), err)
//		}
//	case 1:
//		ok, err := t.inventoryRepository.CheckQuantityAndValue(tx, tr)
//		if err != nil {
//			return fmt.Errorf("check inventory quantity and value: [%w]", err)
//		}
//
//		if !ok {
//			return ErrInvalidQuantity
//		}
//
//		if err := t.transferRepository.Save(tx, tr); err != nil {
//			return fmt.Errorf("transfer svc save: [%w]", err)
//		}
//
//		if err := t.inventoryRepository.Save(tx, tr); err != nil {
//			return fmt.Errorf("transfer svc save: [%w]", err)
//		}
//
//		ok, err = t.documentRepository.IsAuto(tx, tr)
//		if err != nil {
//			return fmt.Errorf("document svc check if is auto: [%w]", err)
//		}
//
//		if !ok {
//			if err := tx.Commit(); err != nil {
//				return fmt.Errorf("commit transaction: [%w]", err)
//			}
//		}
//
//		docIdStr, err := t.documentRepository.Generate(tx, tr)
//		if err != nil {
//			return fmt.Errorf("document svc generate: [%w]", err)
//		}
//
//		if err := t.transferRepository.SaveDocumentID(tx, tr, docIdStr); err != nil {
//			return fmt.Errorf("transfer svc save document id: [%w]", err)
//		}
//	default:
//		if err := t.transferRepository.UpdateStatus(tx, tr.Status()); err != nil {
//			return fmt.Errorf("transfer svc update status: [%w]", err)
//		}
//	}
//
//	if err = tx.Commit(); err != nil {
//		return fmt.Errorf("%s: [%w]", i18n.Localize(t.language, errors.ErrCodeCommitTransaction), err)
//	}
//
//	return nil
//}
