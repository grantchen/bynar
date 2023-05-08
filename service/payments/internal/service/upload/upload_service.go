package svc_upload

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/logger"
	pkg_service "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/service"
	pkg_treegrid "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

type uploadService struct {
	moduleID          int
	accountID         int
	conn              *sql.DB
	approvalService   pkg_service.ApprovalCashPaymentService
	gridRowRepository pkg_treegrid.GridRowRepository
	paymentService    PaymentService
	docSvc            pkg_service.DocumentService
}

type ErrUpload struct {
	Err    error
	ID     string
	Status int
}

func (e *ErrUpload) Error() string {
	return fmt.Sprintf("ErrUpload: %s, ID: %v, Status: %v", e.Err, e.ID, e.Status)
}

type ErrDoc struct {
	Err   error
	ID    int
	DocID int
	DocNo string
}

func (e *ErrDoc) Error() string {
	return fmt.Sprintf("ErrDoc: %s, ID: %v, DocID: %v, DocNo: %v", e.Err, e.ID, e.DocID, e.DocNo)
}

// Handle implements UploadSvc

func NewUploadService(conn *sql.DB,
	approvalService pkg_service.ApprovalCashPaymentService,
	gridRowReppository pkg_treegrid.GridRowRepository,
	procurementService PaymentService,
	moduleID, accoundID int,
	documentService pkg_service.DocumentService,
) (UploadService, error) {

	return &uploadService{
		conn:              conn,
		approvalService:   approvalService,
		gridRowRepository: gridRowReppository,
		paymentService:    procurementService,
		moduleID:          moduleID,
		accountID:         accoundID,
		docSvc:            documentService,
	}, nil
}

var (
	ErrForbiddenAction       = errors.New("forbidden action")
	ErrMissingRequiredParams = errors.New("missing required params")
	ErrInvalidQuantity       = errors.New("invalid quantity")
)

// Handle
func (s *uploadService) Handle(req *pkg_treegrid.PostRequest) (*pkg_treegrid.PostResponse, error) {
	resp := &pkg_treegrid.PostResponse{}

	trList, err := pkg_treegrid.ParseRequestUpload(req, s.gridRowRepository)
	if err != nil {
		return nil, fmt.Errorf("could notparse requst: [%w]", err)
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

func (s *uploadService) handle(tr *pkg_treegrid.MainRow) error {
	// Check Approval Order
	ok, err := s.approvalService.Check(tr, s.moduleID, s.accountID)
	if err != nil {
		return &ErrUpload{Err: err, ID: tr.IDString()}

	}

	if !ok {
		return &ErrUpload{Err: ErrForbiddenAction, Status: tr.Status()}
	}

	// Create new transaction
	tx, err := s.conn.BeginTx(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("could not begin transaction: [%w]", err)
	}
	defer tx.Rollback()

	if err := s.gridRowRepository.Save(tx, tr); err != nil {
		return &ErrUpload{Err: err, ID: tr.IDString()}
	}

	if tr.Status() == 1 {
		logger.Debug("status equal 1 - do calculation, status", tr.Status())

		// working with procurement - calculating and updating.
		entity, err := s.paymentService.GetTx(tx, tr.Fields.GetID())
		if err != nil {
			return fmt.Errorf("could not get procurement service: [%w]", err)
		}

		if err := s.paymentService.Handle(tx, entity, s.moduleID); err != nil {
			return fmt.Errorf("could not handle procurement: [%w]", err)
		}

		if entity.DocumentNo == "" {
			if err := s.docSvc.Handle(tx, entity.ID, entity.DocumentID, entity.DocumentNo); err != nil {
				return &ErrDoc{Err: err, ID: entity.ID, DocID: entity.DocumentID, DocNo: entity.DocumentNo}
			}
		}

	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("could not commit transaction: [%w]", err)
	}

	return nil
}
