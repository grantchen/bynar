package service

import (
	"database/sql"
	"fmt"
	"strings"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/payments/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/i18n"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/logger"
	pkgservice "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/service"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/utils"
)

// UploadService is the service for upload
type UploadService struct {
	db                                 *sql.DB
	updateGRPaymentRepository          treegrid.SimpleGridRowRepository
	updateGRPaymentRepositoryWithChild treegrid.GridRowRepositoryWithChild
	updateGRPaymentLineRepository      treegrid.SimpleGridRowRepository
	language                           string
	approvalService                    pkgservice.ApprovalCashPaymentService
	docSvc                             pkgservice.DocumentService
	accountId                          int
	paymentService                     PaymentService
}

// NewUploadService returns a new upload service
func NewUploadService(db *sql.DB,
	updateGRPaymentRepository treegrid.SimpleGridRowRepository,
	updateGRPaymentRepositoryWithChild treegrid.GridRowRepositoryWithChild,
	updateGRPaymentLineRepository treegrid.SimpleGridRowRepository,
	language string,
	approvalService pkgservice.ApprovalCashPaymentService,
	docSvc pkgservice.DocumentService,
	accountId int,
	paymentService PaymentService,
) *UploadService {
	return &UploadService{
		db:                                 db,
		updateGRPaymentRepository:          updateGRPaymentRepository,
		updateGRPaymentRepositoryWithChild: updateGRPaymentRepositoryWithChild,
		updateGRPaymentLineRepository:      updateGRPaymentLineRepository,
		language:                           language,
		approvalService:                    approvalService,
		docSvc:                             docSvc,
		accountId:                          accountId,
		paymentService:                     paymentService,
	}
}

// Handle handles the upload request
func (u *UploadService) Handle(req *treegrid.PostRequest) (*treegrid.PostResponse, error) {
	trList, err := treegrid.ParseRequestUpload(req, u.updateGRPaymentRepositoryWithChild)
	if err != nil {
		return nil, fmt.Errorf("parse request: [%w]", err)
	}

	resp := treegrid.HandleTreegridWithChildMainRowsLines(
		trList,
		func(mr *treegrid.MainRow) error {
			err = utils.WithTransaction(u.db, func(tx *sql.Tx) error {
				return u.handle(tx, mr)
			})
			return i18n.TranslationErrorToI18n(u.language, err)
		},
	)

	return resp, nil
}

// handle handles the upload request of a main row
func (u *UploadService) handle(tx *sql.Tx, tr *treegrid.MainRow) error {
	//Check Approval Order
	ok, err := u.approvalService.Check(tr, u.accountId, u.language)
	if err != nil {
		return err
	}

	if !ok {
		return i18n.TranslationI18n(u.language, "ForbiddenAction", map[string]string{
			"Message": err.Error(),
		})
	}
	if err := u.save(tx, tr); err != nil {
		return err
	}

	if tr.Status() == 1 {
		logger.Debug("status equal 1 - do calculation, status", tr.Status())

		// working with procurement - calculating and updating.
		entity, err := u.paymentService.GetTx(tx, tr.Fields.GetID())
		if err != nil {
			return i18n.TranslationI18n(u.language, "FailedToGetDataFrom", map[string]string{
				"Message": err.Error(),
			})
		}

		if err := u.paymentService.Handle(tx, entity); err != nil {
			templateData := map[string]string{
				"Field": "procurement",
			}
			return i18n.TranslationI18n(u.language, "FailedToHandle", templateData)
		}
		if entity.DocumentNo == "" {
			if err := u.docSvc.Handle(tx, entity.ID, entity.DocumentID, entity.DocumentNo); err != nil {
				templateData := map[string]string{
					"Field": "document",
				}
				return i18n.TranslationI18n(u.language, "FailedToHandle", templateData)
			}
		}

	}

	return nil
}

// save saves payment and payment lines
func (u *UploadService) save(tx *sql.Tx, tr *treegrid.MainRow) error {
	if err := u.savePayment(tx, tr); err != nil {
		return i18n.TranslationI18n(u.language, "SavePayment", map[string]string{
			"Message": err.Error(),
		})
	}

	if err := u.savePaymentLine(tx, tr, tr.Fields.GetID()); err != nil {
		return i18n.TranslationI18n(u.language, "SavePaymentLine", map[string]string{
			"Message": err.Error(),
		})
	}

	return nil
}

// savePayment saves payment
func (u *UploadService) savePayment(tx *sql.Tx, tr *treegrid.MainRow) error {
	fieldsValidating := []string{"document_id"}

	var err error
	switch tr.Fields.GetActionType() {
	case treegrid.GridRowActionAdd:
		err = tr.Fields.ValidateOnRequiredAll(repository.PaymentFieldUploadNames, u.language)
		if err != nil {
			return err
		}
		tr.Fields["status"] = 0
		tr.Fields["paid"] = 0
		tr.Fields["remaining"] = 0
		tr.Fields["paid_status"] = 0
		for _, field := range fieldsValidating {
			ok, err := u.updateGRPaymentRepository.ValidateOnIntegrity(tx, tr.Fields, []string{field})
			if !ok || err != nil {
				templateData := map[string]string{
					"Field": field,
				}
				return i18n.TranslationI18n(u.language, "ValueDuplicated", templateData)
			}
		}
	case treegrid.GridRowActionChanged:
		err = tr.Fields.ValidateOnRequired(repository.PaymentFieldNames, u.language)
		if err != nil {
			return err
		}

		for _, field := range fieldsValidating {
			ok, err := u.updateGRPaymentRepository.ValidateOnIntegrity(tx, tr.Fields, []string{field})
			if !ok || err != nil {
				templateData := map[string]string{
					"Field": field,
				}
				return i18n.TranslationI18n(u.language, "ValueDuplicated", templateData)
			}
		}
	case treegrid.GridRowActionDeleted:
		// ignore id start with CR
		idStr := tr.Fields.GetIDStr()
		if !strings.HasPrefix(idStr, "CR") {
			stmt, err := tx.Prepare("DELETE FROM payment_lines WHERE parent_id = ?")
			if err != nil {
				templateData := map[string]string{
					"Field": "payment_lines",
				}
				return i18n.TranslationI18n(u.language, "FailedToDelete", templateData)
			}

			defer func(stmt *sql.Stmt) {
				_ = stmt.Close()
			}(stmt)

			_, err = stmt.Exec(idStr)
			if err != nil {
				templateData := map[string]string{
					"Field": "payment_lines",
				}
				return i18n.TranslationI18n(u.language, "FailedToDelete", templateData)
			}
		}
	}

	return u.updateGRPaymentRepositoryWithChild.SaveMainRow(tx, tr)
}

// savePaymentLine saves payment lines
func (u *UploadService) savePaymentLine(tx *sql.Tx, tr *treegrid.MainRow, parentID interface{}) error {
	for _, item := range tr.Items {
		logger.Debug("save payment line: ", tr, "parentID: ", parentID)

		var err error
		switch item.GetActionType() {
		case treegrid.GridRowActionAdd:
			err = item.ValidateOnRequiredAll(repository.PaymentLineFieldUploadNames, u.language)
			if err != nil {
				return err
			}

			logger.Debug("add child row")
			err = u.updateGRPaymentRepositoryWithChild.SaveLineAdd(tx, item)
			if err != nil {
				templateData := map[string]string{
					"Field": "payment-line",
				}
				return i18n.TranslationI18n(u.language, "FailedToAdd", templateData)
			}
		case treegrid.GridRowActionChanged:
			err = tr.Fields.ValidateOnRequired(repository.PaymentLineFieldUploadNames, u.language)
			if err != nil {
				return err
			}

			err = u.updateGRPaymentRepositoryWithChild.SaveLineUpdate(tx, item)
			if err != nil {
				templateData := map[string]string{
					"Field": "payment-line",
				}
				return i18n.TranslationI18n(u.language, "FailedToUpdate", templateData)
			}
		case treegrid.GridRowActionDeleted:
			logger.Debug("delete child")

			// re-assign user_group_lines id
			item["id"] = item.GetID()
			err = u.updateGRPaymentRepositoryWithChild.SaveLineDelete(tx, item)
			if err != nil {
				templateData := map[string]string{
					"Field": "payment-line",
				}
				return i18n.TranslationI18n(u.language, "FailedToDelete", templateData)
			}
		default:
			return i18n.TranslationI18n(u.language, "UndefinedActionType", map[string]string{
				"Message": err.Error(),
			})

		}
	}
	return nil
}
