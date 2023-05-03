package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/task2/internal/repository"
)

type uploadSvc struct {
	conn                *sql.DB
	userRepository      repository.UserRepository
	workflowRepository  repository.WorkflowRepository
	transferRepository  repository.TransferRepository
	inventoryRepository repository.InventoryRepository
	documentRepository  repository.DocumentRepository
}

var (
	ErrForbiddenStatus = errors.New("forbidden transfer status")
	ErrInvalidQuantity = errors.New("invalid quantity")
)

func NewUploadService(db *sql.DB, userRepository repository.UserRepository,
	workflowRepository repository.WorkflowRepository,
	transferRepository repository.TransferRepository,
	inventoryRepository repository.InventoryRepository,
	documentRepository repository.DocumentRepository) UploadService {
	return &uploadSvc{
		conn:                db,
		userRepository:      userRepository,
		workflowRepository:  workflowRepository,
		inventoryRepository: inventoryRepository,
		transferRepository:  transferRepository,
		documentRepository:  documentRepository,
	}
}

func (s *uploadSvc) Handle(req *treegrid.PostRequest, accountID int) (*treegrid.PostResponse, error) {
	trList, err := treegrid.ParseRequestUpload2(req)
	if err != nil {
		return nil, fmt.Errorf("parse requst: [%w]", err)
	}

	for _, tr := range trList.MainRows() {
		ok, err := s.workflowRepository.CheckApprovalOrder(accountID, tr.Status())
		if err != nil {
			return nil, fmt.Errorf("check order: [%w], transfer id: %s", err, tr.IDString())
		}

		if !ok {
			return nil, fmt.Errorf("%w: status: %d", ErrForbiddenStatus, tr.Status())
		}

		if err := s.handle(tr); err != nil {
			return nil, fmt.Errorf("handle transfer: [%w], id: %s", err, tr.IDString())
		}
	}

	return &treegrid.PostResponse{}, nil
}

func (s *uploadSvc) handle(tr *treegrid.MainRow) error {
	tx, err := s.conn.BeginTx(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("begin transaction: [%w]", err)
	}
	defer tx.Rollback()

	switch tr.Status() {
	// update/add
	case 0:
		if err := s.transferRepository.Save(tx, tr); err != nil {
			return fmt.Errorf("transfer svc save '%s': [%w]", tr.IDString(), err)
		}
	case 1:
		ok, err := s.inventoryRepository.CheckQuantityAndValue(tx, tr)
		if err != nil {
			return fmt.Errorf("check inventory quantity and value: [%w]", err)
		}

		if !ok {
			return ErrInvalidQuantity
		}

		if err := s.transferRepository.Save(tx, tr); err != nil {
			return fmt.Errorf("transfer svc save: [%w]", err)
		}

		if err := s.inventoryRepository.Save(tx, tr); err != nil {
			return fmt.Errorf("transfer svc save: [%w]", err)
		}

		ok, err = s.documentRepository.IsAuto(tx, tr)
		if err != nil {
			return fmt.Errorf("document svc check if is auto: [%w]", err)
		}

		if !ok {
			if err := tx.Commit(); err != nil {
				return fmt.Errorf("commit transaction: [%w]", err)
			}
		}

		docIdStr, err := s.documentRepository.Generate(tx, tr)
		if err != nil {
			return fmt.Errorf("document svc generate: [%w]", err)
		}

		if err := s.transferRepository.SaveDocumentID(tx, tr, docIdStr); err != nil {
			return fmt.Errorf("transfer svc save document id: [%w]", err)
		}
	default:
		if err := s.transferRepository.UpdateStatus(tx, tr.Status()); err != nil {
			return fmt.Errorf("transfer svc update status: [%w]", err)
		}
	}

	return nil
}
