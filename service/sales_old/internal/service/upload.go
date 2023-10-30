package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/errors"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/logger"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/models"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/service"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

type (
	SaleService interface {
		GetSaleTx(tx *sql.Tx, id interface{}) (*models.Sale, error)
		Handle(tx *sql.Tx, pr *models.Sale) error
	}

	UploadSvc struct {
		accountID         int
		conn              *sql.DB
		approvalService   service.ApprovalService
		gridRowRepository treegrid.GridRowRepositoryWithChild
		saleSvc           SaleService
		docSvc            service.DocumentService
	}
)

func NewService(conn *sql.DB,
	approvalService service.ApprovalService,
	gridRowRepository treegrid.GridRowRepositoryWithChild,
	procurementSvc SaleService,
	accountID int,
	docSvc service.DocumentService,
) (*UploadSvc, error) {
	return &UploadSvc{
		conn:              conn,
		approvalService:   approvalService,
		gridRowRepository: gridRowRepository,
		saleSvc:           procurementSvc,
		accountID:         accountID,
		docSvc:            docSvc,
	}, nil
}

// Handle
func (s *UploadSvc) Handle(req *treegrid.PostRequest) (*treegrid.PostResponse, error) {
	resp := &treegrid.PostResponse{}

	trList, err := treegrid.ParseRequestUpload(req, s.gridRowRepository)
	if err != nil {
		return nil, fmt.Errorf("parse requst: [%w]", err)
	}

	// handle all transfer, check error and make proper response
	for _, tr := range trList.MainRows() {
		if err := s.handle(tr); err != nil {
			log.Println("Err", err)

			resp.IO.Result = -1
			resp.IO.Message += err.Error() + "\n"
		}

		resp.Changes = append(resp.Changes, tr.Fields)
		for k := range tr.Items {
			resp.Changes = append(resp.Changes, tr.Items[k])
		}
	}

	return resp, nil
}

func (s *UploadSvc) handle(tr *treegrid.MainRow) error {
	// Check Approval Order
	ok, err := s.approvalService.Check(tr, s.accountID, "")
	if err != nil {
		return fmt.Errorf("check order: [%w], transfer id: %s", err, tr.IDString())
	}

	if !ok {
		return fmt.Errorf("invalid approval order: [%w]: status: %d", errors.ErrForbiddenAction, tr.Status())
	}

	// Create new transaction
	tx, err := s.conn.BeginTx(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("begin transaction: [%w]", err)
	}
	defer tx.Rollback()

	if err := s.gridRowRepository.Save(tx, tr); err != nil {
		return fmt.Errorf("transfer svc save '%s': [%w]", tr.IDString(), err)
	}

	if tr.Status() == 1 {
		logger.Debug("status equal 1 - do calculation, status", tr.Status())

		// working with procurement - calculating and updating.
		entity, err := s.saleSvc.GetSaleTx(tx, tr.Fields.GetID())
		if err != nil {
			return fmt.Errorf("get sale service: [%w]", err)
		}

		if err := s.saleSvc.Handle(tx, entity); err != nil {
			return fmt.Errorf("handle sale: [%w]", err)
		}

		if err := s.docSvc.Handle(tx, entity.ID, entity.DocumentID, entity.DocumentNo); err != nil {
			return fmt.Errorf("handle document: [%w], modelID: %d, docID: %d, docNo: '%s'", err, entity.ID, entity.DocumentID, entity.DocumentNo)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: [%w]", err)
	}

	return nil
}
