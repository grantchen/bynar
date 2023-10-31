package service

import (
	"fmt"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/i18n"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/logger"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
	"strings"

	errpkg "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/errors"
)

type approvalCashPaymentService struct {
	storage ApprovalStorage
}

func NewApprovalCashPaymentService(storage ApprovalStorage) ApprovalCashPaymentService {
	return &approvalCashPaymentService{storage}
}

func (s *approvalCashPaymentService) Check(tr *treegrid.MainRow, accountID int, language string) (bool, error) {
	logger.Debug("check", accountID, tr.Fields.GetActionType())
	switch tr.Fields.GetActionType() {
	case treegrid.GridRowActionAdd:
		return s.checkActionAdded(tr, accountID, language)
	case treegrid.GridRowActionChanged:
		strID := tr.Fields.GetIDStr()
		if !strings.HasPrefix(strID, "AR") {
			return s.checkActionUpdated(tr, accountID, language)
		}
		return true, nil
	case treegrid.GridRowActionDeleted:
		strID := tr.Fields.GetIDStr()
		if !strings.HasPrefix(strID, "CR") {
			return s.checkActionDeleted(tr, language)
		}
		return true, nil
	}

	return false, fmt.Errorf("%s : %s",
		i18n.Localize(language, errpkg.ErrCodeUndefinedActionType),
		string(tr.Fields.GetActionType()))
}

func (s *approvalCashPaymentService) checkActionAdded(tr *treegrid.MainRow, accountID int, language string) (bool, error) {
	logger.Debug("check added action")

	docID, ok := tr.Fields.GetValInt("document_id")
	if !ok {
		return false, fmt.Errorf("%s : %s",
			i18n.Localize(language, "missing"), "document_id")
	}

	wrkItem, err := s.storage.GetWorkflowItem(accountID, docID)
	if err != nil {
		return false, fmt.Errorf("%s",
			i18n.Localize(language, "failed-to-get-by", i18n.Localize(language, "workflow-item"), "document_id"))
	}
	tr.Fields["status"] = wrkItem.Status

	// all operations are allowed
	if wrkItem.Status == 0 {
		return true, nil
	}

	// when adding new transfer then approval order must be 1
	return wrkItem.ApprovalOrder == 1, nil
}

func (s *approvalCashPaymentService) checkActionUpdated(tr *treegrid.MainRow, accountID int, language string) (bool, error) {
	logger.Debug("row id", tr.Fields.GetID())
	currStatus, err := s.storage.GetStatus(tr.Fields.GetID())
	if err != nil {
		return false, fmt.Errorf("%s : %v, [%w]", i18n.Localize(language, "failed-to-get-by", "status", "id"), tr.Fields.GetID(), err)
	}

	// can be updated only lines
	if currStatus == 1 && len(tr.Fields.UpdatedFields()) > 0 {
		if len(tr.Fields.UpdatedFields()) > 0 {
			return false, fmt.Errorf("%s.%s", i18n.Localize(language, "data-is", i18n.Localize(language, "current-status"), "1"), i18n.Localize(language, "only-lines-can-be-updated"))
		}
		tr.Fields["status"] = currStatus

		return true, nil
	}

	// if newDocID isn't in update set, then status stays the same. This kind updates only can be with status = 0
	newDocID, ok := tr.Fields.GetValInt("document_id")
	if !ok {
		logger.Debug("document_id isn't updated")

		// document_id hasn't change - valid only with status = 0
		// other type of updates aren't allowed
		return currStatus == 0, nil
	}

	logger.Debug("new document_id", newDocID)

	currDocID, err := s.storage.GetDocID(tr.Fields.GetID())
	if err != nil {
		return false, fmt.Errorf("%s : %v, [%w]", i18n.Localize(language, "failed-to-get-by", "document_id", "id"), tr.Fields.GetID(), err)
	}

	nextWrkItem, err := s.storage.GetWorkflowItem(accountID, newDocID)
	if err != nil {
		return false, fmt.Errorf("%s : %s:%d,%s:%d,[%w]",
			i18n.Localize(language, "failed-to-get-data-from", "next-workflow-item"), "account_id", accountID, "document_id", newDocID, err)
	}

	currentWrkItem, err := s.storage.GetWorkflowItem(accountID, currDocID)
	if err != nil {
		return false, fmt.Errorf("%s : %s:%d,%s:%d,[%w]",
			i18n.Localize(language, "failed-to-get-data-from", "current-workflow-item"), "account_id", accountID, "document_id", currDocID, err)
	}

	logger.Debug("Current approval_order", currentWrkItem.ApprovalOrder, "Next apploval_order", nextWrkItem.ApprovalOrder)

	if (nextWrkItem.ApprovalOrder - currentWrkItem.ApprovalOrder) != 1 {
		return false, fmt.Errorf("invalid approval order, current: %d, got: %d", currentWrkItem.ApprovalOrder, nextWrkItem.ApprovalOrder)
	}

	// No updates allowed. Only document_id can be modified
	if len(tr.Items) > 0 {
		logger.Debug(tr.Items)

		return false, fmt.Errorf("%s. %s: %d", i18n.Localize(language, "only-document_id-can-be-updated"), i18n.Localize(language, "items-length"), len(tr.Items))
	}

	updatedFields := tr.Fields.UpdatedFields()
	if len(updatedFields) != 1 {
		return false, fmt.Errorf("%s. %s: [+%v]", i18n.Localize(language, "only-document_id-can-be-updated"), i18n.Localize(language, "received-fields-to-update"), updatedFields)
	}

	for _, v := range updatedFields {
		if v == "document_id" || v == "status" {
			continue
		}

		return false, fmt.Errorf("%s. %s: [+%v]", i18n.Localize(language, "only-document_id-can-be-updated"), i18n.Localize(language, "received-fields-to-update"), updatedFields)
	}

	tr.Fields["status"] = nextWrkItem.Status

	return true, nil
}

func (s *approvalCashPaymentService) checkActionDeleted(tr *treegrid.MainRow, language string) (bool, error) {
	status, err := s.storage.GetStatus(tr.Fields.GetID())
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return true, nil
		}
		return false, fmt.Errorf("%s : %v, [%w]", i18n.Localize(language, "failed-to-get-by", "status", "id"), tr.Fields.GetID(), err)
	}

	return status != 1, nil
}
