package service

import (
	"context"
	"database/sql"
	"fmt"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/errors"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/i18n"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/logger"
	pkgservice "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/service"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/procurements/internal/repository"
)

type UploadService struct {
	db                  *sql.DB
	language            string
	approvalService     pkgservice.ApprovalCashPaymentService
	docSvc              pkgservice.DocumentService
	gridRowRepository   treegrid.GridRowRepositoryWithChild
	accountId           int
	procurementsService ProcurementsService
}

func NewUploadService(db *sql.DB,
	language string,
	approvalService pkgservice.ApprovalCashPaymentService,
	docSvc pkgservice.DocumentService,
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

// Handle upload handle
func (s *UploadService) Handle(req *treegrid.PostRequest) (*treegrid.PostResponse, error) {
	trList, err := treegrid.ParseRequestUpload(req, s.gridRowRepository)
	if err != nil {
		return nil, fmt.Errorf("parse requst: [%w]", err)
	}

	resp := treegrid.HandleTreegridWithChildMainRowsLines(
		trList,
		func(mr *treegrid.MainRow) error {
			return i18n.TranslationErrorToI18n(s.language, s.handle(mr))
		},
	)

	return resp, nil
}

func (s *UploadService) handle(tr *treegrid.MainRow) error {
	switch tr.Fields.GetActionType() {
	case treegrid.GridRowActionAdd:
		err := tr.Fields.ValidateOnRequiredAll(repository.ProcurementFieldNames, s.language)
		if err != nil {
			return err
		}
		err = tr.Fields.ValidateOnNotNegativeNumber(repository.ProcurementFieldNames, s.language)
		if err != nil {
			return err
		}
	case treegrid.GridRowActionChanged:
		err := tr.Fields.ValidateOnNotNegativeNumber(repository.ProcurementFieldNames, s.language)
		if err != nil {
			return err
		}
	}
	for _, item := range tr.Items {
		switch item.GetActionType() {
		case treegrid.GridRowActionAdd:
			err := item.ValidateOnNotNegativeNumber(repository.ProcurementLineFieldNames, s.language)
			if err != nil {
				return err
			}
			err = s.procurementsService.ValidateParams(s.db, item)
			if err != nil {
				return err
			}
		case treegrid.GridRowActionChanged:
			err := item.ValidateOnNotNegativeNumber(repository.ProcurementLineFieldNames, s.language)
			if err != nil {
				return err
			}
			err = s.procurementsService.ValidateParams(s.db, item)
			if err != nil {
				return err
			}
		}
	}
	// Check Approval Order
	ok, err := s.approvalService.Check(tr, s.accountId, s.language)
	if err != nil {
		return fmt.Errorf("%s: [%w]", i18n.TranslationI18n(s.language, "CheckOrder", map[string]string{}).Error(), err)
	}
	if !ok {
		return fmt.Errorf("invalid approval order: [%w]: status: %d", errors.ErrForbiddenAction, tr.Status())
	}
	// Create new transaction
	tx, err := s.db.BeginTx(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("begin transaction: [%w]", err)
	}
	defer func(tx *sql.Tx) {
		_ = tx.Rollback()
	}(tx)

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
