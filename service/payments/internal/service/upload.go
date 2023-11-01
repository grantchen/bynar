package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/payments/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/errors"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/i18n"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/logger"
	pkg_service "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/service"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
	"log"
	"strings"
)

type UploadService struct {
	db                                 *sql.DB
	updateGRPaymentRepository          treegrid.SimpleGridRowRepository
	updateGRPaymentRepositoryWithChild treegrid.GridRowRepositoryWithChild
	updateGRPaymentLineRepository      treegrid.SimpleGridRowRepository
	language                           string
	approvalService                    pkg_service.ApprovalCashPaymentService
	docSvc                             pkg_service.DocumentService
	accountId                          int
	paymentService                     PaymentService
}

func NewUploadService(db *sql.DB,
	updateGRPaymentRepository treegrid.SimpleGridRowRepository,
	updateGRPaymentRepositoryWithChild treegrid.GridRowRepositoryWithChild,
	updateGRPaymentLineRepository treegrid.SimpleGridRowRepository,
	language string,
	approvalService pkg_service.ApprovalCashPaymentService,
	docSvc pkg_service.DocumentService,
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

func (u *UploadService) Handle(req *treegrid.PostRequest) (*treegrid.PostResponse, error) {
	resp := &treegrid.PostResponse{}
	// Create new transaction
	b, _ := json.Marshal(req)
	logger.Debug("request: ", string(b))
	trList, err := treegrid.ParseRequestUpload(req, u.updateGRPaymentRepositoryWithChild)

	if err != nil {
		return nil, fmt.Errorf("parse request: [%w]", err)
	}

	tx, err := u.db.BeginTx(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf(i18n.Localize(u.language, errors.ErrCodeBeginTransaction))
	}
	defer tx.Rollback()
	m := make(map[string]interface{}, 0)
	var handleErr error
	for _, tr := range trList.MainRows() {
		// todo refactor group
		if tr.Fields["id"] == "Group" {
			continue
		}

		if handleErr = u.handle(tx, tr); handleErr != nil {
			log.Println("Err", handleErr)

			resp.IO.Result = -1
			resp.IO.Message += handleErr.Error() + "\n"
			resp.Changes = append(resp.Changes, treegrid.GenMapColorChangeError(tr.Fields))
			break
		}
		resp.Changes = append(resp.Changes, tr.Fields)
		resp.Changes = append(resp.Changes, treegrid.GenMapColorChangeSuccess(tr.Fields))
		resp.Changes = append(resp.Changes, m)

		for k := range tr.Items {
			resp.Changes = append(resp.Changes, tr.Items[k])
			resp.Changes = append(resp.Changes, treegrid.GenMapColorChangeSuccess(tr.Items[k]))
		}
	}
	if handleErr == nil {
		if err := tx.Commit(); err != nil {
			return nil, fmt.Errorf("%s: [%w]", i18n.Localize(u.language, "commit-transaction"), err)
		}
	}
	return resp, nil
}

func (u *UploadService) handle(tx *sql.Tx, tr *treegrid.MainRow) error {
	//Check Approval Order
	ok, err := u.approvalService.Check(tr, u.accountId, u.language)
	if err != nil {
		return err
	}

	if !ok {
		return fmt.Errorf("%s",
			i18n.Localize(u.language, "forbidden-action"))
	}
	if err := u.save(tx, tr); err != nil {
		return err
	}

	if tr.Status() == 1 {
		logger.Debug("status equal 1 - do calculation, status", tr.Status())

		// working with procurement - calculating and updating.
		entity, err := u.paymentService.GetTx(tx, tr.Fields.GetID())
		if err != nil {
			return fmt.Errorf("%s: [%w]", i18n.Localize(u.language, "failed-to-get-data-from", "procurement"), err)
		}

		if err := u.paymentService.Handle(tx, entity); err != nil {
			return fmt.Errorf("%s: [%w]", i18n.Localize(u.language, "failed-to-handle", "procurement"), err)
		}
		if entity.DocumentNo == "" {
			if err := u.docSvc.Handle(tx, entity.ID, entity.DocumentID, entity.DocumentNo); err != nil {
				return fmt.Errorf("%s: [%w]", i18n.Localize(u.language, "failed-to-handle", "document"), err)
			}
		}

	}

	return nil
}

func (u *UploadService) save(tx *sql.Tx, tr *treegrid.MainRow) error {
	if err := u.savePayment(tx, tr); err != nil {
		return fmt.Errorf("%s %s: [%w]",
			i18n.Localize(u.language, errors.ErrCodeSave),
			i18n.Localize(u.language, errors.ErrCodePayment),
			err)
	}

	if err := u.savePaymentLine(tx, tr, tr.Fields.GetID()); err != nil {
		return fmt.Errorf("%s %s: [%w]",
			i18n.Localize(u.language, errors.ErrCodeSave),
			i18n.Localize(u.language, errors.ErrCodePaymentLine),
			err)
	}

	return nil
}

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
				return fmt.Errorf("%s: %s: %s", field, i18n.Localize(u.language, errors.ErrCodeValueDuplicated), tr.Fields[field])
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
				return fmt.Errorf("%s: %s: %s", field, i18n.Localize(u.language, errors.ErrCodeValueDuplicated), tr.Fields[field])
			}
		}
	case treegrid.GridRowActionDeleted:
		// ignore id start with CR
		idStr := tr.Fields.GetIDStr()
		if !strings.HasPrefix(idStr, "CR") {
			stmt, err := tx.Prepare("DELETE FROM payment_lines WHERE parent_id = ?")
			if err != nil {
				return fmt.Errorf("%s: [%w]", i18n.Localize(u.language, "failed-to-delete", "payment_lines"), err)
			}

			defer stmt.Close()

			_, err = stmt.Exec(idStr)
			if err != nil {
				return fmt.Errorf("%s: [%w]", i18n.Localize(u.language, "failed-to-delete", "payment_lines"), err)
			}
		}

		fmt.Println(tr.Fields.GetID())
	}

	return u.updateGRPaymentRepositoryWithChild.SaveMainRow(tx, tr)
}

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
				return fmt.Errorf("%s: [%w]", i18n.Localize(u.language, "failed-to-add", "payment-line"), err)
			}
		case treegrid.GridRowActionChanged:
			err = tr.Fields.ValidateOnRequired(repository.PaymentLineFieldUploadNames, u.language)
			if err != nil {
				return err
			}

			err = u.updateGRPaymentRepositoryWithChild.SaveLineUpdate(tx, item)
			if err != nil {
				return fmt.Errorf("%s: [%w]", i18n.Localize(u.language, "failed-to-update", "payment-line"), err)
			}
		case treegrid.GridRowActionDeleted:
			logger.Debug("delete child")

			// re-assign user_group_lines id
			item["id"] = item.GetID()
			err = u.updateGRPaymentRepositoryWithChild.SaveLineDelete(tx, item)
			if err != nil {
				return fmt.Errorf("%s: [%w]", i18n.Localize(u.language, "failed-to-delete", "payment-line"), err)
			}
		default:
			return fmt.Errorf("%s: %s", i18n.Localize(u.language, "undefined-action-type"), tr.Fields.GetActionType())

		}
	}
	return nil
}
