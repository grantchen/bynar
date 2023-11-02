package service

import (
	"context"
	"database/sql"
	"fmt"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/errors"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/logger"
	pkg_service "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/service"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
	"log"
)

type UploadService struct {
	db                  *sql.DB
	language            string
	approvalService     pkg_service.ApprovalCashPaymentService
	docSvc              pkg_service.DocumentService
	gridRowRepository   treegrid.GridRowRepositoryWithChild
	accountId           int
	procurementsService ProcurementsService
}

func NewUploadService(db *sql.DB,
	language string,
	approvalService pkg_service.ApprovalCashPaymentService,
	docSvc pkg_service.DocumentService,
	gridRowRepository treegrid.GridRowRepositoryWithChild,
	accountId int,
	procurementsService ProcurementsService,
) *UploadService {
	return &UploadService{
		db:                  db,
		language:            language,
		approvalService:     approvalService,
		docSvc:              docSvc,
		gridRowRepository:   gridRowRepository,
		accountId:           accountId,
		procurementsService: procurementsService,
	}
}

// Handle
func (s *UploadService) Handle(req *treegrid.PostRequest) (*treegrid.PostResponse, error) {
	resp := &treegrid.PostResponse{}

	trList, err := treegrid.ParseRequestUpload(req, s.gridRowRepository)
	if err != nil {
		return nil, fmt.Errorf("parse requst: [%w]", err)
	}

	// handle all transfer, check error and make proper response
	for _, tr := range trList.MainRows() {
		if tr.Fields["id"] == "Group" {
			continue
		}
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

func (s *UploadService) handle(tr *treegrid.MainRow) error {
	// Check Approval Order
	ok, err := s.approvalService.Check(tr, s.accountId, s.language)
	if err != nil {
		return fmt.Errorf("check order: [%w]", err)
	}

	if !ok {
		return fmt.Errorf("invalid approval order: [%w]: status: %d", errors.ErrForbiddenAction, tr.Status())
	}

	// Create new transaction
	tx, err := s.db.BeginTx(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("begin transaction: [%w]", err)
	}
	defer tx.Rollback()

	if err := s.gridRowRepository.Save(tx, tr); err != nil {
		return fmt.Errorf("procecurement svc save '%s': [%w]", tr.IDString(), err)
	}
	//todo refactor
	if tr.Status() == 1 {
		logger.Debug("status equal 1 - do calculation, status", tr.Status())

		// working with procurement - calculating and updating.
		entity, err := s.procurementsService.GetTx(tx, tr.Fields.GetID())
		if err != nil {
			return fmt.Errorf("get procurement service: [%w]", err)
		}

		if err := s.procurementsService.Handle(tx, entity); err != nil {
			return fmt.Errorf("handle procurement: [%w]", err)
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
