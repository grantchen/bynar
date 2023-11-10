package service

import (
	"fmt"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/i18n"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/logger"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
	"strconv"
	"strings"
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
	templateData := map[string]string{
		"Type": string(tr.Fields.GetActionType()),
	}
	return false, i18n.TranslationI18n(language, "UndefinedActionType", templateData)
}

func (s *approvalCashPaymentService) checkActionAdded(tr *treegrid.MainRow, accountID int, language string) (bool, error) {
	logger.Debug("check added action")

	docID, ok := tr.Fields.GetValInt("document_id")
	if !ok {
		templateData := map[string]string{
			"DocumentId": string(tr.Fields.GetActionType()),
		}
		return false, i18n.TranslationI18n(language, "missing", templateData)
	}

	if docID <= 0 {
		return false, i18n.TranslationI18n(language, "ValidateOnPositiveNumber", map[string]string{"Field": "document_id"})
	}

	wrkItem, err := s.storage.GetWorkflowItem(accountID, docID)
	if err != nil {
		templateData := map[string]string{
			"AccountId":  fmt.Sprintf("%d", accountID),
			"DocumentId": fmt.Sprintf("%d", docID),
			"Field":      "current-workflow-item",
		}
		return false, i18n.TranslationI18n(language, "FailedToFrom", templateData)
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

	// nothing changed
	if len(tr.Fields.UpdatedFields()) == 0 {
		return true, nil
	}

	currStatus, err := s.storage.GetStatus(tr.Fields.GetID())
	if err != nil {
		templateData := map[string]string{
			"Id":    fmt.Sprintf("%s", tr.Fields.GetID()),
			"Field": "status",
		}
		return false, i18n.TranslationI18n(language, "FailedToGetBy", templateData)
	}

	// can be updated only lines
	if currStatus == 1 && len(tr.Fields.UpdatedFields()) > 0 {
		if len(tr.Fields.UpdatedFields()) > 0 {
			return false, i18n.TranslationI18n(language, "OnlyLinesCanBeUpdated", nil)
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
		templateData := map[string]string{
			"Id":    fmt.Sprintf("%s", tr.Fields.GetID()),
			"Field": "document_id",
		}
		return false, i18n.TranslationI18n(language, "FailedToGetBy", templateData)
	}

	nextWrkItem, err := s.storage.GetWorkflowItem(accountID, newDocID)
	if err != nil {
		templateData := map[string]string{
			"AccountId":  fmt.Sprintf("%d", accountID),
			"DocumentId": fmt.Sprintf("%d", newDocID),
			"Field":      "next-workflow-item",
		}
		return false, i18n.TranslationI18n(language, "FailedToFrom", templateData)
	}

	currentWrkItem, err := s.storage.GetWorkflowItem(accountID, currDocID)
	if err != nil {
		templateData := map[string]string{
			"AccountId":  fmt.Sprintf("%d", accountID),
			"DocumentId": fmt.Sprintf("%d", newDocID),
			"Field":      "current-workflow-item",
		}
		return false, i18n.TranslationI18n(language, "FailedToFrom", templateData)
	}

	logger.Debug("Current approval_order", currentWrkItem.ApprovalOrder, "Next apploval_order", nextWrkItem.ApprovalOrder)

	if (nextWrkItem.ApprovalOrder - currentWrkItem.ApprovalOrder) != 1 {
		return false, fmt.Errorf("invalid approval order, current: %d, got: %d", currentWrkItem.ApprovalOrder, nextWrkItem.ApprovalOrder)
	}

	// No updates allowed. Only document_id can be modified
	if len(tr.Items) > 0 {
		logger.Debug(tr.Items)
		templateData := map[string]string{
			"Id": fmt.Sprintf("%d", len(tr.Items)),
		}
		return false, i18n.TranslationI18n(language, "OnlyDocumentIdCanUpdated", templateData)
	}

	updatedFields := tr.Fields.UpdatedFields()
	if len(updatedFields) != 1 {
		templateData := map[string]string{
			"Field": strings.Join(updatedFields, " "),
		}
		return false, i18n.TranslationI18n(language, "OnlyDocumentIdCanUpdatedReceived", templateData)
	}

	for _, v := range updatedFields {
		if v == "document_id" || v == "status" {
			continue
		}
		templateData := map[string]string{
			"Field": strings.Join(updatedFields, " "),
		}
		return false, i18n.TranslationI18n(language, "OnlyDocumentIdCanUpdatedReceived", templateData)
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
		templateData := map[string]string{
			"Id":    fmt.Sprintf("%s", tr.Fields.GetID()),
			"Field": "status",
		}
		return false, i18n.TranslationI18n(language, "FailedToGetBy", templateData)
	}

	if status == 1 {
		return false, i18n.TranslationI18n(language, "RecordCanNotDeleteWithStatus", map[string]string{"Status": strconv.Itoa(status)})
	}

	return true, nil
}
