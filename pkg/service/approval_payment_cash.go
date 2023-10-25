package service

import (
	"errors"
	"fmt"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/logger"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

type approvalCashPaymentService struct {
	storage ApprovalStorage
}

func NewApprovalCashPaymentService(storage ApprovalStorage) ApprovalCashPaymentService {
	return &approvalCashPaymentService{storage}
}

func (s *approvalCashPaymentService) Check(tr *treegrid.MainRow, accountID int) (bool, error) {
	logger.Debug("check", accountID, tr.Fields.GetActionType())
	switch tr.Fields.GetActionType() {
	case treegrid.GridRowActionAdd:
		return s.checkActionAdded(tr, accountID)
	case treegrid.GridRowActionChanged:
		return s.checkActionUpdated(tr, accountID)
	case treegrid.GridRowActionDeleted:
		return s.checkActionDeleted(tr)
	}

	return false, errors.New("undefined action type :" + string(tr.Fields.GetActionType()))
}

func (s *approvalCashPaymentService) checkActionAdded(tr *treegrid.MainRow, accountID int) (bool, error) {
	logger.Debug("check added action")

	docID, ok := tr.Fields.GetValInt("document_id")
	if !ok {
		return false, errors.New("missing document_id")
	}

	wrkItem, err := s.storage.GetWorkflowItem(accountID, docID)
	if err != nil {
		return false, fmt.Errorf("get workflow item: [%w]", err)
	}
	tr.Fields["status"] = wrkItem.Status

	// all operations are allowed
	if wrkItem.Status == 0 {
		return true, nil
	}

	// when adding new transfer then approval order must be 1
	return wrkItem.ApprovalOrder == 1, nil
}

func (s *approvalCashPaymentService) checkActionUpdated(tr *treegrid.MainRow, accountID int) (bool, error) {
	logger.Debug("row id", tr.Fields.GetID())
	currStatus, err := s.storage.GetStatus(tr.Fields.GetID())
	if err != nil {
		return false, fmt.Errorf("get status by id: %v, [%w]", tr.Fields.GetID(), err)
	}

	// can be updated only lines
	if currStatus == 1 && len(tr.Fields.UpdatedFields()) > 0 {
		if len(tr.Fields.UpdatedFields()) > 0 {
			return false, errors.New("current status is 1. Only lines can be updated")
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
		return false, fmt.Errorf("get doc id: [%w]", err)
	}

	nextWrkItem, err := s.storage.GetWorkflowItem(accountID, newDocID)
	if err != nil {
		return false, fmt.Errorf("get next workflow items: [%w], accountID: %d, docID: %d", err, accountID, newDocID)
	}

	currentWrkItem, err := s.storage.GetWorkflowItem(accountID, currDocID)
	if err != nil {
		return false, fmt.Errorf("get current workflow item: [%w], accountID: %d, docID: %d", err, accountID, currDocID)
	}

	logger.Debug("Current approval_order", currentWrkItem.ApprovalOrder, "Next apploval_order", nextWrkItem.ApprovalOrder)

	if (nextWrkItem.ApprovalOrder - currentWrkItem.ApprovalOrder) != 1 {
		return false, fmt.Errorf("invalid approval order, current: %d, got: %d", currentWrkItem.ApprovalOrder, nextWrkItem.ApprovalOrder)
	}

	// No updates allowed. Only document_id can be modified
	if len(tr.Items) > 0 {
		logger.Debug(tr.Items)

		return false, fmt.Errorf("only document_id can be updated. Len items: %d", len(tr.Items))
	}

	updatedFields := tr.Fields.UpdatedFields()
	if len(updatedFields) != 1 {
		return false, fmt.Errorf("only document_id can be updated. Received fields to update: [+%v]", updatedFields)
	}

	for _, v := range updatedFields {
		if v == "document_id" || v == "status" {
			continue
		}

		return false, fmt.Errorf("only document_id can be updated. Received fields to update: [+%v]", updatedFields)
	}

	tr.Fields["status"] = nextWrkItem.Status

	return true, nil
}

func (s *approvalCashPaymentService) checkActionDeleted(tr *treegrid.MainRow) (bool, error) {
	status, err := s.storage.GetStatus(tr.Fields.GetID())
	if err != nil {
		return false, err
	}

	return status != 1, nil
}
